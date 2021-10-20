package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableGitHubMyOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_my_organization",
		Description: "GitHub Organizations that you are a member of. GitHub Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubMyOrganizationList,
		},
		Columns: gitHubOrganizationColumns(),
	}
}

//// list ////

func tableGitHubMyOrganizationList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

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

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		orgs := listResponse.org
		resp := listResponse.resp

		for _, i := range orgs {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
