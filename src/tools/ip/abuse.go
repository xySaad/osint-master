package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type pulse struct {
	Name    string   `json:"name"`
	Created string   `json:"created"`
	Tags    []string `json:"tags"`
}

type alienVaultResponse struct {
	Detail    string `json:"detail"`
	PulseInfo struct {
		Count  int     `json:"count"`
		Pulses []pulse `json:"pulses"`
	} `json:"pulse_info"`
}

func GetAbuseHistory(ip string) {
	ipType := "IPv4"
	if strings.Contains(ip, ":") {
		ipType = "IPv6"
	}

	url := fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/%s/%s/general?sections=general", ipType, ip)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating abuse request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching abuse data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp alienVaultResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Detail != "" {
			fmt.Fprintf(os.Stderr, "Error fetching abuse data: %s\n", errResp.Detail)
			return
		}
		fmt.Fprintf(os.Stderr, "Error fetching abuse data (status %d)\n", resp.StatusCode)
		return
	}

	var data alienVaultResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing abuse data: %v\n", err)
		return
	}

	fmt.Println("  - Historical Abuse Data -")
	fmt.Printf("Historical Reports: %d\n", data.PulseInfo.Count)
	if data.PulseInfo.Count > 0 {
		fmt.Println("Recent Reports:")
		max := len(data.PulseInfo.Pulses)
		if max > 5 {
			max = 5
		}
		for i := 0; i < max; i++ {
			p := data.PulseInfo.Pulses[i]
			date := p.Created
			if len(date) >= 10 {
				date = date[:10]
			}
			if len(p.Tags) > 0 {
				fmt.Printf(" - [%s] %s (Tags: %s)\n", date, p.Name, strings.Join(p.Tags, ", "))
			} else {
				fmt.Printf(" - [%s] %s\n", date, p.Name)
			}
		}
	}
}
