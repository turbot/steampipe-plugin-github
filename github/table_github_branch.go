package github

import (
	"context"

	"github.com/google/go-github/v33/github"

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
	type ListPageResponse struct {
		branches []*github.Branch
		resp     *github.Response
	}
	opts := &github.BranchListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		branches, resp, err := client.Repositories.ListBranches(ctx, owner, repo, opts)
		return ListPageResponse{
			branches: branches,
			resp:     resp,
		}, err
	}
	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})
		if err != nil {
			return nil, err
		}
		
		listResponse := listPageResponse.(ListPageResponse)
		branches := listResponse.branches
		resp := listResponse.resp

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
