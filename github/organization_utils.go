package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractOrganizationExternalIdentityFromHydrateItem(h *plugin.HydrateData) (models.OrganizationExternalIdentity, error) {
	if orgExternalIdentity, ok := h.Item.(models.OrganizationExternalIdentity); ok {
		return orgExternalIdentity, nil
	} else {
		return models.OrganizationExternalIdentity{}, fmt.Errorf("unable to parse hydrate item %v as a OrganizationExternalIdentity", h.Item)
	}
}

func appendOrganizationExternalIdentityColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeOrgExternalIdentityGuid"] = githubv4.Boolean(slices.Contains(cols, "guid"))
	(*m)["includeOrgExternalIdentityUser"] = githubv4.Boolean(slices.Contains(cols, "user_detail") || slices.Contains(cols, "user_login"))
	(*m)["includeOrgExternalIdentitySamlIdentity"] = githubv4.Boolean(slices.Contains(cols, "saml_identity"))
	(*m)["includeOrgExternalIdentityScimIdentity"] = githubv4.Boolean(slices.Contains(cols, "scim_identity"))
	(*m)["includeOrgExternalIdentityOrganizationInvitation"] = githubv4.Boolean(slices.Contains(cols, "organization_invitation"))
}

func appendOrgCollaboratorColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeOCPermission"] = githubv4.Boolean(slices.Contains(cols, "permission"))
	(*m)["includeOCNode"] = githubv4.Boolean(slices.Contains(cols, "user_login"))
}

func orgExternalIdentityHydrateGuid(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.Guid, nil
}

func orgExternalIdentityHydrateUserDetail(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.User, nil
}

func orgExternalIdentityHydrateUserLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.User.Login, nil
}

func orgExternalIdentityHydrateSamlIdentity(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.SamlIdentity, nil
}

func orgExternalIdentityHydrateScimIdentity(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.ScimIdentity, nil
}

func orgExternalIdentityHydrateOrganizationInvitation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgExternalIdentity, err := extractOrganizationExternalIdentityFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return orgExternalIdentity.OrganizationInvitation, nil
}

func ocHydratePermission(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	oc, err := extractOrgCollaboratorFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return oc.Permission, nil
}

func ocHydrateRepository(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	oc, err := extractOrgCollaboratorFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return oc.RepositoryName, nil
}

func ocHydrateUserLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	oc, err := extractOrgCollaboratorFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return oc.Node, nil
}

func extractOrganizationFromHydrateItem(h *plugin.HydrateData) (models.OrganizationWithCounts, error) {
	if org, ok := h.Item.(models.OrganizationWithCounts); ok {
		return org, nil
	} else {
		return models.OrganizationWithCounts{}, fmt.Errorf("unable to parse hydrate item %v as a Organization", h.Item)
	}
}

func extractOrgCollaboratorFromHydrateItem(h *plugin.HydrateData) (OrgCollaborators, error) {
	if oc, ok := h.Item.(OrgCollaborators); ok {
		return oc, nil
	} else {
		return OrgCollaborators{}, fmt.Errorf("unable to parse hydrate item %v as a OrgCollaborators", h.Item)
	}
}

func appendOrganizationColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeAnnouncementBanner"] = githubv4.Boolean(slices.Contains(cols, "announcement") || slices.Contains(cols, "announcement_expires_at") || slices.Contains(cols, "announcement_user_dismissible"))
	(*m)["includeAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeIsVerified"] = githubv4.Boolean(slices.Contains(cols, "is_verified"))
	(*m)["includeLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeNewTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "new_team_url"))
	(*m)["includePinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeSamlIdentityProvider"] = githubv4.Boolean(slices.Contains(cols, "saml_identity_provider"))
	(*m)["includeSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeTeamsUrl"] = githubv4.Boolean(slices.Contains(cols, "teams_url"))
	(*m)["includeTotalSponsorshipAmountAsSponsorInCents"] = githubv4.Boolean(slices.Contains(cols, "total_sponsorship_amount_as_sponsor_in_cents"))
	(*m)["includeTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeOrgViewer"] = githubv4.Boolean(slices.Contains(cols, "can_administer") || slices.Contains(cols, "can_changed_pinned_items") || slices.Contains(cols, "can_create_projects") || slices.Contains(cols, "can_create_repositories") || slices.Contains(cols, "can_create_teams") || slices.Contains(cols, "can_sponsor"))
	(*m)["includeIsAMember"] = githubv4.Boolean(slices.Contains(cols, "is_a_member"))
	(*m)["includeIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))
	(*m)["includeMembersWithRole"] = githubv4.Boolean(slices.Contains(cols, "members_with_role_total_count"))
	(*m)["includePackages"] = githubv4.Boolean(slices.Contains(cols, "packages_total_count"))
	(*m)["includePinnableItems"] = githubv4.Boolean(slices.Contains(cols, "pinnable_items_total_count"))
	(*m)["includePinnedItems"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_total_count"))
	(*m)["includeProjectsV2"] = githubv4.Boolean(slices.Contains(cols, "projects_v2_total_count"))
	(*m)["includeSponsoring"] = githubv4.Boolean(slices.Contains(cols, "sponsoring_total_count"))
	(*m)["includeSponsors"] = githubv4.Boolean(slices.Contains(cols, "sponsors_total_count"))
	(*m)["includeTeams"] = githubv4.Boolean(slices.Contains(cols, "teams_total_count"))
	(*m)["includePrivateRepositories"] = githubv4.Boolean(slices.Contains(cols, "private_repositories_total_count"))
	(*m)["includePublicRepositories"] = githubv4.Boolean(slices.Contains(cols, "public_repositories_total_count"))
	(*m)["includeRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_count"))
	(*m)["includeRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_disk_usage"))
}

func orgHydrateMembersWithRoleTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.MembersWithRole.TotalCount, nil
}

func orgHydratePackagesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Packages.TotalCount, nil
}

func orgHydratePinnableItemsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.PinnableItems.TotalCount, nil
}

func orgHydratePinnedItemsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.PinnedItems.TotalCount, nil
}
func orgHydrateProjectsV2TotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.ProjectsV2.TotalCount, nil
}

func orgHydrateRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Repositories.TotalCount, nil
}

func orgHydrateSponsoringTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Sponsoring.TotalCount, nil
}

func orgHydrateSponsorsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Sponsors.TotalCount, nil
}

func orgHydrateTeamsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Teams.TotalCount, nil
}

func orgHydrateRepositoriesTotalDiskUsage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Repositories.TotalDiskUsage, nil
}

func orgHydratePrivateRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.PrivateRepositories.TotalCount, nil
}

func orgHydratePublicRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.PublicRepositories.TotalCount, nil
}

func orgHydrateAnnouncement(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnnouncementBanner.Message, nil
}

func orgHydrateAnnouncementExpiresAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnnouncementBanner.ExpiresAt, nil
}

func orgHydrateAnnouncementUserDismissible(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnnouncementBanner.IsUserDismissible, nil
}

func orgHydrateAnyPinnableItems(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnyPinnableItems, nil
}

func orgHydrateAvatarUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AvatarUrl, nil
}

func orgHydrateEstimatedNextSponsorsPayoutInCents(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.EstimatedNextSponsorsPayoutInCents, nil
}

func orgHydrateHasSponsorsListing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.HasSponsorsListing, nil
}

func orgHydrateInteractionAbility(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.InteractionAbility, nil
}

func orgHydrateIsSponsoringYou(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.IsSponsoringYou, nil
}

func orgHydrateIsVerified(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.IsVerified, nil
}

func orgHydrateLocation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Location, nil
}

func orgHydrateMonthlyEstimatedSponsorsIncomeInCents(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.MonthlyEstimatedSponsorsIncomeInCents, nil
}

func orgHydrateNewTeamUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.NewTeamUrl, nil
}

func orgHydratePinnedItemsRemaining(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.PinnedItemsRemaining, nil
}

func orgHydrateProjectsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.ProjectsUrl, nil
}

func orgHydrateSamlIdentityProvider(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.SamlIdentityProvider, nil
}

func orgHydrateSponsorsListing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.SponsorsListing, nil
}

func orgHydrateTeamsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.TeamsUrl, nil
}

func orgHydrateTotalSponsorshipAmountAsSponsorInCents(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.TotalSponsorshipAmountAsSponsorInCents, nil
}

func orgHydrateTwitterUsername(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.TwitterUsername, nil
}

func orgHydrateCanAdminister(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanAdminister, nil
}

func orgHydrateCanChangedPinnedItems(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanChangedPinnedItems, nil
}

func orgHydrateCanCreateProjects(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanCreateProjects, nil
}

func orgHydrateCanCreateRepositories(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanCreateRepositories, nil
}

func orgHydrateCanCreateTeams(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanCreateTeams, nil
}

func orgHydrateCanSponsor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.CanSponsor, nil
}

func orgHydrateIsAMember(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.IsAMember, nil
}

func orgHydrateIsFollowing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.IsFollowing, nil
}

func orgHydrateIsSponsoring(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.IsSponsoring, nil
}

func orgHydrateWebsiteUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.WebsiteUrl, nil
}
