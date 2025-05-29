package cmd

import (
	"fmt"
	"os"

	structs "github.com/hedibertosilva/pgdump-mapper/structures"
)

func HandleOptions(args []string, options *structs.Options) bool {
	hasOptions := false

	mapOptions := map[string]*bool{
		"--help": &options.Help,
		"-h":     &options.Help,
		"--raw":  &options.Raw,
		"-r":     &options.Raw,
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
