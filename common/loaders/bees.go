package loaders

import (
	"image"
	"sort"

	"github.com/fogleman/gg"
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

func GetBeeAbilities(beeName string) []string {
	return bees[GetBeeId(beeName)].Abilities
}

func GetBeeAbilityList() []string {
	abilities := make([]string, 0)
	for _, v := range bees {
		for _, a := range v.Abilities {
			abilities = append(abilities, a)
		}
	}
	sort.Strings(abilities)
	// remove duplicates
	for i := 0; i < len(abilities)-1; i++ {
		if abilities[i] == abilities[i+1] {
			abilities = append(abilities[:i], abilities[i+1:]...)
			i--
		}
	}
	return abilities
}
