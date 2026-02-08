package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/gameplugin"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
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

	r := handler.New()
	r.Use(middleware.Go)
	//r.Use(middleware.Logger)

	//go statsplugin.Initialize(r, hiveBot, devMode)
	gameplugin.Initialize(r, hiveBot)
	//hiveService := hiveplugin.NewHiveService()
	//hiveplugin.Initialize(r, hiveBot, hiveService)
	//guideplugin.Initialize(r, hiveBot)
	//miscplugin.Initialize(r, hiveBot)

	//r.NotFound()

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
		bot.WithEventListeners(r),
	); err != nil {
		logger.Error("Failed to create disgo client")
		panic(err)
	}

	//if !devMode && syncCommands {
	//	h.SyncCommands(hiveBot.Client)
	//}
	//adminplugin.Initialize(r, hiveBot, hiveService, devMode)
	if devMode && syncCommands {
		//_, _ = hiveBot.Client.Rest().SetGlobalCommands(hiveBot.Client.ApplicationID(), []discord.ApplicationCommandCreate{})
		logger.Info("Syncing commands...")
		if err = handler.SyncCommands(hiveBot.Client, []discord.ApplicationCommandCreate{
			gameplugin.GameCommandCreate,
		}, []snowflake.ID{snowflake.GetEnv("GUILD_ID")}); err != nil {
			logger.Error("error while syncing commands", slog.Any("err", err))
			return
		}
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
