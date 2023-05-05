package github

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/sethvargo/go-retry"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Create Rest API (v3) client
func connect(ctx context.Context, d *plugin.QueryData) *github.Client {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "github_v3"
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
		uv4, err := url.Parse(baseURL)
		if err != nil {
			panic(fmt.Sprintf("github.base_url is invalid: %s", baseURL))
		}

		if uv4.String() != "https://api.github.com/" {
			uv4.Path = uv4.Path + "api/v3/"
		}

		// The upload URL is not set as it's not currently required
		conn, err = github.NewEnterpriseClient(uv4.String(), "", tc)
		if err != nil {
			panic(fmt.Sprintf("error creating GitHub client: %v", err))
		}

		conn.BaseURL = uv4
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn
}

// Create GraphQL API (v4) client
func connectV4(ctx context.Context, d *plugin.QueryData) *githubv4.Client {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "github_v4"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*githubv4.Client)
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
	conn := githubv4.NewClient(tc)

	// If the base URL was provided then set it on the client. Used for
	// enterprise installs.
	if baseURL != "" {
		uv4, err := url.Parse(baseURL)
		if err != nil {
			panic(fmt.Sprintf("github.base_url is invalid: %s", baseURL))
		}

		if uv4.String() != "https://api.github.com/" {
			uv4.Path = uv4.Path + "api/v4/"
		}

		conn = githubv4.NewEnterpriseClient(uv4.String(), tc)
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

func adjustPageSize(pageSize int, limit *int64) int {
	if limit != nil && *limit < int64(pageSize) {
		return int(*limit)
	}
	return pageSize
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

func retryHydrate(ctx context.Context, d *plugin.QueryData, hydrateData *plugin.HydrateData, hydrateFunc plugin.HydrateFunc) (interface{}, error) {

	// Retry configs
	maxRetries := 10
	interval := time.Duration(1)

	// Create the backoff based on the given mode
	// Use exponential instead of fibonacci due to GitHub's aggressive throttling
	backoff := retry.NewExponential(interval * time.Second)

	// Ensure the maximum value is 30s. In this scenario, the sleep values would be
	// 1s, 2s, 4s, 16s, 30s, 30s...
	backoff = retry.WithCappedDuration(30*time.Second, backoff)

	var hydrateResult interface{}

	err := retry.Do(ctx, retry.WithMaxRetries(uint64(maxRetries), backoff), func(ctx context.Context) error {
		var err error
		hydrateResult, err = hydrateFunc(ctx, d, hydrateData)
		if err != nil {
			if shouldRetryError(ctx, err) {
				err = retry.RetryableError(err)
			}
		}
		return err
	})

	return hydrateResult, err
}
