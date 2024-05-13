package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func sharedOrganizationColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Login", "Node.Login"), Description: "The login name of the organization."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID number of the organization."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the organization."},
		{Name: "name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name", "Node.Name"), Description: "The display name of the organization."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt", "Node.CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the organization was created."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt", "Node.UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the organization was last updated."},
		{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description", "Node.Description"), Description: "The description of the organization."},
		{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Email", "Node.Email"), Description: "The email address associated with the organization."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url", "Node.Url"), Description: "The URL for this organization."},
		{Name: "announcement", Type: proto.ColumnType_STRING, Hydrate: orgHydrateAnnouncement, Transform: transform.FromValue(), Description: "The text of the announcement."},
		{Name: "announcement_expires_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: orgHydrateAnnouncementExpiresAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "The expiration date of the announcement, if any."},
		{Name: "announcement_user_dismissible", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateAnnouncementUserDismissible, Transform: transform.FromValue(), Description: "If true, the announcement can be dismissed by the user."},
		{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateAnyPinnableItems, Transform: transform.FromValue(), Description: "If true, this organization has items that can be pinned to their profile."},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Hydrate: orgHydrateAvatarUrl, Transform: transform.FromValue(), Description: "URL pointing to the organization's public avatar."},
		{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Hydrate: orgHydrateEstimatedNextSponsorsPayoutInCents, Transform: transform.FromValue(), Description: "The estimated next GitHub Sponsors payout for this organization in cents (USD)."},
		{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateHasSponsorsListing, Transform: transform.FromValue(), Description: "If true, this organization has a GitHub Sponsors listing."},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Hydrate: orgHydrateInteractionAbility, Transform: transform.FromValue().NullIfZero(), Description: "The interaction ability settings for this organization."},
		{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateIsSponsoringYou, Transform: transform.FromValue(), Description: "If true, you are sponsored by this organization."},
		{Name: "is_verified", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateIsVerified, Transform: transform.FromValue(), Description: "If true, the organization has verified its profile email and website."},
		{Name: "location", Type: proto.ColumnType_STRING, Hydrate: orgHydrateLocation, Transform: transform.FromValue(), Description: "The organization's public profile location."},
		{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Hydrate: orgHydrateMonthlyEstimatedSponsorsIncomeInCents, Transform: transform.FromValue(), Description: "The estimated monthly GitHub Sponsors income for this organization in cents (USD)."},
		{Name: "new_team_url", Type: proto.ColumnType_STRING, Hydrate: orgHydrateNewTeamUrl, Transform: transform.FromValue(), Description: "URL for creating a new team."},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Hydrate: orgHydratePinnedItemsRemaining, Transform: transform.FromValue(), Description: "Returns how many more items this organization can pin to their profile."},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Hydrate: orgHydrateProjectsUrl, Transform: transform.FromValue(), Description: "URL listing organization's projects."},
		{Name: "saml_identity_provider", Type: proto.ColumnType_JSON, Hydrate: orgHydrateSamlIdentityProvider, Transform: transform.FromValue().NullIfZero(), Description: "The Organization's SAML identity provider. Visible to (1) organization owners, (2) organization owners' personal access tokens (classic) with read:org or admin:org scope, (3) GitHub App with an installation token with read or write access to members, else null."},
		{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Hydrate: orgHydrateSponsorsListing, Transform: transform.FromValue().NullIfZero(), Description: "The GitHub sponsors listing for this organization."},
		{Name: "teams_url", Type: proto.ColumnType_STRING, Hydrate: orgHydrateTeamsUrl, Transform: transform.FromValue(), Description: "URL listing organization's teams."},
		{Name: "total_sponsorship_amount_as_sponsor_in_cents", Type: proto.ColumnType_INT, Hydrate: orgHydrateTotalSponsorshipAmountAsSponsorInCents, Transform: transform.FromValue().NullIfZero(), Description: "The amount in United States cents (e.g., 500 = $5.00 USD) that this entity has spent on GitHub to fund sponsorships. Only returns a value when viewed by the user themselves or by a user who can manage sponsorships for the requested organization."},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Hydrate: orgHydrateTwitterUsername, Transform: transform.FromValue(), Description: "The organization's Twitter username."},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanAdminister, Transform: transform.FromValue(), Description: "If true, you can administer the organization."},
		{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanChangedPinnedItems, Transform: transform.FromValue(), Description: "If true, you can change the pinned items on the organization's profile."},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanCreateProjects, Transform: transform.FromValue(), Description: "If true, you can create projects for the organization."},
		{Name: "can_create_repositories", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanCreateRepositories, Transform: transform.FromValue(), Description: "If true, you can create repositories for the organization."},
		{Name: "can_create_teams", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanCreateTeams, Transform: transform.FromValue(), Description: "If true, you can create teams within the organization."},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateCanSponsor, Transform: transform.FromValue(), Description: "If true, you can sponsor this organization."},
		{Name: "is_a_member", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateIsAMember, Transform: transform.FromValue(), Description: "If true, you are an active member of the organization."},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateIsFollowing, Transform: transform.FromValue(), Description: "If true, you are following the organization."},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Hydrate: orgHydrateIsSponsoring, Transform: transform.FromValue(), Description: "If true, you are sponsoring the organization."},
		{Name: "website_url", Type: proto.ColumnType_STRING, Hydrate: orgHydrateWebsiteUrl, Transform: transform.FromValue(), Description: "URL for the organization's public website."},
		// Columns from v3 api - hydrates
		{Name: "hooks", Type: proto.ColumnType_JSON, Description: "The Hooks of the organization.", Hydrate: hydrateOrganizationHooksFromV3, Transform: transform.FromValue()},
		{Name: "billing_email", Type: proto.ColumnType_STRING, Description: "The email address for billing.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "two_factor_requirement_enabled", Type: proto.ColumnType_BOOL, Description: "If true, all members in the organization must have two factor authentication enabled.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "default_repo_permission", Type: proto.ColumnType_STRING, Description: "The default repository permissions for the organization.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_allowed_repository_creation_type", Type: proto.ColumnType_STRING, Description: "Specifies which types of repositories non-admin organization members can create", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_create_internal_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create internal repositories.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_create_pages", Type: proto.ColumnType_BOOL, Description: "If true, members can create pages.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_create_private_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create private repositories.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_create_public_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create public repositories.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_create_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create repositories.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "members_can_fork_private_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can fork private organization repositories.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "plan_filled_seats", Type: proto.ColumnType_INT, Description: "The number of used seats for the plan.", Hydrate: hydrateOrganizationDataFromV3, Transform: transform.FromField("Plan.FilledSeats")},
		{Name: "plan_name", Type: proto.ColumnType_STRING, Description: "The name of the GitHub plan.", Hydrate: hydrateOrganizationDataFromV3, Transform: transform.FromField("Plan.Name")},
		{Name: "plan_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories for the plan.", Hydrate: hydrateOrganizationDataFromV3, Transform: transform.FromField("Plan.PrivateRepos")},
		{Name: "plan_seats", Type: proto.ColumnType_INT, Description: "The number of available seats for the plan", Hydrate: hydrateOrganizationDataFromV3, Transform: transform.FromField("Plan.Seats")},
		{Name: "plan_space", Type: proto.ColumnType_INT, Description: "The total space allocated for the plan.", Hydrate: hydrateOrganizationDataFromV3, Transform: transform.FromField("Plan.Space")},
		{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the organization.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the organization.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators for the organization.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "has_organization_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use organization projects.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "has_repository_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use repository projects.", Hydrate: hydrateOrganizationDataFromV3},
		{Name: "web_commit_signoff_required", Type: proto.ColumnType_BOOL, Description: "If true, contributors are required to sign off on web-based commits for repositories in this organization.", Hydrate: hydrateOrganizationDataFromV3},
	}
}

func sharedOrganizationCountColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "members_with_role_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateMembersWithRoleTotalCount, Transform: transform.FromValue(), Description: "Count of members with a role within the organization."},
		{Name: "packages_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydratePackagesTotalCount, Transform: transform.FromValue(), Description: "Count of packages within the organization."},
		{Name: "pinnable_items_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydratePinnableItemsTotalCount, Transform: transform.FromValue(), Description: "Count of pinnable items within the organization."},
		{Name: "pinned_items_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydratePinnedItemsTotalCount, Transform: transform.FromValue(), Description: "Count of itesm pinned to the organization's profile."},
		{Name: "projects_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateProjectsTotalCount, Transform: transform.FromValue(), Description: "Count of projects within the organization."},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateProjectsV2TotalCount, Transform: transform.FromValue(), Description: "Count of V2 projects within the organization."},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateRepositoriesTotalCount, Transform: transform.FromValue(), Description: "Count of all repositories within the organization."},
		{Name: "sponsoring_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateSponsoringTotalCount, Transform: transform.FromValue(), Description: "Count of users the organization is sponsoring."},
		{Name: "sponsors_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateSponsorsTotalCount, Transform: transform.FromValue(), Description: "Count of sponsors the organization has."},
		{Name: "teams_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydrateTeamsTotalCount, Transform: transform.FromValue(), Description: "Count of teams within the organization."},
		{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT, Hydrate: orgHydrateRepositoriesTotalDiskUsage, Transform: transform.FromValue(), Description: "Total disk usage for all repositories within the organization."},
		{Name: "private_repositories_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydratePrivateRepositoriesTotalCount, Transform: transform.FromValue(), Description: "Count of private repositories within the organization."},
		{Name: "public_repositories_total_count", Type: proto.ColumnType_INT, Hydrate: orgHydratePublicRepositoriesTotalCount, Transform: transform.FromValue(), Description: "Count of public repositories within the organization."},
	}
}

func gitHubOrganizationColumns() []*plugin.Column {
	return append(sharedOrganizationColumns(), sharedOrganizationCountColumns()...)
}

func tableGitHubOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization",
		Description: "GitHub Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("login"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubOrganizationList,
		},
		Columns: commonColumns(gitHubOrganizationColumns()),
	}
}

func tableGitHubOrganizationList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	login := d.EqualsQuals["login"].GetStringValue()

	plugin.Logger(ctx).Debug("github_organization", login)
	var query struct {
		RateLimit    models.RateLimit
		Organization models.OrganizationWithCounts `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	appendOrganizationColumnIncludes(&variables, d.QueryContext.Columns)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_organization", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_organization", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, query.Organization)

	return nil, nil
}

func hydrateOrganizationHooksFromV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(models.OrganizationWithCounts)
	login := org.Login
	var orgHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	client := connect(ctx, d)

	for {
		hooks, resp, err := client.Organizations.ListHooks(ctx, login, opt)
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		orgHooks = append(orgHooks, hooks...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return orgHooks, nil
}

func hydrateOrganizationDataFromV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(models.OrganizationWithCounts)
	login := org.Login

	client := connect(ctx, d)
	organization, _, err := client.Organizations.Get(ctx, login)
	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationDetailV3", err)
		return nil, err
	}

	return organization, nil
}
