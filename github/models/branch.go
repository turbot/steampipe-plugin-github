package models

// Branch returns the Name and BranchProtectionRule by default.
// Pass includeCommits [bool] GraphQL variable to toggle inclusion of BasicCommit information.
type Branch struct {
	Name   string
	Target struct {
		Commit BasicCommit `graphql:"... on Commit @include(if: $includeCommits)"`
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
	PushAllowances PushAllowances `graphql:"pushAllowances(first: $pushAllowancePageSize, after: $pushAllowanceCursor)"`
	// BranchProtectionRuleConflicts
	// BypassForcePushAllowances
	// BypassPullRequestAllowances
	// Repository
	// RestrictsReviewDismissals      bool
	// ReviewDismissalAllowances
}

type PushAllowances struct {
	TotalCount int
	PageInfo   PageInfo
	Nodes      []struct {
		Actor struct {
			Type string `graphql:"type: __typename" json:"type,omitempty"`
			App  struct {
				AppId   int    `graphql:"appId: databaseId" json:"id,omitempty"`
				AppName string `graphql:"appName: name" json:"name,omitempty"`
				AppSlug string `graphql:"appSlug: slug" json:"slug,omitempty"`
			} `graphql:"... on App" json:"app,omitempty"`
			Team struct {
				TeamId   int    `graphql:"teamId: databaseId" json:"id,omitempty"`
				TeamName string `graphql:"teamName: name" json:"name,omitempty"`
				TeamSlug string `graphql:"appSlug: slug" json:"slug,omitempty"`
			} `graphql:"... on Team" json:"team,omitempty"`
			User struct {
				Id    int    `graphql:"id: databaseId" json:"id,omitempty"`
				Name  string `json:"name,omitempty"`
				Login string `json:"login,omitempty"`
			} `graphql:"... on User" json:"user,omitempty"`
		}
	}
}
