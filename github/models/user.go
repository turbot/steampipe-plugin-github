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
	}
	Followers           Count
	Following           Count
	PublicRepositories  Count `graphql:"publicRepositories: repositories(privacy: PUBLIC)"`
	PrivateRepositories Count `graphql:"privateRepositories: repositories(privacy: PRIVATE)"`
	PublicGists         Count `graphql:"publicGists: gists(privacy: PUBLIC)"`
	Issues              Count
	Organizations       Count
	PublicKeys          Count
	OpenPullRequests    Count `graphql:"openPullRequests: pullRequests(states: OPEN)"`
	MergedPullRequests  Count `graphql:"mergedPullRequests: pullRequests(states: MERGED)"`
	ClosedPullRequests  Count `graphql:"closedPullRequests: pullRequests(states: CLOSED)"`
	Packages            Count
	PinnedItems         Count
	Sponsoring          Count
	Sponsors            Count
	StarredRepositories Count
	Watching            Count
}

type userStatus struct {
	CreatedAt                    NullableTime `json:"created_at,omitempty"`
	UpdatedAt                    NullableTime `json:"updated_at,omitempty"`
	ExpiresAt                    NullableTime `json:"expires_at,omitempty"`
	Emoji                        string       `json:"emoji,omitempty"`
	Message                      string       `json:"message,omitempty"`
	IndicatesLimitedAvailability bool         `json:"indicates_limited_availability,omitempty"`
}
