package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func sharedCommentsColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromQual("number"), Description: "The issue/pr number."},

		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromValue(), Hydrate: issueCommentHydrateId, Description: "The ID of the comment."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateNodeId, Description: "The node ID of the comment."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: issueCommentHydrateAuthor, Description: "The actor who authored the comment."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateAuthorLogin, Description: "The login of the comment author."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateAuthorAssociation, Description: "Author's association with the subject of the issue/pr the comment was raised on."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateBody, Description: "The contents of the comment as markdown."},
		{Name: "body_text", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateBodyText, Description: "The contents of the comment as text."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: issueCommentHydrateCreatedAt, Description: "Timestamp when comment was created."},
		{Name: "created_via_email", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateCreatedViaEmail, Description: "If true, comment was created via email."},
		{Name: "editor", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: issueCommentHydrateEditor, Description: "The actor who edited the comment."},
		{Name: "editor_login", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateEditorLogin, Description: "The login of the comment editor."},
		{Name: "includes_created_edit", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateIncludesCreatedEdit, Description: "If true, comment was edited and includes an edit with the creation data."},
		{Name: "is_minimized", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateIsMinimized, Description: "If true, comment has been minimized."},
		{Name: "minimized_reason", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateMinimizedReason, Description: "The reason for comment being minimized."},
		{Name: "last_edited_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: issueCommentHydrateLastEditedAt, Description: "Timestamp when comment was last edited."},
		{Name: "published_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: issueCommentHydratePublishedAt, Description: "Timestamp when comment was published."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Hydrate: issueCommentHydrateUpdatedAt, Description: "Timestamp when comment was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: issueCommentHydrateUrl, Description: "URL for the comment."},
		{Name: "can_delete", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateCanDelete, Description: "If true, user can delete the comment."},
		{Name: "can_minimize", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateCanMinimize, Description: "If true, user can minimize the comment."},
		{Name: "can_react", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateCanReact, Description: "If true, user can react to the comment."},
		{Name: "can_update", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateCanUpdate, Description: "If true, user can update the comment."},
		{Name: "cannot_update_reasons", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: issueCommentHydrateCannotUpdateReasons, Description: "A list of reasons why user cannot update the comment."},
		{Name: "did_author", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: issueCommentHydrateDidAuthor, Description: "If true, user authored the comment."},
	}
}

func tableGitHubIssueComment() *plugin.Table {
	return &plugin.Table{
		Name:        "github_issue_comment",
		Description: "GitHub Issue Comments are the responses/comments on GitHub Issues.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryIssueCommentList,
		},
		Columns: sharedCommentsColumns(),
	}
}

func tableGitHubRepositoryIssueCommentList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	issueNumber := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Issue struct {
				Comments struct {
					PageInfo   models.PageInfo
					TotalCount int
					Nodes      []models.IssueComment
				} `graphql:"comments(first: $pageSize, after: $cursor)"`
			} `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(repoName),
		"issueNumber": githubv4.Int(issueNumber),
		"pageSize":    githubv4.Int(pageSize),
		"cursor":      (*githubv4.String)(nil),
	}
	appendIssuePRCommentColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_issue_comment", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_issue_comment", "api_error", err)
			return nil, err
		}

		for _, comment := range query.Repository.Issue.Comments.Nodes {
			d.StreamListItem(ctx, comment)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Issue.Comments.PageInfo.EndCursor)
	}

	return nil, nil
}
