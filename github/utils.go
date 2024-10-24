package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"

	"github.com/google/go-github/v55/github"
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
	appId := os.Getenv("GITHUB_APP_ID")
	installationId := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	privateKeyPath := os.Getenv("GITHUB_APP_PEM_FILE")

	var githubAppId, githubInstallationId, githubPrivateKeyPath string

	if appId != "" {
		githubAppId = appId
	}
	if installationId != "" {
		githubInstallationId = installationId
	}
	if privateKeyPath != "" {
		githubPrivateKeyPath = privateKeyPath
	}

	// Get connection config for plugin
	githubConfig := GetConfig(d.Connection)
	if githubConfig.Token != nil {
		token = *githubConfig.Token
	}
	if githubConfig.BaseURL != nil {
		baseURL = *githubConfig.BaseURL
	}

	//// Github App authentication.
	if githubConfig.AppId != nil {
		githubAppId = *githubConfig.AppId
	}
	if githubConfig.InstallationId != nil {
		githubInstallationId = *githubConfig.InstallationId
	}
	if githubConfig.PrivateKey != nil {
		githubPrivateKeyPath = *githubConfig.PrivateKey
	}

	if token == "" && (githubAppId == "" || githubInstallationId == "" || githubPrivateKeyPath == "") {
		panic("'token' or 'app_id', 'installation_id' and 'private_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	// Return error for unsupported token by prefix
	// https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#githubs-token-formats
	if token != "" && !strings.HasPrefix(token, "ghs_") && !strings.HasPrefix(token, "ghp_") && !strings.HasPrefix(token, "gho_") && !strings.HasPrefix(token, "github_pat") {
		panic("Supported token formats are 'ghs_', 'gho_', 'ghp_', and 'github_pat'")
	}

	var client *github.Client

	// Authentication with Github access token
	if token != "" && (strings.HasPrefix(token, "ghp_") || strings.HasPrefix(token, "github_pat")){
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	// Authentication Using App Installation Access Token or OAuth Access token
	if token != "" && (strings.HasPrefix(token, "ghs_") || strings.HasPrefix(token, "gho_")) {
		client = github.NewClient(&http.Client{Transport: &oauth2Transport{
			Token: token,
		}})
	}

	// Authentication as Github APP Installation authentication
	if githubAppId != "" && githubInstallationId != "" && githubPrivateKeyPath != "" && token == "" {
		ghAppId, err := strconv.ParseInt(githubAppId, 10, 64)
		if err != nil {
			panic(err)
		}
		ghInstallationId, err := strconv.ParseInt(githubInstallationId, 10, 64)
		if err != nil {
			panic(err)
		}
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, ghAppId, ghInstallationId, githubPrivateKeyPath)
		if err != nil {
			panic("Error occurred in 'connect()' during GitHub App Installation client creation: " + err.Error())
		}

		client = github.NewClient(&http.Client{Transport: itr})
	}

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
		conn, err := github.NewClient(client.Client()).WithEnterpriseURLs(uv4.String(), "")
		if err != nil {
			panic(fmt.Sprintf("error creating GitHub client: %v", err))
		}
		conn.BaseURL = uv4
		client = conn
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client
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
	appId := os.Getenv("GITHUB_APP_ID")
	installationId := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	privateKeyPath := os.Getenv("GITHUB_APP_PEM_FILE")

	var githubAppId, githubInstallationId, githubPrivateKeyPath string

	if appId != "" {
		githubAppId = appId
	}
	if installationId != "" {
		githubInstallationId = installationId
	}
	if privateKeyPath != "" {
		githubPrivateKeyPath = privateKeyPath
	}

	// Get connection config for plugin
	githubConfig := GetConfig(d.Connection)
	if githubConfig.Token != nil {
		token = *githubConfig.Token
	}
	if githubConfig.BaseURL != nil {
		baseURL = *githubConfig.BaseURL
	}

	// Github App authentication.
	if githubConfig.AppId != nil {
		githubAppId = *githubConfig.AppId
	}
	if githubConfig.InstallationId != nil {
		githubInstallationId = *githubConfig.InstallationId
	}
	if githubConfig.PrivateKey != nil {
		githubPrivateKeyPath = *githubConfig.PrivateKey
	}

	// Return error for unsupported token by prefix
	// https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-authentication-to-github#githubs-token-formats
	if token != "" && !strings.HasPrefix(token, "ghs_") && !strings.HasPrefix(token, "ghp_") && !strings.HasPrefix(token, "gho_") && !strings.HasPrefix(token, "github_pat") {
		panic("Supported token formats are 'ghs_', 'gho_', 'ghp_', and 'github_pat'")
	}

	var client *githubv4.Client

	// Authentication Using App Installation Access Token or OAuth Access token
	if token != "" && (strings.HasPrefix(token, "ghs_") || strings.HasPrefix(token, "gho_")) {
		return githubv4.NewClient(&http.Client{Transport: &oauth2Transport{
			Token: token,
		}})
	}

	var transport *ghinstallation.Transport

	// Authentication with Github access token
	if token != "" && (strings.HasPrefix(token, "ghp_") || strings.HasPrefix(token, "github_pat")) {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = githubv4.NewClient(tc)
	}

	if token == "" && (githubAppId == "" || githubInstallationId == "" || githubPrivateKeyPath == "") {
		panic("'token' or 'app_id', 'installation_id' and 'private_key' must be set in the connection configuration. Edit your connection configuration file and then restart Steampipe")
	}

	// Authentication as Github APP Installation
	if githubAppId != "" && githubInstallationId != "" && githubPrivateKeyPath != "" && token == "" {
		ghAppId, err := strconv.ParseInt(githubAppId, 10, 64)
		if err != nil {
			panic(err)
		}
		ghInstallationId, err := strconv.ParseInt(githubInstallationId, 10, 64)
		if err != nil {
			panic(err)
		}
		itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, ghAppId, ghInstallationId, githubPrivateKeyPath)
		if err != nil {
			panic("Error occurred in 'connectV4()' during GitHub App Installation client creation" + err.Error())
		}
		transport = itr

		client = githubv4.NewClient(&http.Client{Transport: itr})
	}

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
		if token != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)
			client = githubv4.NewEnterpriseClient(uv4.String(), tc)
		} else {
			client = githubv4.NewEnterpriseClient(uv4.String(), &http.Client{Transport: transport})
		}
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client
}

// oauth2Transport is an http.RoundTripper that authenticates all requests
type oauth2Transport struct {
	Token string
}

func (t *oauth2Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header.Set("Authorization", "Bearer "+t.Token)
	return http.DefaultTransport.RoundTrip(clone)
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
