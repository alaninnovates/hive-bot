package hiveplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"context"
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
	"io"
	"strconv"
	"strings"
	"time"
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

var InvalidSlotsMessage = discord.MessageCreate{
	Content: "Invalid slot range. Slots must be whole numbers between 1 and 50.",
}

func ValidateRange(rangeStr string, min int, max int) bool {
	for _, num := range GetRangeNumbers(rangeStr) {
		if num < min || num > max {
			return false
		}
	}
	return true
}

func RenderHiveImage(h *hive.Hive, showHiveNumbers bool) *io.PipeReader {
	dc := gg.NewContext(410, 950)
	h.Draw(dc, showHiveNumbers)
	img := dc.Image()
	bg, _ := gg.LoadImage("assets/bg.png")
	hiveImage := gg.NewContextForImage(bg)
	hiveImage.DrawImageAnchored(img, hiveImage.Width()/2, hiveImage.Height()/2, 0.5, 0.5)
	tmImage, _ := gg.LoadPNG("assets/trademark.png")
	hiveImage.DrawImageAnchored(tmImage, hiveImage.Width()/2, 8, 0.5, 0)
	return common.ImageToPipe(hiveImage.Image())
}

func HiveCommand(b *common.Bot, hiveService *State) handler.Command {
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
							Description: "The level of the bee (0 = no level)",
							Required:    true,
							MinValue:    json.Ptr(0),
							MaxValue:    json.Ptr(25),
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
					Name:        "setlevel",
					Description: "Set the level of ALL your bees",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "level",
							Description: "The level you want to set your bees to",
							Required:    true,
							MinValue:    json.Ptr(0),
							MaxValue:    json.Ptr(25),
						},
					},
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
					Name:        "setmutation",
					Description: "Set the mutation of a bee",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "slots",
							Description: "The slot(s) of the bees you want to set the mutation of",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:         "name",
							Description:  "The mutation you want to set",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "info",
					Description: "Get a summary of the bees in your hive",
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
				hiveService.CreateHive(event.User().ID)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Created new hive. You can now add bees with the `/hive add` command.",
				})
			},
			"add": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				name, _ := data.OptString("name")
				if !slices.Contains(loaders.GetBeeNames(), name) {
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
				if !ValidateRange(slots, 1, 50) {
					return event.CreateMessage(InvalidSlotsMessage)
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.AddBee(hive.NewBee(level, loaders.GetBeeId(name), gifted), slot)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Added bee(s) to hive.\nYou can now see your hive by using `/hive view` command",
				})
			},
			"remove": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				slots, _ := data.OptString("slots")
				if !ValidateRange(slots, 1, 50) {
					return event.CreateMessage(InvalidSlotsMessage)
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.RemoveBee(slot)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Removed bee(s) from hive.",
				})
			},
			"giftall": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
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
			"setlevel": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				level, _ := data.OptInt("level")
				for _, bee := range h.GetBees() {
					bee.SetLevel(level)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Set level of all bees in hive to " + strconv.Itoa(level) + ".",
				})
			},
			"setbeequip": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
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
				if !ValidateRange(slots, 1, 50) {
					return event.CreateMessage(InvalidSlotsMessage)
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.GetBee(slot).SetBeequip(name)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Set beequip of bee(s).",
				})
			},
			"setmutation": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				data := event.SlashCommandInteractionData()
				slots, _ := data.OptString("slots")
				name, _ := data.OptString("name")
				if !slices.Contains(loaders.GetMutations(), name) {
					return event.CreateMessage(discord.MessageCreate{
						Content: "That mutation doesn't exist.",
					})
				}
				if !ValidateRange(slots, 1, 50) {
					return event.CreateMessage(InvalidSlotsMessage)
				}
				for _, slot := range GetRangeNumbers(slots) {
					h.GetBee(slot).SetMutation(name)
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Set mutation of bee(s).",
				})
			},
			"info": func(event *events.ApplicationCommandInteractionCreate) error {
				bees := make(map[string]int)
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				for _, bee := range h.GetBees() {
					bees[bee.Name()]++
				}
				beesStr := ""
				for name, count := range bees {
					beesStr += fmt.Sprintf("%s: %d\n", name, count)
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Hive Info",
							Description: beesStr,
							Color:       0x00FF00,
						},
					},
				})
			},
			"view": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
				if h == nil {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have a hive. Create one with the `/hive create` command.",
					})
				}
				err := event.DeferCreateMessage(false)
				if err != nil {
					return err
				}
				data := event.SlashCommandInteractionData()
				showHiveNumbers, provided := data.OptBool("show_hive_numbers")
				if !provided {
					showHiveNumbers = true
				}
				r := RenderHiveImage(h, showHiveNumbers)
				hn := ""
				if showHiveNumbers {
					hn = "1"
				} else {
					hn = "0"
				}
				_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
					Embeds: &[]discord.Embed{
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
					Components: &[]discord.ContainerComponent{
						discord.ActionRowComponent{
							discord.NewPrimaryButton("Gift All", "handler:giftall:"+event.User().ID.String()+":"+hn),
							//discord.NewPrimaryButton("Set Level All", "handler:setlevelall:"+event.User().ID.String()+":"+hn),
							discord.NewSuccessButton("Hive Info", "handler:info:"+event.User().ID.String()),
						},
					},
				})
				return err
			},
			"save": func(event *events.ApplicationCommandInteractionCreate) error {
				h := hiveService.GetHive(event.User().ID)
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
				if int(userSaveCount) >= common.MaxFreeSaves {
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
				userHive := hiveService.CreateHive(event.User().ID)
				for _, bh := range h.Map()["bees"].(primitive.D) {
					be := bh.Value.(primitive.D)
					id := be.Map()["id"].(string)
					level := be.Map()["level"].(int32)
					gifted := be.Map()["gifted"].(bool)
					beequip, ok := be.Map()["beequip"].(string)
					mutation, ok2 := be.Map()["mutation"].(string)
					if !ok {
						beequip = ""
					}
					if !ok2 {
						mutation = "None"
					}
					bee := hive.NewBee(int(level), id, gifted)
					bee.SetBeequip(beequip)
					bee.SetMutation(mutation)
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
			"add":         makeAutocompleteHandler(loaders.GetBeeNames()),
			"setbeequip":  makeAutocompleteHandler(loaders.GetBeequips()),
			"setmutation": makeAutocompleteHandler(loaders.GetMutations()),
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

func GiftAllButton(b *common.Bot, hiveService *State) handler.Component {
	return handler.Component{
		Name:  "giftall",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			err := event.DeferUpdateMessage()
			if err != nil {
				return err
			}
			data := strings.Split(event.ButtonInteractionData().CustomID(), ":")
			uid, shn := data[2], data[3]
			userId, _ := snowflake.Parse(uid)
			h := hiveService.GetHive(userId)
			if h == nil {
				return event.UpdateMessage(discord.MessageUpdate{
					Content:     json.Ptr("Your hive seems to have gone missing... Create a new one with `/hive create`"),
					Embeds:      &[]discord.Embed{},
					Components:  &[]discord.ContainerComponent{},
					Attachments: &[]discord.AttachmentUpdate{},
				})
			}
			for _, b := range h.GetBees() {
				b.SetGifted(true)
			}
			_, err = b.Client.Rest().CreateMessage(event.ChannelID(), discord.MessageCreate{
				Content: "Gifted all bees in your hive!",
			})
			if err != nil {
				b.Logger.Error(err)
				return err
			}
			showHiveNumbers := false
			if shn == "1" {
				showHiveNumbers = true
			} else if shn == "0" {
				showHiveNumbers = false
			}
			r := RenderHiveImage(h, showHiveNumbers)
			_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
				Files: []*discord.File{
					{
						Name:   "hive.png",
						Reader: r,
					},
				},
			})
			return err
		},
	}
}

