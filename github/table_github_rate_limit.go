package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
	client := connect(ctx, d)
	var rateLimits *github.RateLimits
	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		rateLimits, _, err = client.RateLimits(ctx)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, rateLimits)
	return nil, nil
}
