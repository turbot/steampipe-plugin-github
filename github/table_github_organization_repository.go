package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func gitHubOrganizationRepositoryColumns() []*plugin.Column {
	repoColumns := gitHubRepositoryColumns()
	orgColumns := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "permissions", Type: proto.ColumnType_JSON, Description: "The team's permissions for a repository.", Transform: transform.From(perissionsFromMap)},
	}

	return append(repoColumns, orgColumns...)
}

func tableGitHubOrganizationRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_repository",
		Description: "GitHub Repositories that a given organization is associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
			},
			Hydrate:           tableGitHubOrganizationRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubOrganizationRepositoryColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubOrganizationRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: 100}}

	org := d.KeyColumnQuals["organization"].GetStringValue()

	type ListPageResponse struct {
		repos []*github.Repository
		resp  *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		return ListPageResponse{
			repos: repos,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repos := listResponse.repos
		resp := listResponse.resp

		for _, i := range repos {
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

		opt.Page = resp.NextPage
	}

	return nil, nil
}
