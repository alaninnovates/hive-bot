package main

import (
	"alaninnovates.com/hive-bot/hive"
	"github.com/disgoorg/snowflake/v2"
)

type State struct {
	users map[snowflake.ID]*hive.Hive
}

func NewState() *State {
	return &State{users: make(map[snowflake.ID]*hive.Hive)}
}

func (s *State) CreateHive(userID snowflake.ID) *hive.Hive {
	h := hive.NewHive()
	s.users[userID] = h
	return h
}

func (s *State) GetHive(userID snowflake.ID) *hive.Hive {
	return s.users[userID]
}
