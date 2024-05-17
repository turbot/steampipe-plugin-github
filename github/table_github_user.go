package github

import (
	"context"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubUser() *plugin.Table {
	return &plugin.Table{
		Name:        "github_user",
		Description: "GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("login"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubUserGet,
		},
		Columns: commonColumns(tableGitHubUserColumns()),
	}
}

func tableGitHubUserColumns() []*plugin.Column {
	cols := sharedUserColumns()

	counts := []*plugin.Column{
		{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT, Description: "Total disk spaced used by the users repositories.", Transform: transform.FromValue(), Hydrate: userHydrateRepositoriesTotalDiskUsage},
		{Name: "followers_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users this user follows.", Transform: transform.FromValue(), Hydrate: userHydrateFollowersTotalCount},
		{Name: "following_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users follow this user.", Transform: transform.FromValue(), Hydrate: userHydrateFollowingTotalCount},
		{Name: "public_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of public repositories for the user.", Transform: transform.FromValue(), Hydrate: userHydratePublicRepositoriesTotalCount},
		{Name: "private_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of private repositories for the user.", Transform: transform.FromValue(), Hydrate: userHydratePrivateRepositoriesTotalCount},
		{Name: "public_gists_total_count", Type: proto.ColumnType_INT, Description: "Count of public gists for the user.", Transform: transform.FromValue(), Hydrate: userHydratePublicGistsTotalCount},
		{Name: "issues_total_count", Type: proto.ColumnType_INT, Description: "Count of issues associated with the user.", Transform: transform.FromValue(), Hydrate: userHydrateIssuesTotalCount},
		{Name: "organizations_total_count", Type: proto.ColumnType_INT, Description: "Count of organizations the user belongs to.", Transform: transform.FromValue(), Hydrate: userHydrateOrganizationsTotalCount},
		{Name: "public_keys_total_count", Type: proto.ColumnType_INT, Description: "Count of public keys associated with the user.", Transform: transform.FromValue(), Hydrate: userHydratePublicKeysTotalCount},
		{Name: "open_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of open pull requests associated with the user.", Transform: transform.FromValue(), Hydrate: userHydrateOpenPullRequestsTotalCount},
		{Name: "merged_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of merged pull requests associated with the user.", Transform: transform.FromValue(), Hydrate: userHydrateMergedPullRequestsTotalCount},
		{Name: "closed_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of closed pull requests associated with the user.", Transform: transform.FromValue(), Hydrate: userHydrateClosedPullRequestsTotalCount},
		{Name: "packages_total_count", Type: proto.ColumnType_INT, Description: "Count of packages hosted by the user.", Transform: transform.FromValue(), Hydrate: userHydratePackagesTotalCount},
		{Name: "pinned_items_total_count", Type: proto.ColumnType_INT, Description: "Count of items pinned on the users profile.", Transform: transform.FromValue(), Hydrate: userHydratePinnedItemsTotalCount},
		{Name: "sponsoring_total_count", Type: proto.ColumnType_INT, Description: "Count of users that this user is sponsoring.", Transform: transform.FromValue(), Hydrate: userHydrateSponsoringTotalCount},
		{Name: "sponsors_total_count", Type: proto.ColumnType_INT, Description: "Count of users sponsoring this user.", Transform: transform.FromValue(), Hydrate: userHydrateSponsorsTotalCount},
		{Name: "starred_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the user has starred.", Transform: transform.FromValue(), Hydrate: userHydrateStarredRepositoriesTotalCount},
		{Name: "watching_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories being watched by the user.", Transform: transform.FromValue(), Hydrate: userHydrateWatchingTotalCount},
	}

	cols = append(cols, counts...)

	return cols
}

func sharedUserColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user.", Transform: transform.FromField("Login", "Node.Login")},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("Id", "Node.Id")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user.", Transform: transform.FromField("Name", "Node.Name")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("NodeId", "Node.NodeId")},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user.", Transform: transform.FromField("Email", "Node.Email")},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user's GitHub page.", Transform: transform.FromField("Url", "Node.Url")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created.", Transform: transform.FromField("CreatedAt", "Node.CreatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated.", Transform: transform.FromField("UpdatedAt", "Node.UpdatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Description: "If true, user has pinnable items.", Transform: transform.FromValue(), Hydrate: userHydrateAnyPinnableItems},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar.", Transform: transform.FromValue(), Hydrate: userHydrateAvatarUrl},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user.", Transform: transform.FromValue(), Hydrate: userHydrateBio},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company on the users profile.", Transform: transform.FromValue(), Hydrate: userHydrateCompany},
		{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Description: "The estimated next GitHub sponsors payout for this user in cents (USD).", Transform: transform.FromValue(), Hydrate: userHydrateEstimatedNextSponsorsPayoutInCents},
		{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Description: "If true, user has a GitHub sponsors listing.", Transform: transform.FromValue(), Hydrate: userHydrateHasSponsorsListing},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Description: "The interaction ability settings for this user.", Transform: transform.FromValue().NullIfZero(), Hydrate: userHydrateInteractionAbility},
		{Name: "is_bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub security bug bounty.", Transform: transform.FromValue(), Hydrate: userHydrateIsBountyHunter},
		{Name: "is_campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program.", Transform: transform.FromValue(), Hydrate: userHydrateIsCampusExpert},
		{Name: "is_developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub developer program member.", Transform: transform.FromValue(), Hydrate: userHydrateIsDeveloperProgramMember},
		{Name: "is_employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee.", Transform: transform.FromValue(), Hydrate: userHydrateIsEmployee},
		{Name: "is_following_you", Type: proto.ColumnType_BOOL, Description: "If true, user follows you.", Transform: transform.FromValue(), Hydrate: userHydrateIsFollowingYou},
		{Name: "is_github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a member of the GitHub Stars Program.", Transform: transform.FromValue(), Hydrate: userHydrateIsGitHubStar},
		{Name: "is_hireable", Type: proto.ColumnType_BOOL, Description: "If true, user has marked themselves as for hire.", Transform: transform.FromValue(), Hydrate: userHydrateIsHireable},
		{Name: "is_site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is a site administrator.", Transform: transform.FromValue(), Hydrate: userHydrateIsSiteAdmin},
		{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Description: "If true, this user is sponsoring you.", Transform: transform.FromValue(), Hydrate: userHydrateIsSponsoringYou},
		{Name: "is_you", Type: proto.ColumnType_BOOL, Description: "If true, user is you.", Transform: transform.FromValue(), Hydrate: userHydrateIsYou},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the user.", Transform: transform.FromValue(), Hydrate: userHydrateLocation},
		{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Description: "The estimated monthly GitHub sponsors income for this user in cents (USD).", Transform: transform.FromValue(), Hydrate: userHydrateMonthlyEstimatedSponsorsIncomeInCents},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Description: "How many more items this user can pin to their profile.", Transform: transform.FromValue(), Hydrate: userHydratePinnedItemsRemaining},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL listing user's projects.", Transform: transform.FromValue(), Hydrate: userHydrateProjectsUrl},
		{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The user's pronouns.", Transform: transform.FromValue(), Hydrate: userHydratePronouns},
		{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Description: "The GitHub sponsors listing for this user.", Transform: transform.FromValue().NullIfZero(), Hydrate: userHydrateSponsorsListing},
		{Name: "status", Type: proto.ColumnType_JSON, Description: "The user's status.", Transform: transform.FromValue().NullIfZero(), Hydrate: userHydrateStatus},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "Twitter username of the user.", Transform: transform.FromValue(), Hydrate: userHydrateTwitterUsername},
		{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Description: "If true, you can change the pinned items for this user.", Transform: transform.FromValue(), Hydrate: userHydrateCanChangedPinnedItems},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Description: "If true, you can create projects for this user.", Transform: transform.FromValue(), Hydrate: userHydrateCanCreateProjects},
		{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user.", Transform: transform.FromValue(), Hydrate: userHydrateCanFollow},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user.", Transform: transform.FromValue(), Hydrate: userHydrateCanSponsor},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user.", Transform: transform.FromValue(), Hydrate: userHydrateIsFollowing},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user.", Transform: transform.FromValue(), Hydrate: userHydrateIsSponsoring},
		{Name: "website_url", Type: proto.ColumnType_STRING, Description: "The URL pointing to the user's public website/blog.", Transform: transform.FromValue(), Hydrate: userHydrateWebsiteUrl},
	}
}

func tableGitHubUserGet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	login := d.EqualsQuals["login"].GetStringValue()

	var query struct {
		RateLimit models.RateLimit
		User      models.UserWithCounts `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}
	appendUserWithCountColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_user", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_user", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, query.User)

	return nil, nil
}
