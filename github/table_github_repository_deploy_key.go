package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableGitHubRepositoryDeployKey() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_deploy_key",
		Description: "Deploy keys are SSH keys that grant access to a single repository.",
		List: &plugin.ListConfig{
			KeyColumns:   plugin.SingleColumn("repository_full_name"),
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"})},
			Hydrate:      tableGitHubRepositoryDeployKeyList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:   plugin.AllColumns([]string{"repository_full_name", "id"}),
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"})},
			Hydrate:      tableGitHubRepositoryDeployKeyGet,
		},
		// The public key material (`key`) is omitted, matching github_user, which exposes only a count of public keys.
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "The unique identifier of the deploy key."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The name of the deploy key."},
			{Name: "read_only", Type: proto.ColumnType_BOOL, Description: "If true, the key can only read the repository. If false, it can also push to it."},
			{Name: "verified", Type: proto.ColumnType_BOOL, Description: "If true, the key is verified."},
			{Name: "added_by", Type: proto.ColumnType_STRING, Description: "The login of the user who added the key."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Time when the deploy key was created."},
			{Name: "last_used", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastUsed").NullIfZero().Transform(convertTimestamp), Description: "Time when the deploy key was last used."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The REST API URL of the deploy key."},
		}),
	}
}

func tableGitHubRepositoryDeployKeyList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.ListOptions{PerPage: 100}

	client := connect(ctx, d)
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		keys, resp, err := client.Repositories.ListKeys(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range keys {
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

func tableGitHubRepositoryDeployKeyGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if id == 0 || fullName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubRepositoryDeployKeyGet", "owner", owner, "repo", repo, "id", id)

	client := connect(ctx, d)
	key, _, err := client.Repositories.GetKey(ctx, owner, repo, id)
	if err != nil {
		return nil, err
	}

	return key, nil
}
