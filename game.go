package main

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
)

func GameCommand(b *Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "game",
			Description: "Play games",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "identify",
					Description: "Play a game of identify-the-bee",
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"identify": func(event *events.ApplicationCommandInteractionCreate) error {
				return event.CreateMessage(discord.MessageCreate{
					Content: "This command is not yet implemented",
				})
			},
		},
	}
}

func InitializeGameCommands(h *handler.Handler, b *Bot) {
	h.AddCommands(GameCommand(b))
}
