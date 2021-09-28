package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableGitHubMyIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_issue",
		Description: "GitHub Issues owned by you. GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyIssueList,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "state",
					Require: plugin.Optional,
				},
				{
					Name:      "created_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">="},
				},
			},
		},
		Columns: gitHubIssueColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubMyIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	opt := &github.IssueListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	// Additional filters
	if d.KeyColumnQuals["state"] != nil {
		opt.State = d.KeyColumnQuals["state"].GetStringValue()
	}

	if d.Quals["created_at"] != nil {
		for _, q := range d.Quals["created_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				opt.Since = afterTime
			case ">=":
				opt.Since = givenTime
			}
		}
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	client := connect(ctx, d)

	type ListPageResponse struct {
		issues []*github.Issue
		resp   *github.Response
	}
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		issues, resp, err := client.Issues.List(ctx, true, opt)
		return ListPageResponse{
			issues: issues,
			resp:   resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		issues := listResponse.issues
		resp := listResponse.resp

		for _, i := range issues {
			// Only issues, not PRs (those are in the pull_request table...)
			if !i.IsPullRequest() {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
