package loaders

import (
	"github.com/disgoorg/disgo/discord"
	"io/fs"
	"os"
	"strings"
)

func GetHiveBackgroundImagePath(name string) string {
	return "assets/backgrounds/" + name + ".png"
}

func GetHiveBackgrounds() []string {
	var files []string
	dir, _ := os.Open("assets/backgrounds")
	defer dir.Close()
	fileInfos, err := dir.ReadDir(-1)
	if err != nil {
		panic(err)
	}
	for _, fi := range fileInfos {
		if fi.Name() == ".DS_Store" {
			continue
		}
		if fi.Type() == fs.FileMode(0) {
			files = append(files, strings.Split(fi.Name(), ".")[0])
		}
	}
	return files
}

func GetHiveBackgroundsChoices() []discord.ApplicationCommandOptionChoiceString {
	var choices []discord.ApplicationCommandOptionChoiceString
	for _, background := range GetHiveBackgrounds() {
		choices = append(choices, discord.ApplicationCommandOptionChoiceString{
			Name:  background,
			Value: background,
		})
	}
	return choices
}
