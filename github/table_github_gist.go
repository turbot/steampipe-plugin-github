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

func tableGitHubGist() *plugin.Table {
	return &plugin.Table{
		Name:        "github_gist",
		Description: "Github Gist is a simple way to share snippets and pastes with others.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubGistList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    tableGitHubGistGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "id", Type: pb.ColumnType_STRING, Description: "The unique id of the gist."},
			{Name: "description", Type: pb.ColumnType_STRING, Description: "The gist description."},
			{Name: "public", Type: pb.ColumnType_BOOL, Description: "If true, the gist is public, otherwise it is private."},
			{Name: "html_url", Type: pb.ColumnType_STRING, Description: "The HTML URL of the gist."},

			{Name: "comments", Type: pb.ColumnType_INT, Description: "The number of comments for the gist."},
			{Name: "created_at", Type: pb.ColumnType_TIMESTAMP, Description: "The timestamp when the gist was created."},
			{Name: "git_pull_url", Type: pb.ColumnType_STRING, Description: "The https url to pull or clone the gist."},
			{Name: "git_push_url", Type: pb.ColumnType_STRING, Description: "The https url to push the gist."},
			{Name: "node_id", Type: pb.ColumnType_STRING, Description: "The Node ID of the gist."},
			// Only load relevant fields from the owner
			{Name: "owner_id", Type: pb.ColumnType_INT, Description: "The user id (number) of the gist owner.", Transform: transform.FromField("Owner.ID")},
			{Name: "owner_login", Type: pb.ColumnType_STRING, Description: "The user login name of the gist owner.", Transform: transform.FromField("Owner.Login")},
			{Name: "owner_type", Type: pb.ColumnType_STRING, Description: "The type of the gist owner (User or Organization).", Transform: transform.FromField("Owner.Type")},
			{Name: "updated_at", Type: pb.ColumnType_TIMESTAMP, Description: "The timestamp when the gist was last updated."},
		},
	}
}

//// list ////

func tableGitHubGistList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	client := connect(ctx, d)

	opt := &github.GistListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {

		var repos []*github.Gist
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			repos, resp, err = client.Gists.List(ctx, "", opt)
			logger.Error("tableGitHubGistList", "resp", resp)
			logger.Error("tableGitHubGistList", "repos", repos)
			logger.Error("tableGitHubGistList", "err", err)

			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range repos {
			logger.Error("tableGitHubGistList", "i", i)
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

func tableGitHubGistGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	var id string

	if h.Item != nil {
		gist := h.Item.(*github.Gist)
		id = *gist.ID
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
	}

	detail, _, err := client.Gists.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return detail, nil
}
