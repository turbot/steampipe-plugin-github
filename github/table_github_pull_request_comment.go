package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubPullRequestComment() *plugin.Table {
	return &plugin.Table{
		Name:        "github_pull_request_comment",
		Description: "Comments are the responses/comments on Pull Requests.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryPullRequestCommentList,
		},
		Columns: sharedCommentsColumns(),
	}
}

func tableGitHubRepositoryPullRequestCommentList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	prNumber := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			PullRequest struct {
				Comments struct {
					PageInfo   models.PageInfo
					TotalCount int
					Nodes      []models.IssueComment
				} `graphql:"comments(first: $pageSize, after: $cursor)"`
			} `graphql:"pullRequest(number: $prNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"prNumber": githubv4.Int(prNumber),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_pull_request_comment", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_pull_request_comment", "api_error", err)
			return nil, err
		}

		for _, comment := range query.Repository.PullRequest.Comments.Nodes {
			d.StreamListItem(ctx, comment)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.PullRequest.Comments.PageInfo.EndCursor)
	}

	return nil, nil
}
