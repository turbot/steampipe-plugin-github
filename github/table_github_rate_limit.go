package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableGitHubRateLimit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_rate_limit",
		Description: "Rate limit of github.",
		List: &plugin.ListConfig{
			Hydrate: listGitHubRateLimit,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "core_limit", Type: proto.ColumnType_INT, Transform: transform.FromField("Core.Limit"), Description: "The number of requests per hour the client is currently limited to."},
			{Name: "core_remaining", Type: proto.ColumnType_INT, Transform: transform.FromField("Core.Remaining"), Description: "The number of remaining requests the client can make this hour."},
			{Name: "core_reset", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Core.Reset").Transform(convertTimestamp), Description: "The time at which the current rate limit will reset."},
			{Name: "search_limit", Type: proto.ColumnType_INT, Transform: transform.FromField("Search.Limit"), Description: "The number of requests per hour the client is currently limited to."},
			{Name: "search_remaining", Type: proto.ColumnType_INT, Transform: transform.FromField("Search.Remaining"), Description: "The number of remaining requests the client can make this hour."},
			{Name: "search_reset", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Search.Reset").Transform(convertTimestamp), Description: "The time at which the current rate limit will reset."},
		},
	}
}

func listGitHubRateLimit(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client) (interface{}, error) {
		detail, _, err := client.RateLimits(ctx)
		return detail, err
	}

	return streamGitHubListOrItem(ctx, d, h, getDetails)
}
