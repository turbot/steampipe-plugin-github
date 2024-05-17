package github

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRateLimitGraphQL() *plugin.Table {
	return &plugin.Table{
		Name:        "github_rate_limit_graphql",
		Description: "Rate limit information for GitHub GraphQL API endpoint.",
		List: &plugin.ListConfig{
			Hydrate: listGitHubRateLimitGraphQL,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "cost", Type: proto.ColumnType_INT, Description: "Number of points used to return this query.", Transform: transform.FromValue(), Hydrate: rateLimitHydrateCost},
			{Name: "used", Type: proto.ColumnType_INT, Description: "Number of points used from current allocation.", Transform: transform.FromValue(), Hydrate: rateLimitHydrateUsed},
			{Name: "remaining", Type: proto.ColumnType_INT, Description: "Number of points remaining in current allocation.", Transform: transform.FromValue(), Hydrate: rateLimitHydrateRemaining},
			{Name: "limit", Type: proto.ColumnType_INT, Description: "Maximum number of points used that can be used in current allocation.", Transform: transform.FromValue(), Hydrate: rateLimitHydrateLimit},
			{Name: "reset_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the allocation resets.", Transform: transform.FromValue().NullIfZero(), Hydrate: rateLimitHydrateResetAt},
			{Name: "node_count", Type: proto.ColumnType_INT, Description: "Number of nodes returned by this query.", Transform: transform.FromValue(), Hydrate: rateLimitHydrateNodeCount},
		}),
	}
}

func listGitHubRateLimitGraphQL(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit models.BaseRateLimit
	}

	variables := map[string]interface{}{}
	appendRateLimitColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, variables)
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
