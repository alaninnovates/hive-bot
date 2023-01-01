package miscplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
)

func HelpCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "help",
			Description: "Get help with the bot",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": func(event *events.ApplicationCommandInteractionCreate) error {
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Hive Bot Help",
							Description: "this is a todo for now",
							Footer: &discord.EmbedFooter{
								Text: "Made by alaninnovates#0123",
							},
						},
					},
				})
			},
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(HelpCommand(b))
}
