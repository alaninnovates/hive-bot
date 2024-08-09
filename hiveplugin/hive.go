package hiveplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"context"
	"fmt"
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

func RenderHiveImage(h *hive.Hive, showHiveNumbers bool, slotsOnTop bool, skipHiveNumbers []int, background string) *io.PipeReader {
	dc := gg.NewContext(410, 950)
	hive.DrawHive(h, dc, showHiveNumbers, slotsOnTop, skipHiveNumbers)
	img := dc.Image()
	bg, _ := gg.LoadImage(loaders.GetHiveBackgroundImagePath(background))
	hiveImage := gg.NewContextForImage(bg)
	hiveImage.DrawImageAnchored(img, hiveImage.Width()/2, hiveImage.Height()/2, 0.5, 0.5)
	//tmImage, _ := gg.LoadPNG("assets/trademark.png")
	//hiveImage.DrawImageAnchored(tmImage, hiveImage.Width()/2, 8, 0.5, 0)
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
					Name:        "import",
					Description: "Import a hive from Natro Macro Hive Generation",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "input",
							Description: "The output from the Natro Macro Hive Generation",
							Required:    true,
						},
					},
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
						discord.ApplicationCommandOptionBool{
							Name:        "slots_on_top",
							Description: "Whether or not to show slot numbers on top of bees",
							Required:    false,
						},
						discord.ApplicationCommandOptionString{
							Name:         "ability",
							Description:  "Show only bees that have a certain ability",
							Required:     false,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionString{
							Name:         "beequip",
							Description:  "Show only bees that have a certain beequip",
							Required:     false,
							Autocomplete: true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "slots",
							Description: "Show only bees that are in the selected slots",
							Required:    false,
						},
						discord.ApplicationCommandOptionString{
							Name:        "background",
							Description: "What you want the background of your hive to be",
							Required:    false,
							Choices:     loaders.GetHiveBackgroundsChoices(),
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
							MaxLength:   json.Ptr(30),
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
			"import": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				jsonData := data.String("input")
				if !ValidateHiveJson(jsonData) {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid hive input. Please make sure the data is exactly what is outputted by Natro.",
					})
				}
				parsedJson := ParseHiveJson(jsonData)
				h := hiveService.CreateHive(event.User().ID)
				slotNum := 1
				for beeName, beeData := range parsedJson {
					var bd BeeData
					switch beeData.(type) {
					case map[string]interface{}:
						bd.Gifted = beeData.(map[string]interface{})["gifted"].(bool)
						bd.Amount = int(beeData.(map[string]interface{})["amount"].(float64))
					case string:
						//fmt.Println("type is string")
						continue
					}
					for i := 0; i < bd.Amount; i++ {
						h.AddBee(hive.NewBee(0, beeName, bd.Gifted), slotNum)
						slotNum++
					}
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Imported hive.",
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
				name := getBeeFromAbbr(data.String("name"))
				if name == "" {
					return event.CreateMessage(discord.MessageCreate{
						Content: "Invalid bee name.",
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
				slotClearWarning := false
				for _, slot := range GetRangeNumbers(slots) {
					if len(h.GetBeesAt(slot)) > 1 {
						slotClearWarning = true
						// clear the first bee
						h.RemoveBeeAt(slot, 0)
					}
					h.AddBee(hive.NewBee(level, loaders.GetBeeId(name), gifted), slot)
				}
				messageAppend := ""
				if slotClearWarning {
					messageAppend = "\nWarning: Some slots had more than two bees specified! The first bee was removed. To clear slots, use the `/hive remove` command"
				}
				return event.CreateMessage(discord.MessageCreate{
					Content: "Added bee(s) to hive.\nYou can now see your hive by using `/hive view` command" + messageAppend,
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
				for _, bees := range h.GetBees() {
					for _, bee := range bees {
						bee.SetGifted(true)
					}
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
				for _, bees := range h.GetBees() {
					for _, bee := range bees {
						bee.SetLevel(level)
					}
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
				if !(slices.Contains(loaders.GetBeequips(), name) || name == "None") {
					return event.CreateMessage(discord.MessageCreate{
						Content: "That beequip doesn't exist.",
					})
				}
				if !ValidateRange(slots, 1, 50) {
					return event.CreateMessage(InvalidSlotsMessage)
				}
				for _, slot := range GetRangeNumbers(slots) {
					bees := h.GetBeesAt(slot)
					for _, bee := range bees {
						bee.SetBeequip(name)
					}
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
					bees := h.GetBeesAt(slot)
					for _, bee := range bees {
						bee.SetMutation(name)
					}
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
				for _, beeslist := range h.GetBees() {
					for _, bee := range beeslist {
						bees[bee.Name()]++
					}
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
					Components: []discord.ContainerComponent{
						discord.ActionRowComponent{
							discord.NewSuccessButton("Hive Info", "handler:info:"+event.User().ID.String()),
							discord.NewSuccessButton("Mutation Info", "handler:mutationinfo:"+event.User().ID.String()),
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
				slotsOnTop, provided := data.OptBool("slots_on_top")
				if !provided {
					slotsOnTop = false
				}
				var skipHiveNumbers []int
				var includedHiveNumbers []int
				// START: render by ability
				renderAbility, provided := data.OptString("ability")
				if provided {
					premiumLevel, err := common.GetPremiumLevel(b.Db, event.User().ID)
					if err != nil {
						return err
					}
					if premiumLevel < common.PremiumLevelBuilder {
						_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
							Embeds: &[]discord.Embed{
								{
									Title:       "Premium Only Feature",
									Description: "This feature is only available to :sparkles: premium users. You can get premium by donating [here](https://meta-bee.my.to/donate)!",
									Color:       0x800080,
								},
							},
						})
						return err
					}
					for i, bees := range h.GetBees() {
						oneBeeFlag := false
						for _, bee := range bees {
							if !common.ArrayIncludes(loaders.GetBeeAbilities(bee.Name()), renderAbility) {
								oneBeeFlag = false
							} else {
								oneBeeFlag = true
							}
						}
						if oneBeeFlag {
							includedHiveNumbers = append(includedHiveNumbers, i)
						} else {
							skipHiveNumbers = append(skipHiveNumbers, i)
						}
					}
					showHiveNumbers = false
				} // END: render by ability
				// START: render by beequip
				renderBeequip, provided := data.OptString("beequip")
				if provided {
					premiumLevel, err := common.GetPremiumLevel(b.Db, event.User().ID)
					if err != nil {
						return err
					}
					if premiumLevel < common.PremiumLevelBuilder {
						_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
							Embeds: &[]discord.Embed{
								{
									Title:       "Premium Only Feature",
									Description: "This feature is only available to :sparkles: premium users. You can get premium by donating [here](https://meta-bee.my.to/donate)!",
									Color:       0x800080,
								},
							},
						})
						return err
					}
					for i, bees := range h.GetBees() {
						for _, bee := range bees {
							if renderBeequip != bee.Beequip() {
								skipHiveNumbers = append(skipHiveNumbers, i)
							} else {
								includedHiveNumbers = append(includedHiveNumbers, i)
							}
						}
					}
					showHiveNumbers = false
				} // END: render by beequip
				// START: render by slots
				slots, provided := data.OptString("slots")
				if provided {
					premiumLevel, err := common.GetPremiumLevel(b.Db, event.User().ID)
					if err != nil {
						return err
					}
					if premiumLevel < common.PremiumLevelBuilder {
						_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
							Embeds: &[]discord.Embed{
								{
									Title:       "Premium Only Feature",
									Description: "This feature is only available to :sparkles: premium users. You can get premium by donating [here](https://meta-bee.my.to/donate)!",
									Color:       0x800080,
								},
							},
						})
						return err
					}
					if !ValidateRange(slots, 1, 50) {
						return event.CreateMessage(InvalidSlotsMessage)
					}
					slotRange := GetRangeNumbers(slots)
					for i := 1; i <= 50; i++ {
						if !common.ArrayIncludes(slotRange, i) && !common.ArrayIncludes(includedHiveNumbers, i) {
							skipHiveNumbers = append(skipHiveNumbers, i)
						}
					}
					showHiveNumbers = false
				} // END: render by slots
				// START: render background
				background, provided := data.OptString("background")
				if !provided {
					background = "default"
				}
				if provided {
					premiumLevel, err := common.GetPremiumLevel(b.Db, event.User().ID)
					if err != nil {
						return err
					}
					if premiumLevel < common.PremiumLevelBuilder {
						_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
							Embeds: &[]discord.Embed{
								{
									Title:       "Premium Only Feature",
									Description: "This feature is only available to :sparkles: premium users. You can get premium by donating [here](https://meta-bee.my.to/donate)!",
									Color:       0x800080,
								},
							},
						})
						return err
					}
				} // END: render background
				r := RenderHiveImage(h, showHiveNumbers, slotsOnTop, skipHiveNumbers, background)
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
							discord.NewPrimaryButton("Add Bee", "handler:addbee:"+event.User().ID.String()),
							discord.NewPrimaryButton("Gift All", "handler:giftall:"+event.User().ID.String()+":"+hn),
							discord.NewPrimaryButton("Set Level All", "handler:setlevelbutton:"+event.User().ID.String()),
							discord.NewSuccessButton("Hive Info", "handler:info:"+event.User().ID.String()),
							discord.NewSecondaryButton("Rerender Hive", "handler:rerender:"+event.User().ID.String()),
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
				premiumLevel, err := common.GetPremiumLevel(b.Db, event.User().ID)
				if err != nil {
					return err
				}
				userSaveCount, _ := b.Db.Collection("hives").CountDocuments(context.Background(), bson.M{"user_id": event.User().ID})
				maxSaves := common.MaxFreeSaves
				if premiumLevel >= common.PremiumLevelBuilder {
					maxSaves = common.MaxPremiumSaves
				}
				if int(userSaveCount) >= maxSaves {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You have reached the maximum number of free saves.",
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
				for _, bl := range h.Map()["bees"].(primitive.D) {
					key := bl.Key
					for _, bh := range bl.Value.(primitive.A) {
						be := bh.(primitive.D)
						id := be.Map()["id"].(string)
						level := be.Map()["level"].(int32)
						gifted := be.Map()["gifted"].(bool)
						beequip, ok := be.Map()["beequip"].(string)
						mutation, ok2 := be.Map()["mutation"].(string)
						if !ok || beequip == "" {
							beequip = "None"
						}
						if !ok2 {
							mutation = "None"
						}
						bee := hive.NewBee(int(level), id, gifted)
						bee.SetBeequip(beequip)
						bee.SetMutation(mutation)
						i, _ := strconv.ParseInt(key, 10, 64)
						userHive.AddBee(bee, int(i))
					}
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
					b.Logger.Error("Failed to list hive saves for user: %v", err)
				}
				if len(results) == 0 {
					return event.CreateMessage(discord.MessageCreate{
						Content: "You don't have any saves.",
					})
				}
				var saves []string
				var rows []discord.ContainerComponent
				row := discord.ActionRowComponent{}
				for i, result := range results {
					id, _ := result.Map()["_id"].(primitive.ObjectID).MarshalText()
					name := result.Map()["name"].(string)
					saves = append(saves, fmt.Sprintf("%d. %s (`%s`)", i+1, name, id))
					row = row.AddComponents(discord.NewPrimaryButton(name, fmt.Sprintf("handler:save-id:%s:%s", event.User().ID.String(), id)))
					if i%5 == 0 && i != 0 {
						rows = append(rows, row)
						row = discord.ActionRowComponent{}
					}
				}
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title:       "Your Saves",
							Description: strings.Join(saves, "\n"),
							Footer: &discord.EmbedFooter{
								Text: "Press the buttons below to get mobile friendly ids",
							},
						},
					},
					Components: rows,
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
					b.Logger.Error("Error deleting hive save: ", err)
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
			"setbeequip":  makeAutocompleteHandler(append(loaders.GetBeequips(), "None")),
			"setmutation": makeAutocompleteHandler(loaders.GetMutations()),
			"view":        hiveViewAutocompleteHandler,
		},
	}
}

