package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubSearchPullRequestColumns() []*plugin.Column {
	return append(defaultSearchColumns(), sharedPullRequestColumns()...)
}

func tableGitHubSearchPullRequest(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_pull_request",
		Description: "Find pull requests by state and keyword.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchPullRequestList,
		},
		Columns: gitHubSearchPullRequestColumns(),
	}
}

func tableGitHubSearchPullRequestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	input += " is:pr"

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			PageInfo models.PageInfo
			Edges    []struct {
				TextMatches []models.TextMatch
				Node        struct {
					models.BasicPullRequest `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(type: ISSUE, first: $pageSize, after: $cursor, query: $query)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
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
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_pull_request", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_pull_request", "api_error", err)
			return nil, err
		}

		for _, pr := range query.Search.Edges {
			d.StreamListItem(ctx, pr)

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
