package hiveplugin

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func UserHasHiveCheck(hiveService *State) handler.Middleware {
	return func(next handler.Handler) handler.Handler {
		return func(event *handler.InteractionEvent) error {
			h := hiveService.GetHive(event.User().ID)
			if h == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You don't have a hive. Create one with the `/hive create` command.",
				})
			}
			return next(event)
		}
	}
}

func UserOwnsHiveCheck(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		if event.Type() != discord.InteractionTypeComponent {
			return next(event)
		}

		btnEvent := event.Interaction.(discord.ComponentInteraction)
		uid := strings.Split(btnEvent.ButtonInteractionData().CustomID(), "/")

		if event.User().ID.String() != uid[4] {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This is not your hive!").
				SetEphemeral(true).
				Build())
		}
		return next(event)
	}
}
