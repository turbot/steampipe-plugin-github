package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"
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
				{
					Name:    "visibility",
					Require: plugin.Optional,
				},
			},
		},
		Columns: gitHubRepositoryColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubMyRepositoryList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		var repos []*github.Repository
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			repos, resp, err = client.Repositories.List(ctx, "", opt)
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
