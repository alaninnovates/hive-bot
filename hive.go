package main

import (
	"alaninnovates.com/hive-bot/hive"
	"alaninnovates.com/hive-bot/loaders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/fogleman/gg"
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
							Description: "The slot(s) of the loaders you want to remove",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "giftall",
					Description: "Gift all loaders in your hive",
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
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"create": func(event *events.ApplicationCommandInteractionCreate) error {
				b.State.CreateHive(event.User().ID)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Created new hive. You can now add loaders with the `/hive add` command.",
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
					Content: "Gifted all loaders in hive.",
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
				dc := gg.NewContext(400, 900)
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
