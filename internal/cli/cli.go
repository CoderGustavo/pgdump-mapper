package cli

import (
	"fmt"
	"os"
	"strings"

	models "github.com/hedibertosilva/pgdump-mapper/models"
)

var Options = &models.Options{}
var Filters = &models.FilterOptions{}

func HandleOptions(args []string) {
	hasOptions := false
	mapOptions := map[string]*bool{
		"-h":            &Options.Help,
		"--help":        &Options.Help,
		"--json":        &Options.JSON,
		"--json-pretty": &Options.JSONPretty,
		"--yaml":        &Options.YAML,
		"--html":        &Options.HTML,
		"--sqlite":      &Options.SQLite,
		"--cli":         &Options.CLI,
		"--cache":       &Options.Cache,
	}

	for _, arg := range args {
		if opt, exist := mapOptions[arg]; exist {
			*opt = true
			hasOptions = true
			continue
		}

		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			value := parts[1]

			switch key {
			case "--schema":
				Filters.Schema = value
				hasOptions = true
			case "--table":
				Filters.TableName = value
				hasOptions = true
			case "--columns":
				Filters.Columns = strings.Split(value, ",")
				hasOptions = true
			}
		}
	}

	// Set HTML as default
	if !hasOptions {
		Options.HTML = true
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
