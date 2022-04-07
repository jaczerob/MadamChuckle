package admin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaczerob/madamchuckle/internal/middlewares"
	"github.com/jaczerob/madamchuckle/internal/services/database"
	"github.com/jaczerob/madamchuckle/internal/static"
	"github.com/zekrotja/ken"
)

type EventCommand struct{}

var (
	_ ken.SlashCommand                      = (*EventCommand)(nil)
	_ ken.DmCapable                         = (*EventCommand)(nil)
	_ middlewares.RequiresPermissionCommand = (*EventCommand)(nil)
)

func (c *EventCommand) Name() string {
	return "events"
}

func (c *EventCommand) Description() string {
	return "Handles (un)registering messages to/from events"
}

func (c *EventCommand) Version() string {
	return "1.0.0"
}

func (c *EventCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *EventCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "register",
			Description: "Registers a message to an event",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "event",
					Required:    true,
					Description: "The event to register to",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "fieldoffices",
							Value: static.FieldOfficesEventID,
						},
						{
							Name:  "invasions",
							Value: static.InvasionsEventID,
						},
						{
							Name:  "sillymeter",
							Value: static.SillyMeterEventID,
						},
						{
							Name:  "status",
							Value: static.StatusEventID,
						},
						{
							Name:  "doodles",
							Value: static.DoodlesEventID,
						},
						{
							Name:  "population",
							Value: static.PopulationEventID,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message_id",
					Required:    true,
					Description: "The ID of the message (MUST BE SENT BY MADAM CHUCKLE)",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "unregister",
			Description: "Unregisters a message from an event",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message_id",
					Required:    true,
					Description: "The ID of the message (MUST BE SENT BY MADAM CHUCKLE)",
				},
			},
		},
	}
}

func (c *EventCommand) IsDmCapable() bool {
	return false
}

func (c *EventCommand) RequiresPermission() int64 {
	return discordgo.PermissionAdministrator
}

func (c *EventCommand) Run(ctx *ken.Ctx) (err error) {
	if ctx.Event.Message.Author != ctx.Session.State.User {
		ctx.FollowUpError("Message was not written by me!", "")
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{
			Name: "register",
			Run:  c.register,
		},
		ken.SubCommandHandler{
			Name: "unregister",
			Run:  c.unregister,
		},
	)

	return
}

func (c *EventCommand) register(ctx *ken.SubCommandCtx) (err error) {
	eventID := ctx.Options().GetByName("event").IntValue()
	messageID := ctx.Options().GetByName("messageID").StringValue()
	channel, err := ctx.Channel()
	if err != nil {
		return
	}

	channelID := channel.ID

	db := ctx.Get(static.DiDatabase).(*database.Database)
	if err = db.RegisterEvent(eventID, messageID, channelID); err != nil {
		ctx.FollowUpError("Error registering event", "")
		return
	}

	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Event registered",
		},
	})
}

func (c *EventCommand) unregister(ctx *ken.SubCommandCtx) (err error) {
	messageID := ctx.Options().GetByName("messageID").StringValue()
	channel, err := ctx.Channel()
	if err != nil {
		return
	}

	channelID := channel.ID

	db := ctx.Get(static.DiDatabase).(*database.Database)
	if err = db.UnregisterEvent(messageID, channelID); err != nil {
		ctx.FollowUpError("Error unregistering event", "")
		return
	}

	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Event unregistered",
		},
	})
}
