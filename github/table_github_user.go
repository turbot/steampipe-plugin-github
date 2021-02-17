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
		Description: "Github Users are user accounts in Github.",
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
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user."},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar"},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user."},
			{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The user's gravatar ID"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user.", Hydrate: tableGitHubUserGet},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for.", Hydrate: tableGitHubUserGet},
			{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user.", Hydrate: tableGitHubUserGet},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographic location of the user.", Hydrate: tableGitHubUserGet},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user.", Hydrate: tableGitHubUserGet},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Description: "Whether the user currently hireable.", Hydrate: tableGitHubUserGet},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user.", Hydrate: tableGitHubUserGet},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user.", Hydrate: tableGitHubUserGet},
			{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user.", Hydrate: tableGitHubUserGet},
			{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user.", Hydrate: tableGitHubUserGet},
			{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user.", Hydrate: tableGitHubUserGet},
			{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user.", Hydrate: tableGitHubUserGet},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was created.", Hydrate: tableGitHubUserGet, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was last updated.", Hydrate: tableGitHubUserGet, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
			//{Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubUserGet},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of account."},
			{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator."},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories.", Hydrate: tableGitHubUserGet},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories.", Hydrate: tableGitHubUserGet},
			{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user."},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user."},
			{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
			{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
			{Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user."},
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

	client := connect(ctx, d)

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
