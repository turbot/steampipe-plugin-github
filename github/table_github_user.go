package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"strings"

	"github.com/google/go-github/v48/github"

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
		Columns: []*plugin.Column{
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("Id")},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("NodeId")},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user's GitHub page.", Transform: transform.FromField("Url")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated."},
			{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Description: "If true, user has pinnable items."},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar.", Transform: transform.FromField("AvatarUrl")},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company on the users profile."},
			{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Description: "The estimated next GitHub sponsors payout for this user in cents (USD)."},
			{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Description: "If true, user has a GitHub sponsors listing."},
			{Name: "interaction_ability", Type: proto.ColumnType_JSON, Description: "The interaction ability settings for this user.", Transform: transform.FromField("InteractionAbility").NullIfZero()},
			{Name: "is_bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub security bug bounty."},
			{Name: "is_campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program."},
			{Name: "is_developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub developer program member."},
			{Name: "is_employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee."},
			{Name: "is_following_you", Type: proto.ColumnType_BOOL, Description: "If true, user follows you."},
			{Name: "is_github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a member of the GitHub Stars Program.", Transform: transform.FromField("IsGitHubStar")},
			{Name: "is_hireable", Type: proto.ColumnType_BOOL, Description: "If true, user has marked themselves as for hire."},
			{Name: "is_site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is a site administrator."},
			{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Description: "If true, this user is sponsoring you."},
			{Name: "is_you", Type: proto.ColumnType_BOOL, Description: "If true, user is you."},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the user."},
			{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Description: "The estimated monthly GitHub sponsors income for this user in cents (USD)."},
			{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Description: "How many more items this user can pin to their profile."},
			{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL listing user's projects.", Transform: transform.FromField("ProjectsUrl")},
			{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The user's pronouns."},
			{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Description: "The user's sponsors listing.", Transform: transform.FromField("SponsorsListing").NullIfZero()},
			{Name: "status", Type: proto.ColumnType_JSON, Description: "The user's status.", Transform: transform.FromField("Status").NullIfZero()},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "Twitter username of the user."},
			{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Description: "If true, you can change the pinned items for this user."},
			{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Description: "If true, you can create projects for this user."},
			{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user."},
			{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user."},
			{Name: "is_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user."},
			{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user."},
			{Name: "website_url", Type: proto.ColumnType_STRING, Description: "The URL pointing to the user's public website/blog.", Transform: transform.FromField("WebsiteUrl")},
			// Counts
			{Name: "repositories_total_disk_usage", Type: proto.ColumnType_INT, Description: "Total disk spaced used by the users repositories.", Transform: transform.FromField("Repositories.TotalDiskUsage")},
			{Name: "followers_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users this user follows.", Transform: transform.FromField("Followers.TotalCount")},
			{Name: "following_total_count", Type: proto.ColumnType_INT, Description: "Count of how many users follow this user.", Transform: transform.FromField("Following.TotalCount")},
			{Name: "public_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of public repositories for the user.", Transform: transform.FromField("PublicRepositories.TotalCount")},
			{Name: "private_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of private repositories for the user.", Transform: transform.FromField("PrivateRepositories.TotalCount")},
			{Name: "public_gists_total_count", Type: proto.ColumnType_INT, Description: "Count of public gists for the user.", Transform: transform.FromField("PublicGists.TotalCount")},
			{Name: "issues_total_count", Type: proto.ColumnType_INT, Description: "Count of issues associated with the user.", Transform: transform.FromField("Issues.TotalCount")},
			{Name: "organizations_total_count", Type: proto.ColumnType_INT, Description: "Count of organizations the user belongs to.", Transform: transform.FromField("Organizations.TotalCount")},
			{Name: "public_keys_total_count", Type: proto.ColumnType_INT, Description: "Count of public keys associated with the user.", Transform: transform.FromField("PublicKeys.TotalCount")},
			{Name: "open_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of open pull requests associated with the user.", Transform: transform.FromField("OpenPullRequests.TotalCount")},
			{Name: "merged_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of merged pull requests associated with the user.", Transform: transform.FromField("MergedPullRequests.TotalCount")},
			{Name: "closed_pull_requests_total_count", Type: proto.ColumnType_INT, Description: "Count of closed pull requests associated with the user.", Transform: transform.FromField("ClosedPullRequests.TotalCount")},
			{Name: "packages_total_count", Type: proto.ColumnType_INT, Description: "Count of packages hosted by the user.", Transform: transform.FromField("Packages.TotalCount")},
			{Name: "pinned_items_total_count", Type: proto.ColumnType_INT, Description: "Count of items pinned on the users profile.", Transform: transform.FromField("PinnedItems.TotalCount")},
			{Name: "sponsoring_total_count", Type: proto.ColumnType_INT, Description: "Count of users that this user is sponsoring.", Transform: transform.FromField("Sponsoring.TotalCount")},
			{Name: "sponsors_total_count", Type: proto.ColumnType_INT, Description: "Count of users sponsoring this user.", Transform: transform.FromField("Sponsors.TotalCount")},
			{Name: "starred_repositories_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories the user has starred.", Transform: transform.FromField("StarredRepositories.TotalCount")},
			{Name: "watching_total_count", Type: proto.ColumnType_INT, Description: "Count of repositories being watched by the user.", Transform: transform.FromField("Watching.TotalCount")},
		},
	}
}

func tableGitHubUserGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string
	if h.Item != nil {
		item := h.Item.(*github.User)
		plugin.Logger(ctx).Trace("tableGitHubUserGet", item.String())
		login = *item.Login
	} else {
		login = d.EqualsQuals["login"].GetStringValue()
	}

	if login == "" {
		return nil, nil
	}

	client := connectV4(ctx, d)

	var query struct {
		RateLimit models.RateLimit
		User      models.UserWithCounts `graphql:"user(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

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
