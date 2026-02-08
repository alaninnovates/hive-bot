package common

import "github.com/disgoorg/disgo/discord"

var (
	MaxFreeSaves    = 5
	MaxPremiumSaves = 25
	LinksActionRow  = discord.ActionRowComponent{
		discord.NewLinkButton("Documentation", "https://meta-bee.com/hive-builder/"),
		discord.NewLinkButton("Support server", "https://discord.gg/2BgUMfCsHM"),
	}
)
