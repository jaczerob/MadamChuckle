package admin

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaczerob/madamchuckle/internal/middlewares"
	"github.com/zekrotja/ken"
)

type SayCommand struct{}

var (
	_ ken.SlashCommand                      = (*SayCommand)(nil)
	_ ken.DmCapable                         = (*SayCommand)(nil)
	_ middlewares.RequiresPermissionCommand = (*SayCommand)(nil)
)

func (c *SayCommand) Name() string {
	return "say"
}

func (c *SayCommand) Description() string {
	return "Makes Madam Chuckle say something in the current channel"
}

func (c *SayCommand) Version() string {
	return "1.0.0"
}

func (c *SayCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *SayCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "phrase",
			Required:    true,
			Description: "The phrase that Madam Chuckle will say",
		},
	}
}

func (c *SayCommand) IsDmCapable() bool {
	return false
}

func (c *SayCommand) RequiresPermission() int64 {
	return discordgo.PermissionAdministrator
}

func (c *SayCommand) Run(ctx *ken.Ctx) (err error) {
	phrase := ctx.Options().GetByName("phrase").StringValue()
	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: phrase,
		},
	})
}
