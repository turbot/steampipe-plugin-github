package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamRepositoryColumns() []*plugin.Column {
	teamColumns := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
		{Name: "permission", Type: proto.ColumnType_STRING, Description: "The permission level the team has on the repository."},
	}

	return append(teamColumns, sharedRepositoryColumns()...)
}

func tableGitHubTeamRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_repository",
		Description: "GitHub Repositories that a given team is associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
			},
			Hydrate:           tableGitHubTeamRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
				{Name: "name", Require: plugin.Required},
			},
			Hydrate:           tableGitHubTeamRepositoryGet,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubTeamRepositoryColumns(),
	}
}

func tableGitHubTeamRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()
	pageSize := adjustPageSize(50, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Team struct {
				Repositories struct {
					TotalCount int
					PageInfo   models.PageInfo
					Edges      []models.TeamRepositoryWithPermission
				} `graphql:"repositories(first: $pageSize, after: $cursor)"`
			} `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"slug":     githubv4.String(slug),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_team_repository", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_team_repository", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, repo := range query.Organization.Team.Repositories.Edges {
			d.StreamListItem(ctx, repo)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Team.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Team.Repositories.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubTeamRepositoryGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Team struct {
				Repositories struct {
					TotalCount int
					PageInfo   models.PageInfo
					Edges      []models.TeamRepositoryWithPermission
				} `graphql:"repositories(first: $pageSize, query: $name)"`
			} `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"slug":     githubv4.String(slug),
		"name":     githubv4.String(name),
		"pageSize": githubv4.Int(1),
	}

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_team_repository", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_team_repository", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	if len(query.Organization.Team.Repositories.Edges) == 1 && query.Organization.Team.Repositories.Edges[0].Node.Name == name {
		return query.Organization.Team.Repositories.Edges[0], nil
	}

	return nil, nil
}
