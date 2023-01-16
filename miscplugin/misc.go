package miscplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"strconv"
)

var helpMenus = map[string]discord.Embed{
	"home": {
		Title:       "Hive Builder Help",
		Description: "Visit the documentation below for a list of all commands. Join the support server if you have any more questions.",
		Footer: &discord.EmbedFooter{
			Text: "Made by alaninnovates#0123",
		},
		Color: 0xffffff,
	},
	"hive": {
		Title: ":honey_pot: Hive building",
		Description: `
			▸ /hive create
			▸ /hive add <name> <slots> <level> [gifted]
			▸ /hive remove <slots>
			▸ /hive setmutation <slots> <name>
			▸ /hive setbeequip <slots> <name>
			▸ /hive giftall
			▸ /hive setlevel <level>
			▸ /hive view [show_hive_numbers]
			▸ /hive info
			▸ /hive save <name>
			▸ /hive saves list
			▸ /hive saves load <id>
			▸ /hive saves delete <id>

			[] = Optional | <> = Required`,
		Color: 0xfcba03,
	},
	"game": {
		Title: ":video_game: Games",
		Description: `Note: More games are coming in the future!

			▸ /game identify-the-bee

			[] = Optional | <> = Required`,
		Color: 0x03b1fc,
	},
	"guide": {
		Title: ":bee: Guides",
		Description: `Note: More guides are coming in the future!

			▸ /guides hive

			[] = Optional | <> = Required`,
		Color: 0x03fc73,
	},
}

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
						helpMenus["home"],
					},
					Components: []discord.ContainerComponent{
						discord.ActionRowComponent{
							discord.NewLinkButton("Documentation", "https://hive-builder.alaninnovates.com/"),
							discord.NewLinkButton("Support server", "https://discord.gg/hive-builder-community-995988457136603147"),
						},
						discord.ActionRowComponent{
							discord.NewStringSelectMenu(
								"handler:help",
								"Select a category",
								discord.StringSelectMenuOption{
									Label: "Home",
									Value: "home",
									Emoji: &discord.ComponentEmoji{Name: "🏠"},
								},
								discord.StringSelectMenuOption{
									Label: "Hive Building",
									Value: "hive",
									Emoji: &discord.ComponentEmoji{Name: "🍯"},
								},
								discord.StringSelectMenuOption{
									Label: "Games",
									Value: "game",
									Emoji: &discord.ComponentEmoji{Name: "🎮"},
								},
								discord.StringSelectMenuOption{
									Label: "Guides",
									Value: "guide",
									Emoji: &discord.ComponentEmoji{Name: "🐝"},
								},
							),
						},
					},
				})
			},
		},
	}
}

func HelpComponent(b *common.Bot) handler.Component {
	return handler.Component{
		Name: "help",
		Handler: func(event *events.ComponentInteractionCreate) error {
			sectionName := event.StringSelectMenuInteractionData().Values[0]
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					helpMenus[sectionName],
				},
			})
		},
	}
}

func StatsCommand(b *common.Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "stats",
			Description: "Get statistics about the bot",
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"": func(event *events.ApplicationCommandInteractionCreate) error {
				members := 0
				guilds := 0
				b.Client.Caches().GuildsForEach(func(e discord.Guild) {
					guilds++
					members += e.MemberCount
				})
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Description: "Some statistics may be inaccurate due to caching.",
							Fields: []discord.EmbedField{
								{
									Name:   "Guilds",
									Value:  strconv.Itoa(guilds),
									Inline: json.Ptr(true),
								},
								{
									Name:   "Members",
									Value:  strconv.Itoa(members),
									Inline: json.Ptr(true),
								},
							},
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
	h.AddCommands(HelpCommand(b), StatsCommand(b))
	h.AddComponents(HelpComponent(b))
}
