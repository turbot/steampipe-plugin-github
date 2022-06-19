package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

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

//// LIST FUNCTION

func tableGitHubGitignoreList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client) (interface{}, error) {
		raw, _, err := client.Gitignores.List(ctx)

		var list []github.Gitignore
		for _, v := range raw {
			list = append(list, github.Gitignore{Name: github.String(v)})
		}

		return list, err
	}

	return streamGitHubListOrItem(ctx, d, h, listPage)
}

//// HYDRATE FUNCTIONS

func tableGitHubGitignoreGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client) (interface{}, error) {
		var name string
		if h.Item != nil {
			item := h.Item.(github.Gitignore)
			name = *item.Name
		} else {
			name = d.KeyColumnQuals["name"].GetStringValue()
		}

		// Return nil, if no input provided
		if name == "" {
			return nil, nil
		}
		detail, _, err := client.Gitignores.Get(ctx, name)

		return detail, err
	}

	return getGitHubItem(ctx, d, h, getDetails)
}