func hiveViewAutocompleteHandler(event *events.AutocompleteInteractionCreate) error {
	focusedField := ""
	for key, option := range event.Data.Options {
		if option.Focused {
			focusedField = key
			break
		}
	}
	if focusedField == "ability" {
		ability := event.Data.String("ability")
		return getMatches(event, loaders.GetBeeAbilityList(), ability)
	} else if focusedField == "beequip" {
		beequip := event.Data.String("beequip")
		return getMatches(event, loaders.GetBeequips(), beequip)
	}
	return nil
}

func makeAutocompleteHandler(b []string) func(*events.AutocompleteInteractionCreate) error {
	return func(event *events.AutocompleteInteractionCreate) error {
		//fmt.Printf("evt: %d now: %d", event.ID().Time().UnixMilli(), time.Now().UnixMilli())
		name := event.Data.String("name")
		return getMatches(event, b, name)
	}
}

func getMatches(event *events.AutocompleteInteractionCreate, options []string, text string) error {
	text = strings.ToLower(text)
	matches := make([]discord.AutocompleteChoice, 0)
	i := 0
	for _, opt := range options {
		if i >= 25 {
			break
		}
		if strings.Contains(strings.ToLower(opt), text) {
			matches = append(matches, discord.AutocompleteChoiceString{
				Name:  opt,
				Value: opt,
			})
			i++
		}
	}
	return event.AutocompleteResult(matches)
}

