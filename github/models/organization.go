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
	Announcement                           string                       `graphql:"announcement @include(if:$includeAnnouncement)" json:"announcement"`
	AnnouncementExpiresAt                  NullableTime                 `graphql:"announcementExpiresAt @include(if:$includeAnnouncementExpiresAt)" json:"announcement_expires_at"`
	AnnouncementUserDismissible            bool                         `graphql:"announcementUserDismissible @include(if:$includeAnnouncementUserDismissible)" json:"announcement_user_dismissible"`
	AnyPinnableItems                       bool                         `graphql:"anyPinnableItems @include(if:$includeAnyPinnableItems)" json:"any_pinnable_items"`
	AvatarUrl                              string                       `graphql:"avatarUrl @include(if:$includeAvatarUrl)" json:"avatar_url"`
	EstimatedNextSponsorsPayoutInCents     int                          `graphql:"estimatedNextSponsorsPayoutInCents @include(if:$includeEstimatedNextSponsorsPayoutInCents)" json:"estimated_next_sponsors_payout_in_cents"`
	HasSponsorsListing                     bool                         `graphql:"hasSponsorsListing @include(if:$includeHasSponsorsListing)" json:"has_sponsors_listing"`
	InteractionAbility                     RepositoryInteractionAbility `graphql:"interactionAbility @include(if:$includeInteractionAbility)" json:"interaction_ability"`
	IsSponsoringYou                        bool                         `graphql:"isSponsoringYou: isSponsoringViewer @include(if:$includeIsSponsoringYou)" json:"is_sponsoring_you"`
	IsVerified                             bool                         `graphql:"isVerified @include(if:$includeIsVerified)" json:"is_verified"`
	Location                               string                       `graphql:"location @include(if:$includeLocation)" json:"location"`
	MonthlyEstimatedSponsorsIncomeInCents  int                          `graphql:"monthlyEstimatedSponsorsIncomeInCents @include(if:$includeMonthlyEstimatedSponsorsIncomeInCents)" json:"monthly_estimated_sponsors_income_in_cents"`
	NewTeamUrl                             string                       `graphql:"newTeamUrl @include(if:$includeNewTeamUrl)" json:"new_team_url"`
	PinnedItemsRemaining                   int                          `graphql:"pinnedItemsRemaining @include(if:$includePinnedItemsRemaining)" json:"pinned_items_remaining"`
	ProjectsUrl                            string                       `graphql:"projectsUrl @include(if:$includeProjectsUrl)" json:"projects_url"`
	SamlIdentityProvider                   OrganizationIdentityProvider `graphql:"samlIdentityProvider @include(if:$includeSamlIdentityProvider)" json:"saml_identity_provider"`
	SponsorsListing                        SponsorsListing              `graphql:"sponsorsListing @include(if:$includeSponsorsListing)" json:"sponsors_listing"`
	TeamsUrl                               string                       `graphql:"teamsUrl @include(if:$includeTeamsUrl)" json:"teams_url"`
	TotalSponsorshipAmountAsSponsorInCents int                          `graphql:"totalSponsorshipAmountAsSponsorInCents @include(if:$includeTotalSponsorshipAmountAsSponsorInCents)" json:"total_sponsorship_amount_as_sponsor_in_cents"`
	TwitterUsername                        string                       `graphql:"twitterUsername @include(if:$includeTwitterUsername)" json:"twitter_username"`
	CanAdminister                          bool                         `graphql:"canAdminister: viewerCanAdminister @include(if:$includeOrgViewer)" json:"can_administer"`
	CanChangedPinnedItems                  bool                         `graphql:"canChangedPinnedItems: viewerCanChangePinnedItems @include(if:$includeOrgViewer)" json:"can_changed_pinned_items"`
	CanCreateProjects                      bool                         `graphql:"canCreateProjects: viewerCanCreateProjects @include(if:$includeOrgViewer)" json:"can_create_projects"`
	CanCreateRepositories                  bool                         `graphql:"canCreateRepositories: viewerCanCreateRepositories @include(if:$includeOrgViewer)" json:"can_create_repositories"`
	CanCreateTeams                         bool                         `graphql:"canCreateTeams: viewerCanCreateTeams @include(if:$includeOrgViewer)" json:"can_create_teams"`
	CanSponsor                             bool                         `graphql:"canSponsor: viewerCanSponsor @include(if:$includeOrgViewer)" json:"can_sponsor"`
	IsAMember                              bool                         `graphql:"isAMember: viewerIsAMember @include(if:$includeIsAMember)" json:"is_a_member"`
	IsFollowing                            bool                         `graphql:"isFollowing: viewerIsFollowing @include(if:$includeIsFollowing)" json:"is_following"`
	IsSponsoring                           bool                         `graphql:"isSponsoring: viewerIsSponsoring @include(if:$includeIsSponsoring)" json:"is_sponsoring"`
	WebsiteUrl                             string                       `graphql:"websiteUrl @include(if:$includeWebsiteUrl)" json:"website_url"`
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
	DigestMethod    string `json:"digest_method"`
	Issuer          string `json:"issuer"`
	SignatureMethod string `json:"signature_method"`
	SsoUrl          string `json:"sso_url"`
	// ExternalIdentities [pageable]
}

type OrganizationExternalIdentity struct {
	Guid                   string               `json:"guid"`
	User                   BasicUser            `json:"user"`
	SamlIdentity           externalIdentitySaml `json:"saml_identity,omitempty"`
	ScimIdentity           externalIdentityBase `json:"scim_identity,omitempty"`
	OrganizationInvitation struct {
		CreatedAt      NullableTime                        `json:"created_at"`
		Email          string                              `json:"email"`
		InvitationType githubv4.OrganizationInvitationType `json:"invitation_type"`
		Invitee        BasicUser                           `json:"invitee"`
		Inviter        BasicUser                           `json:"inviter"`
		Organization   BasicOrganization                   `json:"organization"`
		Role           githubv4.OrganizationInvitationRole `json:"role"`
	} `json:"organization_invitation"`
}

type externalIdentityBase struct {
	Username   string          `json:"username"`
	GivenName  string          `json:"given_name"`
	FamilyName string          `json:"family_name"`
	Groups     []string        `json:"groups,omitempty"`
	Emails     []emailMetadata `json:"emails,omitempty"`
}

type externalIdentitySaml struct {
	externalIdentityBase
	NameId     string `json:"name_id"`
	Attributes []struct {
		Name     string `json:"name"`
		Value    string `json:"value"`
		Metadata string `json:"metadata"`
	} `json:"attributes,omitempty"`
}

type emailMetadata struct {
	Primary bool   `json:"primary"`
	Type    string `json:"type"`
	Value   string `json:"value"`
}
