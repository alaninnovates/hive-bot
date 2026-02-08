package gameplugin

import (
	"strings"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var GameCommandCreate = discord.SlashCommandCreate{
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
}

func EndGameButton(gameService *State) handler.ComponentHandler {
	return func(event *handler.ComponentEvent) error {
		gameService.EndGame(event.User().ID)
		return event.CreateMessage(discord.MessageCreate{
			Content: "Ended game. Good job!",
		})
	}
}

func StopCommand(gameService *State) handler.CommandHandler {
	return func(event *handler.CommandEvent) error {
		gameService.EndGame(event.User().ID)
		return event.CreateMessage(discord.MessageCreate{
			Content: "Stopped your game",
		})
	}
}

func userIDCheck(next handler.Handler) handler.Handler {
	return func(event *handler.InteractionEvent) error {
		if event.Type() != discord.InteractionTypeComponent {
			return next(event)
		}

		btnEvent := event.Interaction.(discord.ComponentInteraction)
		uid := strings.Split(btnEvent.ButtonInteractionData().CustomID(), "/")

		if event.User().ID.String() != uid[4] {
			return event.CreateMessage(discord.NewMessageCreateBuilder().
				SetContent("This isn't your game!").
				SetEphemeral(true).
				Build())
		}
		return next(event)
	}
}

func Initialize(r *handler.Mux, b *common.Bot) {
	gameService := NewGameService()
	r.Route("/game", func(r handler.Router) {
		r.Use(userIDCheck)
		r.Command("/leaderboard", LeaderboardCommand(b, gameService))
		r.Route("/identify-the-bee", func(r handler.Router) {
			r.Command("/", IdentifyTheBeeCommand(gameService))
			r.Component("/correct/{uid}/{difficulty}", IdentifyCorrectButton(b, gameService))
			r.Component("/incorrect/{uid}/{difficulty}/{correctIdx}", IdentifyIncorrectButton(b, gameService))
		})
		r.Route("/trivia", func(r handler.Router) {
			r.Command("/", TriviaCommand(b, gameService))
			r.Component("/trivia/{uid}/{correct}/{i}", TriviaButton(gameService))
			r.Component("/trivia-review/{uid}", TriviaReviewButton(gameService))
			r.Component("/end/{uid}", EndGameButton(gameService))
		})
		r.Route("/pop-bubbles", func(r handler.Router) {
			r.Command("/", BubbleCommand(b, gameService))
			r.Component("/bubble/{uid}/{i}/{j}", BubbleButton(gameService))
		})
		r.Command("/stop", StopCommand(gameService))
	})
}
