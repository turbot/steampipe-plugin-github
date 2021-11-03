package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableGitHubMyRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository",
		Description: "GitHub Repositories that you are associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyRepositoryList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "visibility", Require: plugin.Optional},
			},
		},
		Columns: gitHubRepositoryColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubMyRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	// Additional filters
	if d.KeyColumnQuals["visibility"] != nil {
		opt.Visibility = d.KeyColumnQuals["visibility"].GetStringValue()
	} else {
		// Will cause a 422 error if 'type' used in the same request as visibility or
		// affiliation.
		opt.Type = "all"
	}
	type ListPageResponse struct {
		repo []*github.Repository
		resp *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		return ListPageResponse{
			repo: repos,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

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
