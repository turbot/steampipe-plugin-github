package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTrafficCloneWeekly() *plugin.Table {
	return &plugin.Table{
		Name:        "github_traffic_clone_weekly",
		Description: "Weekly traffic clone stats over the last 14 days for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTrafficCloneWeeklyList,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(convertTimestamp), Description: "Date for the clone data."},
			{Name: "count", Type: proto.ColumnType_INT, Description: "Clone count for the week."},
			{Name: "uniques", Type: proto.ColumnType_INT, Description: "Unique clone count for the week."},
		}),
	}
}

func tableGitHubTrafficCloneWeeklyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.TrafficBreakdownOptions{Per: "week"}

	cloneStats, _, err := client.Repositories.ListTrafficClones(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	for _, i := range cloneStats.Clones {
		if i != nil {
			d.StreamListItem(ctx, i)
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}