package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubProjectV2Columns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The organization name."},
	}

	return append(tableCols, sharedProjectV2Columns()...)
}

func sharedProjectV2Columns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number", "Node.Number"), Description: "The project number."},
		{Name: "id", Type: proto.ColumnType_INT, Hydrate: projectV2HydrateId, Transform: transform.FromValue(), Description: "The ID of the project."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateNodeId, Transform: transform.FromValue(), Description: "The node ID of the project."},
		{Name: "owner", Type: proto.ColumnType_JSON, Hydrate: projectV2HydrateOwner, Transform: transform.FromValue().NullIfZero(), Description: "The owner of the project."},
		{Name: "creator", Type: proto.ColumnType_JSON, Hydrate: projectV2HydrateCreator, Transform: transform.FromValue().NullIfZero(), Description: "The creator of the project."},
		{Name: "title", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateTitle, Transform: transform.FromValue(), Description: "The title of the project."},
		{Name: "description", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateDescription, Transform: transform.FromValue(), Description: "The description of the project (maps to shortDescription in GraphQL)."},
		{Name: "is_public", Type: proto.ColumnType_BOOL, Hydrate: projectV2HydrateIsPublic, Transform: transform.FromValue(), Description: "If true, the project is public."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: projectV2HydrateClosedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "The time when the project was closed."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: projectV2HydrateCreatedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "The time when the project was created."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: projectV2HydrateUpdatedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "The time when the project was last updated."},
		{Name: "state", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateState, Transform: transform.FromValue(), Description: "The state of the project (open or closed). Derived from the GraphQL closed boolean."},
		{Name: "latest_status_update", Type: proto.ColumnType_JSON, Hydrate: projectV2HydrateLatestStatusUpdate, Transform: transform.FromValue().NullIfZero(), Description: "The latest status update of the project."},
		{Name: "is_template", Type: proto.ColumnType_BOOL, Hydrate: projectV2HydrateIsTemplate, Transform: transform.FromValue(), Description: "If true, the project is a template."},
	}
}

func tableGitHubProjectV2() *plugin.Table {
	return &plugin.Table{
		Name:        "github_project_v2",
		Description: "GitHub Projects are used to organize and manage work on GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "organization",
					Require: plugin.Required,
				},
				{
					Name:      "updated_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">="},
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubProjectV2List,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubProjectV2Get,
		},
		Columns: commonColumns(gitHubProjectV2Columns()),
	}
}

func tableGitHubProjectV2List(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	organization := quals["organization"].GetStringValue()

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			ProjectsV2 struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.ProjectV2
			} `graphql:"projectsV2(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $organization)"`
	}

	variables := map[string]interface{}{
		"organization": githubv4.String(organization),
		"pageSize":     githubv4.Int(pageSize),
		"cursor":       (*githubv4.String)(nil),
	}
	appendProjectV2ColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
			return nil, err
		}

		for _, project := range query.Organization.ProjectsV2.Nodes {
			d.StreamListItem(ctx, project)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.ProjectsV2.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.ProjectsV2.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubProjectV2Get(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	projectId := int(quals["id"].GetInt64Value())
	organization := quals["organization"].GetStringValue()

	client := connectV4(ctx, d)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			ProjectV2 models.ProjectV2 `graphql:"projectV2(id: $projectId)"`
		} `graphql:"organization(login: $organization)"`
	}

	variables := map[string]interface{}{
		"organization": githubv4.String(organization),
		"projectId": githubv4.Int(projectId),
	}
	appendProjectV2ColumnIncludes(&variables, d.QueryContext.Columns)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
		return nil, err
	}

	return query.Organization.ProjectV2, nil
}
