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
	"strconv"
	"strings"
	"time"
)

const (
	instagramEpoch    = 1314220021000 // Instagram snowflake epoch in ms (Aug 24, 2011 UTC)
	instagramAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	instagramURL      = "https://www.instagram.com"
	instagramBotUA    = "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"
)

var (
	metaTagRe     = regexp.MustCompile(`(?i)<meta\s+[^>]+>`)
	metaPropRe    = regexp.MustCompile(`(?i)(?:property|name)=["']([^"']+)["']`)
	metaContentRe = regexp.MustCompile(`(?i)content=["']([^"']*)["']`)

	countsRe = regexp.MustCompile(`([0-9.,KMkm]+)\s*Followers,\s*([0-9.,KMkm]+)\s*Following,\s*([0-9.,KMkm]+)\s*Posts`)
	bioRe    = regexp.MustCompile(`(?s)on Instagram:\s*["“](.*)["”]?$`)
)

type InstagramProfile struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	AvatarURL    string `json:"profile_pic_url"`
	Bio          string `json:"biography"`
	Followers    string `json:"follower_count"`
	Following    string `json:"following_count"`
	Posts        string
	IsPrivate    bool `json:"is_private"`
	IsVerified   bool `json:"is_verified"`
	LastActivity *Activity
}

type instagramTimelineJSON struct {
	Edges []struct {
		Node struct {
			Typename  string `json:"__typename"`
			PK        string `json:"pk"`
			MediaType int    `json:"media_type"`
		} `json:"node"`
	} `json:"edges"`
}

func pkToShortcode(pk int64) string {
	if pk <= 0 {
		return ""
	}
	var res []byte
	for pk > 0 {
		res = append([]byte{instagramAlphabet[pk%64]}, res...)
		pk /= 64
	}
	return string(res)
}

func parseMetaTags(htmlContent string) map[string]string {
	meta := make(map[string]string)
	for _, tag := range metaTagRe.FindAllString(htmlContent, -1) {
		nameMatch := metaPropRe.FindStringSubmatch(tag)
		contentMatch := metaContentRe.FindStringSubmatch(tag)
		if len(nameMatch) > 1 && len(contentMatch) > 1 {
			meta[strings.ToLower(nameMatch[1])] = html.UnescapeString(contentMatch[1])
		}
	}
	return meta
}

func extractUserJSON(htmlContent string) *InstagramProfile {
	for _, key := range []string{`"xig_user_by_igid_v2":`, `"user":`} {
		_, after, ok := strings.Cut(htmlContent, key)
		if !ok {
			continue
		}
		decoder := json.NewDecoder(strings.NewReader(after))
		var u InstagramProfile
		if decoder.Decode(&u) == nil && u.Username != "" {
			return &u
		}
	}
	return nil
}

func extractLastActivity(htmlContent string) *Activity {
	key := `"polaris_timeline_connection":`
	_, after, ok := strings.Cut(htmlContent, key)
	if !ok {
		return nil
	}
	decoder := json.NewDecoder(strings.NewReader(after))
	var timeline instagramTimelineJSON
	if err := decoder.Decode(&timeline); err != nil || len(timeline.Edges) == 0 {
		return nil
	}

	node := timeline.Edges[0].Node
	pk, err := strconv.ParseInt(node.PK, 10, 64)
	if err != nil || pk <= 0 {
		return nil
	}

	activityType := "Post"
	if strings.Contains(strings.ToLower(node.Typename), "video") || node.MediaType == 2 {
		activityType = "Video"
	}

	return &Activity{
		Type: activityType,
		Url:  fmt.Sprintf("%s/p/%s", instagramURL, pkToShortcode(pk)),
		Time: time.UnixMilli((pk >> 23) + instagramEpoch),
	}
}

func parseInstagramHTML(username, htmlContent string) (*InstagramProfile, error) {
	profile := extractUserJSON(htmlContent)
	if profile == nil || profile.Username == "" {
		return nil, fmt.Errorf("user %q not found on instagram", username)
	}

	meta := parseMetaTags(htmlContent)
	profile.AvatarURL = meta["og:image"]
	ogTitle := meta["og:title"]
	ogDesc := meta["og:description"]

	if parts := strings.Split(ogTitle, "(@"+profile.Username); len(parts) > 1 {
		profile.FullName = strings.TrimSpace(parts[0])
	}
	if cm := countsRe.FindStringSubmatch(ogDesc); len(cm) > 3 {
		profile.Followers = cm[1]
		profile.Following = cm[2]
		profile.Posts = cm[3]
	}
	if desc, ok := meta["description"]; ok {
		if bm := bioRe.FindStringSubmatch(desc); len(bm) > 1 {
			profile.Bio = strings.TrimSpace(strings.TrimRight(bm[1], "\"”"))
		}
	}

	profile.LastActivity = extractLastActivity(htmlContent)
	return profile, nil
}

func instagramFetchHTML(username string) (string, error) {
	url := fmt.Sprintf("%s/%s/", instagramURL, username)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", instagramBotUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch instagram data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("user %q not found on instagram", username)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("instagram returned status %d for user %q", resp.StatusCode, username)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(bodyBytes), nil
}

func GetInstagramProfileAndActivity(username string) (*InstagramProfile, error) {
	htmlContent, err := instagramFetchHTML(username)
	if err != nil {
		return nil, err
	}
	return parseInstagramHTML(username, htmlContent)
}

func InstagramByUsername(username string) {
	profile, err := GetInstagramProfileAndActivity(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("  - Instagram -")
	fmt.Println("Username:", profile.Username)
	lib.PrintIfSet("Full Name:", profile.FullName)
	lib.PrintIfSet("Bio:", profile.Bio)
	fmt.Println("Verified", profile.IsVerified)
	fmt.Println("Private:", profile.IsPrivate)
	fmt.Println("Followers:", profile.Followers)
	fmt.Println("Following:", profile.Following)
	fmt.Println("Posts:", profile.Posts)

	if profile.LastActivity == nil {
		fmt.Println("Last Activity: N/A")
	} else {
		fmt.Println("Last Activity:")
		fmt.Println(" Type:", profile.LastActivity.Type)
		fmt.Println(" Time:", profile.LastActivity.Time)
		fmt.Println(" URL:", profile.LastActivity.Url)

	}
}
