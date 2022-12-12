package hive

import (
	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"image/color"
	"strconv"
)

type Bee struct {
	level   int
	name    string
	gifted  bool
	beequip string
}

func NewBee(level int, name string, gifted bool) *Bee {
	return &Bee{level, name, gifted, ""}
}

func (b *Bee) SetGifted(state bool) {
	b.gifted = state
}

func (b *Bee) SetBeequip(name string) {
	b.beequip = name
}

func (b *Bee) ToBson() bson.D {
	return bson.D{{"name", b.name}, {"level", b.level}, {"gifted", b.gifted}}
}

func (b *Bee) Draw(dc *gg.Context, x int, y int) func() {
	if b.gifted {
		dc.DrawRegularPolygon(6, float64(x), float64(y), 50, 0)
		dc.SetHexColor("#ffff00")
		dc.Fill()
	}
	dd := gg.NewContext(410, 900)
	dd.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
	dd.Fill()
	err := dc.SetMask(dd.AsMask())
	if err != nil {
		panic(err)
	}
	dc.DrawImageAnchored(loaders.GetBeeImage(b.name), x, y, 0.5, 0.5)
	dd = gg.NewContext(410, 900)
	err = dc.SetMask(dd.AsMask())
	if err != nil {
		panic(err)
	}
	dc.InvertMask()
	return func() {
		if b.gifted {
			dc.SetHexColor("#ffd085")
		} else {
			dc.SetColor(color.White)
		}
		dc.DrawStringAnchored(strconv.Itoa(b.level), float64(x-60), float64(y-10), 0, 0.5)
		if b.beequip != "" {
			dc.DrawImageAnchored(loaders.GetBeequipImage(b.beequip), x+15, y+15, 0.5, 0)
		}
	}
}
