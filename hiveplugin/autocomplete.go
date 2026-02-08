package hiveplugin

import (
	"strings"

	"alaninnovates.com/hive-bot/common/loaders"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var AddAutocomplete = makeAutocompleteHandler(loaders.GetBeeNames())
var SetBeequipAutocomplete = makeAutocompleteHandler(append(loaders.GetBeequips(), "None"))
var SetMutationAutocomplete = makeAutocompleteHandler(loaders.GetMutations())

func HiveViewAutocomplete(event *handler.AutocompleteEvent) error {
	focusedField := ""
	for key, option := range event.Data.Options {
		if option.Focused {
			focusedField = key
			break
		}
	}
	if focusedField == "ability" {
		ability := event.Data.String("ability")
		return getMatches(event, loaders.GetBeeAbilityList(), ability)
	} else if focusedField == "beequip" {
		beequip := event.Data.String("beequip")
		return getMatches(event, loaders.GetBeequips(), beequip)
	}
	return nil
}

func makeAutocompleteHandler(b []string) handler.AutocompleteHandler {
	return func(event *handler.AutocompleteEvent) error {
		//fmt.Printf("evt: %d now: %d", event.ID().Time().UnixMilli(), time.Now().UnixMilli())
		name := event.Data.String("name")
		return getMatches(event, b, name)
	}
}

func getMatches(event *handler.AutocompleteEvent, options []string, text string) error {
	text = strings.ToLower(text)
	matches := make([]discord.AutocompleteChoice, 0)
	i := 0
	for _, opt := range options {
		if i >= 25 {
			break
		}
		if strings.Contains(strings.ToLower(opt), text) {
			matches = append(matches, discord.AutocompleteChoiceString{
				Name:  opt,
				Value: opt,
			})
			i++
		}
	}
	return event.AutocompleteResult(matches)
}
