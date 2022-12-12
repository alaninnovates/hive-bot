package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"image"
	"image/png"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

type AnswerChoice struct {
	Name      string
	IsCorrect bool
}

func GetAnswerSet() ([]AnswerChoice, int) {
	var answerSet []AnswerChoice
	bees := loaders.GetBees()
	correctI := rand.Intn(4)
	for i := 0; i < 4; i++ {
		randomIndex := rand.Intn(len(bees))
		answerSet = append(answerSet, AnswerChoice{
			Name:      bees[randomIndex],
			IsCorrect: i == correctI,
		})
	}
	return answerSet, correctI
}

func GetAnswerSetButtons(userId string) ([]discord.InteractiveComponent, io.Reader) {
	answerSet, correctI := GetAnswerSet()
	r, w := io.Pipe()
	go func(i image.Image) {
		defer w.Close()
		if err := png.Encode(w, i); err != nil {
			panic(err)
		}
	}(loaders.GetBeeImage(answerSet[correctI].Name))
	buttons := make([]discord.InteractiveComponent, 0)
	for i, answer := range answerSet {
		id := ""
		if answer.IsCorrect {
			id = "handler:correct:" + userId
		} else {
			id = "handler:incorrect:" + userId + ":" + string(i)
		}
		buttons = append(buttons, discord.ButtonComponent{
			Label:    answer.Name,
			Style:    discord.ButtonStylePrimary,
			CustomID: id,
		})
	}
	return buttons, r
}

func GameCommand(b *common.Bot, gameService *State) handler.Command {
	return handler.Command{
		Create: discord.SlashCommandCreate{
			Name:        "game",
			Description: "Play games",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "identify",
					Description: "Play a game of identify-the-bee",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "questions",
							Description: "Number of questions to answer",
							Required:    false,
						},
					},
				},
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"identify": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				questions, ok := data.OptInt("questions")
				if !ok {
					questions = 5
				}
				gameService.StartTriviaGame(event.User().ID, questions)
				buttons, r := GetAnswerSetButtons(event.User().ID.String())
				return event.CreateMessage(discord.MessageCreate{
					Embeds: []discord.Embed{
						{
							Title: "What bee is this?",
							Image: &discord.EmbedResource{
								URL: "attachment://bee.png",
							},
						},
					},
					Components: []discord.ContainerComponent{
						discord.ActionRowComponent{}.AddComponents(buttons...),
					},
					Files: []*discord.File{
						{
							Name:   "bee.png",
							Reader: r,
						},
					},
				})
			},
		},
	}
}

func TriviaButtonHandler(event *events.ComponentInteractionCreate) error {
	buttons, r := GetAnswerSetButtons(event.User().ID.String())
	return event.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			{
				Title: "What bee is this?",
				Image: &discord.EmbedResource{
					URL: "attachment://bee.png",
				},
			},
		},
		Components: &[]discord.ContainerComponent{
			discord.ActionRowComponent{}.AddComponents(buttons...),
		},
		Files: []*discord.File{
			{
				Name:   "bee.png",
				Reader: r,
			},
		},
	})
}

func TriviaSummary(event *events.ComponentInteractionCreate, gameService *State) error {
	user := gameService.GetTriviaUser(event.User().ID)
	gameService.EndTriviaGame(event.User().ID)
	return event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Trivia Summary",
				Description: "You got " + strconv.Itoa(user.Correct) + " correct out of " + strconv.Itoa(user.QuestionsAnswered()),
			},
		},
	})
}

func CorrectButton(b *common.Bot, gameService *State) handler.Component {
	return handler.Component{
		Name:  "correct",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			tu := gameService.GetTriviaUser(event.User().ID)
			if tu == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not in a trivia game",
				})
			}
			tu.IncrementCorrect()
			if tu.QuestionsAnswered() >= tu.QuestionAmount {
				return TriviaSummary(event, gameService)
			}
			return TriviaButtonHandler(event)
		},
	}
}

func IncorrectButton(b *common.Bot, gameService *State) handler.Component {
	return handler.Component{
		Name:  "incorrect",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			tu := gameService.GetTriviaUser(event.User().ID)
			if tu == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not in a trivia game",
				})
			}
			tu.IncrementIncorrect()
			if tu.QuestionsAnswered() >= tu.QuestionAmount {
				return TriviaSummary(event, gameService)
			}
			return TriviaButtonHandler(event)
		},
	}
}

func userIDCheck() handler.Check[*events.ComponentInteractionCreate] {
	return func(event *events.ComponentInteractionCreate) bool {
		return event.User().ID.String() == strings.Split(event.ButtonInteractionData().CustomID(), ":")[2]
	}
}

func Initialize(h *handler.Handler, b *common.Bot) {
	gameSerivce := NewGameService()
	h.AddCommands(GameCommand(b, gameSerivce))
	h.AddComponents(CorrectButton(b, gameSerivce), IncorrectButton(b, gameSerivce))
}
