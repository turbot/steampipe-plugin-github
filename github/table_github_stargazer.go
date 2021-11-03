package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubStargazer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_stargazer",
		Description: "Stargazers are users who have starred the repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubStargazerList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the stargazer."},
			{Name: "starred_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("StarredAt").Transform(convertTimestamp), Description: "Time when the stargazer was created."},
			{Name: "user_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("User.Login"), Description: "The login name of the user who starred the repository."},
			// No extra useful data over login - {Name: "user", Type: proto.ColumnType_JSON, Transform: transform.FromField("User"), Description: "Details of the user who starred the repository."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubStargazerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	type ListPageResponse struct {
		stargazers []*github.Stargazer
		resp       *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		stargazers, resp, err := client.Activity.ListStargazers(ctx, owner, repo, opts)
		return ListPageResponse{
			stargazers: stargazers,
			resp:       resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		stargazers := listResponse.stargazers
		resp := listResponse.resp

		for _, i := range stargazers {
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

		opts.Page = resp.NextPage
	}

	return nil, nil
}
