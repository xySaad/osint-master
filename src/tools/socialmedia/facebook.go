package socialmedia

import (
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
	facebookBaseURL = "https://www.facebook.com"
	facebookBotUA   = "Twitterbot/1.0"
)

var (
	fbMetaTagRe     = regexp.MustCompile(`(?i)<meta\s+[^>]+>`)
	fbMetaPropRe    = regexp.MustCompile(`(?i)(?:property|name)=["']([^"']+)["']`)
	fbMetaContentRe = regexp.MustCompile(`(?i)content=["']([^"']*)["']`)
	fbTitleRe       = regexp.MustCompile(`(?i)<title>(.*?)</title>`)
	fbFollowersRe   = regexp.MustCompile(`(?i)([0-9.,KMkm]+)\s*followers`)
	fbTalkingRe     = regexp.MustCompile(`(?i)([0-9.,KMkm]+)\s*talking about this`)
)

type FacebookProfile struct {
	Username     string
	FullName     string
	Bio          string
	Followers    string
	TalkingAbout string
	URL          string
}

func parseFacebookMetaTags(htmlContent string) map[string]string {
	meta := make(map[string]string)
	for _, tag := range fbMetaTagRe.FindAllString(htmlContent, -1) {
		nameMatch := fbMetaPropRe.FindStringSubmatch(tag)
		contentMatch := fbMetaContentRe.FindStringSubmatch(tag)
		if len(nameMatch) > 1 && len(contentMatch) > 1 {
			meta[strings.ToLower(nameMatch[1])] = html.UnescapeString(contentMatch[1])
		}
	}
	return meta
}

func parseFacebookHTML(username, htmlContent string) (*FacebookProfile, error) {
	meta := parseFacebookMetaTags(htmlContent)

	profile := &FacebookProfile{
		Username: username,
		URL:      fmt.Sprintf("%s/%s", facebookBaseURL, username),
	}

	if ogURL, ok := meta["og:url"]; ok && ogURL != "" {
		profile.URL = ogURL
	}

	if ogTitle, ok := meta["og:title"]; ok && strings.TrimSpace(ogTitle) != "" {
		profile.FullName = strings.TrimSpace(ogTitle)
	} else if titleMatch := fbTitleRe.FindStringSubmatch(htmlContent); len(titleMatch) > 1 {
		rawTitle := html.UnescapeString(strings.TrimSpace(titleMatch[1]))
		lower := strings.ToLower(rawTitle)
		if lower != "facebook" && lower != "log in" && !strings.HasPrefix(lower, "error") && !strings.Contains(lower, "page not found") {
			if idx := strings.Index(rawTitle, " ("); idx != -1 {
				profile.FullName = strings.TrimSpace(rawTitle[:idx])
			} else if idx := strings.Index(rawTitle, " •"); idx != -1 {
				profile.FullName = strings.TrimSpace(rawTitle[:idx])
			} else {
				profile.FullName = rawTitle
			}
		}
	}

	if profile.FullName == "" || strings.EqualFold(profile.FullName, "Facebook") || strings.EqualFold(profile.FullName, "Log in") {
		return nil, fmt.Errorf("user %q not found on facebook", username)
	}

	if ogDesc, ok := meta["og:description"]; ok && ogDesc != "" {
		desc := html.UnescapeString(ogDesc)

		if fm := fbFollowersRe.FindStringSubmatch(desc); len(fm) > 1 {
			profile.Followers = fm[1]
		}
		if tm := fbTalkingRe.FindStringSubmatch(desc); len(tm) > 1 {
			profile.TalkingAbout = tm[1]
		}

		if idx := strings.Index(desc, "talking about this."); idx != -1 {
			rem := strings.TrimSpace(desc[idx+len("talking about this."):])
			if rem != "" {
				profile.Bio = rem
			}
		} else if strings.Contains(desc, "·") {
			parts := strings.Split(desc, "·")
			if len(parts) > 1 {
				last := strings.TrimSpace(parts[len(parts)-1])
				if !strings.Contains(strings.ToLower(last), "followers") && !strings.Contains(strings.ToLower(last), "talking about") {
					profile.Bio = last
				}
			}
		} else if !strings.Contains(desc, "is on Facebook") && !strings.Contains(desc, "Join Facebook") {
			profile.Bio = desc
		}
	}

	return profile, nil
}

func fetchFacebookPublicHTML(username string) (string, error) {
	url := fmt.Sprintf("%s/%s", facebookBaseURL, username)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", facebookBotUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.Header.Set("User-Agent", facebookBotUA)
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch facebook data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("user %q not found on facebook", username)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

func GetFacebookProfile(username string) (*FacebookProfile, error) {
	htmlContent, err := fetchFacebookPublicHTML(username)
	if err != nil {
		return nil, err
	}
	return parseFacebookHTML(username, htmlContent)
}

func FacebookByUsername(username string) {
	profile, err := GetFacebookProfile(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("  - Facebook -")
	fmt.Println("Username:", profile.Username)
	lib.PrintIfSet("Name:", profile.FullName)
	lib.PrintIfSet("Bio:", profile.Bio)
	lib.PrintIfSet("Followers:", profile.Followers)
	lib.PrintIfSet("Talking About:", profile.TalkingAbout)
	lib.PrintIfSet("URL:", profile.URL)
}
