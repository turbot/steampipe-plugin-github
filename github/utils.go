package github

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// create service client
func connect(ctx context.Context, d *plugin.QueryData) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")

	// Get connection config for plugin
	githubConfig := GetConfig(d.Connection)
	if &githubConfig != nil {
		if githubConfig.Token != nil {
			token = *githubConfig.Token
		}
	}

	if token == "" {
		panic("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

//// HELPER FUNCTIONS

func parseRepoFullName(fullName string) (string, string) {
	owner := ""
	repo := ""
	s := strings.Split(fullName, "/")
	owner = s[0]
	if len(s) > 1 {
		repo = s[1]
	}
	return owner, repo
}

// transforms

func convertTimestamp(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	switch t := input.Value.(type) {
	case *github.Timestamp:
		return t.Format(time.RFC3339), nil
	default:
		return nil, nil
	}
}

func filterUserLogins(_ context.Context, input *transform.TransformData) (interface{}, error) {
	user_logins := make([]string, 0)
	if input.Value == nil {
		return user_logins, nil
	}
	users := input.Value.([]*github.User)

	for _, u := range users {
		user_logins = append(user_logins, *u.Login)
	}
	return user_logins, nil
}
