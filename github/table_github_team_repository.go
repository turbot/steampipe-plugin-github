package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	pb "github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func gitHubTeamRepositoryColumns() []*plugin.Column {
	repoColumns := gitHubRepositoryColumns()
	teamColumns := []*plugin.Column{
		{Name: "organization", Type: pb.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: pb.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
	}

	return append(repoColumns, teamColumns...)
}

func tableGitHubTeamRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_repository",
		Description: "GitHub Repositories that a given team is associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
			},
			Hydrate:           tableGitHubTeamRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubTeamRepositoryColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubTeamRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

	org := d.KeyColumnQuals["organization"].GetStringValue()
	slug := d.KeyColumnQuals["slug"].GetStringValue()

	type ListPageResponse struct {
		repo []*github.Repository
		resp *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Teams.ListTeamReposBySlug(ctx, org, slug, opt)
		return ListPageResponse{
			repo: repos,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repos := listResponse.repo
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
