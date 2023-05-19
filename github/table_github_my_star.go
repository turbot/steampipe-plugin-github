package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubMyStar() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_star",
		Description: "GitHub stars owned by you. GitHub stars are repositories.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyStarredRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromField("Node.NameWithOwner")},
			{Name: "starred_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("StarredAt").Transform(convertTimestamp), Description: "The timestamp when the repository was starred."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Url"), Description: "URL of the repository."},
		},
	}
}

func tableGitHubMyStarredRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			StarredRepositories struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					StarredAt models.NullableTime
					Node      struct {
						NameWithOwner string
						Url           string
					}
				}
			} `graphql:"starredRepositories(first: $pageSize, after: $cursor)"`
		}
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_star", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_star", "api_error", err)
			return nil, err
		}

		for _, star := range query.Viewer.StarredRepositories.Edges {
			d.StreamListItem(ctx, star)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.StarredRepositories.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubMyStarredRepositoryListOld(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ActivityListStarredOptions{ListOptions: github.ListOptions{PerPage: 100}}

	type ListPageResponse struct {
		starredRepos []*github.StarredRepository
		resp         *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		starredRepos, resp, err := client.Activity.ListStarred(ctx, "", opt)
		return ListPageResponse{
			starredRepos: starredRepos,
			resp:         resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		starredRepos := listResponse.starredRepos
		resp := listResponse.resp

		for _, i := range starredRepos {
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
