package domain

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"
)

const DEFAULT_CONCURRENCY = 30

func scanDomain(domain string, wordList *string, fingerprints []Fingerprint) {
	printPotentialTakeover(domain, fingerprints)
	addrs, err := net.LookupHost(domain)
	if err != nil {
		if _, ok := err.(*net.DNSError); !ok {
			fmt.Printf("domain lookup error (%s): %s\n", domain, reflect.TypeOf(err))
		}
		return
	}
	details, err := getSSLDetails(domain)
	if err == nil {
		fmt.Printf("%s %s\n%s\n", domain, addrs, details)
	}

	if wordList == nil {
		return
	}

	file, err := os.Open(*wordList)
	if err != nil {
		fmt.Printf("Error opening wordlist: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lock := make(chan struct{}, DEFAULT_CONCURRENCY)
	var wg sync.WaitGroup
	for scanner.Scan() {
		if scanner.Err() != nil {
			fmt.Println(err)
			return
		}
		subdomain := scanner.Text() + "." + domain
		wg.Add(1)
		go func() {
			lock <- struct{}{}
			defer func() { wg.Done(); <-lock }()
			scanDomain(subdomain, nil, fingerprints)
		}()
	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading wordlist:", err)
	}
}
