package hiveplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var HiveCommandCreate = discord.SlashCommandCreate{
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
}

func Initialize(r *handler.Mux, b *common.Bot, hiveService *State) {
	r.Route("/hive", func(r handler.Router) {
		r.Command("/create", CreateCommand(b, hiveService))
		r.Command("/import", ImportCommand(b, hiveService))
		r.Group(func(r handler.Router) {
			r.Use(UserHasHiveCheck(hiveService))
			r.Command("/add", AddCommand(b, hiveService))
			r.Command("/remove", RemoveCommand(b, hiveService))
			r.Command("/giftall", GiftAllCommand(b, hiveService))
			r.Command("/setlevel", SetLevelCommand(b, hiveService))
			r.Command("/setbeequip", SetBeequipCommand(b, hiveService))
			r.Command("/setmutation", SetMutationCommand(b, hiveService))
			r.Command("/info", InfoCommand(b, hiveService))
			r.Command("/view", ViewCommand(b, hiveService))
			r.Command("/save", SaveCommand(b, hiveService))
		})
		r.Route("/saves", func(r handler.Router) {
			r.Command("/load", SavesLoadCommand(b, hiveService))
			r.Command("/list", SavesListCommand(b, hiveService))
			r.Command("/delete", SavesDeleteCommand(b, hiveService))
			r.ButtonComponent("/saveid/{id}", SaveIdButton())
		})
		r.Route("/buttons", func(r handler.Router) {
			r.Use(UserOwnsHiveCheck)
			r.ButtonComponent("/addbee/{uid}", AddBeeButton())
			r.ButtonComponent("/giftall/{uid}/{showHiveNumbers}", GiftAllButton(b, hiveService))
			r.ButtonComponent("/setlevel/{uid}", SetLevelButton())
			r.ButtonComponent("/hiveinfo/{uid}", HiveInfoButton(hiveService))
			r.ButtonComponent("/mutationinfo/{uid}", MutationInfoButton(hiveService))
			r.ButtonComponent("/rerender/{uid}", HiveRerenderButton(b, hiveService))
		})
		r.Route("/modals", func(r handler.Router) {
			r.Use(UserOwnsHiveCheck)
			r.Modal("/addbee/{uid}", AddBeeModal(hiveService))
			r.Modal("/setlevel/{uid}", SetLevelModal(hiveService))
		})
		r.Group(func(r handler.Router) {
			r.Autocomplete("/add", AddAutocomplete)
			r.Autocomplete("/setbeequip", SetBeequipAutocomplete)
			r.Autocomplete("/setmutation", SetMutationAutocomplete)
			r.Autocomplete("/view", HiveViewAutocomplete)
		})
	})
}
