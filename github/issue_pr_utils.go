package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

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
	(*m)["includeIssueAssigneeCount"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count"))
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
	(*m)["includePRAssigneeCount"] = githubv4.Boolean(slices.Contains(cols, "assignees_total_count"))
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

func prHydrateAuthorAssociation(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.AuthorAssociation, nil
}

func prHydrateBaseRefName(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.BaseRefName, nil
}

func prHydrateMaintainerCanModify(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MaintainerCanModify, nil
}

func prHydrateMergedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.MergedAt, nil
}

func prHydrateMergeable(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Mergeable, nil
}

func prHydrateMerged(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Merged, nil
}

func prHydratePermalink(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Permalink, nil
}

func prHydratePublishedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.PublishedAt, nil
}

func prHydrateRevertUrl(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.RevertUrl, nil
}

func prHydrateId(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Id, nil
}

func prHydrateNodeId(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.NodeId, nil
}

func prHydrateAdditions(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Additions, nil
}

func prHydrateChangedFiles(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ChangedFiles, nil
}

func prHydrateChecksUrl(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ChecksUrl, nil
}

func prHydrateClosed(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Closed, nil
}

func prHydrateClosedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ClosedAt, nil
}

func prHydrateCreatedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CreatedAt, nil
}

func prHydrateCreatedViaEmail(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.CreatedViaEmail, nil
}

func prHydrateDeletions(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Deletions, nil
}

func prHydrateHeadRefOid(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.HeadRefOid, nil
}

func prHydrateIncludesCreatedEdit(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IncludesCreatedEdit, nil
}

func prHydrateIsCrossRepository(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsCrossRepository, nil
}

func prHydrateIsDraft(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsDraft, nil
}

func prHydrateIsReadByUser(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.IsReadByUser, nil
}

func prHydrateLastEditedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.LastEditedAt, nil
}

func prHydrateLocked(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Locked, nil
}

func prHydrateReviewDecision(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.ReviewDecision, nil
}

func prHydrateState(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.State, nil
}

func prHydrateTitle(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Title, nil
}

func prHydrateTotalCommentsCount(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.TotalCommentsCount, nil
}

func prHydrateUpdatedAt(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.UpdatedAt, nil
}

func prHydrateUrl(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.Url, nil
}

func prHydrateHeadRefName(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	pr, err := extractPullRequestFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return pr.HeadRefName, nil
}

func prHydrateActiveLockReason(ctx context.Context, queryData *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
