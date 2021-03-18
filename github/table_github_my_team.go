package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

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
		//{Name: "html_url", Type: pb.ColumnType_STRING, Description: "The URL of the team page in GitHub."},
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

func tableGitHubMyTeam() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_team",
		Description: "Github Teams in your organizations.  Github Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyTeamList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"organization_id", "id"}),
			Hydrate:    tableGitHubTeamGet,
		},
		Columns: gitHubTeamColumns(),
	}
}

//// list ////

func tableGitHubMyTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

	for {

		var items []*github.Team
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			items, resp, err = client.Teams.ListUserTeams(ctx, opt)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range items {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

//// hydrate functions ////

func tableGitHubTeamGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var orgID, teamID int64

	if h.Item != nil {
		team := h.Item.(*github.Team)
		orgID = *team.Organization.ID
		teamID = *team.ID
	} else {
		orgID = d.KeyColumnQuals["organization_id"].GetInt64Value()
		teamID = d.KeyColumnQuals["id"].GetInt64Value()
	}

	client := connect(ctx, d)

	var detail *github.Team
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		detail, resp, err = client.Teams.GetTeamByID(ctx, orgID, teamID)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}
