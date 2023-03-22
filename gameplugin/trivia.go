package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

type TriviaDifficulty int

const (
	TriviaDifficultyBeginner TriviaDifficulty = iota
	TriviaDifficultyMidgame
	TriviaDifficultyEndgame
)

var b *common.Bot

func GetDifficulty(difficulty TriviaDifficulty) string {
	switch difficulty {
	case 0:
		return "Beginner"
	case 1:
		return "Midgame"
	case 2:
		return "Endgame"
	default:
		panic("Invalid difficulty")
	}
}

func GetQuestion(difficulty TriviaDifficulty, userId snowflake.ID) (TriviaQuestionAnswer, discord.Embed, discord.ActionRowComponent) {
	var question database.TriviaQuestion
	cursor, err := b.Db.Collection("trivia").Aggregate(
		context.TODO(),
		mongo.Pipeline{
			bson.D{{"$match", bson.D{{"difficulty", difficulty}}}},
			bson.D{{"$sample", bson.D{{"size", 1}}}},
		})
	if err == nil {
		cursor.Next(context.TODO())
		err = cursor.Decode(&question)
	}
	if err != nil {
		return TriviaQuestionAnswer{}, discord.Embed{
				Title:       "Error",
				Description: fmt.Sprintf("Error getting question: %s", err.Error()),
				Footer: &discord.EmbedFooter{
					Text: "Please report this to alaninnovates#0123",
				},
				Color: 0xff0000,
			}, discord.ActionRowComponent{}.AddComponents(discord.ButtonComponent{
				CustomID: "ignore",
				Label:    "Error",
				Disabled: true,
				Style:    discord.ButtonStyleDanger,
			})
	}
	answerArr := []struct {
		Answer  string
		Correct bool
	}{{
		Answer:  question.Answer,
		Correct: true,
	}}
	var answerArrStr []string
	for _, wrongAnswer := range question.Incorrect {
		answerArr = append(answerArr, struct {
			Answer  string
			Correct bool
		}{
			Answer:  wrongAnswer,
			Correct: false,
		})
	}
	answerArr = common.ShuffleArray(answerArr)
	var buttons []discord.InteractiveComponent
	correctIdx := 0
	for i, v := range answerArr {
		answerArrStr = append(answerArrStr, v.Answer)
		correct := "wrong"
		if v.Correct {
			correct = "correct"
			correctIdx = i
		}
		buttons = append(buttons, discord.ButtonComponent{
			CustomID: fmt.Sprintf("handler:trivia:%s:%s:%d", userId, correct, i),
			Label:    strconv.Itoa(i + 1),
			Style:    discord.ButtonStylePrimary,
		})
	}
	answersStr := ""
	for i, v := range answerArrStr {
		answersStr += fmt.Sprintf("%d. %s\n", i+1, v)
	}
	return TriviaQuestionAnswer{
			Question:     question.Question,
			Choices:      answerArrStr,
			CorrectIndex: correctIdx,
			StartTime:    time.Now().UnixMilli(),
		}, discord.Embed{
			Title:       "Trivia",
			Description: fmt.Sprintf("Difficulty: %s", GetDifficulty(difficulty)),
			Fields: []discord.EmbedField{
				{
					Name:  "Question",
					Value: question.Question,
				},
				{
					Name:  "Answers",
					Value: answersStr,
				},
			},
		}, discord.ActionRowComponent{}.AddComponents(buttons...)
}

func TriviaCommand(bot *common.Bot, gameService *State) func(event *events.ApplicationCommandInteractionCreate) error {
	b = bot
	return func(event *events.ApplicationCommandInteractionCreate) error {
		//return event.CreateMessage(discord.MessageCreate{
		//	Embeds: []discord.Embed{
		//		{
		//			Title: "Not Implemented... Yet",
		//			Description: "This command is not implemented yet due to a lack of trivia commands.\n" +
		//				"Want to contribute? Add your trivia questions [here](https://forms.gle/5R5ouBdr4gkh7etH8)!",
		//			Color: 0xff0000,
		//		},
		//	},
		//})
		if gameService.IsPlayingGame(event.User().ID) {
			return event.CreateMessage(discord.MessageCreate{
				Content: "You are already playing a game!",
			})
		}
		data := event.SlashCommandInteractionData()
		difficulty := TriviaDifficulty(data.Int("difficulty"))
		gameService.StartTriviaGame(event.User().ID, difficulty)
		triviaAnswer, emb, buttons := GetQuestion(difficulty, event.User().ID)
		tu := gameService.GetGameUser(event.User().ID, GameTypeTrivia).TriviaGameUser
		tu.AddQuestion(triviaAnswer)
		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				emb,
			},
			Components: []discord.ContainerComponent{
				buttons,
			},
		})
	}
}

