package adminplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func LoadHives(b *common.Bot, hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) {
	hives, err := jsonCacheService.LoadHives("hives.json")
	if err != nil {
		b.Logger.Error(err)
		return
	}
	b.Logger.Infof("Loading %d hives", len(hives))
	for _, cachedUser := range hives {
		h := hiveService.CreateHive(snowflake.MustParse(cachedUser.Id))
		for idx, cachedBee := range cachedUser.Hive {
			h.AddBee(hive.NewBee(cachedBee.Level, cachedBee.Id, cachedBee.Gifted), idx)
			h.GetBee(idx).SetBeequip(cachedBee.Beequip)
			h.GetBee(idx).SetMutation(cachedBee.Mutation)
		}
	}
}

func BackupHives(b *common.Bot, hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) {
	cachedUsers := make([]database.CachedUser, 0)
	for id, h := range hiveService.Hives() {
		cachedHive := make(database.CachedHive)
		for idx, bee := range h.GetBees() {
			cachedHive[idx] = database.CachedBee{
				Id:       bee.Id(),
				Level:    bee.Level(),
				Gifted:   bee.Gifted(),
				Beequip:  bee.Beequip(),
				Mutation: bee.Mutation(),
			}
		}
		cachedUsers = append(cachedUsers, database.CachedUser{
			Id:   id.String(),
			Hive: cachedHive,
		})
	}
	err := jsonCacheService.SaveHives("hives.json", cachedUsers)
	if err != nil {
		b.Logger.Error(err)
	}
	b.Logger.Infof("Backed up %d hives", len(cachedUsers))
}

func AdminCommand(b *common.Bot, hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
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
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"add-trivia-question": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				question := data.String("question")
				answer := data.String("answer")
				difficulty := data.Int("difficulty")
				incorrectAnswer1, ok1 := data.OptString("incorrect-answer-1")
				incorrectAnswer2, ok2 := data.OptString("incorrect-answer-2")
				incorrectAnswer3, ok3 := data.OptString("incorrect-answer-3")
				incorrectAnswer4, ok4 := data.OptString("incorrect-answer-4")
				var incorrectAnswers []string
				if ok1 {
					incorrectAnswers = append(incorrectAnswers, incorrectAnswer1)
				}
				if ok2 {
					incorrectAnswers = append(incorrectAnswers, incorrectAnswer2)
				}
				if ok3 {
					incorrectAnswers = append(incorrectAnswers, incorrectAnswer3)
				}
				if ok4 {
					incorrectAnswers = append(incorrectAnswers, incorrectAnswer4)
				}
				_, err := b.Db.Collection("trivia").InsertOne(context.TODO(), database.TriviaQuestion{
					ID:         primitive.NewObjectID(),
					Difficulty: difficulty,
					Question:   question,
					Answer:     answer,
					Incorrect:  incorrectAnswers,
				})
				if err != nil {
					return err
				}
				return event.CreateMessage(discord.MessageCreate{Content: "ok"})
			},
			"active-hives": func(event *events.ApplicationCommandInteractionCreate) error {
				return event.CreateMessage(discord.MessageCreate{Content: strconv.Itoa(hiveService.HiveCount())})
			},
			"json-save-hives": func(event *events.ApplicationCommandInteractionCreate) error {
				cachedUsers := make([]database.CachedUser, 0)
				for id, h := range hiveService.Hives() {
					cachedHive := make(database.CachedHive)
					for idx, bee := range h.GetBees() {
						cachedHive[idx] = database.CachedBee{
							Id:       bee.Id(),
							Level:    bee.Level(),
							Gifted:   bee.Gifted(),
							Beequip:  bee.Beequip(),
							Mutation: bee.Mutation(),
						}
					}
					cachedUsers = append(cachedUsers, database.CachedUser{
						Id:   id.String(),
						Hive: cachedHive,
					})
				}
				err := jsonCacheService.SaveHives("hives.json", cachedUsers)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{Content: err.Error()})
				}
				return event.CreateMessage(discord.MessageCreate{Content: "ok"})
			},
			"json-load-hives": func(event *events.ApplicationCommandInteractionCreate) error {
				hives, err := jsonCacheService.LoadHives("hives.json")
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{Content: err.Error()})
				}
				for _, cachedUser := range hives {
					h := hiveService.CreateHive(snowflake.MustParse(cachedUser.Id))
					for idx, cachedBee := range cachedUser.Hive {
						h.AddBee(hive.NewBee(cachedBee.Level, cachedBee.Id, cachedBee.Gifted), idx)
						h.GetBee(idx).SetBeequip(cachedBee.Beequip)
						h.GetBee(idx).SetMutation(cachedBee.Mutation)
					}
				}
				return event.CreateMessage(discord.MessageCreate{Content: "ok"})
			},
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot, hiveService *hiveplugin.State) {
	jsonCacheService := database.NewJsonCache()
	h.AddCommands(AdminCommand(b, hiveService, jsonCacheService))
	b.Client.AddEventListeners(&events.ListenerAdapter{
		OnReady: func(event *events.Ready) {
			LoadHives(b, hiveService, jsonCacheService)
			b.Logger.Info("Loaded hives from json.")
			ticker := time.NewTicker(10 * time.Minute)
			go func() {
				for {
					select {
					case <-ticker.C:
						BackupHives(b, hiveService, jsonCacheService)
					}
				}
			}()
			b.Logger.Info("Started automated hive backups.")
		},
	})
}
