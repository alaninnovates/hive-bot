package loaders

import (
	"github.com/fogleman/gg"
	"image"
	"io/fs"
	"os"
	"strings"
)

func GetBeeImage(name string) image.Image {
	gd := gg.NewContext(80, 80)
	img, err := gg.LoadImage("assets/bees/" + name + ".png")
	if err != nil {
		panic(err)
	}
	gd.Scale(0.65, 0.65)
	gd.DrawImage(img, 0, 0)
	//gd.SavePNG("test.png")
	return gd.Image()
}

func GetBees() []string {
	var files []string
	dir, _ := os.Open("assets/bees")
	defer dir.Close()
	fileInfos, err := dir.ReadDir(-1)
	if err != nil {
		panic(err)
	}
	for _, fi := range fileInfos {
		if fi.Type() == fs.FileMode(0) {
			files = append(files, strings.Split(fi.Name(), ".")[0])
		}
	}
	return files
}
