package models

import (
	"github.com/shurcooL/githubv4"
	"time"
)

type BasicOrganization struct {
	basicIdentifiers
	Login       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
	Email       string
	Url         string
}

type Organization struct {
	BasicOrganization
	Announcement                                  string
	AnnouncementExpiresAt                         time.Time
	AnnouncementUserDismissible                   bool
	AnyPinnableItems                              bool
	AvatarUrl                                     string
	EstimatedNextSponsorsPayoutInCents            int
	HasSponsorsListing                            bool
	InteractionAbility                            RepositoryInteractionAbility
	IpAllowListEnabledSetting                     githubv4.IpAllowListEnabledSettingValue
	IpAllowListForInstalledAppsEnabledSetting     githubv4.IpAllowListForInstalledAppsEnabledSettingValue
	IsSponsoringYou                               bool `graphql:"isSponsoringYou: isSponsoringViewer"`
	IsVerified                                    bool
	Location                                      string
	MembersCanForkPrivateRepositories             bool
	MonthlyEstimatedSponsorsIncomeInCents         int
	NewTeamUrl                                    string
	NotificationDeliveryRestrictionEnabledSetting githubv4.NotificationRestrictionSettingValue
	OrganizationBillingEmail                      string
	PinnedItemsRemaining                          int
	ProjectsUrl                                   string
	RequiresTwoFactorAuthentication               bool
	SamlIdentityProvider                          OrganizationIdentityProvider
	SponsorsListing                               SponsorsListing
	TeamsUrl                                      string
	TotalSponsorshipAmountAsSponsorInCents        int
	TwitterUsername                               string
	CanAdminister                                 bool `graphql:"canAdminister: viewerCanAdminister"`
	CanChangedPinnedItems                         bool `graphql:"canChangedPinnedItems: viewerCanChangePinnedItems"`
	CanCreateProjects                             bool `graphql:"canCreateProjects: viewerCanCreateProjects"`
	CanCreateRepositories                         bool `graphql:"canCreateRepositories: viewerCanCreateRepositories"`
	CanCreateTeams                                bool `graphql:"canCreateTeams: viewerCanCreateTeams"`
	CanSponsor                                    bool `graphql:"canSponsor: viewerCanSponsor"`
	IsAMember                                     bool `graphql:"isAMember: viewerIsAMember"`
	IsFollowing                                   bool `graphql:"isFollowing: viewerIsFollowing"`
	IsSponsoring                                  bool `graphql:"isSponsoring: viewerIsSponsoring"`
	WebCommitSignoffRequired                      bool
	WebsiteUrl                                    string
	// AuditLog [pageable]
	// Domains [pageable]
	// EnterpriseOwners [pageable]
	// IpAllowListEntries [pageable]
	// ItemShowcase [pageable sub-item]
	// Mannequins [pageable]
	// MemberStatuses [pageable]
	// MembersWithRole [pageable]
	// Packages [pageable]
	// PendingMembers [pageable]
	// PinnableItems [pageable]
	// PinnedItems [pageable]
	// Projects [pageable]
	// ProjectsV2 [pageable]
	// RecentProjects [pageable]
	// Repositories [pageable]
	// RepositoryDiscussionComments [pageable]
	// RepositoryDiscussions [pageable]
	// RepositoryMigrations [pageable]
	// RuleSets [pageable]
	// Sponsoring [pageable]
	// Sponsors [pageable]
	// SponsorsActivities [pageable]
	// Team [find by slug]
	// Teams [pageable]
}

type OrganizationIdentityProvider struct {
	DigestMethod    string
	IdpCertificate  githubv4.X509Certificate
	Issuer          string
	Organization    *Organization
	SignatureMethod string
	SsoUrl          string
	// ExternalIdentities [pageable]
}
