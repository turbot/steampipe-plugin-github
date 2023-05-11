package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"strings"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromField("Organization.Login")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the team."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the team.", Transform: transform.FromField("Id")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the team.", Transform: transform.FromField("NodeId")},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the team."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was created.", Transform: transform.FromField("CreatedAt")},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was last updated.", Transform: transform.FromField("UpdatedAt")},
		{Name: "combined_slug", Type: proto.ColumnType_STRING, Description: "The slug corresponding to the organization and the team."},
		{Name: "parent_team", Type: proto.ColumnType_JSON, Description: "The teams parent team.", Transform: transform.FromField("ParentTeam").NullIfZero()},
		{Name: "privacy", Type: proto.ColumnType_STRING, Description: "The privacy setting of the team (VISIBLE or SECRET)."},
		{Name: "ancestors_total_count", Type: proto.ColumnType_INT, Description: "Count of ancestors this team has.", Transform: transform.FromField("Ancestors.TotalCount")},
		{Name: "child_teams_total_count", Type: proto.ColumnType_INT, Description: "Count of children teams this team has.", Transform: transform.FromField("ChildTeams.TotalCount")},
		{Name: "discussions_total_count", Type: proto.ColumnType_INT, Description: "Count of team discussions.", Transform: transform.FromField("Discussions.TotalCount")},
		{Name: "invitations_total_count", Type: proto.ColumnType_INT, Description: "Count of outstanding team member invitations for the team.", Transform: transform.FromField("Invitations.TotalCount")},
		{Name: "members_total_count", Type: proto.ColumnType_INT, Description: "Count of team members.", Transform: transform.FromField("Members.TotalCount")},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT, Description: "Count of the teams v2 projects.", Transform: transform.FromField("ProjectsV2.TotalCount")},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the team has.", Transform: transform.FromField("Repositories.TotalCount")},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "URL for the team page in GitHub.", Transform: transform.FromField("Url")},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "URL for teams avatar.", Transform: transform.FromField("AvatarUrl")},
		{Name: "discussions_url", Type: proto.ColumnType_STRING, Description: "URL for team discussions.", Transform: transform.FromField("DiscussionsUrl")},
		{Name: "edit_team_url", Type: proto.ColumnType_STRING, Description: "URL for editing this team.", Transform: transform.FromField("EditTeamUrl")},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "URL for team members.", Transform: transform.FromField("MembersUrl")},
		{Name: "repositories_url", Type: proto.ColumnType_STRING, Description: "URL for team repositories.", Transform: transform.FromField("RepositoriesUrl")},
		{Name: "teams_url", Type: proto.ColumnType_STRING, Description: "URL for this team's teams.", Transform: transform.FromField("TeamsUrl")},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Description: "If true, current user can administer the team."},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Description: "If true, current user can subscribe to the team."},
		{Name: "subscription", Type: proto.ColumnType_STRING, Description: "Subscription status of the current user to the team."},
	}
}

func tableGitHubTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team",
		Description: "GitHub Teams in a given organization. GitHub Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.",
		List: &plugin.ListConfig{
			ParentHydrate:     tableGitHubMyOrganizationList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTeamList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "slug"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTeamGet,
		},
		Columns: gitHubTeamColumns(),
	}
}

func tableGitHubTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	org := h.Item.(*github.Organization)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Teams struct {
				TotalCount int
				PageInfo   struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Nodes []models.TeamWithCounts
			} `graphql:"teams(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(*org.Login),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_team", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_team", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, team := range query.Organization.Teams.Nodes {
			d.StreamListItem(ctx, team)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Teams.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(query.Organization.Teams.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubTeamGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()

	client := connectV4(ctx, d)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Team models.TeamWithCounts `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(org),
		"slug":  githubv4.String(slug),
	}

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_team", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_team", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	return query.Organization.Team, nil
}
