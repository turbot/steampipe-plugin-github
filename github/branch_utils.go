package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractBranchFromHydrateItem(h *plugin.HydrateData) (models.Branch, error) {
	if branch, ok := h.Item.(models.Branch); ok {
		return branch, nil
	} else {
		return models.Branch{}, fmt.Errorf("unable to parse hydrate item %v as a Branch", h.Item)
	}
}

func appendBranchColumnIncludes(m *map[string]interface{}, cols []string) {
	protectionIncluded := githubv4.Boolean(slices.Contains(cols, "protected") || slices.Contains(cols, "branch_protection_rule"))

	(*m)["includeBranchProtectionRule"] = protectionIncluded
	(*m)["includeAllowsDeletions"] = protectionIncluded
	(*m)["includeAllowsForcePushes"] = protectionIncluded
	(*m)["includeBlocksCreations"] = protectionIncluded
	(*m)["includeCreator"] = protectionIncluded
	(*m)["includeBranchProtectionRuleId"] = protectionIncluded
	(*m)["includeDismissesStaleReviews"] = protectionIncluded
	(*m)["includeIsAdminEnforced"] = protectionIncluded
	(*m)["includeLockAllowsFetchAndMerge"] = protectionIncluded
	(*m)["includeLockBranch"] = protectionIncluded
	(*m)["includePattern"] = protectionIncluded
	(*m)["includeRequireLastPushApproval"] = protectionIncluded
	(*m)["includeRequiredApprovingReviewCount"] = protectionIncluded
	(*m)["includeRequiredDeploymentEnvironments"] = protectionIncluded
	(*m)["includeRequiredStatusChecks"] = protectionIncluded
	(*m)["includeRequiresApprovingReviews"] = protectionIncluded
	(*m)["includeRequiresConversationResolution"] = protectionIncluded
	(*m)["includeRequiresCodeOwnerReviews"] = protectionIncluded
	(*m)["includeRequiresCommitSignatures"] = protectionIncluded
	(*m)["includeRequiresDeployments"] = protectionIncluded
	(*m)["includeRequiresLinearHistory"] = protectionIncluded
	(*m)["includeRequiresStatusChecks"] = protectionIncluded
	(*m)["includeRequiresStrictStatusChecks"] = protectionIncluded
	(*m)["includeRestrictsPushes"] = protectionIncluded
	(*m)["includeRestrictsReviewDismissals"] = protectionIncluded
	(*m)["includeMatchingBranches"] = protectionIncluded
}

func branchHydrateProtected(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	branch, err := extractBranchFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return branch.BranchProtectionRule.NodeId, nil
}

func branchHydrateBranchProtectionRule(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	branch, err := extractBranchFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return branch.BranchProtectionRule, nil
}

func extractBranchProtectionRuleFromHydrateItem(h *plugin.HydrateData) (branchProtectionRow, error) {
	if branchProtectionRule, ok := h.Item.(branchProtectionRow); ok {
		return branchProtectionRule, nil
	} else {
		return branchProtectionRow{}, fmt.Errorf("unable to parse hydrate item %v as a BranchProtectionRule", h.Item)
	}
}

func appendBranchProtectionRuleColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeAllowsDeletions"] = githubv4.Boolean(slices.Contains(cols, "allows_deletions"))
	(*m)["includeAllowsForcePushes"] = githubv4.Boolean(slices.Contains(cols, "allows_force_pushes"))
	(*m)["includeBlocksCreations"] = githubv4.Boolean(slices.Contains(cols, "blocks_creations"))
	(*m)["includeCreator"] = githubv4.Boolean(slices.Contains(cols, "creator") || slices.Contains(cols, "creator_login"))
	(*m)["includeBranchProtectionRuleId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeDismissesStaleReviews"] = githubv4.Boolean(slices.Contains(cols, "dismisses_stale_reviews"))
	(*m)["includeIsAdminEnforced"] = githubv4.Boolean(slices.Contains(cols, "is_admin_enforced"))
	(*m)["includeLockAllowsFetchAndMerge"] = githubv4.Boolean(slices.Contains(cols, "lock_allows_fetch_and_merge"))
	(*m)["includeLockBranch"] = githubv4.Boolean(slices.Contains(cols, "lock_branch"))
	(*m)["includePattern"] = githubv4.Boolean(slices.Contains(cols, "pattern"))
	(*m)["includeRequireLastPushApproval"] = githubv4.Boolean(slices.Contains(cols, "require_last_push_approval"))
	(*m)["includeRequiredApprovingReviewCount"] = githubv4.Boolean(slices.Contains(cols, "required_approving_review_count"))
	(*m)["includeRequiredDeploymentEnvironments"] = githubv4.Boolean(slices.Contains(cols, "required_deployment_environments"))
	(*m)["includeRequiredStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "required_status_checks"))
	(*m)["includeRequiresApprovingReviews"] = githubv4.Boolean(slices.Contains(cols, "requires_approving_reviews"))
	(*m)["includeRequiresConversationResolution"] = githubv4.Boolean(slices.Contains(cols, "requires_conversation_resolution"))
	(*m)["includeRequiresCodeOwnerReviews"] = githubv4.Boolean(slices.Contains(cols, "requires_code_owner_reviews"))
	(*m)["includeRequiresCommitSignatures"] = githubv4.Boolean(slices.Contains(cols, "requires_commit_signatures"))
	(*m)["includeRequiresDeployments"] = githubv4.Boolean(slices.Contains(cols, "requires_deployments"))
	(*m)["includeRequiresLinearHistory"] = githubv4.Boolean(slices.Contains(cols, "requires_linear_history"))
	(*m)["includeRequiresStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "requires_status_checks"))
	(*m)["includeRequiresStrictStatusChecks"] = githubv4.Boolean(slices.Contains(cols, "requires_strict_status_checks"))
	(*m)["includeRestrictsPushes"] = githubv4.Boolean(slices.Contains(cols, "restricts_pushes"))
	(*m)["includeRestrictsReviewDismissals"] = githubv4.Boolean(slices.Contains(cols, "restricts_review_dismissals"))
	(*m)["includeMatchingBranches"] = githubv4.Boolean(slices.Contains(cols, "matching_branches"))
}

func branchProtectionRuleHydrateAllowsDeletions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.AllowsDeletions, nil
}

func branchProtectionRuleHydrateAllowsForcePushes(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.AllowsForcePushes, nil
}

func branchProtectionRuleHydrateBlocksCreations(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.BlocksCreations, nil
}

func branchProtectionRuleHydrateCreatorLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.CreatorLogin, nil
}

func branchProtectionRuleHydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.ID, nil
}

func branchProtectionRuleHydrateDismissesStaleReviews(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.DismissesStaleReviews, nil
}

func branchProtectionRuleHydrateIsAdminEnforced(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.IsAdminEnforced, nil
}

func branchProtectionRuleHydrateLockAllowsFetchAndMerge(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.LockAllowsFetchAndMerge, nil
}

func branchProtectionRuleHydrateLockBranch(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.LockBranch, nil
}

func branchProtectionRuleHydratePattern(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.Pattern, nil
}

func branchProtectionRuleHydrateRequireLastPushApproval(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequireLastPushApproval, nil
}

func branchProtectionRuleHydrateRequiredApprovingReviewCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiredApprovingReviewCount, nil
}

func branchProtectionRuleHydrateRequiredDeploymentEnvironments(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiredDeploymentEnvironments, nil
}

func branchProtectionRuleHydrateRequiredStatusChecks(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiredStatusChecks, nil
}

func branchProtectionRuleHydrateRequiresApprovingReviews(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresApprovingReviews, nil
}

func branchProtectionRuleHydrateRequiresConversationResolution(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresConversationResolution, nil
}

func branchProtectionRuleHydrateRequiresCodeOwnerReviews(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresCodeOwnerReviews, nil
}

func branchProtectionRuleHydrateRequiresCommitSignatures(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresCommitSignatures, nil
}

func branchProtectionRuleHydrateRequiresDeployments(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresDeployments, nil
}

func branchProtectionRuleHydrateRequiresLinearHistory(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresLinearHistory, nil
}

func branchProtectionRuleHydrateRequiresStatusChecks(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresStatusChecks, nil
}

func branchProtectionRuleHydrateRequiresStrictStatusChecks(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RequiresStrictStatusChecks, nil
}

func branchProtectionRuleHydrateRestrictsPushes(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RestrictsPushes, nil
}

func branchProtectionRuleHydrateRestrictsReviewDismissals(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.RestrictsReviewDismissals, nil
}

func branchProtectionRuleHydrateMatchingBranchesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	brp, err := extractBranchProtectionRuleFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return brp.MatchingBranches, nil
}
