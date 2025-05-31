package cli

import (
	"fmt"
	"os"

	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func HandleOptions(args []string, opts *models.Options) bool {
	hasOptions := false

	mapOptions := map[string]*bool{
		"--help": &opts.Help,
		"-h":     &opts.Help,
		"--raw":  &opts.Raw,
		"-r":     &opts.Raw,
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
