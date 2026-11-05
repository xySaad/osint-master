package ip

import (
	"fmt"
	"os"
	"osintmaster/src/lib"
)

const INVALID_COMMAND_USAGE = `Error: Invalid usage
` + USAGE

const USAGE = `Usage: 
    osintmaster ip <ip_address>
Arguments:
    -o <FileName> Output destination
    --help Show help message
`

func Cmd() {
	if _, helpOption := lib.GetOption("--help"); helpOption {
		fmt.Fprint(os.Stdout, USAGE)
		return
	}

	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
		return
	}

	ip := os.Args[1]
	GetIPInfo(ip)
	fmt.Println()
	GetAbuseHistory(ip)
}
