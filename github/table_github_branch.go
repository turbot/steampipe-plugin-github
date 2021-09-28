package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubBranch(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:      "protected",
					Require:   plugin.Optional,
					Operators: []string{"<>", "="},
				},
			},
			Hydrate: tableGitHubBranchList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the branch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the branch."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.SHA"), Description: "Commit SHA the branch refers to."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.URL"), Description: "Commit URL the branch refers to."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Description: "True if the branch is protected."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := buildGithubBranchListOptions(d.KeyColumnQuals, d.Quals)
	opts.ListOptions = github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.ListOptions.PerPage) {
			opts.ListOptions.PerPage = int(*limit)
		}
	}

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

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return nil, nil
}

func buildGithubBranchListOptions(equalQuals plugin.KeyColumnEqualsQualMap, quals plugin.KeyColumnQualMap) *github.BranchListOptions {
	request := &github.BranchListOptions{}

	if equalQuals["protected"] != nil {
		request.Protected = types.Bool(equalQuals["protected"].GetBoolValue())
	}

	// Non-Equals Qual Map handling
	if quals["protected"] != nil {
		for _, q := range quals["protected"].Quals {
			value := q.Value.GetBoolValue()
			if q.Operator == "<>" {
				request.Protected = types.Bool(false)
				if !value {
					request.Protected = types.Bool(true)
				}
			}
		}
	}
	return request
}
