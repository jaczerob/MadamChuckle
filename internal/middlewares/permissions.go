package middlewares

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type RequiresPermissionCommand interface {
	RequiresPermission() int64
}

type PermissionsMiddleware struct {
	count int
}

var _ ken.MiddlewareBefore = (*PermissionsMiddleware)(nil)

func (c *PermissionsMiddleware) Before(ctx *ken.Ctx) (next bool, err error) {
	cmd, ok := ctx.Command.(RequiresPermissionCommand)
	if !ok || ctx.Event.Member.Permissions&cmd.RequiresPermission() == cmd.RequiresPermission() {
		next = true
	}

	if !next {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't have permissions to use this command!",
			},
		})
	}

	return
}
