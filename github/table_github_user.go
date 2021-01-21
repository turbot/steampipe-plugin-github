package github

import (
	"context"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubUser() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user",
		Description: "Github User",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("login"),
			Hydrate:    tableGitHubUserGet,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("login"),
			Hydrate:    tableGitHubUserGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "login", Type: proto.ColumnType_STRING},
			{Name: "id", Type: proto.ColumnType_INT},
			{Name: "node_id", Type: proto.ColumnType_STRING},
			{Name: "avatar_url", Type: proto.ColumnType_STRING},
			{Name: "html_url", Type: proto.ColumnType_STRING},
			{Name: "gravatar_id", Type: proto.ColumnType_STRING},
			{Name: "name", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "company", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "blog", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "location", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "email", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubUserGet},
			{Name: "bio", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Hydrate: tableGitHubUserGet},
			{Name: "public_repos", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "public_gists", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "followers", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "following", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubUserGet, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubUserGet, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
			//{Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubUserGet},
			{Name: "type", Type: proto.ColumnType_STRING},
			{Name: "site_admin", Type: proto.ColumnType_BOOL},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Hydrate: tableGitHubUserGet},
			{Name: "private_gists", Type: proto.ColumnType_INT},
			{Name: "disk_usage", Type: proto.ColumnType_INT},
			{Name: "collaborators", Type: proto.ColumnType_INT},
			{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL},
			{Name: "ldap_dn", Type: proto.ColumnType_STRING},
		},
	}
}

//// hydrate functions ////

// tableGitHubUserGet is both the Get and List hydrate function
// Listing all users is not terribly useful, so we require a 'login' qual and essentially always
// do a 'get':  from GitHub API docs: https://developer.github.com/v3/users/#list-users:
//     	Lists all users, in the order that they signed up on GitHub. This list includes personal user
//		accounts and organization accounts.
//     	Note: Pagination is powered exclusively by the since parameter. Use the Link header to get
//		the URL for the next page of users.
func tableGitHubUserGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var login string

	if h.Item != nil {
		item := h.Item.(*github.User)
		logger.Trace("tableGitHubUserGet", item.String())
		login = *item.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	client := connect(ctx, d.ConnectionManager)

	var detail *github.User
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		detail, resp, err = client.Users.Get(ctx, login)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}
