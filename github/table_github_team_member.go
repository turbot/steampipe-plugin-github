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

func tableGitHubTeamMember() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_member",
		Description: "GitHub members for a given team. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
				{Name: "role", Require: plugin.Optional},
			},
			Hydrate:           tableGitHubTeamMemberList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubTeamMemberColumns(),
	}
}

func gitHubTeamMemberColumns() []*plugin.Column {
	cols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The team member's role (MEMBER, MAINTAINER)."},
	}

	cols = append(cols, sharedUserColumns()...)
	return cols
}

func tableGitHubTeamMemberList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()
	slug := quals["slug"].GetStringValue()

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Team struct {
				Members struct {
					TotalCount int
					PageInfo   struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
					Edges []models.TeamMemberWithRole
				} `graphql:"members(first: $pageSize, after: $cursor)"`
			} `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"slug":     githubv4.String(slug),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_team_member", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_team_member", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, member := range query.Organization.Team.Members.Edges {
			d.StreamListItem(ctx, member)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Team.Members.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Team.Members.PageInfo.EndCursor)
	}

	return nil, nil
}
