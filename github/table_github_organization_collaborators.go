package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubOrganizationCollaborators() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the member is associated with.", Transform: transform.FromQual("organization")},
		{Name: "affiliation", Type: proto.ColumnType_STRING, Description: "Affiliation filter - valid values 'ALL' (default), 'OUTSIDE', 'DIRECT'.", Transform: transform.FromQual("affiliation"), Default: "ALL"},
		{Name: "repository_name", Type: proto.ColumnType_STRING, Description: "The name of the repository", Transform: transform.FromValue(), Hydrate: ocHydrateRepository},
		{Name: "permission", Type: proto.ColumnType_STRING, Description: "The permission the collaborator has on the repository.", Transform: transform.FromValue(), Hydrate: ocHydratePermission},
		{Name: "user_login", Type: proto.ColumnType_JSON, Description: "The login of the collaborator", Transform: transform.FromValue(), Hydrate: ocHydrateUserLogin},
	}

	return tableCols
}

type OrgCollaborators struct {
	RepositoryName githubv4.String
	Permission     githubv4.RepositoryPermission
	Node           models.CollaboratorLogin
}

type CollaboratorEdge struct {
	Permission githubv4.RepositoryPermission `graphql:"permission @include(if:$includeOCPermission)" json:"permission"`
	Node       models.CollaboratorLogin      `graphql:"node @include(if:$includeOCNode)" json:"node"`
}

func tableGitHubOrganizationCollaborator() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_collaborator",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "organization",
					Require: plugin.Required,
				},
				{
					Name:       "affiliation",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
			Hydrate: listGitHubOrganizationCollaborators,
		},
		Columns: gitHubOrganizationCollaborators(),
	}
}

func listGitHubOrganizationCollaborators(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()
	a := quals["affiliation"].GetStringValue()
	affiliation := githubv4.CollaboratorAffiliationAll

	if a != "" {
		switch strings.ToLower(a) {
		case "direct":
			affiliation = githubv4.CollaboratorAffiliationDirect
		case "outside":
			affiliation = githubv4.CollaboratorAffiliationOutside
		case "all":
			affiliation = githubv4.CollaboratorAffiliationAll
		default:
			return nil, fmt.Errorf("filter must be 'ALL', 'OUTSIDE' or 'DIRECT' you provided '%s' which is invalid", a)
		}
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			URL          githubv4.String
			Login        githubv4.String
			Repositories struct {
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage githubv4.Boolean
				}
				Nodes []struct {
					Name          githubv4.String
					Collaborators struct {
						Edges []CollaboratorEdge
					} `graphql:"collaborators(affiliation: $affiliation)"`
				}
			} `graphql:"repositories(first: $pageSize, after: $cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":       githubv4.String(org),
		"pageSize":    githubv4.Int(pageSize),
		"cursor":      (*githubv4.String)(nil), // Null after argument to get first page.
		"affiliation": affiliation,
	}
	appendOrgCollaboratorColumnIncludes(&variables, d.QueryContext.Columns)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_organization_collaborator", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_collaborator", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, node := range query.Organization.Repositories.Nodes {
			for _, edge := range node.Collaborators.Edges {
				d.StreamListItem(ctx, OrgCollaborators{node.Name, edge.Permission, edge.Node})
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Repositories.PageInfo.EndCursor)
	}

	return nil, nil
}
