package socialmedia

import (
	"fmt"
	"os"
	"osintmaster/src/lib"
	"time"
)

type Activity struct {
	Type string
	Url  string
	Time time.Time
}

const INVALID_COMMAND_USAGE = `Error: Invalid usage
` + USAGE

const USAGE = `Usage: 
    osintmaster username <username>
Arguments:
    -o <FileName> Output destination
    --help Show help message
`

func Cmd() {
	if _, helpOption := lib.GetOption("--help"); helpOption {
		fmt.Fprint(os.Stdout, USAGE)
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
		return
	}

	username := os.Args[1]
	GithubByUsername(username)
	fmt.Println()
	InstagramByUsername(username)
	fmt.Println()
	LinkedinByUsername(username)
	fmt.Println()
	FacebookByUsername(username)
	fmt.Println()
	XByUsername(username)
}
