package models

import "github.com/shurcooL/githubv4"

type Repository struct {
	basicIdentifiers
	AllowUpdateBranch             bool                             `json:"allow_update_branch"`
	ArchivedAt                    NullableTime                     `json:"archived_at"`
	AutoMergeAllowed              bool                             `json:"auto_merge_allowed"`
	CodeOfConduct                 RepositoryCodeOfConduct          `json:"code_of_conduct"`
	ContactLinks                  []RepositoryContactLink          `json:"contact_links"`
	CreatedAt                     NullableTime                     `json:"created_at"`
	DefaultBranchRef              BasicRefWithBranchProtectionRule `json:"default_branch_ref"`
	DeleteBranchOnMerge           bool                             `json:"delete_branch_on_merge"`
	Description                   string                           `json:"description"`
	DiskUsage                     int                              `json:"disk_usage "`
	ForkCount                     int                              `json:"fork_count"`
	ForkingAllowed                bool                             `json:"forking_allowed"`
	FundingLinks                  []RepositoryFundingLinks         `json:"funding_links"`
	HasDiscussionsEnabled         bool                             `json:"has_discussions_enabled"`
	HasIssuesEnabled              bool                             `json:"has_issues_enabled"`
	HasProjectsEnabled            bool                             `json:"has_projects_enabled"`
	HasVulnerabilityAlertsEnabled bool                             `json:"has_vulnerability_alerts_enabled"`
	HasWikiEnabled                bool                             `json:"has_wiki_enabled"`
	HomepageUrl                   string                           `json:"homepage_url"`
	InteractionAbility            RepositoryInteractionAbility     `json:"interaction_ability"`
	IsArchived                    bool                             `json:"is_archived"`
	IsBlankIssuesEnabled          bool                             `json:"is_blank_issues_enabled"`
	IsDisabled                    bool                             `json:"is_disabled"`
	IsEmpty                       bool                             `json:"is_empty"`
	IsFork                        bool                             `json:"is_fork"`
	IsInOrganization              bool                             `json:"is_in_organization"`
	IsLocked                      bool                             `json:"is_locked"`
	IsMirror                      bool                             `json:"is_mirror"`
	IsPrivate                     bool                             `json:"is_private"`
	IsSecurityPolicyEnabled       bool                             `json:"is_security_policy_enabled"`
	IsTemplate                    bool                             `json:"is_template"`
	IsUserConfigurationRepository bool                             `json:"is_user_configuration_repository"`
	IssueTemplates                []IssueTemplate                  `json:"issue_templates"`
	LatestRelease                 Release                          `json:"latest_release"`
	LicenseInfo                   BasicLicense                     `json:"license_info"`
	LockReason                    githubv4.LockReason              `json:"lock_reason"`
	MergeCommitAllowed            bool                             `json:"merge_commit_allowed"`
	MergeCommitMessage            githubv4.MergeCommitMessage      `json:"merge_commit_message"`
	MergeCommitTitle              githubv4.MergeCommitTitle        `json:"merge_commit_title"`
	MirrorUrl                     string                           `json:"mirror_url"`
	NameWithOwner                 string                           `json:"name_with_owner"`
	OpenGraphImageUrl             string                           `json:"open_graph_image_url"`
	Owner                         struct {
		Login string `json:"login"`
	} `json:"owner"`
	PrimaryLanguage          Language                          `json:"primary_language"`
	ProjectsUrl              string                            `json:"projects_url"`
	PullRequestTemplates     []PullRequestTemplate             `json:"pull_request_templates"`
	PushedAt                 NullableTime                      `json:"pushed_at"`
	RebaseMergeAllowed       bool                              `json:"rebase_merge_allowed"`
	SecurityPolicyUrl        string                            `json:"security_policy_url"`
	SquashMergeAllowed       bool                              `json:"squash_merge_allowed"`
	SquashMergeCommitMessage githubv4.SquashMergeCommitMessage `json:"squash_merge_commit_message"`
	SquashMergeCommitTitle   githubv4.SquashMergeCommitTitle   `json:"squash_merge_commit_title"`
	SshUrl                   string                            `json:"ssh_url"`
	StargazerCount           int                               `json:"stargazer_count"`
	TempCloneToken           string                            `json:"temp_clone_token"`
	UpdatedAt                NullableTime                      `json:"updated_at"`
	Url                      string                            `json:"url"`
	UsesCustomOpenGraphImage bool                              `json:"uses_custom_open_graph_image"`
	CanAdminister            bool                              `graphql:"canAdminister: viewerCanAdminister" json:"can_administer"`
	CanCreateProjects        bool                              `graphql:"canCreateProjects: viewerCanCreateProjects" json:"can_create_projects"`
	CanSubscribe             bool                              `graphql:"canSubscribe: viewerCanSubscribe" json:"can_subscribe"`
	CanUpdateTopics          bool                              `graphql:"canUpdateTopics: viewerCanUpdateTopics" json:"can_update_topics"`
	HasStarred               bool                              `graphql:"hasStarred: viewerHasStarred" json:"has_starred"`
	YourPermission           githubv4.RepositoryPermission     `graphql:"yourPermission: viewerPermission" json:"your_permission"`
	PossibleCommitEmails     []string                          `graphql:"possibleCommitEmails: viewerPossibleCommitEmails" json:"possible_commit_emails"`
	Subscription             githubv4.SubscriptionState        `graphql:"subscription: viewerSubscription" json:"subscription"`
	Visibility               githubv4.RepositoryVisibility     `json:"visibility"`
	WebCommitSignoffRequired bool                              `json:"web_commit_signoff_required"`

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
