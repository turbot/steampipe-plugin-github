package github

import (
	"context"
	"time"

	"github.com/google/go-github/v47/github"
	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableGitHubMyRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository",
		Description: "GitHub Repositories that you are associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
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

func retryHydrate(ctx context.Context, d *plugin.QueryData, hydrateData *plugin.HydrateData, hydrateFunc plugin.HydrateFunc) (interface{}, error) {

	// Retry configs
	maxRetries := 10
	interval := time.Duration(500)

	// Create the backoff based on the given mode
	backoff, err := retry.NewFibonacci(interval * time.Millisecond)
	if err != nil {
		return nil, err
	}

	// Ensure the maximum value is 2.5s. In this scenario, the sleep values would be
	// 0.5s, 0.5s, 1s, 1.5s, 2.5s, 2.5s, 2.5s...
	backoff = retry.WithCappedDuration(2500*time.Millisecond, backoff)

	var hydrateResult interface{}

	err = retry.Do(ctx, retry.WithMaxRetries(uint64(maxRetries), backoff), func(ctx context.Context) error {
		hydrateResult, err = hydrateFunc(ctx, d, hydrateData)
		if err != nil {
			if shouldRetryError(err) {
				err = retry.RetryableError(err)
			}
		}
		return err
	})

	return hydrateResult, err
}
