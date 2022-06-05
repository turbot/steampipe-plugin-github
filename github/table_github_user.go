package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func gitHubUserColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user."},
		{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of account."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user."},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar"},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user."},
		{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The user's gravatar ID"},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
		{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user."},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographic location of the user."},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user."},
		{Name: "hireable", Type: proto.ColumnType_BOOL, Description: "Whether the user currently hireable."},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user."},
		{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user."},
		{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user."},
		{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user."},
		{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was created.", Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was last updated.", Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
		// {Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP},
		{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator."},
		{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories."},
		{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories."},
		{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user."},
		{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user."},
		{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
		{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
		{Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user."},
	}
}

//// TABLE DEFINITION

func tableGitHubUser() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user",
		Description: "GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("login"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubUserGet,
		},
		Columns: gitHubUserColumns(),
	}
}

//// HYDRATE FUNCTRIONS

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

	type GetResponse struct {
		user *github.User
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Users.Get(ctx, login)
		return GetResponse{
			user: detail,
			resp: resp,
		}, err
	}
	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	user := getResp.user

	if user != nil {
		d.StreamListItem(ctx, user)
	}

	return nil, nil
}
