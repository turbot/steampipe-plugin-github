package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubGitignore() *plugin.Table {
	return &plugin.Table{
		Name:        "github_gitignore",
		Description: "GitHub defined .gitignore templates that you can associate with your repository.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubGitignoreList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubGitignoreGetData,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the gitignore template."},
			{Name: "source", Type: proto.ColumnType_STRING, Hydrate: tableGitHubGitignoreGetData, Description: "Source code of the gitignore template."},
		},
	}
}

func tableGitHubGitignoreList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	gitIgnores, _, err := client.Gitignores.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range gitIgnores {
		if i != "" {
			d.StreamListItem(ctx, github.Gitignore{Name: github.String(i)})
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func tableGitHubGitignoreGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		item := h.Item.(github.Gitignore)
		name = *item.Name
	} else {
		name = d.EqualsQuals["name"].GetStringValue()
	}

	// Return nil, if no input provided
	if name == "" {
		return nil, nil
	}

	client := connect(ctx, d)

	gitIgnore, _, err := client.Gitignores.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	return gitIgnore, nil
}
