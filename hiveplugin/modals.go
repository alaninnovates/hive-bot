package hiveplugin

import (
	"strconv"
	"strings"

	"alaninnovates.com/hive-bot/common/loaders"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func getBeeFromAbbr(name string) string {
	for _, bee := range loaders.GetBeeNames() {
		if strings.ToLower(bee) == strings.ToLower(name) ||
			strings.ToLower(bee) == strings.ToLower(name)+" bee" {
			return bee
		}
	}
	return ""
}

func AddBeeModal(hiveService *State) handler.ModalHandler {
	return func(event *handler.ModalEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func SetLevelModal(hiveService *State) handler.ModalHandler {
	return func(event *handler.ModalEvent) error {
		levelStr := event.Data.Text("level")
		levelInt, err := strconv.Atoi(levelStr)
		if err != nil || levelInt < 1 || levelInt > 25 {
			return event.CreateMessage(discord.MessageCreate{
				Content: "Level must be an integer between 1 and 25",
				Flags:   discord.MessageFlagEphemeral,
			})
		}
		h := hiveService.GetHive(event.User().ID)
		for _, beesList := range h.GetBees() {
			for _, bee := range beesList {
				bee.SetLevel(levelInt)
			}
		}
		return event.CreateMessage(discord.MessageCreate{
			Content: "Set level of all bees to " + levelStr,
			Flags:   discord.MessageFlagEphemeral,
		})
	}
}
