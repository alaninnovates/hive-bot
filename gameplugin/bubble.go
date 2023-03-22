package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/snowflake/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
)

func GetBubbleComponents(gameService *State, userId snowflake.ID) []discord.ContainerComponent {
	user := gameService.GetGameUser(userId, GameTypePopBubble).BubbleGameUser
	var rows []discord.ContainerComponent
	for i := range user.Bubbles {
		var buttons []discord.InteractiveComponent
		for j := range user.Bubbles[i] {
			style := discord.ButtonStylePrimary
			if user.Bubbles[i][j] {
				style = discord.ButtonStyleSecondary
			}
			buttons = append(buttons, discord.ButtonComponent{
				Style:    discord.ButtonStyle(style),
				CustomID: "handler:bubble:" + userId.String() + ":" + strconv.Itoa(i) + ":" + strconv.Itoa(j),
				Label:    " ",
				Emoji: &discord.ComponentEmoji{
					Name: "🫧",
				},
				Disabled: user.Bubbles[i][j],
			})
		}
		rows = append(rows, discord.ActionRowComponent{}.AddComponents(buttons...))
	}
	return rows
}

func BubbleCommand(b *common.Bot, gameService *State) func(event *events.ApplicationCommandInteractionCreate) error {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		if gameService.IsPlayingGame(event.User().ID) {
			return event.CreateMessage(discord.MessageCreate{
				Content: "You are already playing a game!",
			})
		}
		gameService.StartBubbleGame(event.User().ID)
		rows := GetBubbleComponents(gameService, event.User().ID)
		return event.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				{
					Title:       "Pop the bubbles!",
					Description: "Pop the bubbles by clicking on the buttons!",
				},
			},
			Components: rows,
		})
	}
}

func BubbleButton(gameService *State) handler.Component {
	return handler.Component{
		Name:  "bubble",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			gu := gameService.GetGameUser(event.User().ID, GameTypePopBubble)
			if gu == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not playing a game!",
				})
			}
			bu := gu.BubbleGameUser
			if bu.PopAmount() == 0 {
				bu.StartTime = time.Now().UnixMilli()
			}
			data := event.ButtonInteractionData()
			split := strings.Split(data.CustomID(), ":")
			row, _ := strconv.Atoi(split[3])
			col, _ := strconv.Atoi(split[4])
			bu.Bubbles[row][col] = true
			if bu.PopAmount() == 25 {
				timeTaken := time.Now().UnixMilli() - bu.StartTime
				gameService.EndGame(event.User().ID)
				res, err := b.Db.Collection("leaderboards").UpdateOne(context.TODO(),
					bson.M{"user_id": event.User().ID},
					bson.A{
						bson.M{"$set": bson.M{
							"username": event.User().Username,
						}},
						bson.M{"$set": bson.M{
							"discriminator": event.User().Discriminator,
						}},
						bson.M{"$set": bson.M{
							"bubble_time": bson.M{
								"$cond": bson.M{
									"if": bson.M{
										"$not": "$bubble_time",
									},
									"then": 999999999999,
									"else": "$bubble_time",
								},
							},
						}},
						bson.M{"$set": bson.M{
							"bubble_time": bson.M{
								"$cond": bson.M{
									"if": bson.M{
										"$lt": bson.A{
											"$bubble_time",
											timeTaken,
										},
									},
									"then": "$bubble_time",
									"else": timeTaken,
								},
							},
						}},
					},
					options.Update().SetUpsert(true))
				if err != nil {
					return err
				}
				if res.MatchedCount != 0 {
					fmt.Println("matched and replaced an existing document")
				}
				if res.UpsertedCount != 0 {
					fmt.Printf("inserted a new document with ID %v\n", res.UpsertedID)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "You popped all the bubbles in " + strconv.FormatInt(timeTaken/1000, 10) + " seconds!",
				})
			}
			rows := GetBubbleComponents(gameService, event.User().ID)
			return event.UpdateMessage(discord.MessageUpdate{
				Components: &rows,
			})
		},
	}
}
