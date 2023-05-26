package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubMyRepositoryV3Columns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the repository (public or private)", Hydrate: tableGitHubRepositoryV3Get},
	}
	return append(gitHubRepositoryV3Columns(), tableCols...)
}

func tableGitHubMyRepositoryV3() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository_v3",
		Description: "",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyRepositoryV3List,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "visibility", Require: plugin.Optional},
			},
		},
		Columns: gitHubMyRepositoryV3Columns(),
	}
}

func tableGitHubMyRepositoryV3List(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: pageSize}}

	// Additional filters
	if d.EqualsQuals["visibility"] != nil {
		opt.Visibility = d.EqualsQuals["visibility"].GetStringValue()
	} else {
		// Will cause a 422 error if 'type' used in the same request as visibility or affiliation.
		opt.Type = "all"
	}
	type ListPageResponse struct {
		repo []*github.Repository
		resp *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		return ListPageResponse{
			repo: repos,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

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
			if d.RowsRemaining(ctx) == 0 {
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
