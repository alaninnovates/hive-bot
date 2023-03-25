package statsplugin

import (
	"github.com/disgoorg/snowflake/v2"
	"golang.org/x/exp/slices"
)

type State struct {
	activeUsers     []string
	commandsRun     int
	popularCommands map[string]int
}

func NewStatsService() *State {
	return &State{activeUsers: []string{}, commandsRun: 0, popularCommands: make(map[string]int)}
}

func (s *State) ActiveUsers() []string {
	return s.activeUsers
}

func (s *State) CommandsRun() int {
	return s.commandsRun
}

type PopularCommand struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (s *State) PopularCommands() []PopularCommand {
	var popCmdArr []PopularCommand
	for cmdName, count := range s.popularCommands {
		popCmdArr = append(popCmdArr, PopularCommand{
			Name:  cmdName,
			Count: count,
		})
	}
	return popCmdArr
}

func (s *State) CommandRun(userId snowflake.ID, commandName string) {
	s.commandsRun++
	if !slices.Contains(s.activeUsers, userId.String()) {
		s.activeUsers = append(s.activeUsers, userId.String())
	}
	s.popularCommands[commandName]++
}

func (s *State) ResetStats() {
	s.commandsRun = 0
	s.activeUsers = []string{}
	s.popularCommands = make(map[string]int)
}
