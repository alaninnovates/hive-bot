package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"alaninnovates.com/hive-bot/adminplugin"
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/gameplugin"
	"alaninnovates.com/hive-bot/guideplugin"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/miscplugin"
	"alaninnovates.com/hive-bot/statsplugin"
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

	gameplugin.Initialize(r, hiveBot)
	hiveService := hiveplugin.NewHiveService()
	hiveplugin.Initialize(r, hiveBot, hiveService)
	guideplugin.Initialize(r, hiveBot)
	miscplugin.Initialize(r, hiveBot)

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

	statsplugin.Initialize(r, hiveBot, devMode)
	adminplugin.Initialize(r, hiveBot, hiveService, devMode)

	if syncCommands {
		//_, _ = hiveBot.Client.Rest().SetGlobalCommands(hiveBot.Client.ApplicationID(), []discord.ApplicationCommandCreate{})
		logger.Info("Syncing commands...")
		commands := []discord.ApplicationCommandCreate{
			hiveplugin.HiveCommandCreate,
			gameplugin.GameCommandCreate,
			guideplugin.GuidesCreateCommand,
			miscplugin.HelpCommandCreate,
			miscplugin.StatsCommandCreate,
		}
		var guilds []snowflake.ID
		if devMode {
			logger.Info("Developer mode enabled: Syncing commands to test guild only")
			guilds = append(guilds, snowflake.GetEnv("GUILD_ID"))
			commands = append(commands, adminplugin.AdminCommandCreate)
		}
		if err = handler.SyncCommands(hiveBot.Client, commands, guilds); err != nil {
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
