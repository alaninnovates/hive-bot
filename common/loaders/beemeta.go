package loaders

import "strings"

const (
	Common = iota
	Rare
	Epic
	Legendary
	Mythic
	Event
)

type BeeMeta struct {
	Name string
	Path string
	Kind int64
}

func p(s string) string {
	return "assets/faces/" + s + ".png"
}

func b(name string, kind int64) BeeMeta {
	return BeeMeta{
		Name: strings.ToUpper(name[:1]) + name[1:] + " Bee",
		Path: p(name),
		Kind: kind,
	}
}

var bees = map[string]BeeMeta{
	"basic":     b("basic", Common),
	"bomber":    b("bomber", Rare),
	"brave":     b("brave", Rare),
	"bumble":    b("bumble", Rare),
	"cool":      b("cool", Rare),
	"hasty":     b("hasty", Rare),
	"looker":    b("looker", Rare),
	"rad":       b("rad", Rare),
	"rascal":    b("rascal", Rare),
	"stubborn":  b("stubborn", Rare),
	"bubble":    b("bubble", Epic),
	"bucko":     b("bucko", Epic),
	"commander": b("commander", Epic),
	"demo":      b("demo", Epic),
	"exhausted": b("exhausted", Epic),
	"fire":      b("fire", Epic),
	"frosty":    b("frosty", Epic),
	"honey":     b("honey", Epic),
	"rage":      b("rage", Epic),
	"riley":     b("riley", Epic),
	"shocked":   b("shocked", Epic),
	"baby":      b("baby", Legendary),
	"carpenter": b("carpenter", Legendary),
	"demon":     b("demon", Legendary),
	"diamond":   b("diamond", Legendary),
	"lion":      b("lion", Legendary),
	"music":     b("music", Legendary),
	"ninja":     b("ninja", Legendary),
	"shy":       b("shy", Legendary),
	"buoyant":   b("buoyant", Mythic),
	"fuzzy":     b("fuzzy", Mythic),
	"precise":   b("precise", Mythic),
	"spicy":     b("spicy", Mythic),
	"tadpole":   b("tadpole", Mythic),
	"vector":    b("vector", Mythic),
	"bear":      b("bear", Event),
	"cobalt":    b("cobalt", Event),
	"crimson":   b("crimson", Event),
	"digital":   b("digital", Event),
	"festive":   b("festive", Event),
	"gummy":     b("gummy", Event),
	"photon":    b("photon", Event),
	"puppy":     b("puppy", Event),
	"tabby":     b("tabby", Event),
	"vicious":   b("vicious", Event),
	"windy":     b("windy", Event),
}
