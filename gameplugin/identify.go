package gameplugin

import (
	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/common/loaders"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/handler"
	"github.com/fogleman/gg"
	"image"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

type IdentifyAnswerChoice struct {
	Name      string
	IsCorrect bool
}

type IdentifyDifficulty int

const (
	IdentifyDifficultyNormal IdentifyDifficulty = iota
	IdentifyDifficultyPixelShuffle
	IdentifyDifficultyDiscoloration
	IdentifyDifficultyCombo
)

func GetAnswerSet() ([]IdentifyAnswerChoice, int) {
	var answerSet []IdentifyAnswerChoice
	answersChose := make(map[string]int)
	bees := loaders.GetBeeNames()
	correctI := rand.Intn(4)
	for i := 0; i < 4; {
		randomIndex := rand.Intn(len(bees))
		if _, found := answersChose[bees[randomIndex]]; found {
			continue
		}
		answerSet = append(answerSet, IdentifyAnswerChoice{
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

func GetAnswerSetButtons(userId string, IdentifyDifficulty IdentifyDifficulty) ([]discord.InteractiveComponent, io.Reader) {
	answerSet, correctI := GetAnswerSet()
	//fmt.Println("correctI: "+strconv.Itoa(correctI), answerSet[correctI].Name)
	img, err := gg.LoadImage("assets/bees/" + answerSet[correctI].Name + ".png")
	switch IdentifyDifficulty {
	case IdentifyDifficultyPixelShuffle:
		img = ShuffleImage(img)
	case IdentifyDifficultyDiscoloration:
		img = DiscolorImage(img)
	case IdentifyDifficultyCombo:
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
			id = "handler:correct:" + userId + ":" + strconv.Itoa(int(IdentifyDifficulty))
		} else {
			id = "handler:incorrect:" + userId + ":" + strconv.Itoa(int(IdentifyDifficulty)) + ":" + strconv.Itoa(i)
		}
		buttons = append(buttons, discord.ButtonComponent{
			Label:    answer.Name,
			Style:    discord.ButtonStylePrimary,
			CustomID: id,
		})
	}
	return buttons, r
}

func IdentifyTheBeeCommand(gameService *State) func(event *events.ApplicationCommandInteractionCreate) error {
	return func(event *events.ApplicationCommandInteractionCreate) error {
		if gameService.IsPlayingGame(event.User().ID) {
			return event.CreateMessage(discord.MessageCreate{
				Content: "You are already playing a game!",
			})
		}
		data := event.SlashCommandInteractionData()
		questions, ok := data.OptInt("questions")
		if !ok {
			questions = 5
		}
		difficulty, ok := data.OptInt("difficulty")
		if !ok {
			difficulty = 0
		}
		gameService.StartIdentifyGame(event.User().ID, questions)
		buttons, r := GetAnswerSetButtons(event.User().ID.String(), IdentifyDifficulty(difficulty))
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
	}
}

func TriviaButtonHandler(event *events.ComponentInteractionCreate) error {
	difficulty, _ := strconv.Atoi(strings.Split(event.ButtonInteractionData().CustomID(), ":")[3])
	buttons, r := GetAnswerSetButtons(event.User().ID.String(), IdentifyDifficulty(difficulty))
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
	user := gameService.GetGameUser(event.User().ID, GameTypeIdentifyTheBee).IdentifyGameUser
	gameService.EndGame(event.User().ID)
	return event.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Trivia Summary",
				Description: "You got " + strconv.Itoa(user.Correct) + " out of " + strconv.Itoa(user.QuestionsAnswered()) + " questions correct!",
			},
		},
	})
}

func IdentifyCorrectButton(b *common.Bot, gameService *State) handler.Component {
	return handler.Component{
		Name:  "correct",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			tu := gameService.GetGameUser(event.User().ID, GameTypeIdentifyTheBee)
			if tu == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not in a trivia game",
				})
			}
			iu := tu.IdentifyGameUser
			iu.IncrementCorrect()
			if iu.QuestionsAnswered() >= iu.QuestionAmount {
				return TriviaSummary(event, gameService)
			}
			return TriviaButtonHandler(event)
		},
	}
}

func IdentifyIncorrectButton(b *common.Bot, gameService *State) handler.Component {
	return handler.Component{
		Name:  "incorrect",
		Check: userIDCheck(),
		Handler: func(event *events.ComponentInteractionCreate) error {
			tu := gameService.GetGameUser(event.User().ID, GameTypeIdentifyTheBee)
			if tu == nil {
				return event.CreateMessage(discord.MessageCreate{
					Content: "You are not in a trivia game",
				})
			}
			iu := tu.IdentifyGameUser
			iu.IncrementIncorrect()
			if iu.QuestionsAnswered() >= iu.QuestionAmount {
				return TriviaSummary(event, gameService)
			}
			return TriviaButtonHandler(event)
		},
	}
}
