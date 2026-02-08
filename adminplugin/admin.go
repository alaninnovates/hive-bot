package adminplugin

import (
	"time"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/hiveplugin"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var AdminCommandCreate = discord.SlashCommandCreate{
	Name:        "admin",
	Description: "Stuff for admins",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "add-trivia-question",
			Description: "Add a trivia question",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "question",
					Description: "The question",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "difficulty",
					Description: "The difficulty of the question",
					Required:    true,
					Choices: []discord.ApplicationCommandOptionChoiceInt{
						{
							Name:  "Beginner",
							Value: 0,
						},
						{
							Name:  "Midgame",
							Value: 1,
						},
						{
							Name:  "Endgame",
							Value: 2,
						},
					},
				},
				discord.ApplicationCommandOptionString{
					Name:        "answer",
					Description: "The correct answer",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "incorrect-answer-1",
					Description: "An incorrect answer",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "incorrect-answer-2",
					Description: "Another incorrect answer",
				},
				discord.ApplicationCommandOptionString{
					Name:        "incorrect-answer-3",
					Description: "Another incorrect answer",
				},
				discord.ApplicationCommandOptionString{
					Name:        "incorrect-answer-4",
					Description: "Another incorrect answer",
				},
			},
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "active-hives",
			Description: "List the number of active, cached hives",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "json-save-hives",
			Description: "Save the current hives to a json file",
		},
		discord.ApplicationCommandOptionSubCommand{
			Name:        "json-load-hives",
			Description: "Load hives from a json file",
		},
	},
}

func Initialize(r *handler.Mux, b *common.Bot, hiveService *hiveplugin.State, devMode bool) {
	jsonCacheService := database.NewJsonCache()
	r.Route("/admin", func(r handler.Router) {
		r.Command("/add-trivia-question", AddTriviaQuestionCommand(b))
		r.Command("/active-hives", ActiveHivesCommand(hiveService))
		r.Command("/json-save-hives", JSONSaveHivesCommand(hiveService, jsonCacheService))
		r.Command("/json-load-hives", JSONLoadHivesCommand(hiveService, jsonCacheService))
	})
	b.Client.AddEventListeners(&events.ListenerAdapter{
		OnReady: func(event *events.Ready) {
			if devMode {
				return
			}
			LoadHives(b, hiveService, jsonCacheService)
			b.Logger.Info("Loaded hives from json.")
			ticker := time.NewTicker(10 * time.Minute)
			go func() {
				for {
					select {
					case <-ticker.C:
						PruneHives(b, hiveService)
						BackupHives(b, hiveService, jsonCacheService)
					}
				}
			}()
			b.Logger.Info("Started automated hive backups.")
		},
	})
}
