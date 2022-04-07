package inits

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"

	"github.com/jaczerob/madamchuckle/internal/commands/admin"
	"github.com/jaczerob/madamchuckle/internal/middlewares"
	"github.com/jaczerob/madamchuckle/internal/services/database"
	"github.com/jaczerob/madamchuckle/internal/static"

	log "github.com/sirupsen/logrus"
)

func InitCommandHandler(ctn di.Container) (k *ken.Ken, err error) {
	session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	db := ctn.Get(static.DiDatabase).(*database.Database)

	k, err = ken.New(session, ken.Options{
		CommandStore:       db,
		DependencyProvider: ctn,

		OnSystemError: func(context string, err error, args ...interface{}) {
			log.WithField("ctx", context).WithError(err).Error("ken error")
		},

		OnCommandError: func(err error, ctx *ken.Ctx) {
			ctx.Defer()

			if err == ken.ErrNotDMCapable {
				ctx.FollowUpError("This command can not be used in DMs.", "")
				return
			}

			ctx.FollowUpError(
				fmt.Sprintf("The command execution failed unexpectedly:\n```\n%s\n```", err.Error()),
				"Command execution failed")
		},
	})

	if err != nil {
		return
	}

	err = k.RegisterCommands(
		new(admin.SayCommand),
		new(admin.EventCommand),
	)

	if err != nil {
		return
	}

	err = k.RegisterMiddlewares(new(middlewares.PermissionsMiddleware))
	return
}
