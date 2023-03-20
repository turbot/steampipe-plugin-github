package github

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

//// TABLE DEFINITION

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

//// LIST FUNCTION

func tableGitHubMyOrganizationList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

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
