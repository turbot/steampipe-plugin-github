package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableGitHubActionsCache() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_cache",
		Description: "Caches are a built-in storage mechanism to speed up jobs in a workflow by persisting data for re-use.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "key", Require: plugin.Optional},
				{Name: "ref", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubCacheList,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the cache."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the cache."},
			{Name: "key", Type: proto.ColumnType_STRING, Description: "Developer-defined string identifier of the cache."},

			// Other columns
			{Name: "ref", Type: proto.ColumnType_STRING, Description: "The git reference of the cache."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "Hash generated from combination of compression tool, runner OS, and path."},
			{Name: "last_accessed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastAccessedAt").Transform(convertTimestamp), Description: "Time of the most recent cache access."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the cache was created."},
			{Name: "size_in_bytes", Type: proto.ColumnType_INT, Description: "Size of the cache in bytes."},
		}),
	}
}

func tableGitHubCacheList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.ActionsCacheListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	if key := d.EqualsQualString("key"); key != "" {
		opts.Key = github.String(key)
	}
	if ref := d.EqualsQualString("ref"); ref != "" {
		opts.Ref = github.String(ref)
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		caches, resp, err := client.Actions.ListCaches(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range caches.ActionsCaches {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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
