package github

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubOrganizationMemberColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the member is associated with.", Transform: transform.FromQual("organization")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The role this user has in the organization. Returns null if information is not available to viewer."},
		{Name: "has_two_factor_enabled", Type: proto.ColumnType_BOOL, Description: "Whether the organization member has two factor enabled or not. Returns null if information is not available to viewer."},
	}

	return append(tableCols, sharedUserColumns()...)
}

type memberWithRole struct {
	HasTwoFactorEnabled *bool
	Role                *string
	Node                models.User
}

func tableGitHubOrganizationMember() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_member",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
			},
			Hydrate: tableGitHubOrganizationMemberList,
		},
		Columns: gitHubOrganizationMemberColumns(),
	}
}

func tableGitHubOrganizationMemberList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Login           string
			MembersWithRole struct {
				Edges    []memberWithRole
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"membersWithRole(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil), // Null after argument to get first page.
	}
	appendUserColumnIncludes(&variables, d.QueryContext.Columns)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_organization_member", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_member", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, member := range query.Organization.MembersWithRole.Edges {
			d.StreamListItem(ctx, member)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.MembersWithRole.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.MembersWithRole.PageInfo.EndCursor)
	}

	return nil, nil
}
