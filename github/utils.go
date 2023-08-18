package github

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"

	"github.com/google/go-github/v48/github"
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
			uv4.Path = uv4.Path + "api/graphql"
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

func rateLimitLogString(table string, rateLimits *models.RateLimit) string {
	return fmt.Sprintf("Query for table %s - rate limit cost: %d (used: %d/%d) [Nodes: %d], resets at: %s", table, rateLimits.Cost, rateLimits.Used, rateLimits.Limit, rateLimits.NodeCount, rateLimits.ResetAt.String())
}

// transforms

func convertTimestamp(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	switch t := input.Value.(type) {
	case *github.Timestamp:
		return t.Format(time.RFC3339), nil
	case github.Timestamp:
		return t.Format(time.RFC3339), nil
	case githubv4.DateTime:
		return t.Format(time.RFC3339), nil
	case *githubv4.DateTime:
		return t.Format(time.RFC3339), nil
	case models.NullableTime:
		return t.Format(time.RFC3339), nil
	default:
		return nil, nil
	}
}

func defaultSearchColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query provided for the search."},
		{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
	}
}

func lastPathElement(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return getLastPathElement(types.SafeString(d.Value)), nil
}

func getLastPathElement(path string) string {
	if path == "" {
		return ""
	}

	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}
