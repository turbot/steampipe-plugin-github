package models

type Branch struct {
	Name   string
	Target struct {
		Commit Commit `graphql:"... on Commit"`
	}
	BranchProtectionRule BranchProtectionRuleForBranch `json:"branch_protection_rule"`
}

type BranchProtectionRuleForBranch struct {
	AllowsDeletions                bool     `json:"allows_deletions"`
	AllowsForcePushes              bool     `json:"allows_force_pushes"`
	BlocksCreations                bool     `json:"blocks_creations"`
	Creator                        Actor    `json:"creator"`
	Id                             int      `graphql:"id: databaseId" json:"id"`
	NodeId                         string   `graphql:"nodeId: id" json:"node_id"`
	DismissesStaleReviews          bool     `json:"dismisses_stale_reviews"`
	IsAdminEnforced                bool     `json:"is_admin_enforced"`
	LockAllowsFetchAndMerge        bool     `json:"lock_allows_fetch_and_merge"`
	LockBranch                     bool     `json:"lock_branch"`
	Pattern                        string   `json:"pattern"`
	RequireLastPushApproval        bool     `json:"require_last_push_approval"`
	RequiredApprovingReviewCount   int      `json:"required_approving_review_count"`
	RequiredDeploymentEnvironments []string `json:"required_deployment_environments"`
	RequiredStatusChecks           []string `graphql:"requiredStatusChecks: requiredStatusCheckContexts" json:"required_status_checks"`
	RequiresApprovingReviews       bool     `json:"requires_approving_reviews"`
	RequiresConversationResolution bool     `json:"requires_conversation_resolution"`
	RequiresCodeOwnerReviews       bool     `json:"requires_code_owner_reviews"`
	RequiresCommitSignatures       bool     `json:"requires_commit_signatures"`
	RequiresDeployments            bool     `json:"requires_deployments"`
	RequiresLinearHistory          bool     `json:"requires_linear_history"`
	RequiresStatusChecks           bool     `json:"requires_status_checks"`
	RequiresStrictStatusChecks     bool     `json:"requires_strict_status_checks"`
	RestrictsPushes                bool     `json:"restricts_pushes"`
	RestrictsReviewDismissals      bool     `json:"restricts_review_dismissals"`
	MatchingBranches               struct {
		TotalCount int `json:"total_count"`
	} `graphql:"matchingBranches: matchingRefs" json:"matching_branches"`
	// BranchProtectionRuleConflicts
}