func AddBeeButton() handler.Component {
	return handler.Component{
		Name:  "addbee",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			return event.Modal(discord.NewModalCreateBuilder().
				SetCustomID("handler:addbeemodal").
				AddActionRow(discord.NewShortTextInput("name", "Bee name")).
				AddActionRow(discord.NewShortTextInput("slots", "Slots")).
				AddActionRow(discord.NewShortTextInput("level", "Level")).
				Build())
		},
	}
}

func getBeeFromAbbr(name string) string {
	for _, bee := range loaders.GetBeeNames() {
		if strings.ToLower(bee) == strings.ToLower(name) ||
			strings.ToLower(bee) == strings.ToLower(name)+" bee" {
			return bee
		}
	}
	return ""
}

func AddBeeModal(hiveService *State) handler.Modal {
	return handler.Modal{
		Name: "addbeemodal",
		Handler: func(event *events.ModalSubmitInteractionCreate) error {
			h := hiveService.GetHive(event.User().ID)
			if h == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Your hive seems to have gone missing... Create a new one with `/hive create`",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			name, _ := event.Data.OptText("name")
			levelStr, _ := event.Data.OptText("level")
			slots, _ := event.Data.OptText("slots")
			if getBeeFromAbbr(name) == "" {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Invalid bee name.",
				})
			}
			level, err := strconv.Atoi(levelStr)
			if err != nil || level < 1 || level > 25 {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Invalid level.",
				})
			}
			if !ValidateRange(slots, 1, 50) {
				return event.CreateMessage(InvalidSlotsMessage)
			}
			beeName := getBeeFromAbbr(name)
			for _, slot := range GetRangeNumbers(slots) {
				h.AddBee(hive.NewBee(level, loaders.GetBeeId(beeName), false), slot)
			}
			return event.CreateMessage(discord.MessageCreate{
				Content: "Added bee.",
				Flags:   discord.MessageFlagEphemeral,
			})
		},
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
				_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
					Content:     json.Ptr("Your hive seems to have gone missing... Create a new one with `/hive create`"),
					Embeds:      &[]discord.Embed{},
					Components:  &[]discord.ContainerComponent{},
					Attachments: &[]discord.AttachmentUpdate{},
				})
				return err
			}
			for _, beesList := range h.GetBees() {
				for _, bee := range beesList {
					bee.SetGifted(true)
				}
			}
			showHiveNumbers := false
			if shn == "1" {
				showHiveNumbers = true
			} else if shn == "0" {
				showHiveNumbers = false
			}
			r := RenderHiveImage(h, showHiveNumbers, false, make([]int, 0), "default")
			_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
				Files: []*discord.File{
					{
						Name:   "hive.png",
						Reader: r,
					},
				},
			})
			_, _ = b.Client.Rest().CreateFollowupMessage(b.Client.ApplicationID(), event.Token(), discord.MessageCreate{
				Content: "Gifted all bees in your hive!",
				Flags:   discord.MessageFlagEphemeral,
			})
			return err
		},
	}
}

