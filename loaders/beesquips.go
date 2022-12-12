package loaders

import (
	"github.com/fogleman/gg"
	"image"
	"io/fs"
	"os"
	"strings"
)

func GetBeequipImage(name string) image.Image {
	gd := gg.NewContext(80, 80)
	img, err := gg.LoadImage("assets/beequips/" + name + ".png")
	if err != nil {
		panic(err)
	}
	gd.Scale(0.3, 0.3)
	gd.DrawImage(img, 0, 0)
	return gd.Image()
}

func GetBeequips() []string {
	var files []string
	dir, _ := os.Open("assets/beequips")
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
