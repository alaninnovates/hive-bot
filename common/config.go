package common

import "github.com/disgoorg/disgo/discord"

var (
	MaxFreeSaves   = 5
	LinksActionRow = discord.ActionRowComponent{
		discord.NewLinkButton("Documentation", "https://hive-builder.alaninnovates.com/"),
		discord.NewLinkButton("Support server", "https://discord.gg/hive-builder-community-995988457136603147"),
	}
)
