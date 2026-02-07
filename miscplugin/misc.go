package miscplugin

import (
	"strconv"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
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
			▸ </hive create:1053476146978758666>
			▸ </hive add:1053476146978758666> <name> <slots> <level> [gifted]
			▸ </hive remove:1053476146978758666> <slots>
			▸ </hive setmutation:1053476146978758666> <slots> <name>
			▸ </hive setbeequip:1053476146978758666> <slots> <name>
			▸ </hive giftall:1053476146978758666>
			▸ </hive setlevel:1053476146978758666> <level>
			▸ </hive view:1053476146978758666> [show_hive_numbers]
			▸ </hive info:1053476146978758666>
			▸ </hive save:1053476146978758666> <name>
			▸ </hive saves list:1053476146978758666>
			▸ </hive saves load:1053476146978758666> <id>
			▸ </hive saves delete:1053476146978758666> <id>

			[] = Optional | <> = Required`,
		Color: 0xfcba03,
	},
	"game": {
		Title: ":video_game: Games",
		Description: `Note: More games are coming in the future!

			▸ </game identify-the-bee:1053476146978758667>

			[] = Optional | <> = Required`,
		Color: 0x03b1fc,
	},
	"guide": {
		Title: ":bee: Guides",
		Description: `Note: More guides are coming in the future!

			▸ </guides hive:1055875239281688576>

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
						common.LinksActionRow,
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
				b.Client.Caches().GuildsForEach(func(e discord.Guild) {
					members += e.MemberCount
				})
				guildId, _ := strconv.Atoi(event.GuildID().String())
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Fields: []discord.EmbedField{
								{
									Name:   "Guilds on this shard (not total)",
									Value:  strconv.Itoa(b.Client.Caches().GuildsLen()),
									Inline: json.Ptr(true),
								},
								{
									Name:   "Members",
									Value:  strconv.Itoa(members),
									Inline: json.Ptr(true),
								},
								{
									Name:  "Shard ID",
									Value: strconv.Itoa(guildId >> 22 % len(b.Client.ShardManager().Shards())),
								},
							},
							Footer: &discord.EmbedFooter{
								Text: "Made by alaninnovates#0123",
							},
						},
					},
					Components: []discord.ContainerComponent{
						common.LinksActionRow,
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
