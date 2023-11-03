package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubStargazer() *plugin.Table {
	return &plugin.Table{
		Name:        "github_stargazer",
		Description: "Stargazers are users who have starred the repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubStargazerList,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the stargazer."},
			{Name: "starred_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().Transform(convertTimestamp), Hydrate: strHydrateStarredAt, Description: "Time when the stargazer was created."},
			{Name: "user_login", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: strHydrateUserLogin, Description: "The login name of the user who starred the repository."},
			{Name: "user_detail", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: strHydrateUser, Description: "Details of the user who starred the repository."},
		},
	}
}

type Stargazer struct {
	StarredAt models.NullableTime `graphql:"starredAt @include(if:$includeStargazerStarredAt)" json:"starred_at"`
	Node      models.BasicUser    `graphql:"node @include(if:$includeStargazerNode)" json:"ndoe"`
}

func tableGitHubStargazerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Stargazers struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					Stargazer
				}
			} `graphql:"stargazers(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendStargazerColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_stargazer", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_stargazer", "api_error", err)
			return nil, err
		}

		for _, sg := range query.Repository.Stargazers.Edges {
			d.StreamListItem(ctx, Stargazer{sg.StarredAt, sg.Node})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Stargazers.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Stargazers.PageInfo.EndCursor)
	}

	return nil, nil
}
