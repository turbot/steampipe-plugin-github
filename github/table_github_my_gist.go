package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

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

//// LIST FUNCTION

func tableGitHubMyGistList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client := connect(ctx, d)

	opt := &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		var repos []*github.Gist
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			repos, resp, err = client.Gists.List(ctx, "", opt)
			logger.Error("tableGitHubGistList", "resp", resp)
			logger.Error("tableGitHubGistList", "repos", repos)
			logger.Error("tableGitHubGistList", "err", err)

			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range repos {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
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
