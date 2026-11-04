package domain

import (
	"encoding/json"
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

	if len(os.Args) > 4 {
		fmt.Fprint(os.Stderr, INVALID_COMMAND_USAGE)
		return
	}

	domain := os.Args[1]
	wordlist := "src/assets/subdomains.txt"
	fingerprints := "src/assets/fingerprints.json"
	if len(os.Args) > 2 {
		wordlist = os.Args[2]
	}
	if len(os.Args) > 3 {
		fingerprints = os.Args[3]

	}

	data, err := os.ReadFile(fingerprints)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading fingerprints: %v\n", err)
		return
	}

	var fps []Fingerprint
	if err := json.Unmarshal(data, &fps); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing fingerprints: %v\n", err)
		return
	}

	scanDomain(domain, &wordlist, fps)
}
