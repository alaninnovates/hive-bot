package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type LbUser struct {
	Username      string
	Discriminator string
	Points        float64
	Type          string
}

func LeaderboardCommand(b *common.Bot, gameService *State) func(event *events.ApplicationCommandInteractionCreate) error {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		data := event.SlashCommandInteractionData()
		lbType := data.String("type")
		lbTypeStr := ""
		var users []LbUser
		if lbType == "trivia" {
			lbTypeStr = "Trivia"
			cursor, _ := b.Db.Collection("leaderboards").Aggregate(context.TODO(), mongo.Pipeline{
				{{"$sort", bson.D{{"trivia_points", -1}}}},
				{{"$limit", 10}},
			})
			for cursor.Next(context.TODO()) {
				var user database.LeaderboardUser
				err := cursor.Decode(&user)
				if err != nil {
					return err
				}
				users = append(users, LbUser{
					Username:      user.Username,
					Discriminator: user.Discriminator,
					Points:        float64(user.TriviaPoints),
					Type:          "points",
				})
			}
		} else if lbType == "pop-bubbles" {
			lbTypeStr = "Pop the Bubbles"
			cursor, err := b.Db.Collection("leaderboards").Aggregate(context.TODO(), mongo.Pipeline{
				{{"$match", bson.D{{"bubble_time", bson.D{{"$gt", 0}}}}}},
				{{"$sort", bson.D{{"bubble_time", 1}}}},
				{{"$limit", 10}},
			})
			if err != nil {
				panic(err)
			}
			for cursor.Next(context.TODO()) {
				var user database.LeaderboardUser
				err := cursor.Decode(&user)
				if err != nil {
					return err
				}
				users = append(users, LbUser{
					Username:      user.Username,
					Discriminator: user.Discriminator,
					Points:        float64(user.BubbleTime) / 1000,
					Type:          "seconds",
				})
			}
		}
		usersString := ""
		for i, user := range users {
			usersString += fmt.Sprintf("%d. %s#%s - %.2f %s\n", i+1, user.Username, user.Discriminator, user.Points, user.Type)
		}
		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				{
					Title:       fmt.Sprintf("%s Leaderboard - Top 10", cases.Title(language.English).String(lbTypeStr)),
					Description: fmt.Sprintf("```\n%s```", usersString),
					Color:       0x00ff00,
				},
			},
		})
	}
}
