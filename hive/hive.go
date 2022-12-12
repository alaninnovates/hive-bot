package hive

import (
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"image/color"
	"strconv"
)

type Hive struct {
	bees map[int]*Bee
}

func NewHive() *Hive {
	return &Hive{make(map[int]*Bee)}
}

func (h *Hive) AddBee(b *Bee, index int) {
	h.bees[index] = b
}

func (h *Hive) RemoveBee(index int) {
	delete(h.bees, index)
}

func (h *Hive) GetBee(index int) *Bee {
	return h.bees[index]
}

func (h *Hive) GetBees() map[int]*Bee {
	return h.bees
}

func (h *Hive) ToBson() bson.D {
	bees := bson.D{}
	for i := 0; i < 50; i++ {
		bee := h.bees[i+1]
		if bee != nil {
			bees = append(bees, bson.E{Key: strconv.Itoa(i + 1), Value: bee.ToBson()})
		}
	}
	return bson.D{{"bees", bees}}
}

var (
	slotColor = "#7d5a2f"
	bgColor   = "#ba8441"
	offsetX   = 10
	offsetY   = 20
)

func (h *Hive) Draw(dc *gg.Context, showHiveNumbers bool) {
	dc.SetHexColor(bgColor)
	dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()
	bottom := dc.Height()
	postProcessFuncs := make([]func(), 0)
	for i := 0; i < 10; i++ {
		bottomCnt := 0
		topCnt := 0
		for j := 0; j < 5; j++ {
			bee := h.bees[i*5+j+1]
			if j%2 == 0 {
				x := bottomCnt*50*3 + 50 + offsetX
				y := bottom - (i*80 + 50) - offsetY
				dc.DrawRegularPolygon(6, float64(x), float64(y), 50, 0)
				dc.SetHexColor(slotColor)
				dc.Fill()
				if bee != nil {
					postProcessFuncs = append(postProcessFuncs, bee.Draw(dc, x, y))
				} else {
					dc.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
					dc.SetHexColor(bgColor)
					dc.Fill()
					if showHiveNumbers {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 20)
						dc.SetFontFace(ff)
						dc.SetColor(color.Black)
						dc.DrawStringAnchored(strconv.Itoa(i*5+j+1), float64(x), float64(y), 0.5, 0.5)
					}
				}
				bottomCnt++
			} else {
				x := topCnt*50*3 + 75 + 50 + offsetX
				y := bottom - (i*80 + 15 + 25 + 50) - offsetY
				dc.DrawRegularPolygon(6, float64(x), float64(y), 50, 0)
				dc.SetHexColor(slotColor)
				dc.Fill()
				if bee != nil {
					postProcessFuncs = append(postProcessFuncs, bee.Draw(dc, x, y))
				} else {
					dc.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
					dc.SetHexColor(bgColor)
					dc.Fill()
					if showHiveNumbers {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 20)
						dc.SetFontFace(ff)
						dc.SetColor(color.Black)
						dc.DrawStringAnchored(strconv.Itoa(i*5+j+1), float64(x), float64(y), 0.5, 0.5)
					}
				}
				topCnt++
			}
		}
	}
	dd := gg.NewContext(410, 900)
	err := dc.SetMask(dd.AsMask())
	if err != nil {
		panic(err)
	}
	dc.InvertMask()
	ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Bold.ttf", 40)
	dc.SetFontFace(ff)
	for _, f := range postProcessFuncs {
		f()
	}
}
