package gameplugin

import (
	"github.com/disgoorg/snowflake/v2"
)

type TriviaGameUser struct {
	Correct        int
	Incorrect      int
	QuestionAmount int
}

func (u *TriviaGameUser) QuestionsAnswered() int {
	return u.Correct + u.Incorrect
}

func (u *TriviaGameUser) IncrementCorrect() {
	u.Correct++
}

func (u *TriviaGameUser) IncrementIncorrect() {
	u.Incorrect++
}

type State struct {
	users map[snowflake.ID]*TriviaGameUser
}

func NewGameService() *State {
	return &State{users: make(map[snowflake.ID]*TriviaGameUser)}
}

func (s *State) StartTriviaGame(userID snowflake.ID, questionAmount int) *TriviaGameUser {
	s.users[userID] = &TriviaGameUser{
		QuestionAmount: questionAmount,
	}
	return s.users[userID]
}

func (s *State) GetTriviaUser(userID snowflake.ID) *TriviaGameUser {
	if s.users[userID] == nil {
		return nil
	}
	return s.users[userID]
}

func (s *State) EndTriviaGame(userID snowflake.ID) {
	delete(s.users, userID)
}
