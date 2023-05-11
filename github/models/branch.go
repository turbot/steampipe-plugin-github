package models

type Branch struct {
	Name   string
	Target struct {
		Commit Commit `graphql:"... on Commit"`
	}
	BranchProtectionRule BranchProtectionRule
}

type BranchProtectionRule struct {
	AllowsDeletions                bool
	AllowsForcePushes              bool
	BlocksCreations                bool
	Creator                        Actor
	Id                             int    `graphql:"id: databaseId"`
	NodeId                         string `graphql:"nodeId: id"`
	DismissesStaleReviews          bool
	IsAdminEnforced                bool
	LockAllowsFetchAndMerge        bool
	LockBranch                     bool
	Pattern                        string
	RequireLastPushApproval        bool
	RequiredApprovingReviewCount   int
	RequiredDeploymentEnvironments []string
	RequiredStatusChecks           []string `graphql:"requiredStatusChecks: requiredStatusCheckContexts"`
	RequiresApprovingReviews       bool
	RequiresConversationResolution bool
	RequiresCodeOwnerReviews       bool
	RequiresCommitSignatures       bool
	RequiresDeployments            bool
	RequiresLinearHistory          bool
	RequiresStatusChecks           bool
	RequiresStrictStatusChecks     bool
	RestrictsPushes                bool
	MatchingBranches               struct {
		TotalCount int
	} `graphql:"matchingBranches: matchingRefs"`
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
