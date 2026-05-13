package main

import (
	"fmt"
	"os"

	"github.com/gentleman-programming/gentle-ai/internal/app"
)

var version = "dev"

func main() {
	app.Version = app.ResolveVersion(version)
	app.SetRuntimeFlavor(app.Flavor{Name: "gentle-ai-crisp", ConfigNamespace: "gentle-ai-crisp"})
	if err := app.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
