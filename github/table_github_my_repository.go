package github

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubMyRepository() *plugin.Table {
	// Get shared columns and override hooks column with Fine-Grained PAT handling
	columns := sharedRepositoryColumns()

	// Override hooks column to use the Fine-Grained PAT-aware hydrate function
	for i, col := range columns {
		if col.Name == "hooks" {
			columns[i] = &plugin.Column{
				Name:        "hooks",
				Type:        proto.ColumnType_JSON,
				Description: "The API Hooks URL.",
				Hydrate:     hydrateMyRepositoryHooksFromV3,
				Transform:   transform.FromValue(),
			}
			break
		}
	}

	return &plugin.Table{
		Name:        "github_my_repository",
		Description: "GitHub Repositories that you are associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubMyRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: commonColumns(columns),
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
	appendRepoColumnIncludesWithQueryData(&variables, d.QueryContext.Columns, d)

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

// hydrateMyRepositoryHooksFromV3 is a version of hydrateRepositoryHooksFromV3
// that skips hooks for Fine-Grained PATs (only for github_my_repository table)
func hydrateMyRepositoryHooksFromV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Skip hooks for Fine-Grained PATs to avoid field errors
	// With Fine-grained access token we are getting field error even though we have proper access.
	// https://spec.graphql.org/October2021/#sec-Errors.Field-errors
	// https://spec.graphql.org/October2021/#sec-Handling-Field-Errors
	githubConfig := GetConfig(d.Connection)
	token := os.Getenv("GITHUB_TOKEN")
	if isGitHubPAT(token) || (githubConfig.Token != nil && isGitHubPAT(*githubConfig.Token)) {
		return nil, nil
	}

	repo, err := extractRepoFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	owner := repo.Owner.Login
	repoName := repo.Name

	client := connect(ctx, d)
	var repositoryHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	for {
		hooks, resp, err := client.Repositories.ListHooks(ctx, owner, repoName, opt)
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		repositoryHooks = append(repositoryHooks, hooks...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositoryHooks, nil
}
