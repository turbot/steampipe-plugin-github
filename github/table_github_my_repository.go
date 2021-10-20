package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubMyRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository",
		Description: "GitHub Repositories that you are associated with.  GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyRepositoryList,
		},

		Columns: gitHubRepositoryColumns(),
	}
}

//// list ////

func tableGitHubMyRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	type ListPageResponse struct {
		repo []*github.Repository
		resp *github.Response
	}

	opt := &github.RepositoryListOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 100}}

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
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
