package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubSearchRepositoryColumns() []*plugin.Column {
	return append(defaultSearchColumns(), sharedRepositoryColumns()...)
}

func tableGitHubSearchRepository(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_repository",
		Description: "Find repositories via various criteria.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchRepositoryList,
		},
		Columns: gitHubSearchRepositoryColumns(),
	}
}

func tableGitHubSearchRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			RepositoryCount int
			PageInfo        models.PageInfo
			Edges           []struct {
				TextMatches []models.TextMatch
				Node        struct {
					models.Repository `graphql:"... on Repository"`
				}
			}
		} `graphql:"search(type: REPOSITORY, first: $pageSize, after: $cursor, query: $query)"`
	}

	pageSize := adjustPageSize(75, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"query":    githubv4.String(input),
	}

	client := connectV4(ctx, d)
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := retryHydrate(ctx, d, h, listPage)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_repository", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_repository", "api_error", err)
			return nil, err
		}

		for _, repo := range query.Search.Edges {
			d.StreamListItem(ctx, repo)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	return nil, nil
}
