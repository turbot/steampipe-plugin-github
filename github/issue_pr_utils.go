package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractPRReviewFromHydrateItem(h *plugin.HydrateData) (models.PullRequestReview, error) {
	if prReview, ok := h.Item.(models.PullRequestReview); ok {
		return prReview, nil
	} else {
		return models.PullRequestReview{}, fmt.Errorf("unable to parse hydrate item %v as an PullRequestReview", h.Item)
	}
}
func appendPRReviewColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includePRReviewAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includePRReviewId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includePRReviewNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includePRReviewAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includePRReviewAuthorCanPushToRepository"] = githubv4.Boolean(slices.Contains(cols, "author_can_push_to_repository"))
	(*m)["includePRReviewState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includePRReviewBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includePRReviewUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includePRReviewSubmittedAt"] = githubv4.Boolean(slices.Contains(cols, "submitted_at"))
}

func prReviewHydrateAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.Author, nil
}

func prReviewHydrateAuthorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.Author.Login, nil
}

func prReviewHydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.Id, nil
}

func prReviewHydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.NodeId, nil
}

func prReviewHydrateAuthorAssociation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.AuthorAssociation, nil
}

func prReviewHydrateAuthorCanPushToRepository(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.AuthorCanPushToRepository, nil
}

func prReviewHydrateBody(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.Body, nil
}

func prReviewHydrateState(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.State, nil
}

func prReviewHydrateSubmittedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.SubmittedAt, nil
}

func prReviewHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	prReview, err := extractPRReviewFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return prReview.Url, nil
}

func extractIssueFromHydrateItem(h *plugin.HydrateData) (models.Issue, error) {
	if issue, ok := h.Item.(models.Issue); ok {
		return issue, nil
	} else if searchResult, ok := h.Item.(models.SearchIssueResult); ok {
		return searchResult.Node.Issue, nil
	} else {
		return models.Issue{}, fmt.Errorf("unable to parse hydrate item %v as an Issue", h.Item)
	}
}

func appendIssueColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeIssueAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeIssueBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueEditor"] = githubv4.Boolean(slices.Contains(cols, "editor"))
	(*m)["includeIssueMilestone"] = githubv4.Boolean(slices.Contains(cols, "milestone"))
	(*m)["includeIssueViewer"] = githubv4.Boolean(slices.Contains(cols, "user_can_close") ||
		slices.Contains(cols, "user_can_react") ||
		slices.Contains(cols, "user_can_reopen") ||
		slices.Contains(cols, "user_can_subscribe") ||
		slices.Contains(cols, "user_can_update") ||
		slices.Contains(cols, "user_cannot_update_reasons") ||
		slices.Contains(cols, "user_did_author") ||
		slices.Contains(cols, "user_subscription"))
	(*m)["includeIssueAssignees"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count") || slices.Contains(cols, "assignees"))
	(*m)["includeIssueCommentCount"] = githubv4.Boolean(slices.Contains(cols, "comments_total_count"))
	(*m)["includeIssueLabels"] = githubv4.Boolean(slices.Contains(cols, "labels") ||
		slices.Contains(cols, "labels_src") ||
		slices.Contains(cols, "labels_total_count"))
	(*m)["includeIssueUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeIssueUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeIssueTitle"] = githubv4.Boolean(slices.Contains(cols, "title"))
	(*m)["includeIssueStateReason"] = githubv4.Boolean(slices.Contains(cols, "state_reason"))
	(*m)["includeIssueState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includeIssuePublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includeIssueLocked"] = githubv4.Boolean(slices.Contains(cols, "locked"))
	(*m)["includeIssueLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includeIssueIsPinned"] = githubv4.Boolean(slices.Contains(cols, "is_pinned"))
	(*m)["includeIssueIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includeIssueFullDatabaseId"] = githubv4.Boolean(slices.Contains(cols, "full_database_id"))
	(*m)["includeIssueCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includeIssueCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeIssueClosedAt"] = githubv4.Boolean(slices.Contains(cols, "closed_at"))
	(*m)["includeIssueClosed"] = githubv4.Boolean(slices.Contains(cols, "closed"))
	(*m)["includeIssueBodyUrl"] = githubv4.Boolean(slices.Contains(cols, "body_url"))
	(*m)["includeIssueAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includeIssueActiveLockReason"] = githubv4.Boolean(slices.Contains(cols, "active_lock_reason"))
	(*m)["includeIssueNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeIssueId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeIssueIsReadByUser"] = githubv4.Boolean(slices.Contains(cols, "is_read_by_user"))
}

