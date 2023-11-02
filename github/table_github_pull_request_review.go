package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func pullRequestReviewColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromQual("number"), Description: "The PR number."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromValue(), Hydrate: prReviewHydrateId, Description: "The ID of the review."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateNodeId, Description: "The node ID of the review."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: prReviewHydrateAuthor, Description: "The actor who authored the review."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateAuthorLogin, Description: "The login of the review author."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateAuthorAssociation, Description: "Author's association with the subject of the pr the review was raised on."},
		{Name: "author_can_push_to_repository", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: prReviewHydrateAuthorCanPushToRepository, Description: "Indicates whether the author of this review has push access to the repository."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateBody, Description: "The body of the review."},
		{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateState, Description: "The state of the review."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: prReviewHydrateUrl, Description: "The HTTP URL permalink for this PullRequestReview."},
		{Name: "submitted_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: prReviewHydrateSubmittedAt, Description: "Identifies when the Pull Request Review was submitted."},
	}
}

func tableGitHubPullRequestReview() *plugin.Table {
	return &plugin.Table{
		Name:        "github_pull_request_review",
		Description: "Pull Request Reviews are groups of pull request review comments on a pull request.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryPullRequestReviewList,
		},
		Columns: pullRequestReviewColumns(),
	}
}

func tableGitHubRepositoryPullRequestReviewList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	prNumber := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			PullRequest struct {
				Reviews struct {
					PageInfo   models.PageInfo
					TotalCount int
					Nodes      []models.PullRequestReview
				} `graphql:"reviews(first: $pageSize, after: $cursor)"`
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
	appendPRReviewColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_pull_request_review", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_pull_request_review", "api_error", err)
			return nil, err
		}

		for _, comment := range query.Repository.PullRequest.Reviews.Nodes {
			d.StreamListItem(ctx, comment)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.PullRequest.Reviews.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.PullRequest.Reviews.PageInfo.EndCursor)
	}

	return nil, nil
}
