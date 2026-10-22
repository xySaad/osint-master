package socialmedia

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"osintmaster/src/lib"
	"regexp"
	"strings"
	"time"
)

const (
	linkedinURL   = "https://www.linkedin.com"
	linkedinBotUA = "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"
)

var (
	linkedinMetaTagRe     = regexp.MustCompile(`(?i)<meta\s+[^>]+>`)
	linkedinMetaPropRe    = regexp.MustCompile(`(?i)(?:property|name)=["']([^"']+)["']`)
	linkedinMetaContentRe = regexp.MustCompile(`(?i)content=["']([^"']*)["']`)
	linkedinTitleRe       = regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	linkedinLDJsonRe      = regexp.MustCompile(`(?is)<script\s+[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	connectionsRe         = regexp.MustCompile(`(?i)([0-9.,KMkm+]+)\s*connections`)
	followersRe           = regexp.MustCompile(`(?i)([0-9.,KMkm+]+)\s*followers`)
)

type LinkedinProfile struct {
	Username     string
	FullName     string
	Headline     string
	Company      string
	Location     string
	Education    string
	Followers    string
	Connections  string
	URL          string
	LastActivity *Activity
}

func parseLinkedinMetaTags(htmlContent string) map[string]string {
	meta := make(map[string]string)
	for _, tag := range linkedinMetaTagRe.FindAllString(htmlContent, -1) {
		nameMatch := linkedinMetaPropRe.FindStringSubmatch(tag)
		contentMatch := linkedinMetaContentRe.FindStringSubmatch(tag)
		if len(nameMatch) > 1 && len(contentMatch) > 1 {
			meta[strings.ToLower(nameMatch[1])] = html.UnescapeString(contentMatch[1])
		}
	}
	return meta
}

func parseTimeString(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse time: %s", timeStr)
}

func extractFirstMapString(items any, key string) string {
	slice, ok := items.([]any)
	if !ok {
		return ""
	}
	for _, it := range slice {
		if m, ok := it.(map[string]any); ok {
			if val, ok := m[key].(string); ok && strings.TrimSpace(val) != "" {
				return strings.TrimSpace(val)
			}
		}
	}
	return ""
}

func parseInteractionStats(stats any, profile *LinkedinProfile) {
	parseItem := func(item map[string]any) {
		name, _ := item["name"].(string)
		var userCount string
		switch v := item["userInteractionCount"].(type) {
		case float64:
			userCount = fmt.Sprintf("%.0f", v)
		case string:
			userCount = v
		}

		if strings.EqualFold(name, "Follows") || strings.Contains(strings.ToLower(name), "follow") {
			if profile.Followers == "" && userCount != "" {
				profile.Followers = userCount
			}
		} else if strings.EqualFold(name, "Connections") || strings.Contains(strings.ToLower(name), "connection") {
			if profile.Connections == "" && userCount != "" {
				profile.Connections = userCount
			}
		}
	}

	switch v := stats.(type) {
	case map[string]any:
		parseItem(v)
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				parseItem(m)
			}
		}
	}
}

func extractPersonAndActivities(node any, profile *LinkedinProfile, activities *[]Activity) {
	switch v := node.(type) {
	case map[string]any:
		typeVal, _ := v["@type"].(string)

		if strings.EqualFold(typeVal, "Person") {
			if name, ok := v["name"].(string); ok && profile.FullName == "" {
				profile.FullName = strings.TrimSpace(name)
			}
			if desc, ok := v["description"].(string); ok && profile.Headline == "" {
				desc = strings.ReplaceAll(desc, "<br>", " ")
				desc = strings.ReplaceAll(desc, "<br/>", " ")
				desc = strings.ReplaceAll(desc, "<br />", " ")
				profile.Headline = strings.TrimSpace(html.UnescapeString(desc))
			}
			if sameAs, ok := v["sameAs"].(string); ok && sameAs != "" {
				profile.URL = strings.TrimSpace(sameAs)
			} else if url, ok := v["url"].(string); ok && profile.URL == "" {
				profile.URL = strings.TrimSpace(url)
			}

			// Location
			if addr, ok := v["address"].(map[string]any); ok {
				loc, _ := addr["addressLocality"].(string)
				country, _ := addr["addressCountry"].(string)
				if loc != "" {
					profile.Location = loc
				} else if country != "" {
					profile.Location = country
				}
			}

			// Works for / Company
			if comp := extractFirstMapString(v["worksFor"], "name"); comp != "" && profile.Company == "" {
				profile.Company = comp
			}

			// Education
			if edu := extractFirstMapString(v["alumniOf"], "name"); edu != "" && profile.Education == "" {
				profile.Education = edu
			}

			// Interaction statistics (followers / connections)
			if stats, exists := v["interactionStatistic"]; exists {
				parseInteractionStats(stats, profile)
			}
		}

		if strings.EqualFold(typeVal, "DiscussionForumPosting") || strings.EqualFold(typeVal, "Article") {
			actURL, _ := v["url"].(string)
			if actURL == "" {
				actURL, _ = v["mainEntityOfPage"].(string)
			}
			dateStr, _ := v["datePublished"].(string)
			actTime, _ := parseTimeString(dateStr)

			actType := "Post"
			if strings.EqualFold(typeVal, "Article") {
				actType = "Article"
			}

			*activities = append(*activities, Activity{
				Type: actType,
				Url:  actURL,
				Time: actTime,
			})
		}

		// Traverse graph or mainEntity if nested
		if graph, exists := v["@graph"]; exists {
			extractPersonAndActivities(graph, profile, activities)
		}
		if mainEntity, exists := v["mainEntity"]; exists {
			extractPersonAndActivities(mainEntity, profile, activities)
		}

	case []any:
		for _, item := range v {
			extractPersonAndActivities(item, profile, activities)
		}
	}
}

func parseLinkedinHTML(username, htmlContent string) (*LinkedinProfile, error) {
	profile := &LinkedinProfile{
		Username: username,
		URL:      fmt.Sprintf("%s/in/%s", linkedinURL, username),
	}

	var activities []Activity

	// Parse JSON-LD scripts
	matches := linkedinLDJsonRe.FindAllStringSubmatch(htmlContent, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		var decoded any
		if err := json.Unmarshal([]byte(match[1]), &decoded); err == nil {
			extractPersonAndActivities(decoded, profile, &activities)
		}
	}

	// Meta tags fallback
	meta := parseLinkedinMetaTags(htmlContent)

	if profile.FullName == "" {
		if ogTitle, ok := meta["og:title"]; ok {
			parts := strings.Split(ogTitle, " - ")
			if len(parts) > 0 {
				profile.FullName = strings.TrimSpace(parts[0])
			}
		} else if titleMatch := linkedinTitleRe.FindStringSubmatch(htmlContent); len(titleMatch) > 1 {
			titleText := html.UnescapeString(titleMatch[1])
			parts := strings.Split(titleText, " - ")
			if len(parts) > 0 && !strings.Contains(strings.ToLower(parts[0]), "linkedin") {
				profile.FullName = strings.TrimSpace(parts[0])
			}
		}
	}

	if profile.Headline == "" {
		if ogDesc, ok := meta["og:description"]; ok {
			profile.Headline = strings.TrimSpace(ogDesc)
		} else if desc, ok := meta["description"]; ok {
			profile.Headline = strings.TrimSpace(desc)
		}
	}

	if profile.Connections == "" {
		if desc, ok := meta["description"]; ok {
			if cm := connectionsRe.FindStringSubmatch(desc); len(cm) > 1 {
				profile.Connections = cm[1]
			}
		}
	}

	if profile.Followers == "" {
		if desc, ok := meta["description"]; ok {
			if fm := followersRe.FindStringSubmatch(desc); len(fm) > 1 {
				profile.Followers = fm[1]
			}
		}
	}

	// Pick latest activity
	if len(activities) > 0 {
		latest := activities[0]
		for _, act := range activities[1:] {
			if act.Time.After(latest.Time) {
				latest = act
			}
		}
		profile.LastActivity = &latest
	}

	// Verification check: ensure we actually resolved a person profile
	if profile.FullName == "" && len(activities) == 0 {
		return nil, fmt.Errorf("user %q not found on linkedin", username)
	}

	return profile, nil
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

	bodyStr := string(bodyBytes)
	if resp.StatusCode == http.StatusNotFound || strings.Contains(bodyStr, "Profile Not Found | LinkedIn") || strings.Contains(bodyStr, "Page not found") {
		return "", fmt.Errorf("user %q not found on linkedin", username)
	}
	if resp.StatusCode == 999 && strings.Contains(bodyStr, "/authwall") {
		return "", fmt.Errorf("user %q not found on linkedin", username)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 999 {
		return "", fmt.Errorf("linkedin returned status %d for user %q", resp.StatusCode, username)
	}

	return bodyStr, nil
}

func GetLinkedinProfileAndActivity(username string) (*LinkedinProfile, error) {
	htmlContent, err := linkedinFetchHTML(username)
	if err != nil {
		return nil, err
	}
	return parseLinkedinHTML(username, htmlContent)
}

func LinkedinByUsername(username string) {
	profile, err := GetLinkedinProfileAndActivity(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("  - Linkedin -")
	fmt.Println("Username:", profile.Username)
	lib.PrintIfSet("Full Name:", profile.FullName)
	lib.PrintIfSet("Headline:", profile.Headline)
	lib.PrintIfSet("Company:", profile.Company)
	lib.PrintIfSet("Location:", profile.Location)
	lib.PrintIfSet("Education:", profile.Education)
	lib.PrintIfSet("Followers:", profile.Followers)
	lib.PrintIfSet("Connections:", profile.Connections)
	lib.PrintIfSet("URL:", profile.URL)

	if profile.LastActivity == nil {
		fmt.Println("Last Activity: N/A")
	} else {
		fmt.Println("Last Activity:")
		fmt.Println(" Type:", profile.LastActivity.Type)
		if !profile.LastActivity.Time.IsZero() {
			fmt.Println(" Time:", profile.LastActivity.Time)
		}
		lib.PrintIfSet(" URL:", profile.LastActivity.Url)
	}
}
