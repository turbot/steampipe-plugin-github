package github

import (
	"context"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func gitHubProjectV2Columns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The organization name. Either this or 'user_login' must be provided."},
		{Name: "user_login", Type: proto.ColumnType_STRING, Transform: transform.FromQual("user_login"), Description: "The user login. Either this or 'organization' must be provided."},
	}

	return append(tableCols, sharedProjectV2Columns()...)
}

func sharedProjectV2Columns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number"), Description: "The project number."},
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
		{Name: "readme", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateReadme, Transform: transform.FromValue(), Description: "The readme of the project."},
		{Name: "resource_path", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateResourcePath, Transform: transform.FromValue(), Description: "The HTTP path for this project."},
		{Name: "url", Type: proto.ColumnType_STRING, Hydrate: projectV2HydrateUrl, Transform: transform.FromValue(), Description: "The HTTP URL for this project."},
		{Name: "repositories", Type: proto.ColumnType_JSON, Hydrate: projectV2HydrateRepositories, Transform: transform.FromValue(), Description: "Array of full repository names (owner/repo) linked to the project (first 100 only; see repositories_total_count for the real total)."},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Hydrate: projectV2HydrateRepositoriesTotalCount, Transform: transform.FromValue(), Description: "Count of repositories linked to the project."},
		{Name: "teams", Type: proto.ColumnType_JSON, Hydrate: projectV2HydrateTeams, Transform: transform.FromValue(), Description: "Array of team slugs linked to the project (first 100 only; see teams_total_count for the real total)."},
		{Name: "teams_total_count", Type: proto.ColumnType_INT, Hydrate: projectV2HydrateTeamsTotalCount, Transform: transform.FromValue(), Description: "Count of teams linked to the project."},
	}
}

func tableGitHubProjectV2() *plugin.Table {
	return &plugin.Table{
		Name:        "github_project_v2",
		Description: "GitHub Projects are used to organize and manage work on GitHub. Projects are listed by their owner, which must be provided as either an 'organization' or a 'user_login' qual.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "organization",
					Require: plugin.AnyOf,
				},
				{
					Name:    "user_login",
					Require: plugin.AnyOf,
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
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:    "number",
					Require: plugin.Required,
				},
				{
					Name:    "organization",
					Require: plugin.AnyOf,
				},
				{
					Name:    "user_login",
					Require: plugin.AnyOf,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubProjectV2Get,
		},
		Columns: commonColumns(gitHubProjectV2Columns()),
	}
}

