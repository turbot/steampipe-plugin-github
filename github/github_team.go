package github

import (
	"context"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"

	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubTeam() *plugin.Table {
	return &plugin.Table{
		Name: "github_team",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubTeamList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"organization_id", "id"}),
			Hydrate:    tableGitHubTeamGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "organization", Type: pb.ColumnType_STRING, Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
			{Name: "slug", Type: pb.ColumnType_STRING},
			{Name: "name", Type: pb.ColumnType_STRING},
			{Name: "members_count", Type: pb.ColumnType_INT, Hydrate: tableGitHubTeamGet},

			// Not yet supported by go-github
			//{Name: "created_at", Type: pb.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet, Transform: transform.Transform(convertTimestamp)},
			{Name: "description", Type: pb.ColumnType_STRING},
			{Name: "html_url", Type: pb.ColumnType_STRING},
			{Name: "id", Type: pb.ColumnType_INT},
			{Name: "members_url", Type: pb.ColumnType_STRING},
			{Name: "node_id", Type: pb.ColumnType_STRING},
			// Only load relevant fields from the organization
			{Name: "organization_id", Type: pb.ColumnType_INT, Transform: transform.FromField("Organization.ID"), Hydrate: tableGitHubTeamGet},
			{Name: "organization_login", Type: pb.ColumnType_STRING, Transform: transform.FromField("Organization.Login"), Hydrate: tableGitHubTeamGet},
			{Name: "organization_type", Type: pb.ColumnType_STRING, Transform: transform.FromField("Organization.Type"), Hydrate: tableGitHubTeamGet},
			{Name: "permission", Type: pb.ColumnType_STRING},
			{Name: "privacy", Type: pb.ColumnType_STRING},
			{Name: "repos_count", Type: pb.ColumnType_INT, Hydrate: tableGitHubTeamGet},
			{Name: "repositories_url", Type: pb.ColumnType_STRING},
			// Not yet supported by go-github
			// {Name: "updated_at", Type: pb.ColumnType_TIMESTAMP, Hydrate: tableGitHubTeamGet},
			{Name: "url", Type: pb.ColumnType_STRING},
		},
	}
}

//// list ////

func tableGitHubTeamList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d.ConnectionManager)

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

	client := connect(ctx, d.ConnectionManager)

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
