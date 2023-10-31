package github

import (
	"context"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

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
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the team.", Transform: transform.FromValue(), Hydrate: teamHydrateDescription},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was created.", Transform: transform.FromValue(), Hydrate: teamHydrateCreatedAt},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was last updated.", Transform: transform.FromValue(), Hydrate: teamHydrateUpdatedAt},
		{Name: "combined_slug", Type: proto.ColumnType_STRING, Description: "The slug corresponding to the organization and the team.", Transform: transform.FromValue(), Hydrate: teamHydrateCombinedSlug},
		{Name: "parent_team", Type: proto.ColumnType_JSON, Description: "The teams parent team.", Transform: transform.FromValue().NullIfZero(), Hydrate: teamHydrateParentTeam},
		{Name: "privacy", Type: proto.ColumnType_STRING, Description: "The privacy setting of the team (VISIBLE or SECRET).", Transform: transform.FromValue(), Hydrate: teamHydratePrivacy},
		{Name: "ancestors_total_count", Type: proto.ColumnType_INT, Description: "Count of ancestors this team has.", Transform: transform.FromValue(), Hydrate: teamHydrateAncestorsTotalCount},
		{Name: "child_teams_total_count", Type: proto.ColumnType_INT, Description: "Count of children teams this team has.", Transform: transform.FromValue(), Hydrate: teamHydrateChildTeamsTotalCount},
		{Name: "discussions_total_count", Type: proto.ColumnType_INT, Description: "Count of team discussions.", Transform: transform.FromValue(), Hydrate: teamHydrateDiscussionsTotalCount},
		{Name: "invitations_total_count", Type: proto.ColumnType_INT, Description: "Count of outstanding team member invitations for the team.", Transform: transform.FromValue(), Hydrate: teamHydrateInvitationsTotalCount},
		{Name: "members_total_count", Type: proto.ColumnType_INT, Description: "Count of team members.", Transform: transform.FromValue(), Hydrate: teamHydrateMembersTotalCount},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT, Description: "Count of the teams v2 projects.", Transform: transform.FromValue(), Hydrate: teamHydrateProjectsV2TotalCount},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the team has.", Transform: transform.FromValue(), Hydrate: teamHydrateRepositoriesTotalCount},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "URL for the team page in GitHub.", Transform: transform.FromValue(), Hydrate: teamHydrateUrl},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "URL for teams avatar.", Transform: transform.FromValue(), Hydrate: teamHydrateAvatarUrl},
		{Name: "discussions_url", Type: proto.ColumnType_STRING, Description: "URL for team discussions.", Transform: transform.FromValue(), Hydrate: teamHydrateDiscussionsUrl},
		{Name: "edit_team_url", Type: proto.ColumnType_STRING, Description: "URL for editing this team.", Transform: transform.FromValue(), Hydrate: teamHydrateEditTeamUrl},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "URL for team members.", Transform: transform.FromValue(), Hydrate: teamHydrateMembersUrl},
		{Name: "new_team_url", Type: proto.ColumnType_STRING, Description: "The HTTP URL creating a new team.", Transform: transform.FromValue(), Hydrate: teamHydrateNewTeamUrl},
		{Name: "repositories_url", Type: proto.ColumnType_STRING, Description: "URL for team repositories.", Transform: transform.FromValue(), Hydrate: teamHydrateRepositoriesUrl},
		{Name: "teams_url", Type: proto.ColumnType_STRING, Description: "URL for this team's teams.", Transform: transform.FromValue(), Hydrate: teamHydrateTeamsUrl},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Description: "If true, current user can administer the team.", Transform: transform.FromValue(), Hydrate: teamHydrateCanAdminister},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Description: "If true, current user can subscribe to the team.", Transform: transform.FromValue(), Hydrate: teamHydrateCanSubscribe},
		{Name: "subscription", Type: proto.ColumnType_STRING, Description: "Subscription status of the current user to the team.", Transform: transform.FromValue(), Hydrate: teamHydrateSubscription},
	}
}

func tableGitHubTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team",
		Description: "GitHub Teams in a given organization. GitHub Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("organization"),
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

func tableGitHubTeamList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	org := d.EqualsQuals["organization"].GetStringValue()
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
		"login":    githubv4.String(org),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}
	appendTeamColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
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

func tableGitHubTeamGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()

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
	appendTeamColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
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
