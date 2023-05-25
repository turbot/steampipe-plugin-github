package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubSearchIssueColumns() []*plugin.Column {
	return append(defaultSearchColumns(), gitHubMyIssueColumns()...)
}

func tableGitHubSearchIssue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_issue",
		Description: "Find issues by state and keyword.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchIssueList,
		},
		Columns: gitHubSearchIssueColumns(),
	}
}

func tableGitHubSearchIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	input += " is:issue"

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			PageInfo models.PageInfo
			Edges    []struct {
				TextMatches []models.TextMatch
				Node        struct {
					models.Issue `graphql:"... on Issue"`
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
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_issue", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_issue", "api_error", err)
			return nil, err
		}

		for _, issue := range query.Search.Edges {
			d.StreamListItem(ctx, issue)

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
