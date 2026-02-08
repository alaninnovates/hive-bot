package adminplugin

import (
	"context"
	"strconv"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTriviaQuestionCommand(b *common.Bot) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
	}
}

func ActiveHivesCommand(hiveService *hiveplugin.State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		return event.CreateMessage(discord.MessageCreate{Content: strconv.Itoa(hiveService.HiveCount())})
	}
}

func JSONSaveHivesCommand(hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		cachedUsers := make([]database.CachedUser, 0)
		for id, h := range hiveService.Hives() {
			cachedHive := make(database.CachedHive)
			for idx, bees := range h.GetBees() {
				for _, bee := range bees {
					cachedHive[idx] = append(cachedHive[idx], database.CachedBee{
						Id:       bee.Id(),
						Level:    bee.Level(),
						Gifted:   bee.Gifted(),
						Beequip:  bee.Beequip(),
						Mutation: bee.Mutation(),
					})
				}
			}
			cachedUsers = append(cachedUsers, database.CachedUser{
				Id:   id.String(),
				Hive: cachedHive,
			})
		}
		err := jsonCacheService.SaveHives("data/hives.json", cachedUsers)
		if err != nil {
			return event.CreateMessage(discord.MessageCreate{Content: err.Error()})
		}
		return event.CreateMessage(discord.MessageCreate{Content: "ok"})
	}
}

func JSONLoadHivesCommand(hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		hives, err := jsonCacheService.LoadHives("data/hives.json")
		if err != nil {
			return event.CreateMessage(discord.MessageCreate{Content: err.Error()})
		}
		for _, cachedUser := range hives {
			h := hiveService.CreateHive(snowflake.MustParse(cachedUser.Id))
			for idx, cachedBees := range cachedUser.Hive {
				for _, cachedBee := range cachedBees {
					h.AddBee(hive.NewBee(cachedBee.Level, cachedBee.Id, cachedBee.Gifted), idx)
					pos := len(h.GetBeesAt(idx)) - 1
					h.GetBeesAt(idx)[pos].SetBeequip(cachedBee.Beequip)
					h.GetBeesAt(idx)[pos].SetMutation(cachedBee.Mutation)
				}
			}
		}
		return event.CreateMessage(discord.MessageCreate{Content: "ok"})
	}
}
