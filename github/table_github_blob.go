package github

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubBlob() *plugin.Table {
	return &plugin.Table{
		Name:        "github_blob",
		Description: "Gets a blob from a repository.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubBlobList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "blob_sha", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the blob."},
			{Name: "blob_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Sha"), Description: "SHA1 of the blob."},
			// Other columns
			{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeID"), Description: "The node ID of the blob."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL of the blob."},
			{Name: "content", Type: proto.ColumnType_STRING, Description: "The encoded content of the blob."},
			{Name: "encoding", Type: proto.ColumnType_STRING, Description: "The encoding of the blob."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Size of the blob."},
		},
	}
}

type blobRow struct {
	Sha      string
	NodeID   string
	Url      string
	Content  string
	Encoding string
	Size     int
}

func tableGitHubBlobList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client := connect(ctx, d)

	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	sha := quals["blob_sha"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	blob, _, err := client.Git.GetBlob(ctx, owner, repo, sha)
	if err != nil {
		logger.Error("github_blob.tableGitHubBlobList", "api_error", err)
		return nil, err
	}

	if blob != nil {
		d.StreamListItem(ctx, blobRow{
			Sha:      *blob.SHA,
			NodeID:   *blob.NodeID,
			Content:  *blob.Content,
			Encoding: *blob.Encoding,
			Size:     *blob.Size,
		})
	}

	return nil, nil
}
