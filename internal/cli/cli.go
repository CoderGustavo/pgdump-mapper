package cli

import (
	"fmt"
	"os"

	models "github.com/hedibertosilva/pgdump-mapper/models"
)

var Options *models.Options

func HandleOptions(args []string) {
	hasOptions := false
	mapOptions := map[string]*bool{
		"-h":       &Options.Help,
		"--help":   &Options.Help,
		"--json":   &Options.Json,
		"--yaml":   &Options.Yaml,
		"--html":   &Options.Html,
		"--sqlite": &Options.Sqlite,
	}

	for _, arg := range args {
		if opt, exist := mapOptions[arg]; exist {
			*opt = true
			hasOptions = true
		}
	}

	// Set HTML as default
	if !hasOptions {
		Options.Html = true
	}
}

func ReturnError(err error) {
	if err != nil {
		fmt.Print(err, "\n\n")
	}

	os.Exit(1)
}

func ReturnSuccess(msg string) {
	if msg != "" {
		fmt.Print(msg, "\n")
	}

	os.Exit(0)
}
