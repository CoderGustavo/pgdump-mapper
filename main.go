package main

import (
	"fmt"
	"os"

	consts "github.com/hedibertosilva/pgdump-mapper/config"
)

func captureHelpRequest() bool {
	input := os.Args[1]
	helpOptions := []string{"--help", "-h"}
	for _, option := range helpOptions {
		if input == option {
			return true
		}
	}

	return false
}

func argsSanityCheck() (msg string, err error) {
	argsLength := len(os.Args[1:])
	if argsLength == 0 {
		return consts.ErrorNoInputFileMsg, fmt.Errorf(consts.ErrorNoInputFile)
	} else if argsLength > 1 {
		return consts.ErrorManyArgsMsg, fmt.Errorf(consts.ErrorManyArgs)
	}

	return "", nil
}

func fileSanityCheck(file string) (msg string, err error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return consts.ErrNoSuchFileMsg, fmt.Errorf(consts.ErrNoSuchFile)
	}

	if mode := fileInfo.Mode(); mode.IsDir() {
		return consts.ErrIsDirectoryMsg, fmt.Errorf(consts.ErrIsDirectory)
	}

	return "", nil
}

func main() {

	if msg, err := argsSanityCheck(); err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}

	if isUserLost := captureHelpRequest(); isUserLost {
		fmt.Printf(consts.HelpMsg)
		os.Exit(0)
	}

	file := os.Args[1]

	if msg, err := fileSanityCheck(file); err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}

	fmt.Println("Loading", file)
	fmt.Println("Load completed")
}
