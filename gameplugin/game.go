package gameplugin

import (
	"strings"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
)

func GameCommand(b *common.Bot, gameService *State) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "game",
			Description: "Play games",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "identify-the-bee",
					Description: "Play a game of identify-the-bee",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "questions",
							Description: "Number of questions to answer",
							Required:    false,
							MinValue:    json.Ptr(5),
							MaxValue:    json.Ptr(100),
						},
						discord.ApplicationCommandOptionInt{
							Name:        "difficulty",
							Description: "Difficulty of the game",
							Required:    false,
							Choices: []discord.ApplicationCommandOptionChoiceInt{
								{
									Name:  "Normal",
									Value: 0,
								},
								{
									Name:  "Difficult",
									Value: 1,
								},
								{
									Name:  "Insane",
									Value: 2,
								},
								{
									Name:  "Expert",
									Value: 3,
								},
							},
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "trivia",
					Description: "Play a game of Bee Swarm Simulator trivia",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "difficulty",
							Description: "Difficulty of the trivia questions",
							Choices: []discord.ApplicationCommandOptionChoiceInt{
								{
									Name:  "Beginner",
									Value: 0,
								},
								{
									Name:  "Midgame",
									Value: 1,
								},
								{
									Name:  "Endgame",
									Value: 2,
								},
							},
							Required: true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pop-bubbles",
					Description: "Play a game of pop the bubbles.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "whack-bee",
					Description: "Play a game of Whack-a-bee.",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "leaderboard",
					Description: "View leaderboards",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionString{
							Name:        "type",
							Description: "Type of leaderboard you want to view",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								//{Name: "Identify the Bee", Value: "identify-the-bee"},
								{Name: "Trivia", Value: "trivia"},
								{Name: "Pop the Bubbles", Value: "pop-bubbles"},
								//{Name: "Whack-a-bee", Value: "whack-bee"},
							},
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "stop",
					Description: "Stop any active game you are playing",
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"identify-the-bee": IdentifyTheBeeCommand(gameService),
			"trivia":           TriviaCommand(b, gameService),
			"pop-bubbles":      BubbleCommand(b, gameService),
			"whack-bee":        WhackCommand(b, gameService),
			"leaderboard":      LeaderboardCommand(b, gameService),
			"stop": func(event *events.ApplicationCommandInteractionCreate) error {
				gameService.EndGame(event.User().ID)
				return event.CreateMessage(discord.MessageCreate{
					Content: "Stopped your game",
				})
			},
		},
	}
}

func EndGameButton(gameService *State) handler.Component {
	return handler.Component{
		Name:  "endgame",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			gameService.EndGame(event.User().ID)
			return event.CreateMessage(discord.MessageCreate{
				Content: "Ended game. Good job!",
			})
		},
	}
}

func userIDCheck() handler.Check[*events.ComponentInteractionCreate] {
	return func(event *events.ComponentInteractionCreate) bool {
		return event.User().ID.String() == strings.Split(event.ButtonInteractionData().CustomID(), ":")[2]
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	gameService := NewGameService()
	h.AddCommands(GameCommand(b, gameService))
	h.AddComponents(IdentifyCorrectButton(b, gameService), IdentifyIncorrectButton(b, gameService))
	h.AddComponents(TriviaButton(gameService), TriviaReviewButton(gameService))
	h.AddComponents(BubbleButton(gameService))
	h.AddComponents(EndGameButton(gameService))
}
