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

func tableGitHubBranch(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubBranchList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repositoryFullNameQual, Transform: transform.FromValue(), Description: "Full name of the repository that contains the branch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the branch."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.SHA"), Description: "Commit SHA the branch refers to."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.URL"), Description: "Commit URL the branch refers to."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Description: "True if the branch is protected."},
		},
	}
}

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.BranchListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		var branches []*github.Branch
		var resp *github.Response
		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}
		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			branches, resp, err = client.Repositories.ListBranches(ctx, owner, repo, opts)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		for _, i := range branches {
			d.StreamListItem(ctx, i)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return nil, nil
}
