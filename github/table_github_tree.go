package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubTree() *plugin.Table {
	return &plugin.Table{
		Name:        "github_tree",
		Description: "Lists directories and files in the given repository's git tree.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubTreeList,
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
			// Other columns
			{Name: "recursive", Type: proto.ColumnType_BOOL, Description: "If set to true, return objects or subtrees referenced by the tree. Defaults to false."},
			{Name: "truncated", Type: proto.ColumnType_BOOL, Description: "True if the entires were truncated because the number of items in the tree exceeded Github's maximum limit."},
			{Name: "mode", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.Mode"), Description: "File mode. Valid values are 100644 (blob file), 100755 (blob executable), 040000 (tree subdirectory), 160000 (commit submodule), 120000 (blob that specifies path of a symlink)."},
			{Name: "path", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.Path"), Description: "The file referenced in the tree."},
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.SHA"), Description: "SHA1 checksum ID of the object in the tree."},
			{Name: "size", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.Size"), Description: "Size of the blob."},
			{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.Type"), Description: "Either blob, tree, or commit."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeEntry.URL"), Description: "URL to the file referenced in the tree."},
		},
	}
}

type treeEntry struct {
	TreeEntry *github.TreeEntry
	Recursive bool
	Truncated *bool
}

func tableGitHubTreeList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client := connect(ctx, d)

	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	sha := quals["tree_sha"].GetStringValue()
	recursive := quals["recursive"].GetBoolValue()
	owner, repo := parseRepoFullName(fullName)

	tree, _, err := client.Git.GetTree(ctx, owner, repo, sha, recursive)
	if err != nil {
		logger.Error("github_tree.tableGitHubTreeList", "api_error", err)
		return nil, err
	}

	entries := tree.Entries
	for _, entry := range entries {
		entryRow := treeEntry{
			TreeEntry: entry,
			Recursive: recursive,
			Truncated: tree.Truncated,
		}
		d.StreamListItem(ctx, entryRow)
	}

	return nil, nil
}
