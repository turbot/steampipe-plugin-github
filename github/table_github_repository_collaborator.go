package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryCollaboratorColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromQual("repository_full_name")},
		{Name: "affiliation", Type: proto.ColumnType_STRING, Description: "Affiliation filter - valid values 'ALL' (default), 'OUTSIDE', 'DIRECT'.", Transform: transform.FromQual("affiliation"), Default: "ALL"},
		{Name: "permission", Type: proto.ColumnType_STRING, Description: "The permission the collaborator has on the repository.", Transform: transform.FromValue(), Hydrate: rcHydratePermission},
		{Name: "user_login", Type: proto.ColumnType_STRING, Description: "The login of the collaborator", Transform: transform.FromValue(), Hydrate: rcHydrateUserLogin},
	}
}

func tableGitHubRepositoryCollaborator() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_collaborator",
		Description: "Collaborators are users that have contributed to the repository.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubRepositoryCollaboratorList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:       "affiliation",
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
			},
		},
		Columns: gitHubRepositoryCollaboratorColumns(),
	}
}

type RepositoryCollaborator struct {
	Permission githubv4.RepositoryPermission
	Node       models.BasicUser
}

func tableGitHubRepositoryCollaboratorList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)
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

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Collaborators struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					Permission githubv4.RepositoryPermission `graphql:"permission @include(if:$includeRCPermission)" json:"permission"`
					Node       models.BasicUser              `graphql:"node @include(if:$includeRCNode)" json:"node"`
				}
			} `graphql:"collaborators(first: $pageSize, after: $cursor, affiliation: $affiliation)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"repo":        githubv4.String(repoName),
		"pageSize":    githubv4.Int(pageSize),
		"cursor":      (*githubv4.String)(nil),
		"affiliation": affiliation,
	}
	appendRepoCollaboratorColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_collaborator", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_collaborator", "api_error", err, "repository", fullName)
			return nil, err
		}

		for _, c := range query.Repository.Collaborators.Edges {
			d.StreamListItem(ctx, RepositoryCollaborator{c.Permission, c.Node})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Collaborators.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Collaborators.PageInfo.EndCursor)
	}

	return nil, nil
}
