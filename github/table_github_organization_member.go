package github

import (
	"context"
	"strings"

	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func gitHubOrganizationMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the member is associated with.", Transform: transform.FromQual("organization")},
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The username used to login.", Transform: transform.FromField("Node.Login")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The role this user has in the organization. Returns null if information is not available to viewer."},
		{Name: "has_two_factor_enabled", Type: proto.ColumnType_BOOL, Description: "Whether the organization member has two factor enabled or not. Returns null if information is not available to viewer."},
	}
}

type memberWithRole struct {
	HasTwoFactorEnabled *bool
	Role                *string
	Node                struct {
		Login string
	}
}

var query struct {
	Organization struct {
		Login           string
		MembersWithRole struct {
			Edges    []memberWithRole
			PageInfo struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
		} `graphql:"membersWithRole(first: $membersWithRolePageSize, after: $membersWithRoleCursor)"`
	} `graphql:"organization(login: $login)"`
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

//// LIST FUNCTION

func tableGitHubOrganizationMemberList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()

	pageSize := 100

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"login":                   githubv4.String(org),
		"membersWithRolePageSize": githubv4.Int(pageSize),
		"membersWithRoleCursor":   (*githubv4.String)(nil), // Null after argument to get first page.
	}

	for {
		err := client.Query(ctx, &query, variables)
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
		variables["membersWithRoleCursor"] = githubv4.NewString(query.Organization.MembersWithRole.PageInfo.EndCursor)
	}

	return nil, nil
}