func issueHydrateIsReadByUser(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.IsReadByUser, nil
}

func issueHydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Id, nil
}

func issueHydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.NodeId, nil
}

func issueHydrateActiveLockReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.ActiveLockReason, nil
}

func issueHydrateAuthorAssociation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.AuthorAssociation, nil
}

func issueHydrateBodyUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.BodyUrl, nil
}

func issueHydrateClosed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Closed, nil
}

func issueHydrateClosedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.ClosedAt, nil
}

func issueHydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.CreatedAt, nil
}

func issueHydrateCreatedViaEmail(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.CreatedViaEmail, nil
}

func issueHydrateFullDatabaseId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.FullDatabaseId, nil
}

func issueHydrateIncludesCreatedEdit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.IncludesCreatedEdit, nil
}

func issueHydrateIsPinned(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.IsPinned, nil
}

func issueHydrateLastEditedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.LastEditedAt, nil
}

func issueHydrateLocked(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Locked, nil
}

func issueHydratePublishedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.PublishedAt, nil
}

func issueHydrateState(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.State, nil
}

func issueHydrateStateReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.StateReason, nil
}

func issueHydrateTitle(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Title, nil
}

func issueHydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UpdatedAt, nil
}

func issueHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Url, nil
}

func issueHydrateAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Author, nil
}

func issueHydrateAuthorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Author.Login, nil
}

func issueHydrateBody(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Body, nil
}

func issueHydrateEditor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Editor, nil
}

func issueHydrateMilestone(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Milestone, nil
}

func issueHydrateUserCanClose(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCanClose, nil
}

func issueHydrateUserCanReact(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCanReact, nil
}

func issueHydrateUserCanReopen(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCanReopen, nil
}

func issueHydrateUserCanSubscribe(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCanSubscribe, nil
}

func issueHydrateUserCanUpdate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCanUpdate, nil
}

func issueHydrateUserCannotUpdateReasons(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserCannotUpdateReasons, nil
}

func issueHydrateUserDidAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserDidAuthor, nil
}

func issueHydrateUserSubscription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.UserSubscription, nil
}

func issueHydrateAssigneeCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Assignees.TotalCount, nil
}

func issueHydrateAssignees(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Assignees.Nodes, nil
}

func issueHydrateCommentCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Comments.TotalCount, nil
}

func issueHydrateLabelsCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Labels.TotalCount, nil
}

func issueHydrateLabels(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issue, err := extractIssueFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issue.Labels.Nodes, nil
}

func extractIssueCommentFromHydrateItem(h *plugin.HydrateData) (models.IssueComment, error) {
	if issueComment, ok := h.Item.(models.IssueComment); ok {
		return issueComment, nil
	} else {
		return models.IssueComment{}, fmt.Errorf("unable to parse hydrate item %v as an IssueComment", h.Item)
	}
}

func appendIssuePRCommentColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeIssueCommentAuthor"] = githubv4.Boolean(slices.Contains(cols, "author") || slices.Contains(cols, "author_login"))
	(*m)["includeIssueCommentBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueCommentEditor"] = githubv4.Boolean(slices.Contains(cols, "editor") || slices.Contains(cols, "editor_login"))
	(*m)["includeIssueCommentViewer"] = githubv4.Boolean(slices.Contains(cols, "can_delete") ||
		slices.Contains(cols, "can_react") ||
		slices.Contains(cols, "can_minimize") ||
		slices.Contains(cols, "can_update") ||
		slices.Contains(cols, "cannot_update_reasons") ||
		slices.Contains(cols, "did_author"))
	(*m)["includeIssueCommentUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeIssueCommentUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeIssueCommentPublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includeIssueCommentMinimizedReason"] = githubv4.Boolean(slices.Contains(cols, "minimized_reason"))
	(*m)["includeIssueCommentLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includeIssueCommentIsMinimized"] = githubv4.Boolean(slices.Contains(cols, "is_minimized"))
	(*m)["includeIssueCommentIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includeIssueCommentCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includeIssueCommentCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeIssueCommentBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includeIssueCommentBodyText"] = githubv4.Boolean(slices.Contains(cols, "body_text"))
	(*m)["includeIssueCommentAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includeIssueCommentNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeIssueCommentId"] = githubv4.Boolean(slices.Contains(cols, "id"))
}

func issueCommentHydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Id, nil
}

func issueCommentHydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.NodeId, nil
}

func issueCommentHydrateAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Author, nil
}

func issueCommentHydrateAuthorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Author.Login, nil
}

func issueCommentHydrateAuthorAssociation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.AuthorAssociation, nil
}

func issueCommentHydrateBody(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Body, nil
}

func issueCommentHydrateBodyText(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.BodyText, nil
}

func issueCommentHydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CreatedAt, nil
}

func issueCommentHydrateCreatedViaEmail(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CreatedViaEmail, nil
}

func issueCommentHydrateEditor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Editor, nil
}

func issueCommentHydrateEditorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Editor.Login, nil
}

func issueCommentHydrateIncludesCreatedEdit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.IncludesCreatedEdit, nil
}

func issueCommentHydrateIsMinimized(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.IsMinimized, nil
}

func issueCommentHydrateMinimizedReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.MinimizedReason, nil
}

func issueCommentHydrateLastEditedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.LastEditedAt, nil
}

func issueCommentHydratePublishedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.PublishedAt, nil
}

func issueCommentHydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.UpdatedAt, nil
}

func issueCommentHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.Url, nil
}

func issueCommentHydrateCanDelete(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CanDelete, nil
}

func issueCommentHydrateCanMinimize(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CanMinimize, nil
}

func issueCommentHydrateCanReact(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CanReact, nil
}

func issueCommentHydrateCanUpdate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CanUpdate, nil
}

func issueCommentHydrateCannotUpdateReasons(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.CannotUpdateReasons, nil
}

func issueCommentHydrateDidAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	issueComment, err := extractIssueCommentFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return issueComment.DidAuthor, nil
}

func extractPullRequestFromHydrateItem(h *plugin.HydrateData) (models.PullRequest, error) {
	if pr, ok := h.Item.(models.PullRequest); ok {
		return pr, nil
	} else if sr, ok := h.Item.(models.SearchPullRequestResult); ok {
		return sr.Node.PullRequest, nil
	}
	return models.PullRequest{}, fmt.Errorf("unable to parse hydrate item %v as a PullRequest", h.Item)
}

func appendPullRequestColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includePRAuthor"] = githubv4.Boolean(slices.Contains(cols, "author"))
	(*m)["includePRBody"] = githubv4.Boolean(slices.Contains(cols, "body"))
	(*m)["includePREditor"] = githubv4.Boolean(slices.Contains(cols, "editor"))
	(*m)["includePRMergedBy"] = githubv4.Boolean(slices.Contains(cols, "merged_by"))
	(*m)["includePRMilestone"] = githubv4.Boolean(slices.Contains(cols, "milestone"))

	(*m)["includePRBaseRef"] = githubv4.Boolean(slices.Contains(cols, "base_ref"))
	(*m)["includePRHeadRef"] = githubv4.Boolean(slices.Contains(cols, "head_ref"))
	(*m)["includePRMergeCommit"] = githubv4.Boolean(slices.Contains(cols, "merge_commit"))
	(*m)["includePRSuggested"] = githubv4.Boolean(slices.Contains(cols, "suggested_reviewers"))
	(*m)["includePRViewer"] = githubv4.Boolean(slices.Contains(cols, "can_apply_suggestion") ||
		slices.Contains(cols, "can_close") ||
		slices.Contains(cols, "can_delete_head_ref") ||
		slices.Contains(cols, "can_disable_auto_merge") ||
		slices.Contains(cols, "can_edit_files") ||
		slices.Contains(cols, "can_enable_auto_merge") ||
		slices.Contains(cols, "can_react") ||
		slices.Contains(cols, "can_reopen") ||
		slices.Contains(cols, "can_subscribe") ||
		slices.Contains(cols, "can_update") ||
		slices.Contains(cols, "can_update_branch") ||
		slices.Contains(cols, "did_author") ||
		slices.Contains(cols, "cannot_update_reasons") ||
		slices.Contains(cols, "subscription"))
	(*m)["includePRAssignees"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count") || slices.Contains(cols, "assignees"))
	(*m)["includePRCommitCount"] = githubv4.Boolean(slices.Contains(cols, "commits_total_count"))
	(*m)["includePRReviewRequestCount"] = githubv4.Boolean(slices.Contains(cols, "review_requests_total_count"))
	(*m)["includePRReviewCount"] = githubv4.Boolean(slices.Contains(cols, "reviews_total_count"))
	(*m)["includePRLabels"] = githubv4.Boolean(slices.Contains(cols, "labels") ||
		slices.Contains(cols, "labels_src") ||
		slices.Contains(cols, "labels_total_count"))
	(*m)["includePRId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includePRNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includePRAuthorAssociation"] = githubv4.Boolean(slices.Contains(cols, "author_association"))
	(*m)["includePRBaseRefName"] = githubv4.Boolean(slices.Contains(cols, "base_ref_name"))
	(*m)["includePRActiveLockReason"] = githubv4.Boolean(slices.Contains(cols, "active_lock_reason"))
	(*m)["includePRAdditions"] = githubv4.Boolean(slices.Contains(cols, "additions"))
	(*m)["includePRChangedFiles"] = githubv4.Boolean(slices.Contains(cols, "changed_files"))
	(*m)["includePRChecksUrl"] = githubv4.Boolean(slices.Contains(cols, "checks_url"))
	(*m)["includePRClosed"] = githubv4.Boolean(slices.Contains(cols, "closed"))
	(*m)["includePRClosedAt"] = githubv4.Boolean(slices.Contains(cols, "closed_at"))
	(*m)["includePRCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includePRCreatedViaEmail"] = githubv4.Boolean(slices.Contains(cols, "created_via_email"))
	(*m)["includePRDeletions"] = githubv4.Boolean(slices.Contains(cols, "deletions"))
	(*m)["includePRHeadRefName"] = githubv4.Boolean(slices.Contains(cols, "head_ref_name"))
	(*m)["includePRHeadRefOid"] = githubv4.Boolean(slices.Contains(cols, "head_ref_oid"))
	(*m)["includePRIncludesCreatedEdit"] = githubv4.Boolean(slices.Contains(cols, "includes_created_edit"))
	(*m)["includePRIsCrossRepository"] = githubv4.Boolean(slices.Contains(cols, "is_cross_repository"))
	(*m)["includePRIsDraft"] = githubv4.Boolean(slices.Contains(cols, "is_draft"))
	(*m)["includePRIsReadByUser"] = githubv4.Boolean(slices.Contains(cols, "is_read_by_user"))
	(*m)["includePRLastEditedAt"] = githubv4.Boolean(slices.Contains(cols, "last_edited_at"))
	(*m)["includePRLocked"] = githubv4.Boolean(slices.Contains(cols, "locked"))
	(*m)["includePRMaintainerCanModify"] = githubv4.Boolean(slices.Contains(cols, "maintainer_can_modify"))
	(*m)["includePRMergeable"] = githubv4.Boolean(slices.Contains(cols, "mergeable"))
	(*m)["includePRMerged"] = githubv4.Boolean(slices.Contains(cols, "merged"))
	(*m)["includePRMergedAt"] = githubv4.Boolean(slices.Contains(cols, "merged_at"))
	(*m)["includePRPermalink"] = githubv4.Boolean(slices.Contains(cols, "permalink"))
	(*m)["includePRPublishedAt"] = githubv4.Boolean(slices.Contains(cols, "published_at"))
	(*m)["includePRRevertUrl"] = githubv4.Boolean(slices.Contains(cols, "revert_url"))
	(*m)["includePRReviewDecision"] = githubv4.Boolean(slices.Contains(cols, "review_decision"))
	(*m)["includePRState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includePRTitle"] = githubv4.Boolean(slices.Contains(cols, "title"))
	(*m)["includePRTotalCommentsCount"] = githubv4.Boolean(slices.Contains(cols, "total_comments_count"))
	(*m)["includePRUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includePRUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
}

func prHydrateAuthorAssociation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.AuthorAssociation, nil
}

