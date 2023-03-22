package adminplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"context"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AdminCommand(b *common.Bot) handler.Command {
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
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	h.AddCommands(AdminCommand(b))
}
