package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hibiken/asynq"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/zekrotja/ken"

	"github.com/jaczerob/madamchuckle/internal/inits"
	"github.com/jaczerob/madamchuckle/internal/services/database"
	"github.com/jaczerob/madamchuckle/internal/services/metrics"
	"github.com/jaczerob/madamchuckle/internal/static"
	"github.com/jaczerob/madamchuckle/pkg/toontown"

	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "starts the Madam Chuckle Discord Bot instance",
	Run: func(_ *cobra.Command, _ []string) {
		diBuilder, _ := di.NewBuilder()

		diBuilder.Add(di.Def{
			Name: static.DiToontownClient,
			Build: func(ctn di.Container) (interface{}, error) {
				return toontown.New(), nil
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiDiscordSession,
			Build: func(ctn di.Container) (interface{}, error) {
				return inits.InitDiscordBotSession(ctn)
			},
			Close: func(obj interface{}) error {
				log.Info("madam chuckle leaving toontown... :(")

				session := obj.(*discordgo.Session)
				return session.Close()
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiCommandHandler,
			Build: func(ctn di.Container) (interface{}, error) {
				return inits.InitCommandHandler(ctn)
			},
			Close: func(obj interface{}) error {
				k := obj.(*ken.Ken)
				return k.Unregister()
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiDatabase,
			Build: func(ctn di.Container) (interface{}, error) {
				return database.New()
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiTaskScheduler,
			Build: func(ctn di.Container) (interface{}, error) {
				return inits.InitTaskScheduler(ctn)
			},
			Close: func(obj interface{}) error {
				s := obj.(*asynq.Scheduler)
				s.Shutdown()
				return nil
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiTaskServer,
			Build: func(ctn di.Container) (interface{}, error) {
				return inits.InitTaskServer(ctn)
			},
			Close: func(obj interface{}) error {
				s := obj.(*asynq.Server)
				s.Shutdown()
				return nil
			},
		})

		diBuilder.Add(di.Def{
			Name: static.DiMetricsServer,
			Build: func(ctn di.Container) (interface{}, error) {
				return inits.InitMetricsServer()
			},
			Close: func(obj interface{}) error {
				s := obj.(*metrics.MetricsServer)
				return s.Shutdown()
			},
		})

		ctn := diBuilder.Build()
		defer func() {
			if err := ctn.DeleteWithSubContainers(); err != nil {
				log.WithError(err).Fatal("failed to close containers")
			}
		}()

		ctn.Get(static.DiToontownClient)
		_, err := ctn.SafeGet(static.DiDatabase)
		if err != nil {
			log.WithError(err).Fatal("error initializing database")
		}

		session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
		_, err = ctn.SafeGet(static.DiCommandHandler)
		if err != nil {
			log.WithError(err).Fatal("error initializing command handler")
		}

		if err = session.Open(); err != nil {
			log.WithError(err).Fatal("error opening discord session")
		}

		_, err = ctn.SafeGet(static.DiMetricsServer)
		if err != nil {
			log.WithError(err).Fatal("error starting metrics server")
		}

		ctn.Get(static.DiTaskScheduler)
		ctn.Get(static.DiTaskServer)

		log.Info("madam chuckle has entered toontown, press CTRL+C to leave")

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	},
}
