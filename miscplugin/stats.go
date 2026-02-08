package miscplugin

import (
	"strconv"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var StatsCommandCreate = discord.SlashCommandCreate{
	Name:        "stats",
	Description: "Get statistics about the bot",
}

func StatsCommand(b *common.Bot) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
	}
}
