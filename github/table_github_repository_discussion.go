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

func gitHubRepositoryDiscussionColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id"), Description: "The ID of the discussion."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId"), Description: "The node ID of the discussion."},
		{Name: "number", Type: proto.ColumnType_INT, Description: "The discussion number."},
		{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the discussion."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url"), Description: "The URL of the discussion."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the discussion was created."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the discussion was last updated."},
		{Name: "author", Type: proto.ColumnType_JSON, Description: "The actor who authored the discussion."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.Login"), Description: "The login of the discussion author."},
		{Name: "category", Type: proto.ColumnType_JSON, Description: "The category of the discussion."},
		{Name: "category_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Category.Name"), Description: "The name of the discussion category."},
		{Name: "answer", Type: proto.ColumnType_JSON, Description: "The answer to the discussion, if any."},
		{Name: "comments_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Comments.TotalCount"), Description: "Total count of comments on the discussion."},
		{Name: "comments", Type: proto.ColumnType_JSON, Hydrate: discussionHydrateComments, Transform: transform.FromValue().NullIfZero(), Description: "All comments on the discussion."},
		{Name: "replies", Type: proto.ColumnType_JSON, Hydrate: discussionHydrateReplies, Transform: transform.FromValue().NullIfZero(), Description: "All replies to comments on the discussion."},
	}
}

func tableGitHubRepositoryDiscussion() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_discussion",
		Description: "GitHub Discussions are conversations that can be started by anyone and are organized into categories.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryDiscussionList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryDiscussionGet,
		},
		Columns: commonColumns(gitHubRepositoryDiscussionColumns()),
	}
}

//// LIST FUNCTION

func tableGitHubRepositoryDiscussionList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	// First, let's check if the repository has discussions enabled
	var checkQuery struct {
		RateLimit  models.RateLimit
		Repository struct {
			Id                    int `graphql:"id: databaseId"`
			Name                  string
			HasDiscussionsEnabled bool `graphql:"hasDiscussionsEnabled"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	checkVariables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repoName),
	}

	client := connectV4(ctx, d)

	// Check repository and discussions status
	err := client.Query(ctx, &checkQuery, checkVariables)
	if err != nil {
		plugin.Logger(ctx).Error("github_repository_discussion", "check_api_error", err, "repository", fullName)
		return nil, err
	}

	if !checkQuery.Repository.HasDiscussionsEnabled {
		plugin.Logger(ctx).Info("github_repository_discussion", "discussions_not_enabled", "repository", fullName)
		return nil, nil
	}

	// Now query discussions with minimal fields
	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Discussions struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.Discussion
			} `graphql:"discussions(first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_discussion", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_discussion", "api_error", err, "repository", fullName)
			return nil, err
		}

		for _, discussion := range query.Repository.Discussions.Nodes {
			d.StreamListItem(ctx, discussion)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Discussions.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Discussions.PageInfo.EndCursor)
	}

	return nil, nil
}

/// HYDRATE FUNCTION

func tableGitHubRepositoryDiscussionGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	discussionNumber := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Discussion models.Discussion `graphql:"discussion(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repoName),
		"number": githubv4.Int(discussionNumber),
	}

	client := connectV4(ctx, d)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_discussion", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_repository_discussion", "api_error", err)
		return nil, err
	}

	return query.Repository.Discussion, nil
}

func discussionHydrateComments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	discussion := h.Item.(models.Discussion)

	// check if the the LIST API already fetched all the comments
	if discussion.Comments.TotalCount == len(discussion.Comments.Nodes) {
		return discussion.Comments.Nodes, nil
	}

	// Parse repository info from the discussion URL
	// URL format: https://github.com/owner/repo/discussions/123
	urlParts := strings.Split(discussion.Url, "/")
	if len(urlParts) < 5 {
		return nil, fmt.Errorf("invalid discussion URL format")
	}

	owner := urlParts[3]
	repoName := urlParts[4]

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Discussion struct {
				Comments struct {
					PageInfo   models.PageInfo
					TotalCount int
					Nodes      []models.DiscussionComment
				} `graphql:"comments(first: $pageSize, after: $cursor)"`
			} `graphql:"discussion(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"number":   githubv4.Int(discussion.Number),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	client := connectV4(ctx, d)
	var allComments []models.DiscussionComment

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_discussion_comments", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_repository_discussion_comments", "api_error", err)
			return nil, err
		}

		allComments = append(allComments, query.Repository.Discussion.Comments.Nodes...)

		if !query.Repository.Discussion.Comments.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Discussion.Comments.PageInfo.EndCursor)
	}

	return allComments, nil
}

