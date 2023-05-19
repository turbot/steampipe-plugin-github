package github

import (
	"context"
	"encoding/json"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubSearchPullRequestColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number", "Node.Number"), Description: "The number of the pull request."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID of the pull request."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the pull request."},
	}

	return append(defaultSearchColumns(), tableCols...)
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
					models.PullRequest `graphql:"... on PullRequest"`
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

	qj, _ := json.Marshal(query)
	plugin.Logger(ctx).Debug(string(qj))

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
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
