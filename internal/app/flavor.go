package app

type Flavor struct {
	Name            string
	ConfigNamespace string
}

var RuntimeFlavor = Flavor{Name: "gentle-ai", ConfigNamespace: "gentle-ai"}

func SetRuntimeFlavor(flavor Flavor) {
	if flavor.Name == "" {
		flavor.Name = "gentle-ai"
	}
	if flavor.ConfigNamespace == "" {
		flavor.ConfigNamespace = flavor.Name
	}
	RuntimeFlavor = flavor
}
