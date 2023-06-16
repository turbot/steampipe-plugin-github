package github

import (
	"context"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

func tableGitHubRateLimitGraphQL() *plugin.Table {
	return &plugin.Table{
		Name:        "github_rate_limit_graphql",
		Description: "Rate limit information for GitHub GraphQL API endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listGitHubRateLimitGraphQL,
		},
		Columns: []*plugin.Column{
			{Name: "cost", Type: proto.ColumnType_INT, Description: "Number of points used to return this query."},
			{Name: "used", Type: proto.ColumnType_INT, Description: "Number of points used from current allocation."},
			{Name: "remaining", Type: proto.ColumnType_INT, Description: "Number of points remaining in current allocation."},
			{Name: "limit", Type: proto.ColumnType_INT, Description: "Maximum number of points used that can be used in current allocation."},
			{Name: "reset_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the allocation resets.", Transform: transform.FromField("ResetAt").NullIfZero()},
			{Name: "node_count", Type: proto.ColumnType_INT, Description: "Number of nodes returned by this query."},
		},
	}
}

func listGitHubRateLimitGraphQL(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit models.RateLimit
	}

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, nil)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_rate_limit_graphql", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_rate_limit_graphql", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, query.RateLimit)

	return nil, nil
}
