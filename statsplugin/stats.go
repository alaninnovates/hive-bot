package statsplugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/handler"
)

var (
	baseUrl = "https://api.statcord.com/v3"
)

type StatsRequest struct {
	Id        string           `json:"id"`
	Key       string           `json:"key"`
	Servers   string           `json:"servers"`
	Users     string           `json:"users"`
	Active    []string         `json:"active"`
	Commands  string           `json:"commands"`
	Popular   []PopularCommand `json:"popular"`
	MemActive string           `json:"memactive"`
	MemLoad   string           `json:"memload"`
	CpuLoad   string           `json:"cpuload"`
	Bandwidth string           `json:"bandwidth"`
}

func PostStats(b *common.Bot, state *State) {
	members := 0
	b.Client.Caches().GuildsForEach(func(e discord.Guild) {
		members += e.MemberCount
	})
	popCmds := state.PopularCommands()
	if popCmds == nil {
		popCmds = make([]PopularCommand, 0)
	}
	data, err := json.Marshal(StatsRequest{
		Id:        b.Client.ID().String(),
		Key:       os.Getenv("STATCORD_KEY"),
		Servers:   strconv.Itoa(b.Client.Caches().GuildsLen()),
		Users:     strconv.Itoa(members),
		Active:    state.ActiveUsers(),
		Commands:  strconv.Itoa(state.CommandsRun()),
		Popular:   popCmds,
		MemActive: "0",
		MemLoad:   "0",
		CpuLoad:   "0",
		Bandwidth: "0",
	})
	//fmt.Printf("%s\n", data)
	if err != nil {
		b.Logger.Error("Failed to send stats: Could not marshal data")
		return
	}
	r := bytes.NewReader(data)
	res, err := http.Post(baseUrl+"/stats", "application/json", r)
	if err != nil {
		b.Logger.Error("Failed to send stats: HTTP post errored: %v", err)
		return
	}
	if res.StatusCode != 200 {
		b.Logger.Error("Failed to send stats: HTTP post errored with code %d", res.StatusCode)
		bts, _ := io.ReadAll(res.Body)
		b.Logger.Error("%s", bts)
		return
	}
	b.Logger.Info(fmt.Sprintf("Posted stats. Servers: %d, Users: %d", b.Client.Caches().GuildsLen(), members))
	state.ResetStats()
}

func Initialize(r *handler.Mux, b *common.Bot, devMode bool) {
	if devMode {
		b.Logger.Info("Not posting bot stats: Developer mode is enabled")
		return
	}
	//statsService := NewStatsService()
	b.Client.AddEventListeners(&events.ListenerAdapter{
		OnReady: func(event *events.Ready) {
			//ticker := time.NewTicker(2 * time.Minute)
			//go func() {
			//	for {
			//		select {
			//		case <-ticker.C:
			//			PostStats(b, statsService)
			//		}
			//	}
			//}()
			//b.Logger.Info("Started auto poster.")
		},
		OnGuildsReady: func(event *events.GuildsReady) {
			//PostStats(b, statsService)
		},
		OnApplicationCommandInteraction: func(event *events.ApplicationCommandInteractionCreate) {
			if event.Data.Type() != discord.ApplicationCommandTypeSlash {
				return
			}
			data := event.SlashCommandInteractionData()
			guildName, ok := event.Guild()
			guild := ""
			if ok {
				guild = guildName.Name
			} else {
				guild = "Dms"
			}
			b.Logger.Info(fmt.Sprintf("%s used %s in %s", event.User().Tag(), data.CommandPath(), guild))
			//statsService.CommandRun(event.User().ID, data.CommandPath())
		},
	})
}