func TriviaButton(gameState *State) handler.Component {
	return handler.Component{
		Name:  "trivia",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			data := strings.Split(event.ButtonInteractionData().CustomID(), ":")
			_, idx := data[3], data[4]
			user := gameState.GetGameUser(event.User().ID, GameTypeTrivia)
			if user == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not playing a trivia game!",
				})
			}
			tu := user.TriviaGameUser
			index, _ := strconv.Atoi(idx)
			tu.SetChosenIndex(index)
			if len(tu.Questions) == 10 {
				b.Db.Collection("leaderboards").FindOneAndUpdate(context.TODO(),
					bson.D{{"user_id", event.User().ID}},
					bson.D{
						{"$inc", bson.D{{"trivia_points", tu.GetScore()}}},
						{"$set", bson.D{{
							"username", event.User().Username,
						}, {
							"discriminator", event.User().Discriminator,
						}}},
					},
					options.FindOneAndUpdate().SetUpsert(true))
				return event.UpdateMessage(discord.MessageUpdate{
					Content: json.Ptr(fmt.Sprintf("You got %d out of 10 questions correct! That makes your score %d!", tu.GetCorrect(), tu.GetScore())),
					Embeds: &[]discord.Embed{
						{
							Title:       "Help us!",
							Description: "We are currently looking for more trivia questions. If you have any, please fill out [this form](https://forms.gle/5R5ouBdr4gkh7etH8).",
							Color:       0x00ff00,
						},
					},
					Components: &[]discord.ContainerComponent{
						discord.ActionRowComponent{}.AddComponents(discord.ButtonComponent{
							CustomID: "handler:endgame:" + event.User().ID.String(),
							Label:    "End Game",
							Style:    discord.ButtonStylePrimary,
						}, discord.ButtonComponent{
							CustomID: "handler:trivia-review:" + event.User().ID.String(),
							Label:    "Review Questions",
							Style:    discord.ButtonStyleSuccess,
						}),
					},
				})
			}
			triviaAnswer, emb, buttons := GetQuestion(user.TriviaGameUser.Difficulty, event.User().ID)
			tu.AddQuestion(triviaAnswer)
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					emb,
				},
				Components: &[]discord.ContainerComponent{
					buttons,
				},
			})
		},
	}
}

func TriviaReviewButton(gameState *State) handler.Component {
	return handler.Component{
		Name:  "trivia-review",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			user := gameState.GetGameUser(event.User().ID, GameTypeTrivia)
			if user == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not playing a trivia game!",
				})
			}
			tu := user.TriviaGameUser
			var fields []discord.EmbedField
			for i, v := range tu.Questions {
				timeTaken := float64(v.EndTime-v.StartTime) / 1000
				timeTakenStr := fmt.Sprintf("%.2f seconds", timeTaken)
				fields = append(fields, discord.EmbedField{
					Name:  fmt.Sprintf("Question %d", i+1),
					Value: fmt.Sprintf("Question: %s\nCorrect Answer: %s\nYour Answer: %s\nTime Taken: %s", v.Question, v.Choices[v.CorrectIndex], v.Choices[v.ChosenIndex], timeTakenStr),
				})
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Content: json.Ptr(""),
				Embeds: &[]discord.Embed{
					{
						Title:       "Trivia Review",
						Description: fmt.Sprintf("You got %d out of 10 questions correct!", tu.GetCorrect()),
						Fields:      fields,
					},
				},
				Components: &[]discord.ContainerComponent{
					discord.ActionRowComponent{}.AddComponents(discord.ButtonComponent{
						CustomID: "handler:endgame:" + event.User().ID.String(),
						Label:    "End Game",
						Style:    discord.ButtonStylePrimary,
					}),
				},
			})
		},
	}
}
