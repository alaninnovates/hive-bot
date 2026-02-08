package miscplugin

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/disgoorg/disgo/handler"
)

func Initialize(r *handler.Mux, b *common.Bot) {
	r.Route("/help", func(r handler.Router) {
		r.Command("/", HelpCommand)
		r.Component("/section", HelpSelectMenu)
	})
	r.Command("/stats", StatsCommand(b))
}
