package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/disgoorg/json"
	"github.com/fogleman/gg"
	"image"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

type AnswerChoice struct {
	Name      string
	IsCorrect bool
}

type Difficulty int

const (
	DifficultyNormal Difficulty = iota
	DifficultyPixelShuffle
	DifficultyDiscoloration
	DifficultyCombo
)

func GetAnswerSet() ([]AnswerChoice, int) {
	var answerSet []AnswerChoice
	answersChose := make(map[string]int)
	bees := loaders.GetBeeNames()
	correctI := rand.Intn(4)
	for i := 0; i < 4; {
		randomIndex := rand.Intn(len(bees))
		if _, found := answersChose[bees[randomIndex]]; found {
			continue
		}
		answerSet = append(answerSet, AnswerChoice{
			Name:      bees[randomIndex],
			IsCorrect: i == correctI,
		})
		answersChose[bees[randomIndex]] = 1
		i++
	}
	return answerSet, correctI
}

func DiscolorImage(img image.Image) image.Image {
	dc := gg.NewContextForImage(img)
	dc.SetLineWidth(3)
	for i := 0; i < 80; i++ {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		dc.DrawLine(rand.Float64()*float64(img.Bounds().Max.X), rand.Float64()*float64(img.Bounds().Max.Y), rand.Float64()*float64(img.Bounds().Max.X), rand.Float64()*float64(img.Bounds().Max.Y))
		dc.Stroke()
	}
	return dc.Image()
}

func ShuffleImage(img image.Image) image.Image {
	dc := gg.NewContext(img.Bounds().Max.X, img.Bounds().Max.Y)
	dc.DrawImage(img, 0, 0)
	for i := 0; i < 200; i++ {
		x1 := rand.Intn(img.Bounds().Max.X)
		y1 := rand.Intn(img.Bounds().Max.Y)
		x2 := rand.Intn(img.Bounds().Max.X)
		y2 := rand.Intn(img.Bounds().Max.Y)
		c1 := img.At(x1, y1)
		c2 := img.At(x2, y2)
		dc.SetColor(c1)
		dc.DrawRectangle(float64(x2), float64(y2), float64(rand.Intn(10)), float64(rand.Intn(10)))
		dc.SetColor(c2)
		dc.DrawRectangle(float64(x1), float64(y1), float64(rand.Intn(10)), float64(rand.Intn(10)))
	}
	dc.Fill()
	return dc.Image()
}

func GetAnswerSetButtons(userId string, difficulty Difficulty) ([]discord.InteractiveComponent, io.Reader) {
	answerSet, correctI := GetAnswerSet()
	//fmt.Println("correctI: "+strconv.Itoa(correctI), answerSet[correctI].Name)
	img, err := gg.LoadImage("assets/bees/" + answerSet[correctI].Name + ".png")
	switch difficulty {
	case DifficultyPixelShuffle:
		img = ShuffleImage(img)
	case DifficultyDiscoloration:
		img = DiscolorImage(img)
	case DifficultyCombo:
		img = ShuffleImage(img)
		img = DiscolorImage(img)
	}
	if err != nil {
		fmt.Println(err, answerSet[correctI].Name)
		img, _ = gg.LoadImage("assets/error.png")
		return []discord.InteractiveComponent{
			discord.ButtonComponent{
				Label:    answerSet[correctI].Name,
				Style:    discord.ButtonStylePrimary,
				CustomID: "error",
			},
		}, common.ImageToPipe(img)
	}
	r := common.ImageToPipe(img)
	buttons := make([]discord.InteractiveComponent, 0)
	for i, answer := range answerSet {
		id := ""
		if answer.IsCorrect {
			id = "handler:correct:" + userId + ":" + strconv.Itoa(int(difficulty))
		} else {
			id = "handler:incorrect:" + userId + ":" + strconv.Itoa(int(difficulty)) + ":" + strconv.Itoa(i)
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
					Name:        "identify-the-bee",
					Description: "Play a game of identify-the-bee",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "questions",
							Description: "Number of questions to answer",
							Required:    false,
							MinValue:    json.Ptr(5),
							MaxValue:    json.Ptr(25),
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
			},
		},
		CommandHandlers: map[string]handler.CommandHandler{
			"identify-the-bee": func(event *events.ApplicationCommandInteractionCreate) error {
				data := event.SlashCommandInteractionData()
				questions, ok := data.OptInt("questions")
				if !ok {
					questions = 5
				}
				difficulty, ok := data.OptInt("difficulty")
				if !ok {
					difficulty = 0
				}
				gameService.StartTriviaGame(event.User().ID, questions)
				buttons, r := GetAnswerSetButtons(event.User().ID.String(), Difficulty(difficulty))
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
	difficulty, _ := strconv.Atoi(strings.Split(event.ButtonInteractionData().CustomID(), ":")[3])
	buttons, r := GetAnswerSetButtons(event.User().ID.String(), Difficulty(difficulty))
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
				Description: "You got " + strconv.Itoa(user.Correct) + " out of " + strconv.Itoa(user.QuestionsAnswered()) + " questions correct!",
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
