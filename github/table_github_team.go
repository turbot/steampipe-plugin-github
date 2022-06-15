package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v45/github"

	pb "github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func gitHubTeamColumns() []*plugin.Column {
	return []*plugin.Column{

		// Top columns
		{Name: "organization", Type: pb.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
		{Name: "slug", Type: pb.ColumnType_STRING, Description: "The team slug name."},
		{Name: "name", Type: pb.ColumnType_STRING, Description: "The name of the team."},
		{Name: "members_count", Type: pb.ColumnType_INT, Description: "The number of members.", Hydrate: tableGitHubTeamGet},

		// Not yet supported by go-github
		//{Name: "created_at", Type: pb.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet, Transform: transform.Transform(convertTimestamp)},
		{Name: "description", Type: pb.ColumnType_STRING, Description: "The description of the team."},
		// Not yet supported by go-github
		// {Name: "html_url", Type: pb.ColumnType_STRING, Description: "The URL of the team page in GitHub."},
		{Name: "id", Type: pb.ColumnType_INT, Description: "The ID of the team."},
		{Name: "members_url", Type: pb.ColumnType_STRING, Description: "The API Members URL."},
		{Name: "node_id", Type: pb.ColumnType_STRING, Description: "The node id of the team."},
		// Only load relevant fields from the organization
		{Name: "organization_id", Type: pb.ColumnType_INT, Description: "The user id (number) of the organization.", Transform: transform.FromField("Organization.ID"), Hydrate: tableGitHubTeamGet},
		{Name: "organization_login", Type: pb.ColumnType_STRING, Description: "The login name of the organization.", Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
		{Name: "organization_type", Type: pb.ColumnType_STRING, Description: "The type of the organization).", Transform: transform.FromField("Organization.Type"), Hydrate: tableGitHubTeamGet},
		{Name: "permission", Type: pb.ColumnType_STRING, Description: "The default repository permissions of the team."},
		{Name: "privacy", Type: pb.ColumnType_STRING, Description: "The privacy setting of the team (closed or secret)."},
		{Name: "repos_count", Type: pb.ColumnType_INT, Description: "The number of repositories for the team.", Hydrate: tableGitHubTeamGet},
		{Name: "repositories_url", Type: pb.ColumnType_STRING, Description: "The API Repositories URL."},
		// Not yet supported by go-github
		// {Name: "updated_at", Type: pb.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet},
		{Name: "url", Type: pb.ColumnType_STRING, Description: "The API URL of the team."},
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
	var org, slug string
	if h.Item != nil {
		team := h.Item.(*github.Team)
		slug = *team.Slug

		// Organization login will be available from different sources based on how
		// this function is called
		if d.KeyColumnQuals["organization"].GetStringValue() != "" { // If called with quals from the github_team table, use the qual value
			org = d.KeyColumnQuals["organization"].GetStringValue()
		} else if team.Organization != nil { // If called with quals from github_my_team table, use the Organization value
			org = *team.Organization.Login
		} else { // If called without quals from the github_team table, extract it from the team's HTML URL
			htmlUrl := *team.HTMLURL
			// Split a URL like "https://github.com/orgs/github/teams/justice-league"
			org = strings.Split(htmlUrl, "/")[4]
		}
	} else {
		org = d.KeyColumnQuals["organization"].GetStringValue()
		slug = d.KeyColumnQuals["slug"].GetStringValue()
	}

	client := connect(ctx, d)

	type GetResponse struct {
		team *github.Team
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Teams.GetTeamBySlug(ctx, org, slug)
		return GetResponse{
			team: detail,
			resp: resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}
	getResp := getResponse.(GetResponse)
	team := getResp.team

	return team, nil
}