type BranchProtectionRule struct {
	AllowsDeletions                bool     `graphql:"allowsDeletions @include(if:$includeAllowsDeletions)" json:"allows_deletions"`
	AllowsForcePushes              bool     `graphql:"allowsForcePushes @include(if:$includeAllowsForcePushes)" json:"allows_force_pushes"`
	BlocksCreations                bool     `graphql:"blocksCreations @include(if:$includeBlocksCreations)" json:"blocks_creations"`
	Creator                        Actor    `graphql:"creator @include(if:$includeCreator)" json:"creator"`
	Id                             int      `graphql:"id: databaseId @include(if:$includeBranchProtectionRuleId)" json:"id"`
	NodeId                         string   `graphql:"nodeId: id" json:"node_id"`
	DismissesStaleReviews          bool     `graphql:"dismissesStaleReviews @include(if:$includeDismissesStaleReviews)" json:"dismisses_stale_reviews"`
	IsAdminEnforced                bool     `graphql:"isAdminEnforced @include(if:$includeIsAdminEnforced)" json:"is_admin_enforced"`
	LockAllowsFetchAndMerge        bool     `graphql:"lockAllowsFetchAndMerge @include(if:$includeLockAllowsFetchAndMerge)" json:"lock_allows_fetch_and_merge"`
	LockBranch                     bool     `graphql:"lockBranch @include(if:$includeLockBranch)" json:"lock_branch"`
	Pattern                        string   `graphql:"pattern @include(if:$includePattern)" json:"pattern"`
	RequireLastPushApproval        bool     `graphql:"requireLastPushApproval @include(if:$includeRequireLastPushApproval)" json:"require_last_push_approval"`
	RequiredApprovingReviewCount   int      `graphql:"requiredApprovingReviewCount @include(if:$includeRequiredApprovingReviewCount)" json:"required_approving_review_count"`
	RequiredDeploymentEnvironments []string `graphql:"requiredDeploymentEnvironments @include(if:$includeRequiredDeploymentEnvironments)" json:"required_deployment_environments"`
	RequiredStatusChecks           []string `graphql:"requiredStatusChecks: requiredStatusCheckContexts @include(if:$includeRequiredStatusChecks)" json:"required_status_checks"`
	RequiresApprovingReviews       bool     `graphql:"requiresApprovingReviews @include(if:$includeRequiresApprovingReviews)" json:"requires_approving_reviews"`
	RequiresConversationResolution bool     `graphql:"requiresConversationResolution @include(if:$includeRequiresConversationResolution)" json:"requires_conversation_resolution"`
	RequiresCodeOwnerReviews       bool     `graphql:"requiresCodeOwnerReviews @include(if:$includeRequiresCodeOwnerReviews)" json:"requires_code_owner_reviews"`
	RequiresCommitSignatures       bool     `graphql:"requiresCommitSignatures @include(if:$includeRequiresCommitSignatures)" json:"requires_commit_signatures"`
	RequiresDeployments            bool     `graphql:"requiresDeployments @include(if:$includeRequiresDeployments)" json:"requires_deployments"`
	RequiresLinearHistory          bool     `graphql:"requiresLinearHistory @include(if:$includeRequiresLinearHistory)" json:"requires_linear_history"`
	RequiresStatusChecks           bool     `graphql:"requiresStatusChecks @include(if:$includeRequiresStatusChecks)" json:"requires_status_checks"`
	RequiresStrictStatusChecks     bool     `graphql:"requiresStrictStatusChecks @include(if:$includeRequiresStrictStatusChecks)" json:"requires_strict_status_checks"`
	RestrictsPushes                bool     `graphql:"restrictsPushes @include(if:$includeRestrictsPushes)" json:"restricts_pushes"`
	RestrictsReviewDismissals      bool     `graphql:"restrictsReviewDismissals @include(if:$includeRestrictsReviewDismissals)" json:"restricts_review_dismissals"`
	MatchingBranches               struct {
		TotalCount int `json:"total_count"`
	} `graphql:"matchingBranches: matchingRefs @include(if:$includeMatchingBranches)" json:"matching_branches"`
	// BranchProtectionRuleConflicts
}

type BranchProtectionRuleWithFirstPageEmbeddedItems struct {
	BranchProtectionRule
	PushAllowances              BranchActorAllowances `graphql:"pushAllowances(first: 100)"`
	BypassForcePushAllowances   BranchActorAllowances `graphql:"bypassForcePushAllowances(first: 100)"`
	BypassPullRequestAllowances BranchActorAllowances `graphql:"bypassPullRequestAllowances(first: 100)"`
}

type BranchProtectionRuleWithPushAllowances struct {
	BranchProtectionRule
	PushAllowances BranchActorAllowances `graphql:"pushAllowances(first: $pageSize, after: $cursor)"`
}

type BranchProtectionRuleWithBypassForcePushAllowances struct {
	BranchProtectionRule
	BypassForcePushAllowances BranchActorAllowances `graphql:"bypassForcePushAllowances(first: $pageSize, after: $cursor)"`
}

type BranchProtectionRuleWithBypassPullRequestAllowances struct {
	BranchProtectionRule
	BypassPullRequestAllowances BranchActorAllowances `graphql:"bypassPullRequestAllowances(first: $pageSize, after: $cursor)"`
}

type BranchActorAllowances struct {
	TotalCount int
	PageInfo   PageInfo
	Nodes      []struct {
		Actor struct {
			Type string `graphql:"type: __typename"`
			App  struct {
				Name string
				Slug string
			} `graphql:"... on App"`
			Team struct {
				Name string
				Slug string
			} `graphql:"... on Team"`
			User struct {
				Name  string
				Login string
			} `graphql:"... on User"`
		}
	}
}

// Explode returns 3 collections from BranchActorAllowances by type (in order) Apps, Teams, Users.
func (b *BranchActorAllowances) Explode() ([]NameSlug, []NameSlug, []NameLogin) {
	var apps, teams []NameSlug
	var users []NameLogin

	for _, a := range b.Nodes {
		switch a.Actor.Type {
		case "App":
			apps = append(apps, NameSlug{Name: a.Actor.App.Name, Slug: a.Actor.App.Slug})
		case "Team":
			teams = append(teams, NameSlug{Name: a.Actor.Team.Name, Slug: a.Actor.Team.Slug})
		case "User":
			users = append(users, NameLogin{Name: a.Actor.User.Name, Login: a.Actor.User.Login})
		}
	}

	return apps, teams, users
}
