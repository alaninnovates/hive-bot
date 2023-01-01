package loaders

var mutations = map[string]string{
	"None":             "#FFFFFF",
	"Attack":           "#FF0000",
	"Convert Amount":   "#FCD049",
	"Gather Amount":    "#9DD678",
	"Energy":           "#8FC1CA",
	"Bee Ability Rate": "#A792CA",
}

func GetMutation(name string) string {
	return mutations[name]
}

func GetMutations() []string {
	mutationNames := make([]string, 0)
	for k := range mutations {
		mutationNames = append(mutationNames, k)
	}
	return mutationNames
}