func SetLevelButton() handler.Component {
	return handler.Component{
		Name:  "setlevelbutton",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			return event.Modal(discord.ModalCreate{
				Title:    "Set level of all bees",
				CustomID: "handler:setlevelmodal",
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.TextInputComponent{
							CustomID:    "level",
							Style:       discord.TextInputStyleShort,
							Label:       "Level",
							MinLength:   json.Ptr(1),
							MaxLength:   2,
							Placeholder: "Value from 1-25",
						},
					},
				},
			})
		},
	}
}

func SetLevelModal(hiveService *State) handler.Modal {
	return handler.Modal{
		Name: "setlevelmodal",
		Handler: func(event *events.ModalSubmitInteractionCreate) error {
			levelStr := event.Data.Text("level")
			levelInt, err := strconv.Atoi(levelStr)
			if err != nil || levelInt < 1 || levelInt > 25 {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Level must be an integer between 1 and 25",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			h := hiveService.GetHive(event.User().ID)
			if h == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "Your hive seems to have gone missing... Create a new one with `/hive create`",
					Flags:   discord.MessageFlagEphemeral,
				})
			}
			for _, beesList := range h.GetBees() {
				for _, bee := range beesList {
					bee.SetLevel(levelInt)
				}
			}
			return event.CreateMessage(discord.MessageCreate{
				Content: "Set level of all bees to " + levelStr,
				Flags:   discord.MessageFlagEphemeral,
			})
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
			bees := make(map[string]struct {
				Count    int
				Beequips map[string]int
			})
			for _, beeslist := range h.GetBees() {
				for _, bee := range beeslist {
					if b, ok := bees[bee.Name()]; ok {
						b.Count++
						if bee.Beequip() != "None" {
							b.Beequips[bee.Beequip()]++
						}
						bees[bee.Name()] = b
					} else {
						bees[bee.Name()] = struct {
							Count    int
							Beequips map[string]int
						}{
							Count:    1,
							Beequips: map[string]int{},
						}
						if bee.Beequip() != "None" {
							bees[bee.Name()].Beequips[bee.Beequip()]++
						}
					}
				}
			}
			beesStr := ""
			for name, info := range bees {
				toAppend := fmt.Sprintf("%s: %d", name, info.Count)
				if len(info.Beequips) > 0 {
					toAppend += " ("
					for beequip, count := range info.Beequips {
						toAppend += fmt.Sprintf("%s: %d, ", beequip, count)
					}
					toAppend = toAppend[:len(toAppend)-2] + ")"
				}
				beesStr += toAppend + "\n"
			}
			return event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					{
						Title:       "Hive Info",
						Description: beesStr,
						Color:       0x00FF00,
					},
				},
				Components: []discord.ContainerComponent{
					discord.ActionRowComponent{
						discord.NewSuccessButton("Hive Info", "handler:info:"+event.User().ID.String()),
						discord.NewSuccessButton("Mutation Info", "handler:mutationinfo:"+event.User().ID.String()),
					},
				},
			})
		},
	}
}

