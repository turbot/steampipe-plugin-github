package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubMyGist() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_gist",
		Description: "GitHub Gists owned by you. GitHub Gist is a simple way to share snippets and pastes with others.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyGistList,
		},
		Columns: gitHubGistColumns(),
	}
}

//// list ////
func tableGitHubMyGistList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	type ListPageResponse struct {
		myGist []*github.Gist
		resp   *github.Response
	}

	opt := &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		myGist, resp, err := client.Gists.List(ctx, "", opt)
		return ListPageResponse{
			myGist: myGist,
			resp:   resp,
		}, err
	}

	for {

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repos := listResponse.myGist
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
