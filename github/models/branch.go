package models

type Branch struct {
	Name   string
	Target struct {
		Commit Commit `graphql:"... on Commit"`
	}
	BranchProtectionRule BranchProtectionRule
}

type BranchProtectionRule struct {
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
	MatchingBranches               struct {
		TotalCount int `json:"total_count"`
	} `graphql:"matchingBranches: matchingRefs" json:"matching_branches"`
	// BranchProtectionRuleConflicts
	// BypassForcePushAllowances
	// BypassPullRequestAllowances
	// Repository
	// RestrictsReviewDismissals      bool
	// ReviewDismissalAllowances
}

type BranchProtectionRuleWithPushAllowances struct {
	BranchProtectionRule
	PushAllowances PushAllowances `graphql:"pushAllowances(first: $pushAllowancePageSize, after: $pushAllowanceCursor)"`
}

type PushAllowances struct {
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
