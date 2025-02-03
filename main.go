package main

import (
	"alaninnovates.com/hive-bot/adminplugin"
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/gameplugin"
	"alaninnovates.com/hive-bot/guideplugin"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/miscplugin"
	"alaninnovates.com/hive-bot/statsplugin"
	"context"
	"flag"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	useEnvFilePtr := flag.Bool("env", false, "a bool")
	devPtr := flag.Bool("dev", false, "a bool")
	syncCommandsPtr := flag.Bool("sync", false, "a bool")
	flag.Parse()
	if *useEnvFilePtr {
		if *devPtr {
			err := godotenv.Load(".env.dev")
			if err != nil {
				logger.Error("Failed to load .env.dev")
				panic(err)
			}
		} else {
			err := godotenv.Load(".env")
			if err != nil {
				logger.Error("Failed to load .env")
				panic(err)
			}
		}
	}
	devMode := *devPtr
	syncCommands := *syncCommandsPtr
	var (
		token = os.Getenv("TOKEN")
		dbUri = os.Getenv("MONGODB_URI")
	)

	hiveBot := &common.Bot{
		Logger: logger,
		Db:     *database.NewDatabase(),
	}

	client, err := hiveBot.Db.Connect(dbUri)
	if err != nil {
		logger.Error("Failed to connect to database")
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	h := handler.New(log.New(log.LstdFlags | log.Lshortfile))
	go statsplugin.Initialize(h, hiveBot, devMode)
	gameplugin.Initialize(h, hiveBot)
	hiveService := hiveplugin.NewHiveService()
	hiveplugin.Initialize(h, hiveBot, hiveService)
	guideplugin.Initialize(h, hiveBot)
	miscplugin.Initialize(h, hiveBot)

	if hiveBot.Client, err = disgo.New(token,
		bot.WithShardManagerConfigOpts(
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentGuilds),
				gateway.WithLogger(logger),
			),
			sharding.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds),
		),
		bot.WithEventListeners(h),
	); err != nil {
		logger.Error("Failed to create disgo client")
		panic(err)
	}

	if !devMode && syncCommands {
		h.SyncCommands(hiveBot.Client)
	}
	adminplugin.Initialize(h, hiveBot, hiveService, devMode)
	if devMode && syncCommands {
		h.SyncCommands(hiveBot.Client, snowflake.GetEnv("GUILD_ID"))
		//_, _ = hiveBot.Client.Rest().SetGlobalCommands(hiveBot.Client.ApplicationID(), []discord.ApplicationCommandCreate{})
	}

	if err = hiveBot.Client.OpenShardManager(context.Background()); err != nil {
		logger.Error("Failed to open shard manager")
		panic(err)
	}

	logger.Info("Hive Bot is running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
