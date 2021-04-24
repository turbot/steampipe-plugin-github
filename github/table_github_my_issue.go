package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubMyIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_issue",
		Description: "GitHub Issues owned by you. GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyIssueList,
		},
		Columns: gitHubIssueColumns(),
	}
}

//// HYDRATE FUNCTIONS

func tableGitHubMyIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// TO DO - get state and other filters from the quals
	opt := &github.IssueListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	client := connect(ctx, d)

	for {
		var issues []*github.Issue
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			issues, resp, err = client.Issues.List(ctx, true, opt)

			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range issues {
			// Only issues, not PRs (those are in the pull_request table...)
			if !i.IsPullRequest() {
				d.StreamListItem(ctx, i)
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
