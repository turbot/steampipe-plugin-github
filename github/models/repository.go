package models

import "github.com/shurcooL/githubv4"

type Repository struct {
	basicIdentifiers
	AllowUpdateBranch             bool                             `graphql:"allowUpdateBranch @include(if:$includeAllowUpdateBranch)" json:"allow_update_branch"`
	ArchivedAt                    NullableTime                     `graphql:"archivedAt @include(if:$includeArchivedAt)" json:"archived_at"`
	AutoMergeAllowed              bool                             `graphql:"autoMergeAllowed @include(if:$includeAutoMergeAllowed)" json:"auto_merge_allowed"`
	CodeOfConduct                 RepositoryCodeOfConduct          `graphql:"codeOfConduct @include(if:$includeCodeOfConduct)" json:"code_of_conduct"`
	ContactLinks                  []RepositoryContactLink          `graphql:"contactLinks @include(if:$includeContactLinks)" json:"contact_links"`
	CreatedAt                     NullableTime                     `graphql:"createdAt @include(if:$includeCreatedAt)" json:"created_at"`
	DefaultBranchRef              BasicRefWithBranchProtectionRule `graphql:"defaultBranchRef @include(if:$includeDefaultBranchRef)" json:"default_branch_ref"`
	DeleteBranchOnMerge           bool                             `graphql:"deleteBranchOnMerge @include(if:$includeDeleteBranchOnMerge)" json:"delete_branch_on_merge"`
	Description                   string                           `graphql:"description @include(if:$includeDescription)" json:"description"`
	DiskUsage                     int                              `graphql:"diskUsage @include(if:$includeDiskUsage)" json:"disk_usage "`
	ForkCount                     int                              `graphql:"forkCount @include(if:$includeForkCount)" json:"fork_count"`
	ForkingAllowed                bool                             `graphql:"forkingAllowed @include(if:$includeForkingAllowed)" json:"forking_allowed"`
	FundingLinks                  []RepositoryFundingLinks         `graphql:"fundingLinks @include(if:$includeFundingLinks)" json:"funding_links"`
	HasDiscussionsEnabled         bool                             `graphql:"hasDiscussionsEnabled @include(if:$includeHasDiscussionsEnabled)" json:"has_discussions_enabled"`
	HasIssuesEnabled              bool                             `graphql:"hasIssuesEnabled @include(if:$includeHasIssuesEnabled)" json:"has_issues_enabled"`
	HasProjectsEnabled            bool                             `graphql:"hasProjectsEnabled @include(if:$includeHasProjectsEnabled)" json:"has_projects_enabled"`
	HasVulnerabilityAlertsEnabled bool                             `graphql:"hasVulnerabilityAlertsEnabled @include(if:$includeHasVulnerabilityAlertsEnabled)" json:"has_vulnerability_alerts_enabled"`
	HasWikiEnabled                bool                             `graphql:"hasWikiEnabled @include(if:$includeHasWikiEnabled)" json:"has_wiki_enabled"`
	HomepageUrl                   string                           `graphql:"homepageUrl @include(if:$includeHomepageUrl)" json:"homepage_url"`
	InteractionAbility            RepositoryInteractionAbility     `graphql:"interactionAbility @include(if:$includeUserInteractionAbility)" json:"interaction_ability,omitempty"`
	IsArchived                    bool                             `graphql:"isArchived @include(if:$includeIsArchived)" json:"is_archived"`
	IsBlankIssuesEnabled          bool                             `graphql:"isBlankIssuesEnabled @include(if:$includeIsBlankIssuesEnabled)" json:"is_blank_issues_enabled"`
	IsDisabled                    bool                             `graphql:"isDisabled @include(if:$includeIsDisabled)" json:"is_disabled"`
	IsEmpty                       bool                             `graphql:"isEmpty @include(if:$includeIsEmpty)" json:"is_empty"`
	IsFork                        bool                             `graphql:"isFork @include(if:$includeIsFork)" json:"is_fork"`
	IsInOrganization              bool                             `graphql:"isInOrganization @include(if:$includeIsInOrganization)" json:"is_in_organization"`
	IsLocked                      bool                             `graphql:"isLocked @include(if:$includeIsLocked)" json:"is_locked"`
	IsMirror                      bool                             `graphql:"isMirror @include(if:$includeIsMirror)" json:"is_mirror"`
	IsPrivate                     bool                             `graphql:"isPrivate @include(if:$includeIsPrivate)" json:"is_private"`
	IsSecurityPolicyEnabled       bool                             `graphql:"isSecurityPolicyEnabled @include(if:$includeIsSecurityPolicyEnabled)" json:"is_security_policy_enabled"`
	IsTemplate                    bool                             `graphql:"isTemplate @include(if:$includeIsTemplate)" json:"is_template"`
	IsUserConfigurationRepository bool                             `graphql:"isUserConfigurationRepository @include(if:$includeIsUserConfigurationRepository)" json:"is_user_configuration_repository"`
	IssueTemplates                []IssueTemplate                  `graphql:"issueTemplates @include(if:$includeIssueTemplates)" json:"issue_templates"`
	LicenseInfo                   BasicLicense                     `graphql:"licenseInfo @include(if:$includeLicenseInfo)" json:"license_info"`
	LockReason                    githubv4.LockReason              `graphql:"lockReason @include(if:$includeLockReason)" json:"lock_reason"`
	MergeCommitAllowed            bool                             `graphql:"mergeCommitAllowed @include(if:$includeMergeCommitAllowed)" json:"merge_commit_allowed"`
	MergeCommitMessage            githubv4.MergeCommitMessage      `graphql:"mergeCommitMessage @include(if:$includeMergeCommitMessage)" json:"merge_commit_message"`
	MergeCommitTitle              githubv4.MergeCommitTitle        `graphql:"mergeCommitTitle @include(if:$includeMergeCommitTitle)" json:"merge_commit_title"`
	MirrorUrl                     string                           `graphql:"mirrorUrl @include(if:$includeMirrorUrl)" json:"mirror_url"`
	NameWithOwner                 string                           `json:"name_with_owner"`
	OpenGraphImageUrl             string                           `graphql:"openGraphImageUrl @include(if:$includeOpenGraphImageUrl)" json:"open_graph_image_url"`
	Owner                         struct {
		Login string `json:"login"`
	} `json:"owner"`
	PrimaryLanguage          Language                          `graphql:"primaryLanguage @include(if:$includePrimaryLanguage)" json:"primary_language"`
	ProjectsUrl              string                            `graphql:"projectsUrl @include(if:$includeProjectsUrl)" json:"projects_url"`
	PullRequestTemplates     []PullRequestTemplate             `graphql:"pullRequestTemplates @include(if:$includePullRequestTemplates)" json:"pull_request_templates"`
	PushedAt                 NullableTime                      `graphql:"pushedAt @include(if:$includePushedAt)" json:"pushed_at"`
	RebaseMergeAllowed       bool                              `graphql:"rebaseMergeAllowed @include(if:$includeRebaseMergeAllowed)" json:"rebase_merge_allowed"`
	SecurityPolicyUrl        string                            `graphql:"securityPolicyUrl @include(if:$includeSecurityPolicyUrl)" json:"security_policy_url"`
	SquashMergeAllowed       bool                              `graphql:"squashMergeAllowed @include(if:$includeSquashMergeAllowed)" json:"squash_merge_allowed"`
	SquashMergeCommitMessage githubv4.SquashMergeCommitMessage `graphql:"squashMergeCommitMessage @include(if:$includeSquashMergeCommitMessage)" json:"squash_merge_commit_message"`
	SquashMergeCommitTitle   githubv4.SquashMergeCommitTitle   `graphql:"squashMergeCommitTitle @include(if:$includeSquashMergeCommitTitle)" json:"squash_merge_commit_title"`
	SshUrl                   string                            `graphql:"sshUrl @include(if:$includeSshUrl)" json:"ssh_url"`
	StargazerCount           int                               `graphql:"stargazerCount @include(if:$includeStargazerCount)" json:"stargazer_count"`
	UpdatedAt                NullableTime                      `graphql:"updatedAt @include(if:$includeUpdatedAt)" json:"updated_at"`
	Url                      string                            `graphql:"url @include(if:$includeUrl)" json:"url"`
	UsesCustomOpenGraphImage bool                              `graphql:"usesCustomOpenGraphImage @include(if:$includeUsesCustomOpenGraphImage)" json:"uses_custom_open_graph_image"`
	CanAdminister            bool                              `graphql:"canAdminister: viewerCanAdminister @include(if:$includeCanAdminister)" json:"can_administer"`
	CanCreateProjects        bool                              `graphql:"canCreateProjects: viewerCanCreateProjects @include(if:$includeCanCreateProjects)" json:"can_create_projects"`
	CanSubscribe             bool                              `graphql:"canSubscribe: viewerCanSubscribe @include(if:$includeCanSubscribe)" json:"can_subscribe"`
	CanUpdateTopics          bool                              `graphql:"canUpdateTopics: viewerCanUpdateTopics @include(if:$includeCanUpdateTopics)" json:"can_update_topics"`
	HasStarred               bool                              `graphql:"hasStarred: viewerHasStarred @include(if:$includeHasStarred)" json:"has_starred"`
	YourPermission           githubv4.RepositoryPermission     `graphql:"yourPermission: viewerPermission  @include(if:$includeYourPermission)" json:"your_permission"`
	PossibleCommitEmails     []string                          `graphql:"possibleCommitEmails: viewerPossibleCommitEmails @include(if:$includePossibleCommitEmails)" json:"possible_commit_emails"`
	Subscription             githubv4.SubscriptionState        `graphql:"subscription: viewerSubscription @include(if:$includeSubscription)" json:"subscription"`
	Visibility               githubv4.RepositoryVisibility     `graphql:"visibility @include(if:$includeVisibility)" json:"visibility"`
	WebCommitSignoffRequired bool                              `graphql:"webCommitSignoffRequired @include(if:$includeWebCommitSignoffRequired)" json:"web_commit_signoff_required"`
	RepositoryTopics         Count                             `graphql:"repositoryTopics @include(if:$includeRepositoryTopics)" json:"repository_topics"`
	OpenIssues               Count                             `graphql:"issues(states: OPEN) @include(if:$includeOpenIssues)" json:"open_issues"`
	Watchers                 Count                             `graphql:"watchers @include(if:$includeWatchers)" json:"watchers"`
	// AssignableUsers [pageable]
	// BranchProtectionRules [pageable]
	// CodeOwners [search by refName]
	// Collaborators [pageable]
	// CommitComments [pageable]
	// DependencyGraphManifests [pageable]
	// DeployKeys [pageable]
	// Deployments [pageable]
	// Discussion [search by number]
	// Discussions [pageable]
	// DiscussionCategory [search by slug]
	// DiscussionCategories [pageable]
	// Environment [search by name]
	// Environments [pageable]
	// Forks [pageable]
	// Issue [search by number]
	// Issues [pageable]
	// IssueOrPullRequest [search by number]
	// Label [search by name]
	// Labels [pageable]
	// Languages [pageable]
	// MentionableUsers [pageable]
	// MergeQueue [search by branch]
	// Milestone [search by number]
	// Milestones [pageable]
	// Packages [pageable]
	// PinnedDiscussions [pageable]
	// PinnedIssues [pageable]
	// Project [find by number]
	// Projects [pageable]
	// ProjectV2 [find by number]
	// ProjectsV2 [pageable]
	// PullRequest [search by number]
	// PullRequests [pageable]
	// Ref [search by qualifiedName]
	// Refs [pageable]
	// Release [search by tagName]
	// Releases [pageable]
	// RepositoryTopics [pageable]
	// RuleSets [pageable]
	// Stargazers [pageable]
	// Submodules [pageable]
	// VulnerabilityAlert [search by number]
	// VulnerabilityAlerts [pageable]
	// Watchers [pageable]
}

type RepositoryInteractionAbility struct {
	ExpiresAt NullableTime `json:"expires_at,omitempty"`
	Limit     string       `json:"repository_interaction_limit,omitempty"`
	Origin    string       `json:"repository_interaction_limit_origin,omitempty"`
}

type RepositoryCodeOfConduct struct {
	Id   githubv4.ID `json:"-"`
	Key  string      `json:"key"`
	Name string      `json:"name"`
	Body string      `json:"body"`
	Url  string      `json:"url"`
}

type RepositoryContactLink struct {
	Name  string `json:"name"`
	About string `json:"about"`
	Url   string `json:"url"`
}

type RepositoryFundingLinks struct {
	Url      string                   `json:"url"`
	Platform githubv4.FundingPlatform `json:"platform"`
}
