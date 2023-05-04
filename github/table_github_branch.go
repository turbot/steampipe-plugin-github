package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubBranch(ctx context.Context) *plugin.Table {
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
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the branch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the branch.", Transform: transform.FromField("Node.Name")},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Oid"), Description: "Commit SHA the branch refers to."},
			{Name: "commit_short_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.AbbreviatedOid"), Description: "Commit short SHA the branch refers to."},
			{Name: "commit_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Node.Target.Commit.AuthoredDate"), Description: "Date commit was authored."},
			{Name: "commit_author", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Author.Name"), Description: "Commit author."},
			{Name: "commit_message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Message"), Description: "Commit message."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Url"), Description: "Commit URL the branch refers to."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.Id").Transform(HasValue), Description: "True if the branch is protected."},
			{Name: "protection_rule_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.BranchProtectionRule.Id").NullIfZero(), Description: "Branch protection rule id, null if not protected."},
		},
	}
}

var branchesQuery struct {
	Repository struct {
		Refs struct {
			TotalCount int
			PageInfo   struct {
				EndCursor   githubv4.String
				HasNextPage bool
			}
			Edges []struct {
				Node struct {
					Name   string
					Target struct {
						Commit struct {
							Oid            string
							AbbreviatedOid string
							Url            string
							AuthoredDate   time.Time
							Message        string
							Author         struct {
								Name string
							}
						} `graphql:"... on Commit"`
					}
					BranchProtectionRule struct {
						Id string
					}
				}
			}
		} `graphql:"refs(refPrefix: \"refs/heads/\", first: $pageSize, after: $cursor)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := 100
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &branchesQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_branch", "api_error", err)
			return nil, err
		}

		for _, branch := range branchesQuery.Repository.Refs.Edges {
			d.StreamListItem(ctx, branch)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !branchesQuery.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(branchesQuery.Repository.Refs.PageInfo.EndCursor)
	}

	return nil, nil
}

// Note: if useful to other tables, move to utils.go
func HasValue(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil || input.Value.(string) == "" {
		return false, nil
	}

	return true, nil
}
