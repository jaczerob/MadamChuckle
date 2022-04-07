package tasks

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/hibiken/asynq"
	"github.com/sarulabs/di/v2"

	"github.com/jaczerob/madamchuckle/internal/services/toontown"
	"github.com/jaczerob/madamchuckle/internal/static"

	log "github.com/sirupsen/logrus"
)

func NewGameStatusTask() (*asynq.Task, error) {
	return asynq.NewTask(static.AsynqGamestatusUpdate, nil), nil
}

type GameStatusProcessor struct {
	session *discordgo.Session
	ttr     *toontown.ToontownClient
}

var _ asynq.Handler = (*GameStatusProcessor)(nil)

func NewGameStatusProcessor(ctn di.Container) *GameStatusProcessor {
	return &GameStatusProcessor{
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		ttr:     ctn.Get(static.DiToontownClient).(*toontown.ToontownClient),
	}
}

func (p *GameStatusProcessor) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	population, err := p.ttr.Population()
	if err != nil {
		return
	}

	status := fmt.Sprintf("with %d Toons", population.TotalPopulation)

	log.Info("updating game status: ", status)
	return p.session.UpdateGameStatus(0, status)
}
