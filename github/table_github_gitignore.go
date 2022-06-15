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

	listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}

	listResponse := listPageResponse.(ListPageResponse)
	gitIgnores := listResponse.gitIgnores

	for _, i := range gitIgnores {
		if i != "" {
			d.StreamListItem(ctx, github.Gitignore{Name: github.String(i)})
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func tableGitHubGitignoreGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	gitIgnore := getResp.gitIgnore

	return gitIgnore, nil
}
