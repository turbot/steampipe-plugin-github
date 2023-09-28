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
