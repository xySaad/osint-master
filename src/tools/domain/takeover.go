package domain

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"osintmaster/src/lib"
	"strings"
	"time"
)

type Fingerprint struct {
	Service     string   `json:"service"`
	CNAME       []string `json:"cname"`
	Fingerprint []string `json:"fingerprint"`
	NXDomain    bool     `json:"nxdomain"`
}

var client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func printPotentialTakeover(domain string, fingerprints []Fingerprint) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return
	}

	cname = strings.TrimSuffix(cname, ".")
	if cname == domain || cname == "" {
		return
	}
	for _, fp := range fingerprints {
		for _, pattern := range fp.CNAME {
			cleanPattern := strings.TrimPrefix(pattern, ".")
			if strings.Contains(cname, cleanPattern) {
				if isPotentialTakeOver(fp, cname, domain) {
					fmt.Printf("[POTENTIAL TAKEOVER] %s -> %s (Service: %s)\n", domain, cname, fp.Service)
				}
				break
			}
		}

	}
}

func isPotentialTakeOver(fp Fingerprint, cname, domain string) bool {
	if fp.NXDomain {
		_, err := net.LookupHost(cname)
		if lib.IsNotFoundDnsErr(err) {
			return true
		}
	}

	if len(fp.Fingerprint) < 1 {
		return false
	}
	for _, proto := range []string{"http://", "https://"} {
		req, err := http.NewRequest("GET", proto+domain, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}
		for _, match := range fp.Fingerprint {
			if strings.Contains(string(body), match) {
				return true
			}
		}
	}

	return false
}