func MutationInfoButton(hiveService *State) handler.Component {
	return handler.Component{
		Name:  "mutationinfo",
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
			mutationData := make(map[string]map[string]int)
			for _, beeslist := range h.GetBees() {
				for _, bee := range beeslist {
					if bee.Mutation() == "None" {
						continue
					}
					if mutationData[bee.Mutation()] == nil {
						mutationData[bee.Mutation()] = make(map[string]int)
					}
					mutationData[bee.Mutation()][bee.Name()]++
				}
			}
			content := ""
			for mutationName, info := range mutationData {
				content += mutationName + ": "
				var strs []string
				for name, amt := range info {
					strs = append(strs, fmt.Sprintf(" x%d %s", amt, name))
				}
				content += strings.Join(strs, ",") + "\n"
			}
			return event.UpdateMessage(discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					{
						Title:       "Mutation Info",
						Description: content,
						Color:       0x00FF00,
					},
				},
			})
		},
	}
}

func SaveIdButton() handler.Component {
	return handler.Component{
		Name:  "save-id",
		Check: common.UserIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			id := strings.Split(event.ButtonInteractionData().CustomID(), ":")[3]
			return event.CreateMessage(discord.MessageCreate{Content: id})
		},
	}
}

func HiveRerenderButton(b *common.Bot, hiveService *State) handler.Component {
	return handler.Component{
		Name:  "rerender",
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
			err := event.DeferCreateMessage(true)
			if err != nil {
				return err
			}
			message := event.Message
			r := RenderHiveImage(h, true, false, make([]int, 0), "default")
			_, err = b.Client.Rest().UpdateMessage(message.ChannelID, message.ID, discord.MessageUpdate{
				Files: []*discord.File{
					{
						Name:   "hive.png",
						Reader: r,
					},
				},
			})
			if err != nil {
				cause := ""
				if strings.Contains(err.Error(), "403") {
					cause = "I didn't have permission to edit that message!"
				} else {
					cause = "Something went wrong!"
				}
				_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
					Content: json.Ptr("Failed to re-render hive: " + cause),
				})
				return err
			}
			_, err = b.Client.Rest().UpdateInteractionResponse(b.Client.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr("Rerendered hive"),
			})
			return err
		},
	}
}

func Initialize(h *handler.Handler, b *common.Bot, hiveService *State) {
	h.AddCommands(HiveCommand(b, hiveService))
	h.AddComponents(AddBeeButton(), GiftAllButton(b, hiveService), SetLevelButton(),
		HiveInfoButton(hiveService), MutationInfoButton(hiveService),
		SaveIdButton(), HiveRerenderButton(b, hiveService))
	h.AddModals(AddBeeModal(hiveService), SetLevelModal(hiveService))
}
