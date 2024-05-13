package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryDeploymentColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromValue(), Hydrate: deploymentHydrateId, Description: "The ID of the deployment."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateNodeId, Description: "The node ID of the deployment."},
		{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateCommitSha, Description: "SHA of the commit the deployment is using."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: deploymentHydrateCreatedAt, Description: "Timestamp when the deployment was created."},
		{Name: "creator", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: deploymentHydrateCreator, Description: "The deployment creator."},
		{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateDescription, Description: "The description of the deployment."},
		{Name: "environment", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateEnvironment, Description: "The name of the environment to which the deployment was made."},
		{Name: "latest_environment", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateLatestEnvironment, Description: "The name of the latest environment to which the deployment was made."},
		{Name: "latest_status", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: deploymentHydrateLatestStatus, Description: "The latest status of the deployment."},
		{Name: "original_environment", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateOriginalEnvironment, Description: "The original environment to which this deployment was made."},
		{Name: "payload", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydratePayload, Description: "Extra information that a deployment system might need."},
		{Name: "ref", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: deploymentHydrateRef, Description: "Identifies the Ref of the deployment, if the deployment was created by ref."},
		{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateState, Description: "The current state of the deployment."},
		{Name: "task", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: deploymentHydrateTask, Description: "The deployment task."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: deploymentHydrateUpdatedAt, Description: "Timestamp when the deployment was last updated."},
	}
}

func tableGitHubRepositoryDeployment() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_deployment",
		Description: "GitHub Deployments are releases of the app/service/etc to an environment.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryDeploymentList,
		},
		Columns: commonColumns(gitHubRepositoryDeploymentColumns()),
	}
}

func tableGitHubRepositoryDeploymentList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Deployments struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.Deployment
			} `graphql:"deployments(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendRepoDeploymentColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_deployment", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_deployment", "api_error", err)
			return nil, err
		}

		for _, deployment := range query.Repository.Deployments.Nodes {
			d.StreamListItem(ctx, deployment)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Deployments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Deployments.PageInfo.EndCursor)
	}

	return nil, nil
}
