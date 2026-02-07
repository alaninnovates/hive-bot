package common

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
)

func UserIDCheck() handler.Check[*events.ComponentInteractionCreate] {
	return func(event *events.ComponentInteractionCreate) bool {
		allow := event.User().ID.String() == strings.Split(event.ButtonInteractionData().CustomID(), ":")[2]
		if !allow {
			_ = event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This is not your hive!").
				SetEphemeral(true).
				Build())
		}
		return allow
	}
}
