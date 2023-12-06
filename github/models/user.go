package models

// basicIdentifiers is used to store basic identifying information.
type basicIdentifiers struct {
	Id     int    `graphql:"id: databaseId" json:"id,omitempty"`
	NodeId string `graphql:"nodeId: id" json:"node_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

type BasicUser struct {
	basicIdentifiers
	Login     string       `json:"login"`
	Email     string       `json:"email"`
	CreatedAt NullableTime `json:"created_at"`
	UpdatedAt NullableTime `json:"updated_at"`
	Url       string       `json:"url"`
}

type User struct {
	BasicUser
	AnyPinnableItems                      bool                         `graphql:"anyPinnableItems @include(if:$includeUserAnyPinnableItems)" json:"any_pinnable_items"`
	AvatarUrl                             string                       `graphql:"avatarUrl @include(if:$includeUserAvatarUrl)" json:"avatar_url"`
	Bio                                   string                       `graphql:"bio @include(if:$includeUserBio)" json:"bio"`
	Company                               string                       `graphql:"company @include(if:$includeUserCompany)" json:"company"`
	EstimatedNextSponsorsPayoutInCents    int                          `graphql:"estimatedNextSponsorsPayoutInCents @include(if:$includeUserEstimatedNextSponsorsPayoutInCents)" json:"estimated_next_sponsors_payout_in_cents"`
	HasSponsorsListing                    bool                         `graphql:"hasSponsorsListing @include(if:$includeUserHasSponsorsListing)" json:"has_sponsors_listing"`
	InteractionAbility                    RepositoryInteractionAbility `graphql:"interactionAbility @include(if:$includeUserInteractionAbility)" json:"interaction_ability,omitempty"`
	IsBountyHunter                        bool                         `graphql:"isBountyHunter @include(if:$includeUserIsBountyHunter)" json:"is_bounty_hunter"`
	IsCampusExpert                        bool                         `graphql:"isCampusExpert @include(if:$includeUserIsCampusExpert)" json:"is_campus_expert"`
	IsDeveloperProgramMember              bool                         `graphql:"isDeveloperProgramMember @include(if:$includeUserIsDeveloperProgramMember)" json:"is_developer_program_member"`
	IsEmployee                            bool                         `graphql:"isEmployee @include(if:$includeUserIsEmployee)" json:"is_employee"`
	IsFollowingYou                        bool                         `graphql:"isFollowingYou: isFollowingViewer @include(if:$includeUserIsFollowingYou)" json:"is_following_you"`
	IsGitHubStar                          bool                         `graphql:"isGitHubStar @include(if:$includeUserIsGitHubStar)" json:"is_github_star"`
	IsHireable                            bool                         `graphql:"isHireable @include(if:$includeUserIsHireable)" json:"is_hireable"`
	IsSiteAdmin                           bool                         `graphql:"isSiteAdmin @include(if:$includeUserIsSiteAdmin)" json:"is_site_admin"`
	IsSponsoringYou                       bool                         `graphql:"isSponsoringYou: isSponsoringViewer @include(if:$includeUserIsSponsoringYou)" json:"is_sponsoring_you"`
	IsYou                                 bool                         `graphql:"isYou: isViewer @include(if:$includeUserIsYou)" json:"is_you"`
	Location                              string                       `graphql:"location @include(if:$includeUserLocation)" json:"location"`
	MonthlyEstimatedSponsorsIncomeInCents int                          `graphql:"monthlyEstimatedSponsorsIncomeInCents @include(if:$includeUserMonthlyEstimatedSponsorsIncomeInCents)" json:"monthly_estimated_sponsors_income_in_cents"`
	PinnedItemsRemaining                  int                          `graphql:"pinnedItemsRemaining @include(if:$includeUserPinnedItemsRemaining)" json:"pinned_items_remaining"`
	ProjectsUrl                           string                       `graphql:"projectsUrl @include(if:$includeUserProjectsUrl)" json:"projects_url"`
	Pronouns                              string                       `graphql:"pronouns @include(if:$includeUserPronouns)" json:"pronouns"`
	SponsorsListing                       SponsorsListing              `graphql:"sponsorsListing @include(if:$includeUserSponsorsListing)" json:"sponsors_listing,omitempty"`
	Status                                userStatus                   `graphql:"status @include(if:$includeUserStatus)" json:"status,omitempty"`
	TwitterUsername                       string                       `graphql:"twitterUsername @include(if:$includeUserTwitterUsername)" json:"twitter_username"`
	CanChangedPinnedItems                 bool                         `graphql:"canChangedPinnedItems: viewerCanChangePinnedItems @include(if:$includeUserCanChangedPinnedItems)" json:"can_changed_pinned_items"`
	CanCreateProjects                     bool                         `graphql:"canCreateProjects: viewerCanCreateProjects @include(if:$includeUserCanCreateProjects)" json:"can_create_projects"`
	CanFollow                             bool                         `graphql:"canFollow: viewerCanFollow @include(if:$includeUserCanFollow)" json:"can_follow"`
	CanSponsor                            bool                         `graphql:"canSponsor: viewerCanSponsor @include(if:$includeUserCanSponsor)" json:"can_sponsor"`
	IsFollowing                           bool                         `graphql:"isFollowing: viewerIsFollowing @include(if:$includeUserIsFollowing)" json:"is_following"`
	IsSponsoring                          bool                         `graphql:"isSponsoring: viewerIsSponsoring @include(if:$includeUserIsSponsoring)" json:"is_sponsoring"`
	WebsiteUrl                            string                       `graphql:"websiteUrl @include(if:$includeUserWebsiteUrl)" json:"website_url"`
	// CommitComments [pageable]
	// ContributionsCollection [pageable]
	// CanReceiveOrganizationEmailsWhenNotificationsRestricted [requires login]
	// Followers [pageable]
	// Following [pageable]
	// Gist [find by id]
	// GistComments [pageable]
	// Gists [pageable]
	// Hovercard [requires id]
	// IsSponsoredBy [requires id]
	// IssueComments [pageable]
	// Issues [pageable]
	// ItemShowcase [pageable on sub id]
	// Organization [requires login]
	// Organizations [pageable]
	// OrganizationVerifiedDomains [requires org login & array of string for matching on]
	// Packages [pageable]
	// PinnableItems [pageable]
	// PinnedItems [pageable]
	// Project [find by number]
	// ProjectV2 [find by number]
	// Projects [pageable]
	// ProjectsV2 [pageable]
	// PublicKeys [pageable]
	// PullRequests [pageable]
	// RecentProjects [pageable]
	// Repositories [pageable]
	// RepositoriesContributedTo [pageable - note: nice filtering options]
	// Repository [find by name]
	// RepositoryDiscussionComments [pageable]
	// RepositoryDiscussions [pageable]
	// SavedReplies [pageable]
	// SocialAccounts [pageable]
	// Sponsoring [pageable]
	// Sponsors [pageable]
	// SponsorsActivities [pageable]
	// SponsorshipForViewerAsSponsor [revisit potentially nestable with a default var for `activeOnly`]
	// SponsorshipForViewerAsSponsorable [revisit potentially nestable with a default var for `activeOnly`]
	// SponsorshipNewsletters [pageable]
	// SponsorshipAsMaintainer [pageable]
	// SponsorshipAsSponsor [pageable]
	// StarredRepositories [pageable]
	// TopRepositories [pageable]
	// TotalSponsorshipAmountAsSponsorInCents [requires params & only viewable by user themselves]
	// Watching [pageable]
}

type UserWithCounts struct {
	User
	Repositories struct {
		TotalDiskUsage int
	} `graphql:"repositories @include(if:$includeUserRepositories)" json:"repositories"`
	Followers           Count `graphql:"followers @include(if:$includeUserFollowers)" json:"followers"`
	Following           Count `graphql:"following @include(if:$includeUserFollowing)" json:"following"`
	PublicRepositories  Count `graphql:"publicRepositories: repositories(privacy: PUBLIC) @include(if:$includeUserPublicRepositories)" json:"public_repositories"`
	PrivateRepositories Count `graphql:"privateRepositories: repositories(privacy: PRIVATE) @include(if:$includeUserPrivateRepositories)" json:"private_repositories"`
	PublicGists         Count `graphql:"publicGists: gists(privacy: PUBLIC) @include(if:$includeUserPublicGists)" json:"public_gists"`
	Issues              Count `graphql:"issues @include(if:$includeUserIssues)" json:"issues"`
	Organizations       Count `graphql:"organizations @include(if:$includeUserOrganizations)" json:"organizations"`
	PublicKeys          Count `graphql:"publicKeys @include(if:$includeUserPublicKeys)" json:"public_keys"`
	OpenPullRequests    Count `graphql:"openPullRequests: pullRequests(states: OPEN) @include(if:$includeUserOpenPullRequests)" json:"open_pull_requests"`
	MergedPullRequests  Count `graphql:"mergedPullRequests: pullRequests(states: MERGED) @include(if:$includeUserMergedPullRequests)" json:"merged_pull_requests"`
	ClosedPullRequests  Count `graphql:"closedPullRequests: pullRequests(states: CLOSED) @include(if:$includeUserClosedPullRequests)" json:"closed_pull_requests"`
	Packages            Count `graphql:"packages @include(if:$includeUserPackages)" json:"packages"`
	PinnedItems         Count `graphql:"pinnedItems @include(if:$includeUserPinnedItems)" json:"pinned_items"`
	Sponsoring          Count `graphql:"sponsoring @include(if:$includeUserSponsoring)" json:"sponsoring"`
	Sponsors            Count `graphql:"sponsors @include(if:$includeUserSponsors)" json:"sponsors"`
	StarredRepositories Count `graphql:"starredRepositories @include(if:$includeUserStarredRepositories)" json:"starred_repositories"`
	Watching            Count `graphql:"watching @include(if:$includeUserWatching)" json:"watching"`
}

type userStatus struct {
	CreatedAt                    NullableTime `json:"created_at,omitempty"`
	UpdatedAt                    NullableTime `json:"updated_at,omitempty"`
	ExpiresAt                    NullableTime `json:"expires_at,omitempty"`
	Emoji                        string       `json:"emoji,omitempty"`
	Message                      string       `json:"message,omitempty"`
	IndicatesLimitedAvailability bool         `json:"indicates_limited_availability,omitempty"`
}

type BaseUser struct {
	BasicUser
	AnyPinnableItems                      bool                         `json:"any_pinnable_items"`
	AvatarUrl                             string                       `json:"avatar_url"`
	Bio                                   string                       `json:"bio"`
	Company                               string                       `json:"company"`
	EstimatedNextSponsorsPayoutInCents    int                          `json:"estimated_next_sponsors_payout_in_cents"`
	HasSponsorsListing                    bool                         `json:"has_sponsors_listing"`
	InteractionAbility                    RepositoryInteractionAbility `json:"interaction_ability,omitempty"`
	IsBountyHunter                        bool                         `json:"is_bounty_hunter"`
	IsCampusExpert                        bool                         `json:"is_campus_expert"`
	IsDeveloperProgramMember              bool                         `json:"is_developer_program_member"`
	IsEmployee                            bool                         `json:"is_employee"`
	IsFollowingYou                        bool                         `graphql:"isFollowingYou: isFollowingViewer" json:"is_following_you"`
	IsGitHubStar                          bool                         `json:"is_github_star"`
	IsHireable                            bool                         `json:"is_hireable"`
	IsSiteAdmin                           bool                         `json:"is_site_admin"`
	IsSponsoringYou                       bool                         `graphql:"isSponsoringYou: isSponsoringViewer" json:"is_sponsoring_you"`
	IsYou                                 bool                         `graphql:"isYou: isViewer" json:"is_you"`
	Location                              string                       `json:"location"`
	MonthlyEstimatedSponsorsIncomeInCents int                          `json:"monthly_estimated_sponsors_income_in_cents"`
	PinnedItemsRemaining                  int                          `json:"pinned_items_remaining"`
	ProjectsUrl                           string                       `json:"projects_url"`
	Pronouns                              string                       `json:"pronouns"`
	SponsorsListing                       SponsorsListing              `json:"sponsors_listing,omitempty"`
	Status                                userStatus                   `json:"status,omitempty"`
	TwitterUsername                       string                       `json:"twitter_username"`
	CanChangedPinnedItems                 bool                         `graphql:"canChangedPinnedItems: viewerCanChangePinnedItems" json:"can_changed_pinned_items"`
	CanCreateProjects                     bool                         `graphql:"canCreateProjects: viewerCanCreateProjects" json:"can_create_projects"`
	CanFollow                             bool                         `graphql:"canFollow: viewerCanFollow" json:"can_follow"`
	CanSponsor                            bool                         `graphql:"canSponsor: viewerCanSponsor" json:"can_sponsor"`
	IsFollowing                           bool                         `graphql:"isFollowing: viewerIsFollowing" json:"is_following"`
	IsSponsoring                          bool                         `graphql:"isSponsoring: viewerIsSponsoring" json:"is_sponsoring"`
	WebsiteUrl                            string                       `json:"website_url"`
}
