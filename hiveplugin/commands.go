package hiveplugin

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

func CreateCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		hiveService.CreateHive(event.User().ID)
		return event.CreateMessage(discord.MessageCreate{
			Content: "Created new hive. You can now add bees with the `/hive add` command.",
		})
	}
}

func ImportCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
	}
}

func AddCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func RemoveCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func GiftAllCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
		for _, bees := range h.GetBees() {
			for _, bee := range bees {
				bee.SetGifted(true)
			}
		}
		return event.CreateMessage(discord.MessageCreate{
			Content: "Gifted all bees in hive.",
		})
	}
}

func SetLevelCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func SetBeequipCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func SetMutationCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func InfoCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		bees := make(map[string]int)
		h := hiveService.GetHive(event.User().ID)
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
					discord.NewSuccessButton("Hive Info", fmt.Sprintf("/hive/buttons/hiveinfo/%s", event.User().ID)),
					discord.NewSuccessButton("Mutation Info", fmt.Sprintf("/hive/buttons/mutationinfo/%s", event.User().ID)),
				},
			},
		})
	}
}

func ViewCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
					discord.NewPrimaryButton("Add Bee", fmt.Sprintf("/hive/buttons/addbee/%s", event.User().ID)),
					discord.NewPrimaryButton("Gift All", fmt.Sprintf("/hive/buttons/giftall/%s/%s", event.User().ID, hn)),
					discord.NewPrimaryButton("Set Level All", fmt.Sprintf("/hive/buttons/setlevel/%s", event.User().ID)),
					discord.NewSuccessButton("Hive Info", fmt.Sprintf("/hive/buttons/hiveinfo/%s", event.User().ID)),
					discord.NewSecondaryButton("Rerender Hive", fmt.Sprintf("/hive/buttons/rerender/%s", event.User().ID)),
				},
			},
		})
		return err
	}
}

func SaveCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func SavesLoadCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
	}
}

func SavesListCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
			row = row.AddComponents(discord.NewPrimaryButton(name, fmt.Sprintf("/hive/saves/saveid/%s", id)))
			if i%4 == 0 && i != 0 {
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
	}
}

func SavesDeleteCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
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
	}
}

func PostCommand(b *common.Bot, hiveService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		data := event.SlashCommandInteractionData()
		id, _ := data.OptString("id")
		title, _ := data.OptString("title")
		content, _ := data.OptString("content")
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
		r := RenderHiveImage(userHive, false, false, make([]int, 0), "default")
		link, err := b.R2.UploadImage(fmt.Sprintf("hive-%s.png", id), r)
		if err != nil {
			b.Logger.Error("Error uploading hive image: ", err)
			return event.CreateMessage(discord.MessageCreate{
				Content: "Something went wrong while uploading the hive image. Try again later.",
			})
		}
		post := database.Post{
			ID:        primitive.NewObjectID(),
			Title:     title,
			Content:   content,
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			HiveId:    oid,
			ImageUrl:  link,
		}
		res, err := b.Db.Collection("posts").InsertOne(context.Background(), post)
		if err != nil {
			b.Logger.Error("Error creating post: ", err)
			return event.CreateMessage(discord.MessageCreate{
				Content: "Something went wrong while creating the post.",
			})
		}
		insertedId := res.InsertedID.(primitive.ObjectID)
		insertedStr, _ := insertedId.MarshalText()
		return event.CreateMessage(discord.MessageCreate{
			Content: "Created post. Access it at " + fmt.Sprintf("https://hives.meta-bee.com/posts/%s", insertedStr),
		})
	}
}
