package hive

import (
	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

type Bee struct {
	level    int
	id       string
	name     string
	gifted   bool
	beequip  string
	mutation string
}

func NewBee(level int, id string, gifted bool) *Bee {
	_, meta := loaders.GetBee(id)
	return &Bee{level, id, meta.Name, gifted, "", "None"}
}

func (b *Bee) Name() string {
	return b.name
}

func (b *Bee) SetGifted(state bool) {
	b.gifted = state
}

func (b *Bee) SetBeequip(name string) {
	b.beequip = name
}

func (b *Bee) SetMutation(mutation string) {
	b.mutation = mutation
}

func (b *Bee) ToBson() bson.D {
	return bson.D{
		{"id", b.id},
		{"level", b.level},
		{"gifted", b.gifted},
		{"beequip", b.beequip},
		{"mutation", b.mutation},
	}
}

func (b *Bee) Draw(dc *gg.Context, x int, y int) func() {
	if b.gifted {
		dc.SetLineWidth(4)
		dc.SetHexColor("#ffff00")
		dc.DrawRegularPolygon(6, float64(x), float64(y), 42, 0)
		dc.Stroke()
	}
	face, beeMeta := loaders.GetBee(b.id)
	switch beeMeta.Kind {
	case loaders.Common:
		dc.SetHexColor("#A76F33")
	case loaders.Rare:
		dc.SetHexColor("#9B9B9B")
	case loaders.Epic:
		dc.SetHexColor("#A48B37")
	case loaders.Legendary:
		dc.SetHexColor("#87CFCE")
	case loaders.Mythic:
		dc.SetHexColor("#826FAC")
	case loaders.Event:
		dc.SetHexColor("#74B052")
	}
	dc.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
	dc.Fill()
	dc.DrawImageAnchored(face, x, y, 0.5, 0.5)
	return func() {
		dc.SetHexColor(loaders.GetMutation(b.mutation))
		dc.DrawStringAnchored(strconv.Itoa(b.level), float64(x-60), float64(y-10), 0, 0.5)
		if b.beequip != "" {
			dc.DrawImageAnchored(loaders.GetBeequipImage(b.beequip), x+15, y+15, 0.5, 0)
		}
	}
}
