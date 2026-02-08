package guideplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

var GuidesCreateCommand = discord.SlashCommandCreate{
	Name:        "guides",
	Description: "Various guides for BSS",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "hive",
			Description: "Hive guides",
		},
	},
}

type Guide struct {
	Embed   discord.Embed
	EmojiID snowflake.ID
}

var hiveGuides = map[string]Guide{
	"Red Hive": {
		Embed: discord.Embed{
			Title:       "Red Hive Guides",
			Color:       0xC51E3A,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.com/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.com/ for all of our posts!",
				IconURL: "",
			},
		},
		EmojiID: 1055878225663901827,
	},
	"Blue Hive": {
		Embed: discord.Embed{
			Title:       "Blue Hive Guides",
			Color:       0x318CE7,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.com/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.com/ for all of our posts!",
				IconURL: "",
			},
		},
		EmojiID: 1055878223931654244,
	},
	"White Hive": {
		Embed: discord.Embed{
			Title:       "White Hive Guides",
			Color:       0xFFFDD0,
			Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.com/category/hive-builds/). We have guides for all hive colors, and much more to come!",
			Footer: &discord.EmbedFooter{
				Text:    "Visit https://meta-bee.com/ for all of our posts!",
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

func GuidesCommand(event *handler.CommandEvent) error {
	return event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Hive Guides",
				Color:       0x318CE7,
				Description: "Find all of our hive guides at [Meta Bee's website](https://meta-bee.com/category/hive-builds/). We have guides for all hive colors, and much more to come!",
				Footer: &discord.EmbedFooter{
					Text:    "Visit https://meta-bee.com/ for all of our posts!",
					IconURL: "",
				},
			},
		},
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.NewStringSelectMenu("/guides/hive", "Select a guide", GetComponents(hiveGuides)...),
			},
		},
	})
}

func GuidesComponent(event *handler.ComponentEvent) error {
	guideType := event.Vars["type"]
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
}

func Initialize(r *handler.Mux, b *common.Bot) {
	r.Route("/guides", func(r handler.Router) {
		r.Command("/hive", GuidesCommand)
		r.Component("/{type}", GuidesComponent)
	})
}
