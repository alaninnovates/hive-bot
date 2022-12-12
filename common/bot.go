package common

import (
	"alaninnovates.com/hive-bot/database"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/log"
)

type Bot struct {
	Logger log.Logger
	Client bot.Client
	Db     database.Database
}
