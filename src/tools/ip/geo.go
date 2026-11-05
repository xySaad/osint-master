package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"osintmaster/src/lib"
	"time"
)

type ipApiResponse struct {
	Status     string  `json:"status"`
	Message    string  `json:"message"`
	Country    string  `json:"country"`
	RegionName string  `json:"regionName"`
	City       string  `json:"city"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	Timezone   string  `json:"timezone"`
	ISP        string  `json:"isp"`
	Org        string  `json:"org"`
	AS         string  `json:"as"`
}

func GetIPInfo(ip string) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,regionName,city,lat,lon,timezone,isp,org,as", ip)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching geolocation data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var data ipApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing geolocation data: %v\n", err)
		return
	}

	if data.Status != "success" {
		fmt.Fprintf(os.Stderr, "Error fetching geolocation data: %s\n", data.Message)
		return
	}

	fmt.Println("  - Geolocation -")
	lib.PrintIfSet("Country:", data.Country)
	lib.PrintIfSet("Region:", data.RegionName)
	lib.PrintIfSet("City:", data.City)
	fmt.Printf("Latitude: %g\n", data.Lat)
	fmt.Printf("Longitude: %g\n", data.Lon)
	lib.PrintIfSet("Timezone:", data.Timezone)

	fmt.Println("\n  - ISP Details -")
	lib.PrintIfSet("ISP:", data.ISP)
	lib.PrintIfSet("Organization:", data.Org)
	lib.PrintIfSet("AS:", data.AS)
}
