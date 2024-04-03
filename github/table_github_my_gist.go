package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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

func tableGitHubMyGistList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		gists, resp, err := client.Gists.List(ctx, "", opt)
		if err != nil {
			return nil, err
		}

		for _, i := range gists {
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
