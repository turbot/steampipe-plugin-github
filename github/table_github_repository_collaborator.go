package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

func gitHubRepositoryCollaboratorColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromQual("repository_full_name")},
		{Name: "filter", Type: proto.ColumnType_STRING, Description: "Affiliation filter - valid values 'ALL' (default), 'OUTSIDE', 'DIRECT'.", Transform: transform.FromQual("filter"), Default: "ALL"},
		{Name: "permission", Type: proto.ColumnType_STRING, Description: "The permission the collaborator has on the repository."},
		{Name: "user_login", Type: proto.ColumnType_STRING, Description: "The login of the collaborator", Transform: transform.FromField("Node.Login")},
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
					Name:       "filter",
					Require:    plugin.Optional,
					Operators:  []string{"="},
					CacheMatch: "exact",
				},
			},
		},
		Columns: gitHubRepositoryCollaboratorColumns(),
	}
}

func tableGitHubRepositoryCollaboratorList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)
	filter := quals["filter"].GetStringValue()
	affiliation := githubv4.CollaboratorAffiliationAll

	if filter != "" {
		switch strings.ToLower(filter) {
		case "direct":
			affiliation = githubv4.CollaboratorAffiliationDirect
		case "outside":
			affiliation = githubv4.CollaboratorAffiliationOutside
		case "all":
			affiliation = githubv4.CollaboratorAffiliationAll
		default:
			return nil, fmt.Errorf("filter must be 'ALL', 'OUTSIDE' or 'DIRECT' you provided '%s' which is invalid", filter)
		}
	}

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Collaborators struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					Permission githubv4.RepositoryPermission
					Node       models.BasicUser
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

	client := connectV4(ctx, d)
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_collaborator", &query.RateLimit))
		if err != nil {
			if strings.Contains(err.Error(), "You do not have permission to view repository collaborators") {
				return nil, nil
			}
			plugin.Logger(ctx).Error("github_repository_collaborator", "api_error", err)
			return nil, err
		}

		for _, c := range query.Repository.Collaborators.Edges {
			d.StreamListItem(ctx, c)

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
