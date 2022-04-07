package inits

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/viper"

	"github.com/jaczerob/madamchuckle/internal/static"
	"github.com/jaczerob/madamchuckle/internal/tasks"

	log "github.com/sirupsen/logrus"
)

func InitTaskScheduler(ctn di.Container) (s *asynq.Scheduler, err error) {
	loc, err := time.LoadLocation(static.AsynqLocation)
	if err != nil {
		return
	}

	opts := &asynq.SchedulerOpts{
		Logger:   log.StandardLogger(),
		Location: loc,
	}

	redis := asynq.RedisClientOpt{Addr: viper.GetString("redis.addr")}
	s = asynq.NewScheduler(redis, opts)

	// GAME STATUS TASK
	gamestatusTask, err := tasks.NewGameStatusTask()
	if err != nil {
		return
	}

	info, err := s.Register("@every 10m", gamestatusTask, asynq.MaxRetry(0), asynq.ProcessIn(time.Second*10))
	if err != nil {
		return
	}

	log.Info("registered task: %q", info)

	// EVENTS TASK
	eventsTask, err := tasks.NewEventTask()
	if err != nil {
		return
	}

	info, err = s.Register("@every 10m", eventsTask, asynq.MaxRetry(0), asynq.ProcessIn(time.Second*10))
	if err != nil {
		return
	}

	log.Info("registered task: %q", info)

	err = s.Start()
	return
}
