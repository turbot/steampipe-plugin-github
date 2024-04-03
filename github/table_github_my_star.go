package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubMyStar() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_star",
		Description: "GitHub stars owned by you. GitHub stars are repositories.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyStarredRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromValue(), Hydrate: starHydrateNameWithOwner},
			{Name: "starred_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().Transform(convertTimestamp), Hydrate: starHydrateStarredAt, Description: "The timestamp when the repository was starred."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: starHydrateUrl, Description: "URL of the repository."},
		},
	}
}

func tableGitHubMyStarredRepositoryList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			StarredRepositories struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					StarredAt models.NullableTime
					Node      struct {
						NameWithOwner string
						Url           string
					} `graphql:"node @include(if:$includeStarNode)"`
				} `graphql:"edges @include(if:$includeStarEdges)"`
			} `graphql:"starredRepositories(first: $pageSize, after: $cursor)"`
		}
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendStarColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_star", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_star", "api_error", err)
			return nil, err
		}

		for _, star := range query.Viewer.StarredRepositories.Edges {
			d.StreamListItem(ctx, myStar{star.StarredAt, star.Node.NameWithOwner, star.Node.Url})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.StarredRepositories.PageInfo.EndCursor)
	}

	return nil, nil
}

type myStar struct {
	StarredAt     models.NullableTime
	NameWithOwner string
	Url           string
}
