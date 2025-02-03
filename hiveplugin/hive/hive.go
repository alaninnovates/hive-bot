package hive

import (
	"alaninnovates.com/hive-bot/common"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"image/color"
	"strconv"
)

type Hive struct {
	bees         map[int][]*Bee
	lastModified int64
}

func NewHive() *Hive {
	return &Hive{make(map[int][]*Bee), 0}
}

func (h *Hive) AddBee(b *Bee, index int) {
	h.bees[index] = append(h.bees[index], b)
	h.lastModified = common.CurrentTimeMillis()
}

func (h *Hive) RemoveBee(index int) {
	delete(h.bees, index)
	h.lastModified = common.CurrentTimeMillis()
}

func (h *Hive) RemoveBeeAt(index int, beeIndex int) {
	h.bees[index] = append(h.bees[index][:beeIndex], h.bees[index][beeIndex+1:]...)
	h.lastModified = common.CurrentTimeMillis()
}

func (h *Hive) GetBeesAt(index int) []*Bee {
	return h.bees[index]
}

func (h *Hive) GetBees() map[int][]*Bee {
	return h.bees
}

func (h *Hive) LastModified() int64 {
	return h.lastModified
}

func (h *Hive) ToBson() bson.D {
	bees := bson.D{}
	for i := 0; i < 50; i++ {
		bee := h.bees[i+1]
		if bee != nil {
			beeList := bson.A{}
			for _, b := range bee {
				beeList = append(beeList, b.ToBson())
			}
			bees = append(bees, bson.E{Key: strconv.Itoa(i + 1), Value: beeList})
		}
	}
	return bson.D{{"bees", bees}}
}

var (
	slotColor = "#7d5a2f"
	bgColor   = "#E0B153"
	offsetX   = 18
	offsetY   = 15
)

func DrawHive(h *Hive, dc *gg.Context, showHiveNumbers bool, slotsOnTop bool, skipHiveNumbers []int) {
	h.lastModified = common.CurrentTimeMillis()
	//dc.SetHexColor(bgColor)
	//dc.DrawRectangle(0, 0, float64(dc.Width()), float64(dc.Height()))
	dc.Fill()
	bottom := dc.Height()
	postProcessFuncs := make([]map[string]func(), 0)
	for i := 0; i < 10; i++ {
		bottomCnt := 0
		topCnt := 0
		for j := 0; j < 5; j++ {
			hiveNumber := i*5 + j + 1
			bee := h.bees[hiveNumber]
			var hiveNumFunc func()
			if j%2 == 0 {
				x := bottomCnt*46*3 + 50 + offsetX
				y := bottom - (i*80 + 50) - offsetY
				dc.DrawRegularPolygon(6, float64(x), float64(y), 50, 0)
				dc.SetHexColor(slotColor)
				dc.Fill()
				hiveNumFunc = func() {
					if bee != nil && slotsOnTop {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 30)
						dc.SetFontFace(ff)
						dc.SetColor(color.White)
					} else {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 20)
						dc.SetFontFace(ff)
						dc.SetColor(color.Black)
					}
					dc.DrawStringAnchored(strconv.Itoa(hiveNumber), float64(x), float64(y), 0.5, 0.5)
				}
				if bee != nil && !common.ArrayIncludes(skipHiveNumbers, hiveNumber) {
					postProcessFuncs = append(postProcessFuncs, DrawBees(bee, dc, x, y))
				} else {
					dc.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
					dc.SetHexColor(bgColor)
					dc.Fill()
					if showHiveNumbers && !slotsOnTop {
						hiveNumFunc()
					}
				}
				if slotsOnTop {
					hiveNumFunc()
				}
				bottomCnt++
			} else {
				x := topCnt*46*3 + 70 + 50 + offsetX
				y := bottom - (i*80 + 15 + 25 + 50) - offsetY
				dc.DrawRegularPolygon(6, float64(x), float64(y), 50, 0)
				dc.SetHexColor(slotColor)
				dc.Fill()
				hiveNumFunc = func() {
					if bee != nil && slotsOnTop {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 30)
						dc.SetFontFace(ff)
						dc.SetColor(color.White)
					} else {
						ff, _ := gg.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 20)
						dc.SetFontFace(ff)
						dc.SetColor(color.Black)
					}
					dc.DrawStringAnchored(strconv.Itoa(hiveNumber), float64(x), float64(y), 0.5, 0.5)
				}
				if bee != nil && !common.ArrayIncludes(skipHiveNumbers, hiveNumber) {
					postProcessFuncs = append(postProcessFuncs, DrawBees(bee, dc, x, y))
				} else {
					dc.DrawRegularPolygon(6, float64(x), float64(y), 40, 0)
					dc.SetHexColor(bgColor)
					dc.Fill()
					if showHiveNumbers && !slotsOnTop {
						hiveNumFunc()
					}
				}
				if slotsOnTop {
					hiveNumFunc()
				}
				topCnt++
			}
		}
	}
	// add some credits
	//ff, _ := gg.LoadFontFace("assets/fonts/UniformRnd-Black.ttf", 40)
	//dc.SetFontFace(ff)
	//dc.SetColor(color.White)
	//dc.DrawStringAnchored("Hive Builder", float64(dc.Width()/2), 50, 0.5, 0)
	// simulate "layers" with post-processing functions
	dd := gg.NewContext(410, 950)
	err := dc.SetMask(dd.AsMask())
	if err != nil {
		panic(err)
	}
	dc.InvertMask()
	normalFont, _ := gg.LoadFontFace("assets/fonts/Roboto-Bold.ttf", 30)
	dc.SetFontFace(normalFont)
	for i := 0; i < 3; i++ {
		for _, f := range postProcessFuncs {
			if i == 0 {
				f["gifted"]()
			} else if i == 1 {
				f["beequip"]()
			} else if i == 2 {
				f["level"]()
			}
		}
	}
}
