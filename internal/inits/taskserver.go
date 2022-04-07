package inits

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/viper"

	"github.com/jaczerob/madamchuckle/internal/static"
	"github.com/jaczerob/madamchuckle/internal/tasks"

	log "github.com/sirupsen/logrus"
)

func InitTaskServer(ctn di.Container) (s *asynq.Server, err error) {
	s = asynq.NewServer(
		asynq.RedisClientOpt{Addr: viper.GetString("redis.addr")},
		asynq.Config{
			Concurrency:  10,
			ErrorHandler: asynq.ErrorHandlerFunc(logError),
			Logger:       log.StandardLogger(),
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(static.AsynqGamestatusUpdate, tasks.NewGameStatusProcessor(ctn))
	mux.Handle(static.AsynqToontownEvent, tasks.NewEventProcessor(ctn))

	err = s.Start(mux)
	return
}

func logError(ctx context.Context, task *asynq.Task, err error) {
	log.WithError(err).Error("error handling task: ", task.Type)
}
