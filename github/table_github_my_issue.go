package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubMyIssueColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Repo.NameWithOwner", "Node.Repo.NameWithOwner"), Description: "The full name of the repository (login/repo-name)."},
	}

	return append(tableCols, sharedIssueColumns()...)
}

func tableGitHubMyIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_issue",
		Description: "GitHub Issues owned by you. GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyIssueList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "state", Require: plugin.Optional},
				{Name: "updated_at", Require: plugin.Optional, Operators: []string{">", ">="}},
			},
		},
		Columns: commonColumns(gitHubMyIssueColumns()),
	}
}

func tableGitHubMyIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var filters githubv4.IssueFilters

	quals := d.EqualsQuals
	if quals["state"] != nil {
		state := quals["state"].GetStringValue()
		switch state {
		case "OPEN":
			filters.States = &[]githubv4.IssueState{githubv4.IssueStateOpen}
		case "CLOSED":
			filters.States = &[]githubv4.IssueState{githubv4.IssueStateClosed}
		default:
			plugin.Logger(ctx).Error("github_my_issue", "invalid filter", "state", state)
			return nil, fmt.Errorf("invalid value for 'state' can only filter for 'OPEN' or 'CLOSED' - you attempted to filter for '%s'", state)
		}
	} else {
		filters.States = &[]githubv4.IssueState{githubv4.IssueStateOpen, githubv4.IssueStateClosed}
	}

	if d.Quals["updated_at"] != nil {
		for _, q := range d.Quals["updated_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				filters.Since = githubv4.NewDateTime(githubv4.DateTime{Time: afterTime})
			case ">=":
				filters.Since = githubv4.NewDateTime(githubv4.DateTime{Time: givenTime})
			}
		}
	}

	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			Issues struct {
				TotalCount int
				PageInfo   models.PageInfo
				Nodes      []models.Issue
			} `graphql:"issues(first: $pageSize, after: $cursor, filterBy: $filters)"`
		}
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"filters":  filters,
	}
	appendIssueColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_issue", &query.RateLimit))
		if err != nil && len(query.Viewer.Issues.Nodes)==0{
			plugin.Logger(ctx).Error("github_my_issue", "api_error", err)
			return nil, err
		}

		for _, issue := range query.Viewer.Issues.Nodes {
			d.StreamListItem(ctx, issue)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.Issues.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.Issues.PageInfo.EndCursor)
	}

	return nil, nil
}