func discussionHydrateReplies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	discussion := h.Item.(models.Discussion)
	discussionNumber := discussion.Number

	var commentNodeIds []string
	// Check if the LIST API already fetched all comments to collect all node IDs from the discussion
	if discussion.Comments.TotalCount == len(discussion.Comments.Nodes) {
		for _, nodeId := range discussion.Comments.Nodes {
			commentNodeIds = append(commentNodeIds, nodeId.NodeId)
		}
	}

	// Parse repository info from the discussion URL
	// URL format: https://github.com/owner/repo/discussions/123
	urlParts := strings.Split(discussion.Url, "/")
	if len(urlParts) < 5 {
		return nil, fmt.Errorf("invalid discussion URL format")
	}
	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	client := connectV4(ctx, d)

	// if the commentNodeIds is empty, we need to fetch all the comment node IDs from the discussion
	// Step 1: Get all comment node IDs from the discussion
	if len(commentNodeIds) == 0 {
		owner := urlParts[3]
		repoName := urlParts[4]

		var commentsQuery struct {
			RateLimit  models.RateLimit
			Repository struct {
				Discussion struct {
					Comments struct {
						PageInfo   models.PageInfo
						TotalCount int
						Nodes      []struct {
							NodeId string `graphql:"nodeId: id"`
						}
					} `graphql:"comments(first: $pageSize, after: $cursor)"`
				} `graphql:"discussion(number: $number)"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}

		commentsVariables := map[string]interface{}{
			"owner":    githubv4.String(owner),
			"name":     githubv4.String(repoName),
			"number":   githubv4.Int(discussionNumber),
			"pageSize": githubv4.Int(pageSize),
			"cursor":   (*githubv4.String)(nil),
		}

		// Collect all comment node IDs
		for {
			err := client.Query(ctx, &commentsQuery, commentsVariables)
			plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_discussion_comments_for_replies", &commentsQuery.RateLimit))
			if err != nil {
				plugin.Logger(ctx).Error("github_repository_discussion_comments_for_replies", "api_error", err)
				return nil, err
			}

			for _, comment := range commentsQuery.Repository.Discussion.Comments.Nodes {
				commentNodeIds = append(commentNodeIds, comment.NodeId)
			}

			if !commentsQuery.Repository.Discussion.Comments.PageInfo.HasNextPage {
				break
			}
			commentsVariables["cursor"] = githubv4.NewString(commentsQuery.Repository.Discussion.Comments.PageInfo.EndCursor)
		}
	}
	// Step 2: Get replies for each comment node ID
	var allReplies []models.DiscussionComment

	for _, nodeId := range commentNodeIds {
		var repliesQuery struct {
			RateLimit models.RateLimit
			Node      struct {
				DiscussionComment struct {
					Replies struct {
						PageInfo   models.PageInfo
						TotalCount int
						Nodes      []models.DiscussionComment
					} `graphql:"replies(first: $pageSize, after: $cursor)"`
				} `graphql:"... on DiscussionComment"`
			} `graphql:"node(id: $nodeId)"`
		}

		repliesVariables := map[string]interface{}{
			"nodeId":   githubv4.ID(nodeId),
			"pageSize": githubv4.Int(pageSize),
			"cursor":   (*githubv4.String)(nil),
		}

		// Get replies for this comment with pagination
		for {
			err := client.Query(ctx, &repliesQuery, repliesVariables)
			plugin.Logger(ctx).Debug(rateLimitLogString("github_repository_discussion_replies", &repliesQuery.RateLimit))
			if err != nil {
				plugin.Logger(ctx).Error("github_repository_discussion_replies", "api_error", err)
				return nil, err
			}

			allReplies = append(allReplies, repliesQuery.Node.DiscussionComment.Replies.Nodes...)

			if !repliesQuery.Node.DiscussionComment.Replies.PageInfo.HasNextPage {
				break
			}
			repliesVariables["cursor"] = githubv4.NewString(repliesQuery.Node.DiscussionComment.Replies.PageInfo.EndCursor)
		}
	}

	return allReplies, nil
}
