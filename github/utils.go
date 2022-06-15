package github

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

// create service client
func connect(ctx context.Context, d *plugin.QueryData) *github.Client {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "github"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*github.Client)
	}

	token := os.Getenv("GITHUB_TOKEN")
	baseURL := os.Getenv("GITHUB_BASE_URL")

	// Get connection config for plugin
	githubConfig := GetConfig(d.Connection)
	if githubConfig.Token != nil {
		token = *githubConfig.Token
	}
	if githubConfig.BaseURL != nil {
		baseURL = *githubConfig.BaseURL
	}

	if token == "" {
		panic("'token' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	conn := github.NewClient(tc)

	// If the base URL was provided then set it on the client. Used for
	// enterprise installs.
	if baseURL != "" {
		uv3, err := url.Parse(baseURL)
		if err != nil {
			panic(fmt.Sprintf("github.base_url is invalid: %s", baseURL))
		}

		if uv3.String() != "https://api.github.com/" {
			uv3.Path = uv3.Path + "api/v3/"
		}

		// The upload URL is not set as it's not currently required
		conn, err = github.NewEnterpriseClient(uv3.String(), "", tc)
		if err != nil {
			panic(fmt.Sprintf("error creating GitHub client: %v", err))
		}

		conn.BaseURL = uv3
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn
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
	case github.Timestamp:
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

	var userType []*github.User

	// Check type of the transform values otherwise it is throwing error while type casting the interface to []*github.User type
	if reflect.TypeOf(input.Value) != reflect.TypeOf(userType) {
		return nil, nil
	}

	users := input.Value.([]*github.User)

	if users == nil {
		return user_logins, nil
	}

	for _, u := range users {
		user_logins = append(user_logins, *u.Login)
	}
	return user_logins, nil
}

func gitHubSearchRepositoryColumns(columns []*plugin.Column) []*plugin.Column {
	return append(gitHubRepositoryColumns(), columns...)
}
