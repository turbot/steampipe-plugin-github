package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubTree(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_tree",
		Description: "Tree in the given repository, lists files in the git tree",
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "tree_sha", Require: plugin.Required},
				{Name: "recursive", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubTreeGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository"},
			{Name: "tree_sha", Type: proto.ColumnType_STRING, Transform: transform.FromQual("tree_sha"), Description: "Tree SHA"},
			{Name: "recursive", Type: proto.ColumnType_BOOL, Transform: transform.FromQual("recursive"), Description: "Recursive tree content"},
			{Name: "truncated", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Truncated"), Description: "Whether results are truncated"},
			{Name: "entries", Type: proto.ColumnType_JSON, Transform: transform.FromField("Entries"), Description: "Tree entries"},
		},
	}
}

//// GET FUNCTION

func tableGitHubTreeGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("Connecting to client")
	client := connect(ctx, d)

	if h.Item != nil {
		item := h.Item
		logger.Trace("Got hydrate data", item)
	} else {
		logger.Trace("No hydrate data")
	}
	quals := d.KeyColumnQuals
	logger.Trace("Parsing key column quals", quals)
	fullName := quals["repository_full_name"].GetStringValue()
	sha := quals["tree_sha"].GetStringValue()
	recursive := quals["recursive"].GetBoolValue()
	owner, repo := parseRepoFullName(fullName)

	type GetResponse struct {
		tree *github.Tree
		resp *github.Response
	}

	getTree := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		logger.Trace("Getting tree", "owner", owner, "repo", repo, "sha", sha)
		tree, resp, err := client.Git.GetTree(ctx, owner, repo, sha, recursive)
		return GetResponse{
			tree: tree,
			resp: resp,
		}, err
	}
	getResponse, err := plugin.RetryHydrate(ctx, d, h, getTree, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		logger.Error("Error getting tree", err)
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	tree := getResp.tree
	if tree != nil {
		logger.Trace("Returning tree", tree)
		return tree, nil
	}
	logger.Error("Nothing found")
	return nil, nil
}
