package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func gitHubSearchRepositoryColumns() []*plugin.Column {
	return append(defaultSearchColumns(), sharedRepositoryColumns()...)
}

func tableGitHubSearchRepository(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_repository",
		Description: "Find repositories via various criteria.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchRepositoryList,
		},
		Columns: gitHubSearchRepositoryColumns(),
	}
}

func tableGitHubSearchRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			RepositoryCount int
			PageInfo        models.PageInfo
			Edges           []struct {
				TextMatches []models.TextMatch
				Node        struct {
					models.Repository `graphql:"... on Repository"`
				}
			}
		} `graphql:"search(type: REPOSITORY, first: $pageSize, after: $cursor, query: $query)"`
	}

	pageSize := adjustPageSize(75, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"query":    githubv4.String(input),
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_repository", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_repository", "api_error", err)
			return nil, err
		}

		for _, repo := range query.Search.Edges {
			d.StreamListItem(ctx, repo)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubSearchRepositoryListOld(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchRepositoryList")

	quals := d.EqualsQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.RepositoriesSearchResult
		resp   *github.Response
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		result, resp, err := client.Search.Repositories(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchRepositoryList", "error_Search.Repositories", err)
			return nil, err
		}

		return ListPageResponse{
			result: result,
			resp:   resp,
		}, nil
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			logger.Error("tableGitHubSearchRepositoryList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repoResults := listResponse.result.Repositories
		resp := listResponse.resp

		for _, i := range repoResults {
			d.StreamListItem(ctx, i)

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
