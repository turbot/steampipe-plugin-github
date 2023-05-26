package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubMyOrganizationV3() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_organization_v3",
		Description: "",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyOrganizationV3List,
		},
		Columns: gitHubOrganizationV3Columns(),
	}
}

func tableGitHubMyOrganizationV3List(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	opt := &github.ListOptions{PerPage: pageSize}

	type ListPageResponse struct {
		org  []*github.Organization
		resp *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		orgs, resp, err := client.Organizations.List(ctx, "", opt)
		return ListPageResponse{
			org:  orgs,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		orgs := listResponse.org
		resp := listResponse.resp

		for _, i := range orgs {
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

		opt.Page = resp.NextPage
	}

	return nil, nil
}
