package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubOrganizationExternalIdentityColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the external identity is associated with.", Transform: transform.FromQual("organization")},
		{Name: "guid", Type: proto.ColumnType_STRING, Description: "Guid identifier for the external identity.", Transform: transform.FromField("Guid")},
		{Name: "user_login", Type: proto.ColumnType_STRING, Description: "The GitHub user login.", Transform: transform.FromField("User.Login")},
		{Name: "user_detail", Type: proto.ColumnType_JSON, Description: "The GitHub user details.", Transform: transform.FromField("User")},
		{Name: "saml_identity", Type: proto.ColumnType_JSON, Description: "The external SAML identity."},
		{Name: "scim_identity", Type: proto.ColumnType_JSON, Description: "The external SCIM identity."},
		{Name: "organization_invitation", Type: proto.ColumnType_JSON, Description: "The invitation to the organization."},
	}
}

func tableGitHubOrganizationExternalIdentity() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_external_identity",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "organization",
					Require: plugin.Required,
				},
			},
			Hydrate: tableGitHubOrganizationExternalIdentityList,
		},
		Columns: gitHubOrganizationExternalIdentityColumns(),
	}
}

func tableGitHubOrganizationExternalIdentityList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			SamlIdentityProvider struct {
				ExternalIdentities struct {
					PageInfo   models.PageInfo
					TotalCount int
					Nodes      []models.OrganizationExternalIdentity
				} `graphql:"externalIdentities(first: $pageSize, after: $cursor)"`
			}
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_pull_request_comment", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_pull_request_comment", "api_error", err)
			return nil, err
		}

		for _, eid := range query.Organization.SamlIdentityProvider.ExternalIdentities.Nodes {
			d.StreamListItem(ctx, eid)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor)
	}

	return nil, nil
}
