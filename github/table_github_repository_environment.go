package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryEnvironmentColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID of the environment."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the environment."},
		{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name", "Node.Name"), Description: "The name of the environment."},
	}
}

func tableGitHubRepositoryEnvironment() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_environment",
		Description: "GitHub Environments are named deployment targets, usually isolated for usage such as test, prod, staging, etc.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryEnvironmentList,
		},
		Columns: gitHubRepositoryEnvironmentColumns(),
	}
}

func tableGitHubRepositoryEnvironmentList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Environments struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.Environment
			} `graphql:"environments(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_environment", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_environment", "api_error", err)
			return nil, err
		}

		for _, environment := range query.Repository.Environments.Nodes {
			d.StreamListItem(ctx, environment)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Environments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Environments.PageInfo.EndCursor)
	}

	return nil, nil
}
