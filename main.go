package main

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/gameplugin"
	"alaninnovates.com/hive-bot/guideplugin"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/miscplugin"
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(log.LstdFlags | log.Lshortfile)
	logger.SetLevel(log.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load .env: ", err)
	}

	var (
		token = os.Getenv("TOKEN")
		_     = snowflake.GetEnv("GUILD_ID")
		dbUri = os.Getenv("MONGODB_URI")
	)

	hiveBot := &common.Bot{
		Logger: logger,
		Db:     *database.NewDatabase(),
	}

	client, err := hiveBot.Db.Connect(dbUri)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	h := handler.New(logger)
	gameplugin.Initialize(h, hiveBot)
	hiveplugin.Initialize(h, hiveBot)
	guideplugin.Initialize(h, hiveBot)
	miscplugin.Initialize(h, hiveBot)

	if hiveBot.Client, err = disgo.New(token,
		bot.WithLogger(logger),
		bot.WithDefaultGateway(),
		bot.WithEventListeners(h),
	); err != nil {
		logger.Fatal("Failed to create disgo client: ", err)
	}

	h.SyncCommands(hiveBot.Client)

	if err = hiveBot.Client.OpenGateway(context.TODO()); err != nil {
		logger.Fatal("Failed to open gateway: ", err)
	}

	logger.Info("Hive Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
