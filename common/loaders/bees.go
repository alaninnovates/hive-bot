package loaders

import (
	"github.com/fogleman/gg"
	"image"
	"sort"
)

func GetBee(id string) (image.Image, BeeMeta) {
	img, err := gg.LoadImage(bees[id].Path)
	dc := gg.NewContext(42, 42)
	dc.Scale(0.12, 0.12)
	dc.DrawImage(img, -20, -20)
	if err != nil {
		panic(err)
	}
	return dc.Image(), bees[id]
}

func GetBeeName(id string) string {
	return bees[id].Name
}

func GetBeeId(name string) string {
	for k, v := range bees {
		if v.Name == name {
			return k
		}
	}
	return ""
}

func GetBeeIds() []string {
	beeIds := make([]string, 0)
	for k := range bees {
		beeIds = append(beeIds, k)
	}
	sort.Strings(beeIds)
	return beeIds
}

func GetBeeNames() []string {
	beeNames := make([]string, 0)
	for _, v := range bees {
		beeNames = append(beeNames, v.Name)
	}
	sort.Strings(beeNames)
	return beeNames
}
