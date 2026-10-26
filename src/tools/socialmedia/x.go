package socialmedia

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"osintmaster/src/lib"
	"strconv"
	"strings"
	"time"
)

const (
	xBaseURL = "https://x.com"
	xUA      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

type xGraphQLResponse struct {
	Data struct {
		UserResultByScreenName struct {
			Result xUser
		} `json:"user_result_by_screen_name"`
	}
}

type xUser struct {
	Core               xUserCore
	Location           xUserLocation
	ProfileBio         xUserProfileBio    `json:"profile_bio"`
	RelationshipCounts xUserRelationships `json:"relationship_counts"`
	TweetCounts        xUserTweetCounts   `json:"tweet_counts"`
	Verification       xUserVerification
}

type xUserCore struct {
	CreatedAtMs int64 `json:"created_at_ms"`
	Name        string
	ScreenName  string `json:"screen_name"`
}

type xUserLocation struct {
	Location string
}

type xUserProfileBio struct {
	Description string
}

type xUserRelationships struct {
	Followers int
	Following int
}

type xUserTweetCounts struct {
	Tweets int
}

type xUserVerification struct {
	IsBlueVerified bool   `json:"is_blue_verified"`
	VerifiedType   string `json:"verified_type"`
}

type XProfile struct {
	Username  string
	Name      string
	Bio       string
	Location  string
	Followers string
	Following string
	Posts     string
	Joined    string
	Verified  bool
	URL       string
}

func parseXHTML(username, htmlContent string) (*XProfile, error) {
	_, after, ok := strings.Cut(htmlContent, `value:`)
	if !ok {
		return nil, fmt.Errorf("user %q not found on x", username)
	}

	var jsonStr string
	if err := json.NewDecoder(strings.NewReader(after)).Decode(&jsonStr); err != nil {
		return nil, fmt.Errorf("failed to extract user JSON for %q: %w", username, err)
	}

	var gql xGraphQLResponse
	if err := json.Unmarshal([]byte(jsonStr), &gql); err != nil {
		return nil, fmt.Errorf("failed to parse user JSON for %q: %w", username, err)
	}

	res := gql.Data.UserResultByScreenName.Result
	if res.Core.Name == "" && res.Core.ScreenName == "" {
		return nil, fmt.Errorf("user %q not found on x", username)
	}

	cleanUsername := strings.TrimPrefix(username, "@")
	if res.Core.ScreenName != "" {
		cleanUsername = res.Core.ScreenName
	}

	var joined string
	if res.Core.CreatedAtMs > 0 {
		joined = time.UnixMilli(res.Core.CreatedAtMs).UTC().Format("January 2006")
	}

	isVerified := res.Verification.IsBlueVerified || res.Verification.VerifiedType != ""

	return &XProfile{
		Username:  cleanUsername,
		Name:      res.Core.Name,
		Bio:       res.ProfileBio.Description,
		Location:  res.Location.Location,
		Followers: strconv.Itoa(res.RelationshipCounts.Followers),
		Following: strconv.Itoa(res.RelationshipCounts.Following),
		Posts:     strconv.Itoa(res.TweetCounts.Tweets),
		Joined:    joined,
		Verified:  isVerified,
		URL:       fmt.Sprintf("%s/%s", xBaseURL, cleanUsername),
	}, nil
}

func fetchXPublicHTML(username string) (string, error) {
	cleanUsername := strings.TrimPrefix(username, "@")
	url := fmt.Sprintf("%s/%s", xBaseURL, cleanUsername)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", xUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			req.Header.Set("User-Agent", xUA)
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.9")
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch x data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("user %q not found on x", username)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("x returned status %d for user %q", resp.StatusCode, username)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(bodyBytes), nil
}

func GetXProfile(username string) (*XProfile, error) {
	htmlContent, err := fetchXPublicHTML(username)
	if err != nil {
		return nil, err
	}
	return parseXHTML(username, htmlContent)
}

func XByUsername(username string) {
	profile, err := GetXProfile(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println("  - X -")
	fmt.Println("Username:", profile.Username)
	lib.PrintIfSet("Name:", profile.Name)
	lib.PrintIfSet("Bio:", profile.Bio)
	lib.PrintIfSet("Location:", profile.Location)
	fmt.Println("Verified:", profile.Verified)
	lib.PrintIfSet("Followers:", profile.Followers)
	lib.PrintIfSet("Following:", profile.Following)
	lib.PrintIfSet("Posts:", profile.Posts)
	lib.PrintIfSet("Joined:", profile.Joined)
	lib.PrintIfSet("URL:", profile.URL)
}
