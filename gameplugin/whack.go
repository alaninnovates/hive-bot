package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/events"
)

func WhackCommand(b *common.Bot, gameService *State) func(event *events.ApplicationCommandInteractionCreate) error {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		return nil
	}
}
