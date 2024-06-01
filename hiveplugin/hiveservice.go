package hiveplugin

import (
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/snowflake/v2"
)

type State struct {
	users map[snowflake.ID]*hive.Hive
}

func NewHiveService() *State {
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

func (s *State) HiveCount() int {
	return len(s.users)
}

func (s *State) Hives() map[snowflake.ID]*hive.Hive {
	return s.users
}
