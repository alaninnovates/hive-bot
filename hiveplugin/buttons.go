package hiveplugin

import (
	"fmt"
	"strings"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

func AddBeeButton() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		return event.Modal(discord.NewModalCreateBuilder().
			SetTitle("Add a bee to your hive").
			SetCustomID(fmt.Sprintf("/hive/modals/addbee/%s", event.User().ID)).
			AddActionRow(discord.NewShortTextInput("name", "Bee name")).
			AddActionRow(discord.NewShortTextInput("slots", "Slots")).
			AddActionRow(discord.NewShortTextInput("level", "Level")).
			Build())
	}
}

func GiftAllButton(b *common.Bot, hiveService *State) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		err := event.DeferUpdateMessage()
		if err != nil {
			return err
		}
		shn := event.Vars["showHiveNumbers"]
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func SetLevelButton() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		return event.Modal(discord.ModalCreate{
			Title:    "Set level of all bees",
			CustomID: fmt.Sprintf("/hive/modals/setlevel/%s", event.User().ID.String()),
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
	}
}

func HiveInfoButton(hiveService *State) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
					discord.NewSuccessButton("Hive Info", fmt.Sprintf("/hive/buttons/hiveinfo/%s", event.User().ID)),
					discord.NewSuccessButton("Mutation Info", fmt.Sprintf("/hive/buttons/mutationinfo/%s", event.User().ID)),
				},
			},
		})
	}
}

func MutationInfoButton(hiveService *State) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
	}
}

func HiveRerenderButton(b *common.Bot, hiveService *State) handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		h := hiveService.GetHive(event.User().ID)
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
			if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "Missing Access") {
				cause = "I didn't have permission to edit that message! Please make sure I have edit message permissions in this channel."
			} else {
				cause = "Something went wrong!"
				b.Logger.Error("Error rerendering hive: ", err)
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
	}
}

func SaveIdButton() handler.ButtonComponentHandler {
	return func(data discord.ButtonInteractionData, event *handler.ComponentEvent) error {
		id := event.Vars["id"]
		return event.CreateMessage(discord.MessageCreate{Content: id})
	}
}
