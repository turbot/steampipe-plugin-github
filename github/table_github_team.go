package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamColumns() []*plugin.Column {
	return []*plugin.Column{

		// Top columns
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromField("Organization.Login")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the team."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was created.", Transform: transform.FromField("CreatedAt")}, // .Transform(convertTimestamp)},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the team."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the team page in GitHub.", Transform: transform.FromField("Url")},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the team.", Transform: transform.FromField("DatabaseId")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the team.", Transform: transform.FromField("Id")},
		{Name: "organization_id", Type: proto.ColumnType_INT, Description: "The user id (number) of the organization.", Transform: transform.FromField("Organization.DatabaseId")},
		{Name: "organization_login", Type: proto.ColumnType_STRING, Description: "The login name of the organization.", Transform: transform.FromField("Organization.Login")},
		{Name: "members_count", Type: proto.ColumnType_INT, Description: "The number of members.", Transform: transform.FromField("Members.TotalCount")},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "The API Members URL.", Transform: transform.FromField("MembersUrl")},
		{Name: "parent", Type: proto.ColumnType_JSON, Description: "The parent team of the team.", Transform: transform.FromField("ParentTeam")},
		{Name: "privacy", Type: proto.ColumnType_STRING, Description: "The privacy setting of the team (VISIBLE or SECRET)."},
		{Name: "repos_count", Type: proto.ColumnType_INT, Description: "The number of repositories for the team.", Transform: transform.FromField("Repositories.TotalCount")},
		{Name: "repositories_url", Type: proto.ColumnType_STRING, Description: "The API Repositories URL.", Transform: transform.FromField("RepositoriesUrl")},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when team was last updated.", Transform: transform.FromField("UpdatedAt")},
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

type teamDetail struct {
	Slug                string
	Name                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Description         string
	DatabaseId          int
	Privacy             string
	ViewerCanAdminister bool
	Id                  string
	Members             struct {
		TotalCount int
	}
	Repositories struct {
		TotalCount int
	}
	Organization struct {
		DatabaseId int
		Name       string
		Login      string
	}
	Url             string
	TeamsUrl        string
	RepositoriesUrl string
	MembersUrl      string
	ParentTeam      struct {
		DatabaseId int
		Name       string
		Slug       string
	}
}

var teamsQuery struct {
	Organization struct {
		Teams struct {
			TotalCount int
			PageInfo   struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
			Nodes []teamDetail
		} `graphql:"teams(first: $teamPageSize, after: $teamCursor)"`
	} `graphql:"organization(login: $login)"`
}

var teamQuery struct {
	Organization struct {
		Team teamDetail `graphql:"team(slug: $slug)"`
	} `graphql:"organization(login: $login)"`
}

func tableGitHubTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	org := h.Item.(*github.Organization)

	pageSize := 100
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"login":        githubv4.String(*org.Login),
		"teamPageSize": githubv4.Int(pageSize),
		"teamCursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &teamsQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_team", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, team := range teamsQuery.Organization.Teams.Nodes {
			d.StreamListItem(ctx, team)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !teamsQuery.Organization.Teams.PageInfo.HasNextPage {
			break
		}

		variables["teamCursor"] = githubv4.NewString(teamsQuery.Organization.Teams.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubTeamGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()

	client := connectV4(ctx, d)

	variables := map[string]interface{}{
		"login": githubv4.String(org),
		"slug":  githubv4.String(slug),
	}

	err := client.Query(ctx, &teamQuery, variables)
	if err != nil {
		plugin.Logger(ctx).Error("github_team", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	return teamQuery.Organization.Team, nil
}
