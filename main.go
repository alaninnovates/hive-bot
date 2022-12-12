package main

import (
	"alaninnovates.com/hive-bot/database"
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

type Bot struct {
	Logger log.Logger
	Client bot.Client
	Db     database.Database
	State  State
}

func main() {
	logger := log.New(log.LstdFlags | log.Lshortfile)
	logger.SetLevel(log.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Failed to load .env: ", err)
	}

	var (
		token   = os.Getenv("TOKEN")
		guildID = snowflake.GetEnv("GUILD_ID")
		dbUri   = os.Getenv("MONGODB_URI")
	)

	hiveBot := &Bot{
		Logger: logger,
		Db:     *database.NewDatabase(),
		State:  *NewState(),
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
	InitializeHiveCommands(h, hiveBot)
	InitializeGameCommands(h, hiveBot)

	if hiveBot.Client, err = disgo.New(token,
		bot.WithLogger(logger),
		bot.WithDefaultGateway(),
		bot.WithEventListeners(h),
	); err != nil {
		logger.Fatal("Failed to create disgo client: ", err)
	}

	h.SyncCommands(hiveBot.Client, guildID)

	if err = hiveBot.Client.OpenGateway(context.TODO()); err != nil {
		logger.Fatal("Failed to open gateway: ", err)
	}

	logger.Info("Hive Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
