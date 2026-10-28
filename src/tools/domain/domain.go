package domain

import (
	"fmt"
	"os"
	"osintmaster/src/options"
)

const INVALID_COMMAND_USAGE = `Error: Invalid usage
` + USAGE

const USAGE = `Usage: 
    osintmaster domain <subdomain.domain.tld>
Arguments:
    -o <FileName> Output destination
    --help Show help message
`

func ScanDNSAndSSl() {
	if _, helpOption := options.GetOption("--help"); helpOption {
		fmt.Fprint(os.Stdout, USAGE)
		return
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
		return
	}

	domain := os.Args[1]
	printSSLDetails(domain)

}