func prHydrateBaseRefName(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.BaseRefName, nil
}

func prHydrateMaintainerCanModify(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MaintainerCanModify, nil
}

func prHydrateMergedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MergedAt, nil
}

func prHydrateMergeable(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Mergeable, nil
}

func prHydrateMerged(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Merged, nil
}

func prHydratePermalink(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Permalink, nil
}

func prHydratePublishedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.PublishedAt, nil
}

func prHydrateRevertUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.RevertUrl, nil
}

func prHydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Id, nil
}

func prHydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.NodeId, nil
}

func prHydrateAdditions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Additions, nil
}

func prHydrateChangedFiles(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ChangedFiles, nil
}

func prHydrateChecksUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ChecksUrl, nil
}

func prHydrateClosed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Closed, nil
}

func prHydrateClosedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ClosedAt, nil
}

func prHydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CreatedAt, nil
}

func prHydrateCreatedViaEmail(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CreatedViaEmail, nil
}

func prHydrateDeletions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Deletions, nil
}

func prHydrateHeadRefOid(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.HeadRefOid, nil
}

func prHydrateIncludesCreatedEdit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IncludesCreatedEdit, nil
}

func prHydrateIsCrossRepository(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsCrossRepository, nil
}

func prHydrateIsDraft(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsDraft, nil
}

func prHydrateIsReadByUser(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsReadByUser, nil
}

func prHydrateLastEditedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.LastEditedAt, nil
}

func prHydrateLocked(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Locked, nil
}

func prHydrateReviewDecision(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ReviewDecision, nil
}

func prHydrateState(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.State, nil
}

func prHydrateTitle(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Title, nil
}

func prHydrateTotalCommentsCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.TotalCommentsCount, nil
}

func prHydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.UpdatedAt, nil
}

func prHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Url, nil
}

func prHydrateHeadRefName(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.HeadRefName, nil
}

func prHydrateActiveLockReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ActiveLockReason, nil
}

func prHydrateAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Author, nil
}

func prHydrateBody(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Body, nil
}

func prHydrateEditor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Editor, nil
}

func prHydrateMergeBy(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MergedBy, nil
}

func prHydrateMilestone(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Milestone, nil
}

func prHydrateBaseRef(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.BaseRef, nil
}

func prHydrateHeadRef(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.HeadRef, nil
}

func prHydrateMergeCommit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MergeCommit, nil
}

func prHydrateSuggested(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.SuggestedReviewers, nil
}

func prHydrateLabels(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Labels.Nodes, nil
}

func prHydrateLabelsCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Labels.TotalCount, nil
}

func prHydrateAssigneeCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Assignees.TotalCount, nil
}

func prHydrateAssignees(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Assignees.Nodes, nil
}

func prHydrateCommitCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Commits.TotalCount, nil
}

func prHydrateReviewRequestCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ReviewRequests.TotalCount, nil
}

func prHydrateReviewCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Reviews.TotalCount, nil
}

func prHydrateCanApplySuggestion(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanApplySuggestion, nil
}

func prHydrateCanClose(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanClose, nil
}

func prHydrateCanDeleteHeadRef(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanDeleteHeadRef, nil
}

func prHydrateCanDisableAutoMerge(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanDisableAutoMerge, nil
}

func prHydrateCanEditFiles(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanEditFiles, nil
}

func prHydrateCanEnableAutoMerge(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanEnableAutoMerge, nil
}

func prHydrateCanMergeAsAdmin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanMergeAsAdmin, nil
}

func prHydrateCanReact(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanReact, nil
}

func prHydrateCanReopen(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanReopen, nil
}

func prHydrateCanSubscribe(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanSubscribe, nil
}

func prHydrateCanUpdate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanUpdate, nil
}

func prHydrateCanUpdateBranch(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CanUpdateBranch, nil
}

func prHydrateDidAuthor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.DidAuthor, nil
}

func prHydrateCannotUpdateReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CannotUpdateReasons, nil
}

func prHydrateSubscription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Subscription, nil
}
