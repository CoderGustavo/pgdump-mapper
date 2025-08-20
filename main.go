package main

import (
	"fmt"
	"os"
	"strings"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	errors "github.com/hedibertosilva/pgdump-mapper/internal/cli/errors"
	file "github.com/hedibertosilva/pgdump-mapper/internal/file"

	messages "github.com/hedibertosilva/pgdump-mapper/internal/cli/messages"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func argsSanityCheck(args []string) error {
	numOptions := 0
	for _, arg := range args {
		if _, exist := models.CatalogOptions[arg]; exist {
			numOptions++
			continue
		}

		// If the argument follow --key=value format, check if the key is in CatalogOptions
		if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			if _, exist := models.CatalogOptions[key]; exist {
				numOptions++
				continue
			}
		}
	}

	argsLength := len(args) - numOptions
	if argsLength == 0 {
		return fmt.Errorf(errors.ErrorNoInputFile)
	} else if argsLength > 1 {
		return fmt.Errorf(errors.ErrorManyArgs)
	}

	return nil
}

func inputSanityCheck(input string) error {
	inputInfo, err := os.Stat(input)
	if err != nil {
		return fmt.Errorf(errors.ErrNoSuchFile)
	}

	if mode := inputInfo.Mode(); mode.IsDir() {
		return fmt.Errorf(errors.ErrIsDirectory)
	}

	return nil
}

func main() {
	defer os.Exit(0)

	var args = os.Args[1:]
	var opts = models.Options{}

	cli.Options = &opts
	cli.HandleOptions(args)

	if opts.Help {
		msg := messages.HelpContent
		cli.ReturnSuccess(msg)
	}

	if err := argsSanityCheck(args); err != nil {
		cli.ReturnError(err)
	}

	input := os.Args[1]
	if err := inputSanityCheck(input); err != nil {
		cli.ReturnError(err)
	}

	file.Input = &input
	file.Options = opts

	file.Read()
	file.Export()
}
