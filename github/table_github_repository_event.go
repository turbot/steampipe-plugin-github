package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableGitHubRepositoryEvent(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_event",
		Description: "Rate limit of github.",
		List: &plugin.ListConfig{
			Hydrate: listGitHubRepositoryEvent,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
			},
		},
		Columns: []*plugin.Column{
			// Top columns

			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "The unique identifier for the event"},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the event was created."},
			{Name: "organization_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Org.ID"), Description: "The associated organization id."},
			{Name: "organization_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Org.Login"), Description: "The associated organization login."},
			{Name: "actor_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.ID"), Description: "The associated actor id."},
			{Name: "actor_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.Login"), Description: "The associated actor login."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Repo.Name"), Description: "The associated repository full name."},
			{Name: "raw_payload", Type: proto.ColumnType_JSON, Description: "The raw payload."},
			{Name: "public", Type: proto.ColumnType_BOOL, Description: "The associated repository full name."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The associated repository full name."},
		},
	}
}

func listGitHubRepositoryEvent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	repoName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(repoName)
	plugin.Logger(ctx).Trace("listGitHubRepositoryEvent", "owner", owner, "repo", repo)

	opts := github.ListOptions{PerPage: 100}

	type ListPageResponse struct {
		event []*github.Event
		resp  *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Activity.ListRepositoryEvents(ctx, owner, repo, &opts)

		return ListPageResponse{
			event: detail,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		events := listResponse.event
		resp := listResponse.resp
		for _, event := range events {
			plugin.Logger(ctx).Trace("-------------------------------->>>>>>>>>>>>>>>>>>>", "events", event)
			d.StreamListItem(ctx, event)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}
