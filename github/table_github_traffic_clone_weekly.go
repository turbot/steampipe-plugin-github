package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINTION

func tableGitHubTrafficCloneWeekly(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_traffic_clone_weekly",
		Description: "Weekly traffic clone over the last 14 days for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTrafficCloneWeeklyList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(convertTimestamp), Description: "Date for the clone data."},
			{Name: "count", Type: proto.ColumnType_INT, Description: "Clone count for the week."},
			{Name: "uniques", Type: proto.ColumnType_INT, Description: "Unique clones count for the week."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubTrafficCloneWeeklyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.TrafficBreakdownOptions{Per: "week"}

	type ListResponse struct {
		trafficClones *github.TrafficClones
		resp          *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		trafficClones, resp, err := client.Repositories.ListTrafficClones(ctx, owner, repo, opts)
		return ListResponse{
			trafficClones: trafficClones,
			resp:          resp,
		}, err
	}
	listResponse, err := retryHydrate(ctx, d, h, listPage)

	if err != nil {
		return nil, err
	}

	result := listResponse.(ListResponse)
	trafficClones := result.trafficClones

	for _, i := range trafficClones.Clones {
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
