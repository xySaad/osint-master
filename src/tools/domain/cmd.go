package domain

import (
	"fmt"
	"os"
	"osintmaster/src/lib"
)

const INVALID_COMMAND_USAGE = `Error: Invalid usage
` + USAGE

const USAGE = `Usage: 
    osintmaster domain <subdomain.domain.tld> <wordlist.txt>
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

	domain := os.Args[1]
	var wordlist *string = nil
	if len(os.Args) > 2 {
		tmp := os.Args[2]
		wordlist = &tmp
	}

	ScanDNSAndSSL(domain, wordlist)
}
