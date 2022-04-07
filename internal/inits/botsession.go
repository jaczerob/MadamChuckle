package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/viper"
)

func InitDiscordBotSession(ctn di.Container) (session *discordgo.Session, err error) {
	return discordgo.New("Bot " + viper.GetString("discord.token"))
}
