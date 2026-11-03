package main

import (
	"fmt"
	"os"
	"osintmaster/src/lib"
	"osintmaster/src/tools/domain"
	"osintmaster/src/tools/socialmedia"
	"slices"
)

const INVALID_COMMAND_USAGE = `Error: Invalid usage
` + USAGE

const USAGE = `Usage: 
	osintmaster [COMMAND] ...args
Commands:
	ip, domain, username
Arguments:
	-o <FileName> Output destination
	--help Show help message
`

func main() {
	if len(os.Args) == 2 {
		if _, helpOption := lib.GetOption("--help"); helpOption {
			fmt.Fprint(os.Stdout, USAGE)
			return
		}
	}

	if outputFilePath, dashO := lib.GetOption("-o"); dashO {
		if outputFilePath == "" {
			fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
			return
		}
		close, err := lib.CloneOutput(outputFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer close()
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
		return
	}

	subcommand := os.Args[1]
	os.Args = slices.Delete(os.Args, 1, 2)
	switch subcommand {
	case "ip":
		//
	case "domain":
		domain.Cmd()
	case "username":
		socialmedia.Cmd()
	default:
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
	}
}
