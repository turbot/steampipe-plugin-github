package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func gitHubTeamColumns() []*plugin.Column {
	return []*plugin.Column{

		// Top columns
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the team."},
		{Name: "members_count", Type: proto.ColumnType_INT, Description: "The number of members.", Hydrate: tableGitHubTeamGet},

		// Not yet supported by go-github
		//{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the team."},
		// Not yet supported by go-github
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the team page in GitHub."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the team."},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "The API Members URL."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the team."},
		// Only load relevant fields from the organization
		{Name: "organization_id", Type: proto.ColumnType_INT, Description: "The user id (number) of the organization.", Transform: transform.FromField("Organization.ID"), Hydrate: tableGitHubTeamGet},
		{Name: "organization_login", Type: proto.ColumnType_STRING, Description: "The login name of the organization.", Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
		{Name: "organization_type", Type: proto.ColumnType_STRING, Description: "The type of the organization).", Transform: transform.FromField("Organization.Type"), Hydrate: tableGitHubTeamGet},
		{Name: "permission", Type: proto.ColumnType_STRING, Description: "The default repository permissions of the team."},
		{Name: "privacy", Type: proto.ColumnType_STRING, Description: "The privacy setting of the team (closed or secret)."},
		{Name: "repos_count", Type: proto.ColumnType_INT, Description: "The number of repositories for the team.", Hydrate: tableGitHubTeamGet},
		{Name: "repositories_url", Type: proto.ColumnType_STRING, Description: "The API Repositories URL."},
		// Not yet supported by go-github
		//{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL of the team."},
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

//// LIST FUNCTION

func tableGitHubTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	org := h.Item.(*github.Organization)

	opt := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

	type ListPageResponse struct {
		teams []*github.Team
		resp  *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		teams, resp, err := client.Teams.ListTeams(ctx, *org.Login, opt)
		return ListPageResponse{
			teams: teams,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		teams := listResponse.teams
		resp := listResponse.resp

		for _, i := range teams {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func tableGitHubTeamGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client) (interface{}, error) {
		var org, slug string
		if h.Item != nil {
			team := h.Item.(*github.Team)
			slug = *team.Slug

			// Organization login will be available from different sources based on how
			// this function is called
			if team.Organization != nil { // If called from github_my_team table, use the team's organization login
				org = *team.Organization.Login
			} else if h.ParentItem != nil { // If called from github_team table through parent hydrate, use the parent organization's login
				parentOrg := h.ParentItem.(*github.Organization)
				org = *parentOrg.Login
			} else { // Unknown caller
				plugin.Logger(ctx).Error("tableGitHubTeam.tableGitHubTeamGet", "unknown_caller_error")
				return nil, fmt.Errorf("unknown caller for tableGitHubTeamGet function")
			}
		} else {
			org = d.KeyColumnQuals["organization"].GetStringValue()
			slug = d.KeyColumnQuals["slug"].GetStringValue()
		}

		detail, _, err := client.Teams.GetTeamBySlug(ctx, org, slug)
		return detail, err
	}

	return getGitHubItem(ctx, d, h, getDetails)
}
