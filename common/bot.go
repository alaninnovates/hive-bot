package common

import (
	"log/slog"

	"alaninnovates.com/hive-bot/database"
	"github.com/disgoorg/disgo/bot"
)

type Bot struct {
	Logger *slog.Logger
	Client bot.Client
	Db     database.Database
	R2     database.R2
}
