package models

import (
	"github.com/shurcooL/githubv4"
)

type BasicOrganization struct {
	basicIdentifiers
	Login       string       `json:"login"`
	CreatedAt   NullableTime `json:"created_at"`
	UpdatedAt   NullableTime `json:"updated_at"`
	Description string       `json:"description"`
	Email       string       `json:"email"`
	Url         string       `json:"url"`
}

type Organization struct {
	BasicOrganization
	Announcement                           string                       `json:"announcement"`
	AnnouncementExpiresAt                  NullableTime                 `json:"announcement_expires_at"`
	AnnouncementUserDismissible            bool                         `json:"announcement_user_dismissible"`
	AnyPinnableItems                       bool                         `json:"any_pinnable_items"`
	AvatarUrl                              string                       `json:"avatar_url"`
	EstimatedNextSponsorsPayoutInCents     int                          `json:"estimated_next_sponsors_payout_in_cents"`
	HasSponsorsListing                     bool                         `json:"has_sponsors_listing"`
	InteractionAbility                     RepositoryInteractionAbility `json:"interaction_ability"`
	IsSponsoringYou                        bool                         `graphql:"isSponsoringYou: isSponsoringViewer" json:"is_sponsoring_you"`
	IsVerified                             bool                         `json:"is_verified"`
	Location                               string                       `json:"location"`
	MonthlyEstimatedSponsorsIncomeInCents  int                          `json:"monthly_estimated_sponsors_income_in_cents"`
	NewTeamUrl                             string                       `json:"new_team_url"`
	PinnedItemsRemaining                   int                          `json:"pinned_items_remaining"`
	ProjectsUrl                            string                       `json:"projects_url"`
	SamlIdentityProvider                   OrganizationIdentityProvider `json:"saml_identity_provider"`
	SponsorsListing                        SponsorsListing              `json:"sponsors_listing"`
	TeamsUrl                               string                       `json:"teams_url"`
	TotalSponsorshipAmountAsSponsorInCents int                          `json:"total_sponsorship_amount_as_sponsor_in_cents"`
	TwitterUsername                        string                       `json:"twitter_username"`
	CanAdminister                          bool                         `graphql:"canAdminister: viewerCanAdminister" json:"can_administer"`
	CanChangedPinnedItems                  bool                         `graphql:"canChangedPinnedItems: viewerCanChangePinnedItems" json:"can_changed_pinned_items"`
	CanCreateProjects                      bool                         `graphql:"canCreateProjects: viewerCanCreateProjects" json:"can_create_projects"`
	CanCreateRepositories                  bool                         `graphql:"canCreateRepositories: viewerCanCreateRepositories" json:"can_create_repositories"`
	CanCreateTeams                         bool                         `graphql:"canCreateTeams: viewerCanCreateTeams" json:"can_create_teams"`
	CanSponsor                             bool                         `graphql:"canSponsor: viewerCanSponsor" json:"can_sponsor"`
	IsAMember                              bool                         `graphql:"isAMember: viewerIsAMember" json:"is_a_member"`
	IsFollowing                            bool                         `graphql:"isFollowing: viewerIsFollowing" json:"is_following"`
	IsSponsoring                           bool                         `graphql:"isSponsoring: viewerIsSponsoring" json:"is_sponsoring"`
	WebsiteUrl                             string                       `json:"website_url"`
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

type OrganizationWithOwnerProperties struct {
	Organization
	IpAllowListEnabledSetting                     githubv4.IpAllowListEnabledSettingValue                 `json:"ip_allow_list_enabled_setting"`
	IpAllowListForInstalledAppsEnabledSetting     githubv4.IpAllowListForInstalledAppsEnabledSettingValue `json:"ip_allow_list_for_installed_apps_enabled_setting"`
	MembersCanForkPrivateRepositories             bool                                                    `json:"members_can_fork_private_repositories"`
	OrganizationBillingEmail                      string                                                  `json:"organization_billing_email"`
	NotificationDeliveryRestrictionEnabledSetting githubv4.NotificationRestrictionSettingValue            `json:"notification_delivery_restriction_enabled_setting"`
	RequiresTwoFactorAuthentication               bool                                                    `json:"requires_two_factor_authentication"`
	WebCommitSignoffRequired                      bool                                                    `json:"web_commit_signoff_required"`
}

type OrganizationWithCounts struct {
	Organization
	MembersWithRole     Count `json:"members_with_role"`
	Packages            Count `json:"packages"`
	PinnableItems       Count `json:"pinnable_items"`
	PinnedItems         Count `json:"pinned_items"`
	Projects            Count `json:"projects"`
	ProjectsV2          Count `json:"projects_v2"`
	Sponsoring          Count `json:"sponsoring"`
	Sponsors            Count `json:"sponsors"`
	Teams               Count `json:"teams"`
	PrivateRepositories Count `graphql:"privateRepositories: repositories(privacy: PRIVATE)" json:"private_repositories"`
	PublicRepositories  Count `graphql:"publicRepositories: repositories(privacy: PUBLIC)" json:"public_repositories"`
	Repositories        struct {
		TotalCount     int `json:"total_count"`
		TotalDiskUsage int `json:"total_disk_usage"`
	} `json:"repositories"`
}

type OrganizationWithOwnerPropertiesAndCounts struct {
	OrganizationWithOwnerProperties
	MembersWithRole     Count `json:"members_with_role"`
	Packages            Count `json:"packages"`
	PinnableItems       Count `json:"pinnable_items"`
	PinnedItems         Count `json:"pinned_items"`
	Projects            Count `json:"projects"`
	ProjectsV2          Count `json:"projects_v2"`
	Sponsoring          Count `json:"sponsoring"`
	Sponsors            Count `json:"sponsors"`
	Teams               Count `json:"teams"`
	PrivateRepositories Count `graphql:"privateRepositories: repositories(privacy: PRIVATE)" json:"private_repositories"`
	PublicRepositories  Count `graphql:"publicRepositories: repositories(privacy: PUBLIC)" json:"public_repositories"`
	Repositories        struct {
		TotalCount     int `json:"total_count"`
		TotalDiskUsage int `json:"total_disk_usage"`
	} `json:"repositories"`
}

type OrganizationIdentityProvider struct {
	DigestMethod    string                   `json:"digest_method"`
	IdpCertificate  githubv4.X509Certificate `json:"idp_certificate"`
	Issuer          string                   `json:"issuer"`
	SignatureMethod string                   `json:"signature_method"`
	SsoUrl          string                   `json:"sso_url"`
	// ExternalIdentities [pageable]
}
