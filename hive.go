package main

import (
	"alaninnovates.com/hive-bot/hive"
	"alaninnovates.com/hive-bot/loaders"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"image"
	"image/png"
	"io"
	"strconv"
	"strings"
)

func GetRangeNumbers(rangeStr string) []int {
	var nums []int
	for _, s := range strings.Split(rangeStr, ",") {
		if strings.Contains(s, "-") {
			l, _ := strconv.ParseInt(strings.Split(s, "-")[0], 10, 64)
			r, _ := strconv.ParseInt(strings.Split(s, "-")[1], 10, 64)
			for i := l; i <= r; i++ {
				nums = append(nums, int(i))
			}
		} else {
			n, _ := strconv.ParseInt(s, 10, 64)
			nums = append(nums, int(n))
		}
	}
	return nums
}

func HiveCommand(b *Bot) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "hive",
			Description: "Mange your hive",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "create",
					Description: "Start building your hive",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "add",
					Description: "Add a bee to your hive",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:         "name",
							Description:  "The name of the bee",
							Required:     true,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "slots",
							Description: "The slot(s) where you want to add the bee",
							Required:    true,
						},
						discord.ApplicationCommandOptionInt{
							Name:        "level",
							Description: "The level of the bee",
							Required:    true,
						},
						discord.ApplicationCommandOptionBool{
							Name:        "gifted",
							Description: "Whether or not the bee is gifted",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "remove",
					Description: "Remove a bee from your hive",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "slots",
							Description: "The slot(s) of the bees you want to remove",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "giftall",
					Description: "Gift all bees in your hive",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "setbeequip",
					Description: "Set the beequip of a bee",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "slots",
							Description: "The slot(s) of the bees you want to set the beequip of",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:         "name",
							Description:  "The beequip you want to set",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "view",
					Description: "View your hive",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionBool{
							Name:        "show_hive_numbers",
							Description: "Whether or not to show the hive numbers",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "save",
					Description: "Save your hive",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "name",
							Description: "The name of the hive",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommandGroup{
					Name:        "saves",
					Description: "Manage saves",
					Options: []discord.ApplicationCommandOptionSubCommand{
						{
							Name:        "load",
							Description: "Load a hive save",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "id",
									Description: "The id of the hive save",
									Required:    true,
								},
							},
						},
						{
							Name:        "list",
							Description: "View your hive saves",
						},
						{
							Name:        "delete",
							Description: "Delete a hive save",
							Options: []discord.ApplicationCommandOption{
								discord.ApplicationCommandOptionString{
									Name:        "id",
									Description: "The id of the hive save",
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"create": func(event *events.ApplicationCommandInteractionCreate) error {
				b.State.CreateHive(event.User().ID)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Created new hive. You can now add bees with the `/hive add` command.",
				})
			},
			"add": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				name, _ := data.OptString("name")
				if !slices.Contains(loaders.GetBees(), name) {
					return event.CreateMessage(discord.MessageCreate{
						Content: "That bee doesn't exist.",
					})
				}
				slots, _ := data.OptString("slots")
				level, _ := data.OptInt("level")
				gifted, exists := data.OptBool("gifted")
				if !exists {
					gifted = false
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.AddBee(hive.NewBee(level, name, gifted), slot)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Added bee(s) to hive.",
				})
			},
			"remove": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				slots, _ := data.OptString("slots")
				for _, slot := range GetRangeNumbers(slots) {
					h.RemoveBee(int(slot))
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Removed bee(s) from hive.",
				})
			},
			"giftall": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				for _, bee := range h.GetBees() {
					bee.SetGifted(true)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Gifted all bees in hive.",
				})
			},
			"setbeequip": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				slots, _ := data.OptString("slots")
				name, _ := data.OptString("name")
				if !slices.Contains(loaders.GetBeequips(), name) {
					return event.CreateMessage(discord.MessageCreate{
						Content: "That beequip doesn't exist.",
					})
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.GetBee(slot).SetBeequip(name)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Set beequip of bee(s).",
				})
			},
			"view": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				showHiveNumbers, provided := data.OptBool("show_hive_numbers")
				if !provided {
					showHiveNumbers = true
				}
				dc := gg.NewContext(410, 900)
				h.Draw(dc, showHiveNumbers)
				img := dc.Image()
				r, w := io.Pipe()
				go func(i image.Image) {
					defer w.Close()
					if err := png.Encode(w, i); err != nil {
						panic(err)
					}
				}(img)
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title: event.User().Username + "'s Hive",
							Image: &discord.EmbedResource{
								URL: "attachment://hive.png",
							},
						},
					},
					Files: []*discord.File{
						{
							Name:   "hive.png",
							Reader: r,
						},
					},
				})
			},
			"save": func(event *events.ApplicationCommandInteractionCreate) error {
				h := b.State.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				name, _ := data.OptString("name")
				if name == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide a non-empty name for the hive.",
					})
				}
				userSaveCount, _ := b.Db.Collection("hives").CountDocuments(context.Background(), bson.M{"user_id": event.User().ID})
				if int(userSaveCount) >= MAX_FREE_SAVES {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You have reached the maximum number of free saves. You can get more saves by donating.",
					})
				}
				res, err := b.Db.Collection("hives").UpdateOne(context.Background(), bson.M{
					"user_id": event.User().ID,
					"name":    name,
				}, bson.D{{
					"$set",
					h.ToBson(),
				}}, options.Update().SetUpsert(true))
				if err != nil {
					b.Logger.Error("Error saving hive: ", err)
				}
				id := res.UpsertedID
				if id == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Updated save.",
					})
				}
				oid, _ := id.(primitive.ObjectID)
				hiveId, _ := oid.MarshalText()
				return event.CreateMessage(discord.MessageCreate{
					Content: "Saved hive. ID: `" + string(hiveId) + "`",
				})
			},
			"saves/load": func(event *events.ApplicationCommandInteractionCreate) error {
				//maybe defer?
				data := event.SlashCommandInteractionData()
				id, _ := data.OptString("id")
				if id == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide a non-empty id for the hive.",
					})
				}
				oid, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid id.",
					})
				}
				var h bson.D
				err = b.Db.Collection("hives").FindOne(context.Background(), bson.M{
					"user_id": event.User().ID,
					"_id":     oid,
				}).Decode(&h)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive with that id.",
					})
				}
				userHive := b.State.CreateHive(event.User().ID)
				for _, bh := range h.Map()["bees"].(primitive.D) {
					be := bh.Value.(primitive.D)
					name := be.Map()["name"].(string)
					level := be.Map()["level"].(int32)
					gifted := be.Map()["gifted"].(bool)
					beequip, ok := be.Map()["beequip"].(string)
					if !ok {
						beequip = ""
					}
					bee := hive.NewBee(int(level), name, gifted)
					bee.SetBeequip(beequip)
					i, _ := strconv.ParseInt(bh.Key, 10, 64)
					userHive.AddBee(bee, int(i))
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Loaded hive.",
				})
			},
			"saves/list": func(event *events.ApplicationCommandInteractionCreate) error {
				var results []bson.D
				cur, _ := b.Db.Collection("hives").Find(context.Background(), bson.M{"user_id": event.User().ID})
				err := cur.All(context.Background(), &results)
				if err != nil {
					b.Logger.Error(err)
				}
				if len(results) == 0 {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have any saves.",
					})
				}
				var saves []string
				for i, result := range results {
					id, _ := result.Map()["_id"].(primitive.ObjectID).MarshalText()
					saves = append(saves, fmt.Sprintf("%d. %s (`%s`)", i+1, result.Map()["name"], id))
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Your Saves",
							Description: strings.Join(saves, "\n"),
						},
					},
				})
			},
			"saves/delete": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				id, _ := data.OptString("id")
				if id == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You need to provide the ID of the save you want to delete.",
					})
				}
				oid, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid ID.",
					})
				}
				res, err := b.Db.Collection("hives").DeleteOne(context.Background(), bson.M{
					"_id":     oid,
					"user_id": event.User().ID,
				})
				if err != nil {
					b.Logger.Error(err)
				}
				if res.DeletedCount == 0 {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a save with that ID.",
					})
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Deleted save.",
				})
			},
		},
		AutocompleteHandlers: map[string]handler.AutocompleteHandler{
			"add":        makeAutocompleteHandler(loaders.GetBees()),
			"setbeequip": makeAutocompleteHandler(loaders.GetBeequips()),
		},
	}
}

func makeAutocompleteHandler(b []string) func(*events.AutocompleteInteractionCreate) error {
	return func(event *events.AutocompleteInteractionCreate) error {
		name, _ := event.Data.OptString("name")
		name = strings.ToLower(name)
		matches := make([]discord.AutocompleteChoice, 0)
		i := 0
		for _, bee := range b {
			if i >= 25 {
				break
			}
			if strings.Contains(strings.ToLower(bee), name) {
				matches = append(matches, discord.AutocompleteChoiceString{
					Name:  bee,
					Value: bee,
				})
				i++
			}
		}
		return event.Result(matches)
	}
}

func InitializeHiveCommands(h *handler.Handler, b *Bot) {
	h.AddCommands(HiveCommand(b))
}
