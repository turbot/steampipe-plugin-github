package github

import (
	"context"

	"github.com/google/go-github/v33/github"
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
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})
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
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
