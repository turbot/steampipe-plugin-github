package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubTree(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_tree",
		Description: "Tree in the given repository, lists files in the git tree",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubTreeGet,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "tree_sha", Require: plugin.Required},
				{Name: "recursive", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tree."},
			{Name: "tree_sha", Type: proto.ColumnType_STRING, Transform: transform.FromQual("tree_sha"), Description: "SHA1 of the tree."},
			// TODO: How to handle recursive and truncated columns when splitting items?
			{Name: "recursive", Type: proto.ColumnType_BOOL, Transform: transform.FromQual("recursive"), Description: "If set to true, return objects or subtrees referenced by the tree. Defaults to false."},
			{Name: "truncated", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Truncated"), Description: "True if the entires were truncated because the number of items in the tree exceeded Github's maximum limit."},
			{Name: "mode", Type: proto.ColumnType_STRING, Description: "File mode, valid values are 100644 (file, blob), 100755 (executable, blob), 040000 (subdirectory, tree), 160000 (submodule, commit), 120000 (blob that specifies path of a symlink)."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The file referenced in the tree."},
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA"), Description: "SHA1 checksum ID of the object in the tree."},
			{Name: "size", Type: proto.ColumnType_STRING, Description: "Path."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the entry, valid values are blob, tree, or commit."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "URL to the file referenced in the tree."},
		},
	}
}

//// GET FUNCTION

func tableGitHubTreeGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client := connect(ctx, d)

	quals := d.KeyColumnQuals
	fullName := quals["repository_full_name"].GetStringValue()
	sha := quals["tree_sha"].GetStringValue()

	recursive := quals["recursive"].GetBoolValue()
	owner, repo := parseRepoFullName(fullName)

	type GetResponse struct {
		tree *github.Tree
		resp *github.Response
	}

	getTree := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		tree, resp, err := client.Git.GetTree(ctx, owner, repo, sha, recursive)
		return GetResponse{
			tree: tree,
			resp: resp,
		}, err
	}
	getResponse, err := plugin.RetryHydrate(ctx, d, h, getTree, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		logger.Error("github_tree.tableGitHubTreeGet", "api_error", err)
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	tree := getResp.tree
	entries := tree.Entries

	for _, entry := range entries {
		d.StreamListItem(ctx, entry)
	}

	return nil, nil
}
