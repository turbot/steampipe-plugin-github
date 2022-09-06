package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubBranch(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "protected", Require: plugin.Optional, Operators: []string{"<>", "="}},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubBranchList,
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

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	type ListPageResponse struct {
		branches []*github.Branch
		resp     *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		branches, resp, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
		return ListPageResponse{
			branches: branches,
			resp:     resp,
		}, err
	}
	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
		branches := listResponse.branches
		resp := listResponse.resp

		for _, i := range branches {
			if i != nil {
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
