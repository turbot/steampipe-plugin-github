package github

import (
	"context"

	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func gitHubRateLimitGraphQLColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The username used to login.", Transform: transform.FromField("Viewer.Login")},
		{Name: "rate_limit", Type: proto.ColumnType_INT, Description: "The organization the member is associated with.", Transform: transform.FromField("RateLimit.Limit")},
		{Name: "cost", Type: proto.ColumnType_INT, Description: "The organization the member is associated with.", Transform: transform.FromField("RateLimit.Cost")},
		{Name: "remaining", Type: proto.ColumnType_INT, Description: "The organization the member is associated with.", Transform: transform.FromField("RateLimit.Remaining")},
		{Name: "reset_at", Type: proto.ColumnType_STRING, Description: "The organization the member is associated with.", Transform: transform.FromField("RateLimit.ResetAt")},
	}
}

var rateLimitQuery struct {
	Viewer struct {
		Login githubv4.String
	}
	RateLimit struct {
		Limit     githubv4.Int
		Cost      githubv4.Int
		Remaining githubv4.Int
		ResetAt   githubv4.String
	}
}

func tableGitHubRateLimitGraphQL() *plugin.Table {
	return &plugin.Table{
		Name:        "github_rate_limit_graphql",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubRateLimitGraphQLList,
		},
		Columns: gitHubRateLimitGraphQLColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubRateLimitGraphQLList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	err := client.Query(ctx, &rateLimitQuery, nil)
	if err != nil {
		plugin.Logger(ctx).Error("github_rate_limit_graphql", "api_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, rateLimitQuery)

	return nil, nil
}
