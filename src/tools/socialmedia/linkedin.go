package socialmedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"osintmaster/src/lib"
	"regexp"
	"time"
)

const (
	linkedinURL   = "https://www.linkedin.com"
	linkedinBotUA = "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"
)

var linkedinScriptRe = regexp.MustCompile(`(?is)<script\s+[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)

type LinkedinProfile struct {
	MainEntity ldPerson `json:"mainEntity"`
	Graph      []ldNode `json:"@graph"`
}

type ldPerson struct {
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	SameAs               string          `json:"sameAs"`
	URL                  string          `json:"url"`
	Address              ldAddress       `json:"address"`
	WorksFor             []ldOrg         `json:"worksFor"`
	AlumniOf             []ldOrg         `json:"alumniOf"`
	InteractionStatistic []ldInteraction `json:"interactionStatistic"`
}

type ldAddress struct {
	AddressLocality string `json:"addressLocality"`
	AddressCountry  string `json:"addressCountry"`
}

type ldOrg struct {
	Name string `json:"name"`
}

type ldInteraction struct {
	Name                 string      `json:"name"`
	UserInteractionCount json.Number `json:"userInteractionCount"`
}

type ldNode struct {
	Type          string `json:"@type"`
	URL           string `json:"url"`
	DatePublished string `json:"datePublished"`
}

func linkedinFetchHTML(username string) (string, error) {
	url := fmt.Sprintf("%s/in/%s", linkedinURL, username)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", linkedinBotUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.Header.Set("User-Agent", linkedinBotUA)
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch linkedin data: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	body := string(bodyBytes)
	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("user %q not found on linkedin", username)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 999 {
		return "", fmt.Errorf("linkedin returned status %d for user %q", resp.StatusCode, username)
	}

	return body, nil
}

func GetLinkedinProfileAndActivity(username string) (*LinkedinProfile, error) {
	htmlContent, err := linkedinFetchHTML(username)
	if err != nil {
		return nil, err
	}

	match := linkedinScriptRe.FindStringSubmatch(htmlContent)
	if match == nil {
		return nil, fmt.Errorf("user %q not found on linkedin", username)
	}
	var profile LinkedinProfile
	if err := json.Unmarshal([]byte(match[1]), &profile); err != nil {
		return nil, fmt.Errorf("failed to parse linkedin data for %q: %w", username, err)
	}

	return &profile, nil
}

func LinkedinByUsername(username string) {
	profile, err := GetLinkedinProfileAndActivity(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("  - Linkedin -")
	lib.PrintIfSet("Full Name:", profile.MainEntity.Name)
	lib.PrintIfSet("Description:", profile.MainEntity.Description)
	lib.PrintIfSet("Country:", profile.MainEntity.Address.AddressCountry)
	lib.PrintIfSet("Locality:", profile.MainEntity.Address.AddressLocality)
	for _, inter := range profile.MainEntity.InteractionStatistic {
		fmt.Println(inter.Name, ":", inter.UserInteractionCount)
	}
	var lastActivity *Activity
	for _, node := range profile.Graph {
		t, _ := time.Parse(time.RFC3339, node.DatePublished)
		if lastActivity == nil || t.After(lastActivity.Time) {
			act := Activity{Type: node.Type, Url: node.URL, Time: t}
			lastActivity = &act
		}
	}

	if lastActivity != nil {
		fmt.Println("Last Activity:")
		fmt.Println(" Type:", lastActivity.Type)
		fmt.Println(" Time:", lastActivity.Time)
		lib.PrintIfSet(" URL:", lastActivity.Url)
	} else {
		fmt.Println("Last Activity: N/A")
	}
}
