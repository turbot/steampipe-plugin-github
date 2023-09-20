package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubMyRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_repository",
		Description: "GitHub Repositories that you are associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: sharedRepositoryColumns(),
	}
}

func tableGitHubMyRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	pageSize := adjustPageSize(50, d.QueryContext.Limit)

	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			Repositories struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.Repository
			} `graphql:"repositories(first: $pageSize, after: $cursor, affiliations: [COLLABORATOR, OWNER, ORGANIZATION_MEMBER], ownerAffiliations: [COLLABORATOR, OWNER, ORGANIZATION_MEMBER])"`
		}
	}

	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendRepoColumnIncludes(&variables, d.QueryContext.Columns)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_repository", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_repository", "api_error", err)
			return nil, err
		}

		for _, repo := range query.Viewer.Repositories.Nodes {
			d.StreamListItem(ctx, repo)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.Repositories.PageInfo.EndCursor)
	}

	return nil, nil
}