func SetLevelButton(b *common.Bot, hiveService *State) handler.Component {
	return handler.Component{
		Name:  "setlevelall",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			go func() error {
				data := strings.Split(event.ButtonInteractionData().CustomID(), ":")
				uid, shn := data[2], data[3]
				userId, _ := snowflake.Parse(uid)
				h := hiveService.GetHive(userId)
				if h == nil {
					return event.UpdateMessage(discord.MessageUpdate{
						Content:     json.Ptr("Your hive seems to have gone missing... Create a new one with `/hive create`"),
						Embeds:      &[]discord.Embed{},
						Components:  &[]discord.ContainerComponent{},
						Attachments: &[]discord.AttachmentUpdate{},
					})
				}
				err := event.CreateMessage(discord.MessageCreate{
					Content: "What level do you want to set your hive to?\nYou have 10 seconds to reply",
				})
				if err != nil {
					b.Logger.Error(err)
					return err
				}
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer func() {
					cancel()
					println("cancelled")
				}()
				print("waiting for message")
				level := 0
				bot.WaitForEvent(b.Client, ctx,
					func(e *events.MessageCreate) bool {
						println(e.Message.Content)
						return e.Message.Author.ID == event.User().ID
					}, func(e *events.MessageCreate) {
						println(e.Message.Content)
						l, err := strconv.Atoi(e.Message.Content)
						if err != nil || l < 0 || l > 25 {
							_, _ = b.Client.Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
								Content: "Please input an integer between 0-25.",
							})
							return
						}
						level = l
					}, func() {
						println("timeout")
					})
				print("after: ", level)
				for _, bee := range h.GetBees() {
					bee.SetLevel(level)
				}
				showHiveNumbers := false
				if shn == "1" {
					showHiveNumbers = true
				} else if shn == "0" {
					showHiveNumbers = false
				}
				r := RenderHiveImage(h, showHiveNumbers)
				_, err = b.Client.Rest().UpdateMessage(event.ChannelID(), event.Message.ID, discord.MessageUpdate{
					Files: []*discord.File{
						{
							Name:   "hive.png",
							Reader: r,
						},
					},
				})
				return err
			}()
			return nil
		},
	}
}

func HiveInfoButton(hiveService *State) handler.Component {
	return handler.Component{
		Name:  "info",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			data := strings.Split(event.ButtonInteractionData().CustomID(), ":")
			uid := data[2]
			userId, _ := snowflake.Parse(uid)
			h := hiveService.GetHive(userId)
			if h == nil {
				return event.UpdateMessage(discord.MessageUpdate{
					Content:     json.Ptr("Your hive seems to have gone missing... Create a new one with `/hive create`"),
					Embeds:      &[]discord.Embed{},
					Components:  &[]discord.ContainerComponent{},
					Attachments: &[]discord.AttachmentUpdate{},
				})
			}
			bees := make(map[string]int)
			for _, bee := range h.GetBees() {
				bees[bee.Name()]++
			}
			beesStr := ""
			for name, count := range bees {
				beesStr += fmt.Sprintf("%s: %d\n", name, count)
			}
			return event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					{
						Title:       "Hive Info",
						Description: beesStr,
						Color:       0x00FF00,
					},
				},
			})
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	hiveService := NewHiveService()
	h.AddCommands(HiveCommand(b, hiveService))
	h.AddComponents(GiftAllButton(b, hiveService), SetLevelButton(b, hiveService), HiveInfoButton(hiveService))
}
