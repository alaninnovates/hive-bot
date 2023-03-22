package gameplugin

import (
	"github.com/disgoorg/snowflake/v2"
	"math"
	"time"
)

type IdentifyGameUser struct {
	Correct        int
	Incorrect      int
	QuestionAmount int
}

func (u *IdentifyGameUser) QuestionsAnswered() int {
	return u.Correct + u.Incorrect
}

func (u *IdentifyGameUser) IncrementCorrect() {
	u.Correct++
}

func (u *IdentifyGameUser) IncrementIncorrect() {
	u.Incorrect++
}

type TriviaQuestionAnswer struct {
	Question     string
	Choices      []string
	CorrectIndex int
	ChosenIndex  int
	StartTime    int64
	EndTime      int64
}

func (u *TriviaGameUser) AddQuestion(question TriviaQuestionAnswer) {
	u.Questions = append(u.Questions, question)
}

func (u *TriviaGameUser) SetChosenIndex(index int) {
	u.Questions[len(u.Questions)-1].ChosenIndex = index
	u.Questions[len(u.Questions)-1].EndTime = time.Now().UnixMilli()
}

func (u *TriviaGameUser) GetCorrect() int {
	correct := 0
	for _, question := range u.Questions {
		if question.CorrectIndex == question.ChosenIndex {
			correct++
		}
	}
	return correct
}

func (u *TriviaGameUser) GetScore() int {
	//	compute score based off of correct answers and time
	//	every question is 5 points, and every extra 20 seconds taken is a 1 point reduction
	score := 0
	for _, question := range u.Questions {
		if question.CorrectIndex == question.ChosenIndex {
			timeTaken := float64(question.EndTime-question.StartTime) / 1000
			reduction := math.Floor(timeTaken / 20)
			if reduction > 1 {
				reduction = 4
			}
			score += 5 - int(reduction)
		}
	}
	return score
}

type TriviaGameUser struct {
	Difficulty TriviaDifficulty
	Questions  []TriviaQuestionAnswer
}

type BubbleGameUser struct {
	StartTime int64
	Bubbles   [][]bool
}

func (u *BubbleGameUser) PopAmount() int {
	amount := 0
	for _, row := range u.Bubbles {
		for _, bubble := range row {
			if bubble {
				amount++
			}
		}
	}
	return amount
}

type WhackGameUser struct {
}

type GameType int

const (
	GameTypeIdentifyTheBee GameType = iota
	GameTypeTrivia
	// GameTypePopBubble Idea: 5x5 grid of button emojis with bubble. click as fast as possible.
	GameTypePopBubble
	GameTypeWhackABee
)

type GameUser struct {
	GameType         GameType
	IdentifyGameUser *IdentifyGameUser
	TriviaGameUser   *TriviaGameUser
	BubbleGameUser   *BubbleGameUser
}

type State struct {
	users map[snowflake.ID]*GameUser
}

func NewGameService() *State {
	return &State{users: make(map[snowflake.ID]*GameUser)}
}

func (s *State) IsPlayingGame(userID snowflake.ID) bool {
	return s.users[userID] != nil
}

func (s *State) StartIdentifyGame(userID snowflake.ID, questionAmount int) *IdentifyGameUser {
	s.users[userID] = &GameUser{
		GameType: GameTypeIdentifyTheBee,
		IdentifyGameUser: &IdentifyGameUser{
			QuestionAmount: questionAmount,
		},
	}
	return s.users[userID].IdentifyGameUser
}

func (s *State) StartTriviaGame(userID snowflake.ID, difficulty TriviaDifficulty) *TriviaGameUser {
	s.users[userID] = &GameUser{
		GameType: GameTypeTrivia,
		TriviaGameUser: &TriviaGameUser{
			Difficulty: difficulty,
		},
	}
	return s.users[userID].TriviaGameUser
}

func (s *State) StartBubbleGame(userID snowflake.ID) *BubbleGameUser {
	bubbles := make([][]bool, 5)
	for i := range bubbles {
		bubbles[i] = make([]bool, 5)
	}
	s.users[userID] = &GameUser{
		GameType: GameTypePopBubble,
		BubbleGameUser: &BubbleGameUser{
			Bubbles: bubbles,
		},
	}
	return s.users[userID].BubbleGameUser
}

func (s *State) GetGameUser(userID snowflake.ID, gameType GameType) *GameUser {
	if s.users[userID] == nil || s.users[userID].GameType != gameType {
		return nil
	}
	return s.users[userID]
}

func (s *State) EndGame(userID snowflake.ID) {
	delete(s.users, userID)
}
