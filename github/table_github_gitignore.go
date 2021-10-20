package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubGitignore() *plugin.Table {
	return &plugin.Table{
		Name:        "github_gitignore",
		Description: "GitHub defined .gitignore templates that you can associate with your repository.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubGitignoreList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("name"),
			Hydrate:    tableGitHubGitignoreGetData,
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

	type ListPageResponse struct {
		gitIgnores []string
		resp       *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		gitignore, resp, err := client.Gitignores.List(ctx)
		return ListPageResponse{
			gitIgnores: gitignore,
			resp:       resp,
		}, err
	}

	listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

	if err != nil {
		return nil, err
	}

	listResponse := listPageResponse.(ListPageResponse)
	gitIgnores := listResponse.gitIgnores

	for _, i := range gitIgnores {
		d.StreamListItem(ctx, github.Gitignore{Name: github.String(i)})
	}
	return nil, nil
}

func tableGitHubGitignoreGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var name string
	if h.Item != nil {
		item := h.Item.(github.Gitignore)
		name = *item.Name
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}
	client := connect(ctx, d)

	type GetResponse struct {
		gitIgnore *github.Gitignore
		resp      *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Gitignores.Get(ctx, name)
		return GetResponse{
			gitIgnore: detail,
			resp:      resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	gitIgnore := getResp.gitIgnore

	return gitIgnore, nil
}
