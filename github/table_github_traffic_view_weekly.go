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

//// TABLE DEFINTION

func tableGitHubTrafficViewWeekly(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_traffic_view_weekly",
		Description: "Weekly traffic view over the last 14 days for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubTrafficViewWeeklyList,
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

func tableGitHubTrafficViewWeeklyList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.TrafficBreakdownOptions{Per: "week"}
	var result *github.TrafficViews

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		result, _, err = client.Repositories.ListTrafficViews(ctx, owner, repo, opts)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, i := range result.Views {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if plugin.IsCancelled(ctx) {
			return nil, nil
		}
	}
	return nil, nil
}
