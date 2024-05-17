package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsArtifact() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_artifact",
		Description: "Artifacts allow you to share data between jobs in a workflow and store data once that workflow has completed.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubArtifactList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubArtifactGet,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the artifact."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the artifact."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the artifact."},
			{Name: "size_in_bytes", Type: proto.ColumnType_INT, Description: "Size of the artifact in bytes."},

			// Other columns
			{Name: "archive_download_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ArchiveDownloadURL"), Description: "Archive download URL for the artifact."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the artifact was created."},
			{Name: "expired", Type: proto.ColumnType_BOOL, Description: "It defines whether the artifact is expires or not."},
			{Name: "expires_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ExpiresAt").Transform(convertTimestamp), Description: "Time when the artifact expires."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "Node where GitHub stores this data internally."},
		}),
	}
}

func tableGitHubArtifactList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		artifacts, resp, err := client.Actions.ListArtifacts(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range artifacts.Artifacts {
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

func tableGitHubArtifactGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if id == 0 || fullName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubArtifactGet", "owner", owner, "repo", repo, "id", id)

	client := connect(ctx, d)

	artifact, _, err := client.Actions.GetArtifact(ctx, owner, repo, id)
	if err != nil {
		return nil, err
	}

	return artifact, nil
}
