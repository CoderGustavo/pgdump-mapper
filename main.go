package main

import (
	"fmt"
	"os"

	cli "github.com/hedibertosilva/pgdump-mapper/internal/cli"
	errors "github.com/hedibertosilva/pgdump-mapper/internal/cli/errors"
	models "github.com/hedibertosilva/pgdump-mapper/models"
)

func argsSanityCheck(args []string) error {
	argsLength := len(args)
	if argsLength == 0 {
		return fmt.Errorf(errors.ErrorNoInputFile)
	} else if argsLength > 1 {
		return fmt.Errorf(errors.ErrorManyArgs)
	}

	return nil
}

func fileSanityCheck(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf(errors.ErrNoSuchFile)
	}

	if mode := fileInfo.Mode(); mode.IsDir() {
		return fmt.Errorf(errors.ErrIsDirectory)
	}

	return nil
}

func main() {
	defer os.Exit(0)

	var opts = models.Options{}
	var args = os.Args[1:]
	cli.HandleOptions(args, &opts)

	if opts.Help {
		msg := cli.HelpContent
		cli.ReturnSuccess(msg)
	}

	if err := argsSanityCheck(args); err != nil {
		cli.ReturnError(err)
	}

	file := os.Args[1]
	if err := fileSanityCheck(file); err != nil {
		cli.ReturnError(err)
	}

	fmt.Println("Loading", file)
	fmt.Println("Load completed")
}