func tableGitHubProjectV2List(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	organization := quals["organization"].GetStringValue()
	user := quals["user_login"].GetStringValue()

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	// The projectsV2 GraphQL field has no server-side filter for updatedAt, only an
	// orderBy argument. To still honor the updated_at qual (and avoid scanning every
	// page), we always request results ordered newest-first by updatedAt, then stop
	// paging as soon as we see a project older than the requested threshold.
	var minUpdatedAt *time.Time
	var minUpdatedAtInclusive bool
	if d.Quals["updated_at"] != nil {
		for _, q := range d.Quals["updated_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">":
				if minUpdatedAt == nil || givenTime.After(*minUpdatedAt) {
					minUpdatedAt = &givenTime
					minUpdatedAtInclusive = false
				}
			case ">=":
				if minUpdatedAt == nil || givenTime.After(*minUpdatedAt) {
					minUpdatedAt = &givenTime
					minUpdatedAtInclusive = true
				}
			}
		}
	}

	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"orderBy": githubv4.ProjectV2Order{
			Field:     githubv4.ProjectV2OrderFieldUpdatedAt,
			Direction: githubv4.OrderDirectionDesc,
		},
	}
	appendProjectV2ColumnIncludes(&variables, d.QueryContext.Columns)
	if minUpdatedAt != nil {
		// Force updatedAt to be fetched so we can compare it, even if the column wasn't requested.
		variables["includeUpdatedAt"] = githubv4.Boolean(true)
	}

	client := connectV4(ctx, d)

	// projectsV2 is exposed identically by the Organization and User GraphQL types (both
	// implement ProjectV2Owner), but the root query field ("organization"/"user") differs,
	// so we need a distinct query shape per owner kind.
	if organization != "" {
		var query struct {
			RateLimit    models.RateLimit
			Organization struct {
				ProjectsV2 struct {
					PageInfo models.PageInfo
					Nodes    []models.ProjectV2
				} `graphql:"projectsV2(first: $pageSize, after: $cursor, orderBy: $orderBy)"`
			} `graphql:"organization(login: $owner)"`
		}
		variables["owner"] = githubv4.String(organization)

		for {
			err := client.Query(ctx, &query, variables)
			plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
			if err != nil {
				plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
				if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
					return nil, nil
				}
				return nil, err
			}

			for _, project := range query.Organization.ProjectsV2.Nodes {
				if minUpdatedAt != nil && !project.UpdatedAt.IsZero() {
					updatedAt := project.UpdatedAt.Time
					if updatedAt.Before(*minUpdatedAt) || (!minUpdatedAtInclusive && updatedAt.Equal(*minUpdatedAt)) {
						// Results are ordered newest-first, so once we see a project
						// older than the threshold, every subsequent project (on this
						// page and later pages) will also be too old.
						return nil, nil
					}
				}

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

	var query struct {
		RateLimit models.RateLimit
		User      struct {
			ProjectsV2 struct {
				PageInfo models.PageInfo
				Nodes    []models.ProjectV2
			} `graphql:"projectsV2(first: $pageSize, after: $cursor, orderBy: $orderBy)"`
		} `graphql:"user(login: $owner)"`
	}
	variables["owner"] = githubv4.String(user)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, project := range query.User.ProjectsV2.Nodes {
			if minUpdatedAt != nil && !project.UpdatedAt.IsZero() {
				updatedAt := project.UpdatedAt.Time
				if updatedAt.Before(*minUpdatedAt) || (!minUpdatedAtInclusive && updatedAt.Equal(*minUpdatedAt)) {
					// Results are ordered newest-first, so once we see a project
					// older than the threshold, every subsequent project (on this
					// page and later pages) will also be too old.
					return nil, nil
				}
			}

			d.StreamListItem(ctx, project)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.User.ProjectsV2.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.User.ProjectsV2.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubProjectV2Get(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	number := int(quals["number"].GetInt64Value())
	organization := quals["organization"].GetStringValue()
	user := quals["user_login"].GetStringValue()

	client := connectV4(ctx, d)

	variables := map[string]interface{}{
		"number": githubv4.Int(number),
	}
	appendProjectV2ColumnIncludes(&variables, d.QueryContext.Columns)

	var project models.ProjectV2
	if organization != "" {
		var query struct {
			RateLimit    models.RateLimit
			Organization struct {
				ProjectV2 models.ProjectV2 `graphql:"projectV2(number: $number)"`
			} `graphql:"organization(login: $owner)"`
		}
		variables["owner"] = githubv4.String(organization)

		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
			// GitHub's GraphQL API returns a "Could not resolve to a ProjectV2 ..."
			// error (rather than a 404) when the project number doesn't exist, and a
			// "Could not resolve to an Organization ..." error when the organization
			// doesn't exist. Treat both the same as a not-found result rather than
			// propagating the error, since letting a Get error out (rather than a
			// clean "no rows") can leave the plugin SDK's query cache with an
			// orphaned pending entry that causes subsequent identical queries to hang
			// waiting on it.
			if strings.Contains(err.Error(), "Could not resolve to a ProjectV2") ||
				strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}
		project = query.Organization.ProjectV2
	} else {
		var query struct {
			RateLimit models.RateLimit
			User      struct {
				ProjectV2 models.ProjectV2 `graphql:"projectV2(number: $number)"`
			} `graphql:"user(login: $owner)"`
		}
		variables["owner"] = githubv4.String(user)

		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_project_v2", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to a ProjectV2") ||
				strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
				return nil, nil
			}
			return nil, err
		}
		project = query.User.ProjectV2
	}

	return project, nil
}
