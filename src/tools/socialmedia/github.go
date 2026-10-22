package socialmedia

import (
	"context"
	"fmt"
	"os"
	"osintmaster/src/lib"
	"time"

	"github.com/google/go-github/v90/github"
)

func GetLastGithubActivity(username string) (*Activity, error) {
	client, err := github.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	opts := &github.ListOptions{PerPage: 1}
	events, _, err := client.Activity.ListEventsPerformedByUser(context.Background(), username, false, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user events: %w", err)
	}
	if len(events) == 0 {
		return nil, nil
	}

	lastEvent := events[0]
	repoName := lastEvent.GetRepo().GetName()
	activityUrl := fmt.Sprintf("https://github.com/%s", repoName)
	payload, err := lastEvent.ParsePayload()
	if err == nil {
		switch p := payload.(type) {
		case *github.CommitCommentEvent:
			activityUrl = p.GetComment().GetHTMLURL()
		case *github.PullRequestEvent:
			activityUrl = p.GetPullRequest().GetHTMLURL()
		case *github.PullRequestReviewEvent:
			activityUrl = p.GetReview().GetHTMLURL()
		case *github.PullRequestReviewCommentEvent:
			activityUrl = p.GetComment().GetHTMLURL()
		case *github.PullRequestReviewThreadEvent:
			activityUrl = p.GetPullRequest().GetHTMLURL()
		case *github.IssuesEvent:
			activityUrl = p.GetIssue().GetHTMLURL()
		case *github.IssueCommentEvent:
			activityUrl = p.GetComment().GetHTMLURL()
		case *github.ReleaseEvent:
			activityUrl = p.GetRelease().GetHTMLURL()
		case *github.ForkEvent:
			activityUrl = p.GetForkee().GetHTMLURL()
		case *github.MemberEvent:
			activityUrl = p.GetMember().GetHTMLURL()
		case *github.DiscussionEvent:
			activityUrl = p.GetDiscussion().GetHTMLURL()
		case *github.DiscussionCommentEvent:
			activityUrl = p.GetComment().GetHTMLURL()
		case *github.PushEvent:
			if len(p.Commits) > 0 {
				activityUrl = fmt.Sprintf("https://github.com/%s/commit/%s", repoName, p.Commits[0].GetSHA())
			}
		case *github.CreateEvent:
			if p.GetRefType() != "repository" {
				activityUrl = fmt.Sprintf("https://github.com/%s/tree/%s", repoName, p.GetRef())
			}
		case *github.GollumEvent:
			if len(p.Pages) > 0 {
				activityUrl = p.Pages[0].GetHTMLURL()
			}
		}
	}

	return &Activity{
		Url:  activityUrl,
		Type: lastEvent.GetType(),
		Time: lastEvent.GetCreatedAt().Time,
	}, nil
}

type Profile struct {
	Username    string
	Name        string
	Email       string
	Bio         string
	Company     string
	Location    string
	Blog        string
	AvatarURL   string
	Twitter     string
	Followers   int
	Following   int
	PublicRepos int
	CreatedAt   time.Time
}

func GetGithubProfile(username string) (Profile, error) {
	client, err := github.NewClient()
	if err != nil {
		return Profile{}, fmt.Errorf("failed to create github client: %w", err)
	}

	user, _, err := client.Users.Get(context.Background(), username)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to fetch user profile for %q: %w", username, err)
	}

	return Profile{
		Username:    user.GetLogin(),
		Name:        user.GetName(),
		Email:       user.GetEmail(),
		Bio:         user.GetBio(),
		Company:     user.GetCompany(),
		Location:    user.GetLocation(),
		Blog:        user.GetBlog(),
		AvatarURL:   user.GetAvatarURL(),
		Twitter:     user.GetTwitterUsername(),
		Followers:   user.GetFollowers(),
		Following:   user.GetFollowing(),
		PublicRepos: user.GetPublicRepos(),
		CreatedAt:   user.GetCreatedAt().Time,
	}, nil
}

func GithubByUsername(username string) {
	profile, err := GetGithubProfile(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println("  - Github -")
	fmt.Println("Username:", profile.Username)
	lib.PrintIfSet("Name:", profile.Name)
	lib.PrintIfSet("Email:", profile.Email)
	lib.PrintIfSet("Bio:", profile.Bio)
	lib.PrintIfSet("Company:", profile.Company)
	lib.PrintIfSet("Location:", profile.Location)
	lib.PrintIfSet("Blog:", profile.Blog)
	lib.PrintIfSet("Twitter:", profile.Twitter)
	fmt.Println("Followers:", profile.Followers)
	fmt.Println("Following:", profile.Following)
	fmt.Println("Public Repos:", profile.PublicRepos)
	fmt.Println("Created At:", profile.CreatedAt.Format("Jan 2006"))

	activity, err := GetLastGithubActivity(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if activity == nil {
		fmt.Println("Last Activity: N/A")
	} else {
		fmt.Println("Last Activity:")
		fmt.Println(" Type:", activity.Type)
		fmt.Println(" Time:", activity.Time)
		fmt.Println(" URL:", activity.Url)

	}
}
