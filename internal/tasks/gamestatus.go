package tasks

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sarulabs/di/v2"

	"github.com/jaczerob/madamchuckle/internal/services/metrics"
	"github.com/jaczerob/madamchuckle/internal/static"
	"github.com/jaczerob/madamchuckle/pkg/toontown"

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
	population, err := p.getPopulation()
	if err != nil {
		return
	}

	status := fmt.Sprintf("with %d Toons", population)

	log.Info("updating game status: ", status)
	return p.session.UpdateGameStatus(0, status)
}

func (p *GameStatusProcessor) getPopulation() (population int, err error) {
	populationData, err := p.ttr.Population()
	if err != nil {
		return
	}

	population = populationData.TotalPopulation
	floatPopulation := float64(population)

	hour := strconv.FormatInt(int64(time.Now().Hour()), 10)
	roundedPopulation := math.Round(floatPopulation/100) * 100

	metrics.PopulationTracker.Set(floatPopulation)
	metrics.PopulationByHourTracker.With(prometheus.Labels{"hour": hour}).Observe(roundedPopulation)

	return
}
