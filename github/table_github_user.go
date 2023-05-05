package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubUser() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user",
		Description: "GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("login"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubUserGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("DatabaseId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("Id")},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar", Transform: transform.FromField("AvatarUrl")},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user.", Transform: transform.FromField("Url")},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
			{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user.", Transform: transform.FromField("WebsiteUrl")},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographic location of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user."},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Description: "Whether the user currently hireable.", Transform: transform.FromField("IsHireable")},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user."},
			{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user.", Transform: transform.FromField("PublicRepos.TotalCount")},
			{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user.", Transform: transform.FromField("PublicGists.TotalCount")},
			{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user.", Transform: transform.FromField("Followers.TotalCount")},
			{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user.", Transform: transform.FromField("Following.TotalCount")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was last updated."},
			{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator.", Transform: transform.FromField("IsSiteAdmin")},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories.", Transform: transform.FromField("PrivateRepos.TotalCount")},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories.", Transform: transform.FromField("PrivateRepos.TotalCount")},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user.", Transform: transform.FromField("Repositories.TotalDiskUsage")},
			{Name: "status_message", Type: proto.ColumnType_STRING, Description: "The status message set by the user.", Transform: transform.FromField("Status.Message")},
			{Name: "issues", Type: proto.ColumnType_INT, Description: "Count of issues authored by the user.", Transform: transform.FromField("Issues.TotalCount")},
			{Name: "organizations", Type: proto.ColumnType_INT, Description: "Count of organizations the user is a member of.", Transform: transform.FromField("Organizations.TotalCount")},
			{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The pronouns of the user."},
			// {Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user.", Transform: transform.FromField("PrivateGists.TotalCount")},
			// {Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
			// {Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
			// {Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user."},
		},
	}
}

var userQuery struct {
	User struct {
		Login           string
		DatabaseId      int
		Name            string
		Id              string
		AvatarUrl       string
		Url             string
		Company         string
		WebsiteUrl      string
		Location        string
		Email           string
		IsHireable      bool
		IsSiteAdmin     bool
		Bio             string
		TwitterUsername string
		Followers       struct {
			TotalCount int
		}
		Following struct {
			TotalCount int
		}
		PublicRepos struct {
			TotalCount int
		} `graphql:"publicRepos: repositories(privacy: PUBLIC)"`
		PrivateRepos struct {
			TotalCount int
		} `graphql:"privateRepos: repositories(privacy: PRIVATE)"`
		Repositories struct {
			TotalDiskUsage int
		}
		PublicGists struct {
			TotalCount int
		} `graphql:"publicGists: gists(privacy: PUBLIC)"`
		// PrivateGists struct {
		// 	TotalCount int
		// } `graphql:"privateGists: gists(privacy: SECRET)"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Status    struct {
			Message string
		}
		Issues struct {
			TotalCount int
		}
		Organizations struct {
			TotalCount int
		}
		Pronouns string
	} `graphql:"user(login: $login)"`
}

func tableGitHubUserGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string
	if h.Item != nil {
		item := h.Item.(*github.User)
		plugin.Logger(ctx).Trace("tableGitHubUserGet", item.String())
		login = *item.Login
	} else {
		login = d.EqualsQuals["login"].GetStringValue()
	}

	if login == "" {
		return nil, nil
	}

	client := connectV4(ctx, d)

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	err := client.Query(ctx, &userQuery, variables)
	if err != nil {
		plugin.Logger(ctx).Error("github_user", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, userQuery.User)

	return nil, nil
}
