package tasks

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hibiken/asynq"
	"github.com/sarulabs/di/v2"

	"github.com/jaczerob/madamchuckle/internal/services/database"
	"github.com/jaczerob/madamchuckle/internal/services/toontown"
	"github.com/jaczerob/madamchuckle/internal/static"

	log "github.com/sirupsen/logrus"
)

func NewEventTask() (*asynq.Task, error) {
	return asynq.NewTask(static.AsynqToontownEvent, nil), nil
}

type EventProcessor struct {
	session *discordgo.Session
	ttr     *toontown.ToontownClient
	db      *database.Database
}

var _ asynq.Handler = (*EventProcessor)(nil)

func NewEventProcessor(ctn di.Container) *EventProcessor {
	return &EventProcessor{
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		ttr:     ctn.Get(static.DiToontownClient).(*toontown.ToontownClient),
		db:      ctn.Get(static.DiDatabase).(*database.Database),
	}
}

func (p *EventProcessor) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	events, err := p.db.GetEvents()
	if err != nil {
		return
	}

	for _, event := range events {
		var embed *discordgo.MessageEmbed

		switch event.ID {
		case int64(static.FieldOfficesEventID):
			fieldOffices, err := p.ttr.FieldOffices()
			if err != nil {
				return err
			}

			embed, err = fieldOffices.ToEmbed()
			if err != nil {
				return err
			}
		case int64(static.InvasionsEventID):
			invasions, err := p.ttr.Invasions()
			if err != nil {
				return err
			}

			embed, err = invasions.ToEmbed()
			if err != nil {
				return err
			}
		case int64(static.PopulationEventID):
			population, err := p.ttr.Population()
			if err != nil {
				return err
			}

			embed, err = population.ToEmbed()
			if err != nil {
				return err
			}
		default:
			continue
		}

		log.Info("updating event: ", event.ID, event.MessageID, event.ChannelID)

		_, err := p.session.ChannelMessageEditEmbed(event.MessageID, event.ChannelID, embed)
		if err != nil {
			log.WithError(err).Error("error handling event")
			err = p.db.UnregisterEvent(event.MessageID, event.ChannelID)
			return err
		}

		time.Sleep(5 * time.Second)
	}

	return
}
