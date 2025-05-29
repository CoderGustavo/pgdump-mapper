package main

import (
	"fmt"
	"os"

	cmd "github.com/hedibertosilva/pgdump-mapper/cmd"
	consts "github.com/hedibertosilva/pgdump-mapper/constants"
	structs "github.com/hedibertosilva/pgdump-mapper/structures"
)

func argsSanityCheck(args []string) error {
	argsLength := len(args)
	if argsLength == 0 {
		return fmt.Errorf(consts.ErrorNoInputFile)
	} else if argsLength > 1 {
		return fmt.Errorf(consts.ErrorManyArgs)
	}

	return nil
}

func fileSanityCheck(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf(consts.ErrNoSuchFile)
	}

	if mode := fileInfo.Mode(); mode.IsDir() {
		return fmt.Errorf(consts.ErrIsDirectory)
	}

	return nil
}

func main() {

	defer os.Exit(0)

	var options = structs.Options{}
	var args = os.Args[1:]

	cmd.HandleOptions(args, &options)

	if options.Help {
		cmd.ReturnSuccess(consts.HelpContent)
	}

	if err := argsSanityCheck(args); err != nil {
		cmd.ReturnError(err)
	}

	file := os.Args[1]
	if err := fileSanityCheck(file); err != nil {
		cmd.ReturnError(err)
	}

	fmt.Println("Loading", file)
	fmt.Println("Load completed")
}
