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
	Name      string
	Path      string
	Kind      int64
	Abilities []string
}

func p(s string) string {
	return "assets/faces/" + s + ".png"
}

func b(name string, kind int64, abilities []string) BeeMeta {
	return BeeMeta{
		Name:      strings.ToUpper(name[:1]) + name[1:] + " Bee",
		Path:      p(name),
		Kind:      kind,
		Abilities: abilities,
	}
}

func a(abilities ...string) []string {
	return abilities
}

var bees = map[string]BeeMeta{
	"basic":     b("basic", Common, a()),
	"bomber":    b("bomber", Rare, a("Buzz Bomb")),
	"brave":     b("brave", Rare, a()),
	"bumble":    b("bumble", Rare, a("Blue Bomb")),
	"cool":      b("cool", Rare, a("Blue Boost")),
	"hasty":     b("hasty", Rare, a("Haste")),
	"looker":    b("looker", Rare, a("Focus")),
	"rad":       b("rad", Rare, a("Red Boost")),
	"rascal":    b("rascal", Rare, a("Red Bomb")),
	"stubborn":  b("stubborn", Rare, a("Pollen Mark")),
	"bubble":    b("bubble", Epic, a("Blue Bomb", "Passive: Gathering Bubbles")),
	"bucko":     b("bucko", Epic, a("Blue Boost")),
	"commander": b("commander", Epic, a("Focus", "Buzz Bomb")),
	"demo":      b("demo", Epic, a("Buzz Bomb+")),
	"exhausted": b("exhausted", Epic, a("Buzz Bomb", "Token Link")),
	"fire":      b("fire", Epic, a("Red Bomb+", "Passive: Gathering Flames")),
	"frosty":    b("frosty", Epic, a("Blue Boost", "Blue Bomb+")),
	"honey":     b("honey", Epic, a("Honey Gift", "Honey Mark")),
	"rage":      b("rage", Epic, a("Rage", "Token Link")),
	"riley":     b("riley", Epic, a("Red Boost")),
	"shocked":   b("shocked", Epic, a("Haste", "Token Link")),
	"baby":      b("baby", Legendary, a("Baby Love")),
	"carpenter": b("carpenter", Legendary, a("Honey Mark+", "Pollen Mark")),
	"demon":     b("demon", Legendary, a("Red Bomb", "Red Bomb+", "Passive: Gathering Flames+")),
	"diamond":   b("diamond", Legendary, a("Blue Boost", "Honey Gift+", "Passive: Shimmering Honey")),
	"lion":      b("lion", Legendary, a("Buzz Bomb+")),
	"music":     b("music", Legendary, a("Melody", "Focus", "Token Link")),
	"ninja":     b("ninja", Legendary, a("Haste", "Blue Bomb+")),
	"shy":       b("shy", Legendary, a("Red Boost", "Red Bomb", "Passive: Nectar Lover")),
	"buoyant":   b("buoyant", Mythic, a("Blue Bomb", "Inflate Balloons", "Gifted: Surprise Party", "Passive: Balloon Enthusiast")),
	"fuzzy":     b("fuzzy", Mythic, a("Fuzz Bombs", "Buzz Bomb+", "Gifted: Pollen Haze", "Passive: Fuzzy Coat")),
	"precise":   b("precise", Mythic, a("Target Practice", "Passive: Sniper")),
	"spicy":     b("spicy", Mythic, a("Inferno", "Rage", "Gifted: Flame Fuel", "Passive: Steam Engine")),
	"tadpole":   b("tadpole", Mythic, a("Blue Boost", "Summon Frog", "Gifted: Baby Love", "Passive: Gathering Bubbles+")),
	"vector":    b("vector", Mythic, a("Pollen Mark+", "Triangulate", "Gifted: Mark Surge")),
	"bear":      b("bear", Event, a("Bear Morph")),
	"cobalt":    b("cobalt", Event, a("Blue Pulse", "Blue Bomb Sync")),
	"crimson":   b("crimson", Event, a("Red Pulse", "Red Bomb Sync")),
	"digital":   b("digital", Event, a("Glitch", "Mind Hack", "Gifted: Map Corruption", "Passive: Drive Expansion")),
	"festive":   b("festive", Event, a("Festive Gift", "Honey Mark", "Red Bomb+", "Festive Mark")),
	"gummy":     b("gummy", Event, a("Gumdrop Barrage", "Glob")),
	"photon":    b("photon", Event, a("Beamstorm", "Haste")),
	"puppy":     b("puppy", Event, a("Fetch", "Puppy Love", "Focus")),
	"tabby":     b("tabby", Event, a("Scratch", "Tabby Love")),
	"vicious":   b("vicious", Event, a("Impale", "Blue Bomb+")),
	"windy":     b("windy", Event, a("White Boost", "Rain Cloud", "Tornado")),
}
