package common

import (
	"alaninnovates.com/hive-bot/database"
	"github.com/disgoorg/disgo/bot"
	"log/slog"
)

type Bot struct {
	Logger *slog.Logger
	Client bot.Client
	Db     database.Database
}
