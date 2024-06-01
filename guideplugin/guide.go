package guideplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/snowflake/v2"
	"strings"
)

type Guide struct {
	Embed   discord.Embed
	EmojiID snowflake.ID
}

var hiveGuides = map[string]Guide{
	"Red Hive": {
		Embed: discord.Embed{
			Title:       "Red Hive Guides",
			Color:       0xC51E3A,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.my.to/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.my.to/ for all of our posts!",
				IconURL: "",
			},
		},
		EmojiID: 1055878225663901827,
	},
	"Blue Hive": {
		Embed: discord.Embed{
			Title:       "Blue Hive Guides",
			Color:       0x318CE7,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.my.to/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.my.to/ for all of our posts!",
				IconURL: "",
			},
		},
		EmojiID: 1055878223931654244,
	},
	"White Hive": {
		Embed: discord.Embed{
			Title:       "White Hive Guides",
			Color:       0xFFFDD0,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.my.to/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.my.to/ for all of our posts!",
				IconURL: "",
			},
		},
		EmojiID: 1055878226548883507,
	},
}

func GetComponents(guide map[string]Guide) []discord.StringSelectMenuOption {
	beeNames := make([]discord.StringSelectMenuOption, 0)
	for k, v := range guide {
		beeNames = append(beeNames, discord.StringSelectMenuOption{
			Label: k,
			Value: k,
			Emoji: &discord.ComponentEmoji{
				ID: v.EmojiID,
			},
		})
	}
	return beeNames
}

func GuidesCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "guides",
			Description: "Various guides for BSS",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "hive",
					Description: "Hive guides",
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"hive": func(event *events.ApplicationCommandInteractionCreate) error {
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Hive Guides",
							Color:       0x318CE7,
							Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.my.to/category/hive-builds/). We have guides for all hive colors, and much more to come!",
							Footer: &discord.EmbedFooter{
								Text:    "Visit https://meta-bee.my.to/ for all of our posts!",
								IconURL: "",
							},
						},
					},
					//Components: []discord.ContainerComponent{
					//	discord.ActionRowComponent{
					//		discord.NewStringSelectMenu("handler:guides:hive", "Select a guide", GetComponents(hiveGuides)...),
					//	},
					//},
				})
			},
		},
	}
}

func GuidesComponent(b *common.Bot) handler.Component {
	return handler.Component{
		Name: "guides",
		Handler: func(event *events.ComponentInteractionCreate) error {
			guideType := strings.Split(event.Data.CustomID(), ":")[2]
			guideName := event.StringSelectMenuInteractionData().Values[0]
			var g Guide
			switch guideType {
			case "hive":
				g = hiveGuides[guideName]
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					g.Embed,
				},
			})
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(GuidesCommand(b))
	h.AddComponents(GuidesComponent(b))
}
