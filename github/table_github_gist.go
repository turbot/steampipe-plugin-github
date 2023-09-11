package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubGistColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique id of the gist."},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The gist description."},
		{Name: "public", Type: proto.ColumnType_BOOL, Description: "If true, the gist is public, otherwise it is private."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The HTML URL of the gist."},
		{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments for the gist."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "The timestamp when the gist was created."},
		{Name: "git_pull_url", Type: proto.ColumnType_STRING, Description: "The https url to pull or clone the gist."},
		{Name: "git_push_url", Type: proto.ColumnType_STRING, Description: "The https url to push the gist."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The Node ID of the gist."},
		// Only load relevant fields from the owner
		{Name: "owner_id", Type: proto.ColumnType_INT, Description: "The user id (number) of the gist owner.", Transform: transform.FromField("Owner.ID")},
		{Name: "owner_login", Type: proto.ColumnType_STRING, Description: "The user login name of the gist owner.", Transform: transform.FromField("Owner.Login")},
		{Name: "owner_type", Type: proto.ColumnType_STRING, Description: "The type of the gist owner (User or Organization).", Transform: transform.FromField("Owner.Type")},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "The timestamp when the gist was last updated."},
		{Name: "files", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files").Transform(gistFileMapToArray), Description: "Files in the gist."},
	}
}

func tableGitHubGist() *plugin.Table {
	return &plugin.Table{
		Name:        "github_gist",
		Description: "GitHub Gist is a simple way to share snippets and pastes with others.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubGistList,
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubGistColumns(),
	}
}

func tableGitHubGistList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	var id string
	if h.Item != nil {
		gist := h.Item.(*github.Gist)
		id = *gist.ID
	} else {
		id = d.EqualsQuals["id"].GetStringValue()
	}

	gist, _, err := client.Gists.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if gist != nil {
		d.StreamListItem(ctx, gist)
	}

	return nil, nil
}

func gistFileMapToArray(ctx context.Context, input *transform.TransformData) (interface{}, error) {
	var objectList []github.GistFile
	objectMap := input.Value.(map[github.GistFilename]github.GistFile)
	for _, v := range objectMap {
		objectList = append(objectList, v)
	}
	return objectList, nil
}
