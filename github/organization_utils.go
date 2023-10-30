package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractOrganizationFromHydrateItem(h *plugin.HydrateData) (models.Organization, error) {
	if org, ok := h.Item.(models.Organization); ok {
		return org, nil
	} else {
		return models.Organization{}, fmt.Errorf("unable to parse hydrate item %v as a Organization", h.Item)
	}
}

func appendOrganizationColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeAnnouncement"] = githubv4.Boolean(slices.Contains(cols, "announcement"))
	(*m)["includeAnnouncementExpiresAt"] = githubv4.Boolean(slices.Contains(cols, "announcement_expires_at"))
	(*m)["includeAnnouncementUserDismissible"] = githubv4.Boolean(slices.Contains(cols, "announcement_user_dismissible"))
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
}

func orgHydrateAnnouncement(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.Announcement, nil
}

func orgHydrateAnnouncementExpiresAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnnouncementExpiresAt, nil
}

func orgHydrateAnnouncementUserDismissible(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org, err := extractOrganizationFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return org.AnnouncementUserDismissible, nil
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