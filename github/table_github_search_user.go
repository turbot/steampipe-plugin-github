package github

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubSearchUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_user",
		Description: "Find users via various criteria.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchUserList,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID of the user or organization."},
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user or organization."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the the user or organization."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the user or organization."},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar."},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user or organization."},
			{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user or organization."},
			{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was created."},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user or organization."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user or organization."},
			{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The event URL of the user or organization."},
			{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user or organization."},
			{Name: "followers_url", Type: proto.ColumnType_STRING, Description: "The URL to get list of followers."},
			{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user or organization."},
			{Name: "following_url", Type: proto.ColumnType_STRING, Description: "The URL to get list of users followed by the user or organization."},
			{Name: "gists_url", Type: proto.ColumnType_STRING, Description: "The URL get the gists of the user or organization."},
			{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The gravatar id of the user or organization."},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the user or organization is hireable."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user or organization."},
			{Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user or organization."},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The URL of the user or organization."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user or organization."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user or organization."},
			{Name: "organizations_url", Type: proto.ColumnType_STRING, Description: "The URL to get the organization details of the user or organization."},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories by the user or organization."},
			{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user or organization."},
			{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user or organization."},
			{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user or organization."},
			{Name: "received_events_url", Type: proto.ColumnType_STRING, Description: "The URL to get the received events of the user or organization."},
			{Name: "repos_url", Type: proto.ColumnType_STRING, Description: "The URL to get the repositories that the user or organization is part of."},
			{Name: "site_admin", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the user or organization is an administrator."},
			{Name: "starred_url", Type: proto.ColumnType_STRING, Description: "The URL to get the starred details of the user or organization."},
			{Name: "subscriptions_url", Type: proto.ColumnType_STRING, Description: "The URL to get subscription details of the user or organization."},
			{Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was suspended."},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories of the user or organization."},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user or organization."},
			{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether two-factor authentication is enabled for the user or organization."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was updated."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL to get information regarding the user or organization."},
			{Name: "permissions", Type: proto.ColumnType_JSON, Description: "The permission details."},
			{Name: "plan", Type: proto.ColumnType_JSON, Description: "The plan details."},
			{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubSearchUserList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchUserList")

	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.UsersSearchResult
		resp   *github.Response
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		result, resp, err := client.Search.Users(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchUserList", "error_Search.Users", err)
			return nil, err
		}

		return ListPageResponse{
			result: result,
			resp:   resp,
		}, nil
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			logger.Error("tableGitHubSearchUserList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		users := listResponse.result.Users
		resp := listResponse.resp

		for _, i := range users {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
