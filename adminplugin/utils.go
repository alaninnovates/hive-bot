package adminplugin

import (
	"fmt"
	"sync"

	"alaninnovates.com/hive-bot/common"
	"alaninnovates.com/hive-bot/database"
	"alaninnovates.com/hive-bot/hiveplugin"
	"alaninnovates.com/hive-bot/hiveplugin/hive"
	"github.com/disgoorg/snowflake/v2"
)

var fileMutex = sync.Mutex{}

func LoadHives(b *common.Bot, hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) {
	hives, err := jsonCacheService.LoadHives("data/hives.json")
	if err != nil {
		b.Logger.Error("Failed to load hives: %v", err)
		return
	}
	b.Logger.Info(fmt.Sprintf("Loading %d hives", len(hives)))
	for _, cachedUser := range hives {
		h := hiveService.CreateHive(snowflake.MustParse(cachedUser.Id))
		for idx, cachedBees := range cachedUser.Hive {
			for _, cachedBee := range cachedBees {
				h.AddBee(hive.NewBee(cachedBee.Level, cachedBee.Id, cachedBee.Gifted), idx)
				pos := len(h.GetBeesAt(idx)) - 1
				h.GetBeesAt(idx)[pos].SetBeequip(cachedBee.Beequip)
				h.GetBeesAt(idx)[pos].SetMutation(cachedBee.Mutation)
			}
		}
	}
}

func BackupHives(b *common.Bot, hiveService *hiveplugin.State, jsonCacheService *database.JsonCache) {
	// acquire lock
	fileMutex.Lock()
	cachedUsers := make([]database.CachedUser, 0)
	for id, h := range hiveService.Hives() {
		cachedHive := make(database.CachedHive)
		for idx, bees := range h.GetBees() {
			for _, bee := range bees {
				cachedHive[idx] = append(cachedHive[idx], database.CachedBee{
					Id:       bee.Id(),
					Level:    bee.Level(),
					Gifted:   bee.Gifted(),
					Beequip:  bee.Beequip(),
					Mutation: bee.Mutation(),
				})
			}
		}
		cachedUsers = append(cachedUsers, database.CachedUser{
			Id:           id.String(),
			Hive:         cachedHive,
			LastModified: h.LastModified(),
		})
	}
	err := jsonCacheService.SaveHives("data/hives.json", cachedUsers)
	if err != nil {
		b.Logger.Error("Failed to back up hives: %v", err)
		return
	}
	b.Logger.Info(fmt.Sprintf("Backed up %d hives", len(cachedUsers)))
	// release lock
	fileMutex.Unlock()
}

func PruneHives(b *common.Bot, hiveService *hiveplugin.State) {
	for id, h := range hiveService.Hives() {
		if h.LastModified() < common.CurrentTimeMillis()-1000*60*60 {
			hiveService.DeleteHive(id)
		}
	}
}
