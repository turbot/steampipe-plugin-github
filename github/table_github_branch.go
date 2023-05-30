package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubBranchList,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the branch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the branch.", Transform: transform.FromField("Node.Name")},
			{Name: "commit", Type: proto.ColumnType_JSON, Transform: transform.FromField("Node.Target.Commit"), Description: "Latest commit on the branch."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.NodeId").Transform(HasValue), Description: "If true, the branch is protected."},
			{Name: "branch_protection_rule", Type: proto.ColumnType_JSON, Transform: transform.FromField("Node.BranchProtectionRule").NullIfZero(), Description: "Branch protection rule if protected."},
		},
	}
}

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Refs struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					Node models.Branch
				}
			} `graphql:"refs(refPrefix: \"refs/heads/\", first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := retryHydrate(ctx, d, h, listPage)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_branch", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_branch", "api_error", err)
			return nil, err
		}

		for _, branch := range query.Repository.Refs.Edges {
			d.StreamListItem(ctx, branch)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Refs.PageInfo.EndCursor)
	}

	return nil, nil
}

// HasValue Note: if useful to other tables, move to utils.go
func HasValue(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil || input.Value.(string) == "" {
		return false, nil
	}

	return true, nil
}
