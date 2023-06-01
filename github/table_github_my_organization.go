package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
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

func tableGitHubMyOrganizationList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit models.RateLimit
		Viewer    struct {
			Organizations struct {
				TotalCount int
				PageInfo   models.PageInfo
				Nodes      []models.OrganizationWithCounts
			} `graphql:"organizations(first: $pageSize, after: $cursor)"`
		}
	}

	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_my_organization", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_my_organization", "api_error", err)
			return nil, err
		}

		for _, org := range query.Viewer.Organizations.Nodes {
			d.StreamListItem(ctx, org)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Viewer.Organizations.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Viewer.Organizations.PageInfo.EndCursor)
	}

	return nil, nil
}
