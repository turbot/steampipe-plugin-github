package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINTION

func tableGitHubTrafficViewWeekly(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_traffic_view_weekly",
		Description: "Weekly traffic view over the last 14 days for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTrafficViewWeeklyList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the branch."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(convertTimestamp), Description: "Date for the view data."},
			{Name: "count", Type: proto.ColumnType_INT, Description: "View count for the week."},
			{Name: "uniques", Type: proto.ColumnType_INT, Description: "Unique viewer count for the week."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubTrafficViewWeeklyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.TrafficBreakdownOptions{Per: "week"}

	type ListResponse struct {
		trafficViews *github.TrafficViews
		resp         *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		trafficViews, resp, err := client.Repositories.ListTrafficViews(ctx, owner, repo, opts)
		return ListResponse{
			trafficViews: trafficViews,
			resp:         resp,
		}, err
	}
	listResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}

	result := listResponse.(ListResponse)
	trafficViews := result.trafficViews

	for _, i := range trafficViews.Views {
		if i != nil {
			d.StreamListItem(ctx, i)
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
