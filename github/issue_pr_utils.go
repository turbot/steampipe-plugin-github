package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"slices"
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
