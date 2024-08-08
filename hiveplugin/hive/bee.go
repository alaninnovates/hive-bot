package hive

import (
	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"math"
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
	return &Bee{level, id, meta.Name, gifted, "None", "None"}
}

func (b *Bee) Level() int {
	return b.level
}

func (b *Bee) Gifted() bool {
	return b.gifted
}

func (b *Bee) Id() string {
	return b.id
}

func (b *Bee) Name() string {
	return b.name
}

func (b *Bee) Beequip() string {
	return b.beequip
}

func (b *Bee) Mutation() string {
	return b.mutation
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

func (b *Bee) SetLevel(level int) {
	b.level = level
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

func DrawBees(b []*Bee, dc *gg.Context, x int, y int) map[string]func() {
	if len(b) == 1 {
		return drawOneBee(b[0], dc, x, y)
	}

	/*
		apothem: a=rcos(180/n)
	*/
	a := 40 * math.Cos(180/float64(6)*math.Pi/180)

	n := 6
	r := 40
	centerAng := 2 * math.Pi / float64(n)
	startAng := 0

	cornersX := make([]float64, n)
	cornersY := make([]float64, n)
	for i := 0; i < n; i++ {
		ang := float64(startAng) + (float64(i) * centerAng)
		cornersX[i] = float64(x) + float64(r)*math.Cos(ang)
		cornersY[i] = float64(y) - float64(r)*math.Sin(ang)
	}

	dc.NewSubPath()
	dc.MoveTo(float64(x-r), float64(y)-a)
	dc.LineTo(cornersX[1], cornersY[1])
	dc.LineTo(cornersX[4], cornersY[4])
	dc.LineTo(float64(x-r), float64(y)+a)
	dc.ClosePath()
	dc.Clip()

	funcs := drawOneBee(b[0], dc, x, y)

	dc.ResetClip()

	dc.NewSubPath()
	dc.MoveTo(float64(x+r), float64(y)-a)
	dc.LineTo(cornersX[1], cornersY[1])
	dc.LineTo(cornersX[4], cornersY[4])
	dc.LineTo(float64(x+r), float64(y)+a)
	dc.ClosePath()
	dc.Clip()

	funcs = drawOneBee(b[1], dc, x, y)

	dc.ResetClip()

	return funcs
}

func drawOneBee(b *Bee, dc *gg.Context, x int, y int) map[string]func() {
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
	funcs := make(map[string]func())
	funcs["gifted"] = func() {
		if b.gifted {
			dc.SetLineWidth(4)
			dc.SetHexColor("#ffff00")
			dc.DrawRegularPolygon(6, float64(x), float64(y), 42, 0)
			dc.Stroke()
		}
	}
	funcs["beequip"] = func() {
		if b.beequip != "None" {
			dc.DrawImageAnchored(loaders.GetBeequipImage(b.beequip), x+15, y+15, 0.5, 0)
		}
	}
	funcs["level"] = func() {
		if b.level != 0 {
			dc.SetHexColor("#000000")
			//todo: this border drawing function creates cpu spikes, find a better way to do this
			n := 3
			for dy := -n; dy <= n; dy++ {
				for dx := -n; dx <= n; dx++ {
					if dx*dx+dy*dy >= n*n {
						continue
					}
					//println("drawing", dx, dy)
					xx := float64(x - 60 + dx)
					yy := float64(y - 10 + dy)
					dc.DrawStringAnchored(strconv.Itoa(b.level), xx, yy, 0, 0.5)
				}
			}
			dc.SetHexColor(loaders.GetMutation(b.mutation))
			dc.DrawStringAnchored(strconv.Itoa(b.level), float64(x-60), float64(y-10), 0, 0.5)
		}
	}
	return funcs
}
