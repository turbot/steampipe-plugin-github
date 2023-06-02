package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"strings"

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
		{Name: "announcement", Type: proto.ColumnType_STRING, Transform: transform.FromField("Announcement", "Node.Announcement"), Description: "The text of the announcement."},
		{Name: "announcement_expires_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("AnnouncementExpiresAt", "Node.AnnouncementExpiresAt").NullIfZero().Transform(convertTimestamp), Description: "The expiration date of the announcement, if any."},
		{Name: "announcement_user_dismissible", Type: proto.ColumnType_BOOL, Transform: transform.FromField("AnnouncementUserDismissible", "Node.AnnouncementUserDismissible"), Description: "If true, the announcement can be dismissed by the user."},
		{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Transform: transform.FromField("AnyPinnableItems", "Node.AnyPinnableItems"), Description: "If true, this organization has items that can be pinned to their profile."},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("AvatarUrl", "Node.AvatarUrl"), Description: "URL pointing to the organization's public avatar."},
		{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Transform: transform.FromField("EstimatedNextSponsorsPayoutInCents", "Node.EstimatedNextSponsorsPayoutInCents"), Description: "The estimated next GitHub Sponsors payout for this organization in cents (USD)."},
		{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasSponsorsListing", "Node.HasSponsorsListing"), Description: "If true, this organization has a GitHub Sponsors listing."},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Transform: transform.FromField("InteractionAbility", "Node.InteractionAbility").NullIfZero(), Description: "The interaction ability settings for this organization."},
		{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsSponsoringYou", "Node.IsSponsoringYou"), Description: "If true, you are sponsored by this organization."},
		{Name: "is_verified", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsVerified", "Node.IsVerified"), Description: "If true, the organization has verified its profile email and website."},
		{Name: "location", Type: proto.ColumnType_STRING, Transform: transform.FromField("Location", "Node.Location"), Description: "The organization's public profile location."},
		{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Transform: transform.FromField("MonthlyEstimatedSponsorsIncomeInCents", "Node.MonthlyEstimatedSponsorsIncomeInCents"), Description: "The estimated monthly GitHub Sponsors income for this organization in cents (USD)."},
		{Name: "new_team_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewTeamUrl", "Node.NewTeamUrl"), Description: "URL for creating a new team."},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Transform: transform.FromField("PinnedItemsRemaining", "Node.PinnedItemsRemaining"), Description: "Returns how many more items this organization can pin to their profile."},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ProjectsUrl", "Node.ProjectsUrl"), Description: "URL listing organization's projects."},
		{Name: "saml_identity_provider", Type: proto.ColumnType_JSON, Transform: transform.FromField("SamlIdentityProvider", "Node.SamlIdentityProvider").NullIfZero(), Description: "The Organization's SAML identity provider. Visible to (1) organization owners, (2) organization owners' personal access tokens (classic) with read:org or admin:org scope, (3) GitHub App with an installation token with read or write access to members, else null."},
		{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Transform: transform.FromField("SponsorsListing", "Node.SponsorsListing").NullIfZero(), Description: "The GitHub sponsors listing for this organization."},
		{Name: "teams_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TeamsUrl", "Node.TeamsUrl"), Description: "URL listing organization's teams."},
		{Name: "total_sponsorship_amount_as_sponsor_in_cents", Type: proto.ColumnType_INT, Transform: transform.FromField("TotalSponsorshipAmountAsSponsorInCents", "Node.TotalSponsorshipAmountAsSponsorInCents").NullIfZero(), Description: "The amount in United States cents (e.g., 500 = $5.00 USD) that this entity has spent on GitHub to fund sponsorships. Only returns a value when viewed by the user themselves or by a user who can manage sponsorships for the requested organization."},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Transform: transform.FromField("TwitterUsername", "Node.TwitterUsername"), Description: "The organization's Twitter username."},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanAdminister", "Node.CanAdminister"), Description: "If true, you can administer the organization."},
		{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanChangedPinnedItems", "Node.CanChangedPinnedItems"), Description: "If true, you can change the pinned items on the organization's profile."},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanCreateProjects", "Node.CanCreateProjects"), Description: "If true, you can create projects for the organization."},
		{Name: "can_create_repositories", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanCreateRepositories", "Node.CanCreateRepositories"), Description: "If true, you can create repositories for the organization."},
		{Name: "can_create_teams", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanCreateTeams", "Node.CanCreateTeams"), Description: "If true, you can create teams within the organization."},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanSponsor", "Node.CanSponsor"), Description: "If true, you can sponsor this organization."},
		{Name: "is_a_member", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsAMember", "Node.IsAMember"), Description: "If true, you are an active member of the organization."},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsFollowing", "Node.IsFollowing"), Description: "If true, you are following the organization."},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsSponsoring", "Node.IsSponsoring"), Description: "If true, you are sponsoring the organization."},
		{Name: "website_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("WebsiteUrl", "Node.WebsiteUrl"), Description: "URL for the organization's public website."},
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
	}
}

func sharedOrganizationCountColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "members_with_role_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("MembersWithRole.TotalCount", "Node.MembersWithRole.TotalCount"), Description: "Count of members with a role within the organization."},
		{Name: "packages_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Packages.TotalCount", "Node.Packages.TotalCount"), Description: "Count of packages within the organization."},
		{Name: "pending_members_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("PendingMembers.TotalCount", "Node.PendingMembers.TotalCount"), Description: "Count of pending members within the organization."},
		{Name: "pinnable_items_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("PinnableItems.TotalCount", "Node.PinnableItems.TotalCount"), Description: "Count of pinnable items within the organization."},
		{Name: "pinned_items_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("PinnedItems.TotalCount", "Node.PinnedItems.TotalCount"), Description: "Count of itesm pinned to the organization's profile."},
		{Name: "projects_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Projects.TotalCount", "Node.Projects.TotalCount"), Description: "Count of projects within the organization."},
		{Name: "projects_v2_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("ProjectsV2.TotalCount", "Node.ProjectsV2.TotalCount"), Description: "Count of V2 projects within the organization."},
		{Name: "repositories_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Repositories.TotalCount", "Node.Repositories.TotalCount"), Description: "Count of all repositories within the organization."},
		{Name: "sponsoring_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Sponsoring.TotalCount", "Node.Sponsoring.TotalCount"), Description: "Count of users the organization is sponsoring."},
		{Name: "sponsors_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Sponsors.TotalCount", "Node.Sponsors.TotalCount"), Description: "Count of sponsors the organization has."},
		{Name: "teams_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Teams.TotalCount", "Node.Teams.TotalCount"), Description: "Count of teams within the organization."},
		{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT, Transform: transform.FromField("Repositories.TotalDiskUsage", "Node.Repositories.TotalDiskUsage"), Description: "Total disk usage for all repositories within the organization."},
		{Name: "private_repositories_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("PrivateRepositories.TotalCount", "Node.PrivateRepositories.TotalCount"), Description: "Count of private repositories within the organization."},
		{Name: "public_repositories_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("PublicRepositories.TotalCount", "Node.PublicRepositories.TotalCount"), Description: "Count of public repositories within the organization."},
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
		Columns: gitHubOrganizationColumns(),
	}
}

func tableGitHubOrganizationList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
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

	client := connect(ctx, d)

	var orgHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		hooks, resp, err := client.Organizations.ListHooks(ctx, login, opt)
		return ListHooksResponse{
			hooks: hooks,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		listResponse := listPageResponse.(ListHooksResponse)
		hooks := listResponse.hooks
		resp := listResponse.resp
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

	type GetResponse struct {
		org  *github.Organization
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Organizations.Get(ctx, login)
		return GetResponse{
			org:  detail,
			resp: resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, retryConfig())

	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationDetailV3", err)
		return nil, err
	}

	getResp := getResponse.(GetResponse)

	return getResp.org, nil
}

type ListHooksResponse struct {
	hooks []*github.Hook
	resp  *github.Response
}
