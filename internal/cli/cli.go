package cli

import (
	"fmt"
	"os"

	models "github.com/hedibertosilva/pgdump-mapper/models"
)

var Options *models.Options

func HandleOptions(args []string) bool {
	hasOptions := false
	mapOptions := map[string]*bool{
		"--help": &Options.Help,
		"-h":     &Options.Help,
		"--raw":  &Options.Raw,
		"-r":     &Options.Raw,
		"--json": &Options.Json,
		"--yaml": &Options.Yaml,
		"--html": &Options.Html,
	}

	for _, arg := range args {
		if opt, exist := mapOptions[arg]; exist {
			*opt = true
			hasOptions = true
		}
	}

	return hasOptions
}

func ReturnError(err error) {
	if err != nil {
		fmt.Print(err, "\n\n")
	}

	os.Exit(1)
}

func ReturnSuccess(msg string) {
	if msg != "" {
		fmt.Printf(msg, "\n\n")
	}

	os.Exit(0)
}
