package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

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

func tableGitHubMyOrganizationList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

	for {

		var orgs []*github.Organization
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			orgs, resp, err = client.Organizations.List(ctx, "", opt)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

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
