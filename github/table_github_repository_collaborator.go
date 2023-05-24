package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryCollaboratorColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromQual("repository_full_name")},
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
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
		},
		Columns: gitHubRepositoryCollaboratorColumns(),
	}
}

func tableGitHubRepositoryCollaboratorList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

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
			} `graphql:"collaborators(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repoName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_collaborator", &query.RateLimit))
		if err != nil {
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
