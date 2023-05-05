package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"strings"
	"time"

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
			// Top columns
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("DatabaseId")},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("Id")},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar", Transform: transform.FromField("AvatarUrl")},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user.", Transform: transform.FromField("Url")},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
			{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user.", Transform: transform.FromField("WebsiteUrl")},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographic location of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user."},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Description: "Whether the user currently hireable.", Transform: transform.FromField("IsHireable")},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user."},
			{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user.", Transform: transform.FromField("PublicRepos.TotalCount")},
			{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user.", Transform: transform.FromField("PublicGists.TotalCount")},
			{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user.", Transform: transform.FromField("Followers.TotalCount")},
			{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user.", Transform: transform.FromField("Following.TotalCount")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was last updated."},
			{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator.", Transform: transform.FromField("IsSiteAdmin")},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories.", Transform: transform.FromField("PrivateRepos.TotalCount")},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories.", Transform: transform.FromField("PrivateRepos.TotalCount")},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user.", Transform: transform.FromField("Repositories.TotalDiskUsage")},
			{Name: "status_message", Type: proto.ColumnType_STRING, Description: "The status message set by the user.", Transform: transform.FromField("Status.Message")},
			{Name: "status_expires_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the status message expires.", Transform: transform.FromField("Status.ExpiresAt")},
			{Name: "issues", Type: proto.ColumnType_INT, Description: "Count of issues authored by the user.", Transform: transform.FromField("Issues.TotalCount")},
			{Name: "organizations", Type: proto.ColumnType_INT, Description: "Count of organizations the user is a member of.", Transform: transform.FromField("Organizations.TotalCount")},
			{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The pronouns of the user."},
			{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Description: "If true, user has items that are pinnable."},
			{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Description: "If true, user has a GitHub sponsors listing."},
			{Name: "bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub bounty hunter program.", Transform: transform.FromField("IsBountyHunter")},
			{Name: "campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program.", Transform: transform.FromField("IsCampusExpert")},
			{Name: "developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub Developer Program member..", Transform: transform.FromField("IsDeveloperProgramMember")},
			{Name: "employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee.", Transform: transform.FromField("IsEmployee")},
			{Name: "github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub Stars program.", Transform: transform.FromField("IsGitHubStar")},
			{Name: "following_you", Type: proto.ColumnType_BOOL, Description: "If true, user is following you.", Transform: transform.FromField("IsFollowingViewer")},
			{Name: "sponsoring_you", Type: proto.ColumnType_BOOL, Description: "If true, user is sponsoring you.", Transform: transform.FromField("IsSponsoringViewer")},
			{Name: "can_change_pinned_items", Type: proto.ColumnType_BOOL, Description: "If true, you can change pinned items on the users profile.", Transform: transform.FromField("ViewerCanChangePinnedItems")},
			{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Description: "If true, you can create projects in the users account.", Transform: transform.FromField("ViewerCanCreateProjects")},
			{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user.", Transform: transform.FromField("ViewerCanFollow")},
			{Name: "are_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user.", Transform: transform.FromField("ViewerIsFollowing")},
			{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user.", Transform: transform.FromField("ViewerCanSponsor")},
			{Name: "are_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user.", Transform: transform.FromField("ViewerIsSponsoring")},
			{Name: "issue_comments", Type: proto.ColumnType_INT, Description: "Count of issue comments made by the user.", Transform: transform.FromField("IssueComments.TotalCount")},
			{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Description: "The estimated next GitHub Sponsors payout for this user in cents (USD).", Transform: transform.FromField("EstimatedNextSponsorsPayoutInCents")},
			{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Description: "The estimated monthly GitHub Sponsors income for this user in cents (USD).", Transform: transform.FromField("EstimatedNextSponsorsPayoutInCents")},
			{Name: "packages", Type: proto.ColumnType_INT, Description: "Count of packages owned by the user.", Transform: transform.FromField("Packages.TotalCount")},
			{Name: "pinned_items", Type: proto.ColumnType_INT, Description: "Count of items currently pinned to the users profile.", Transform: transform.FromField("PinnedItems.TotalCount")},
			{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Description: "Count of available slots to pin items on the users profile.", Transform: transform.FromField("PinnedItemsRemaining")},
			{Name: "projects", Type: proto.ColumnType_INT, Description: "Count of projects under the user.", Transform: transform.FromField("Projects.TotalCount")},
			{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL for the users projects.", Transform: transform.FromField("ProjectsUrl")},
			{Name: "public_keys", Type: proto.ColumnType_INT, Description: "Count of public keys for the user.", Transform: transform.FromField("PublicKeys.TotalCount")},
			{Name: "open_pull_requests", Type: proto.ColumnType_INT, Description: "Count of open pull requests the user has.", Transform: transform.FromField("OpenPullRequests.TotalCount")},
			{Name: "merged_pull_requests", Type: proto.ColumnType_INT, Description: "Count of merged pull requests the user has.", Transform: transform.FromField("MergedPullRequests.TotalCount")},
			{Name: "closed_pull_requests", Type: proto.ColumnType_INT, Description: "Count of closed pull requests the user has.", Transform: transform.FromField("ClosedPullRequests.TotalCount")},
			{Name: "social_accounts", Type: proto.ColumnType_JSON, Description: "Array of social accounts associated with the user.", Transform: transform.FromField("SocialAccounts.Nodes")},
			{Name: "sponsoring", Type: proto.ColumnType_INT, Description: "Count of users the user is sponsoring.", Transform: transform.FromField("Sponsoring.TotalCount")},
			{Name: "sponsors", Type: proto.ColumnType_INT, Description: "Count of sponsors the user has.", Transform: transform.FromField("Sponsors.TotalCount")},
			{Name: "starred_repositories", Type: proto.ColumnType_INT, Description: "Count of repositories the user has starred.", Transform: transform.FromField("StarredRepositories.TotalCount")},
			{Name: "watching", Type: proto.ColumnType_INT, Description: "Count of repositories the user is watching.", Transform: transform.FromField("Watching.TotalCount")},
			// {Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user.", Transform: transform.FromField("PrivateGists.TotalCount")},
			// {Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
			// {Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
			// {Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user."},
		},
	}
}

var userQuery struct {
	User struct {
		Login           string
		DatabaseId      int
		Name            string
		Id              string
		AvatarUrl       string
		Url             string
		Company         string
		WebsiteUrl      string
		Location        string
		Email           string
		IsHireable      bool
		IsSiteAdmin     bool
		Bio             string
		TwitterUsername string
		Followers       struct {
			TotalCount int
		}
		Following struct {
			TotalCount int
		}
		PublicRepos struct {
			TotalCount int
		} `graphql:"publicRepos: repositories(privacy: PUBLIC)"`
		PrivateRepos struct {
			TotalCount int
		} `graphql:"privateRepos: repositories(privacy: PRIVATE)"`
		Repositories struct {
			TotalDiskUsage int
		}
		PublicGists struct {
			TotalCount int
		} `graphql:"publicGists: gists(privacy: PUBLIC)"`
		// PrivateGists struct {
		// 	TotalCount int
		// } `graphql:"privateGists: gists(privacy: SECRET)"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Status    struct {
			Message   string
			ExpiresAt time.Time
		}
		Issues struct {
			TotalCount int
		}
		Organizations struct {
			TotalCount int
		}
		Pronouns                   string
		AnyPinnableItems           bool
		HasSponsorsListing         bool
		IsBountyHunter             bool
		IsCampusExpert             bool
		IsDeveloperProgramMember   bool
		IsEmployee                 bool
		IsGitHubStar               bool
		IsFollowingViewer          bool
		IsSponsoringViewer         bool
		ViewerCanChangePinnedItems bool
		ViewerCanCreateProjects    bool
		ViewerCanFollow            bool
		ViewerCanSponsor           bool
		ViewerIsFollowing          bool
		ViewerIsSponsoring         bool
		IssueComments              struct {
			TotalCount int
		}
		EstimatedNextSponsorsPayoutInCents    int
		MonthlyEstimatedSponsorsIncomeInCents int
		Packages                              struct {
			TotalCount int
		}
		PinnedItems struct {
			TotalCount int
		}
		PinnedItemsRemaining int
		Projects             struct {
			TotalCount int
		}
		ProjectsUrl string
		PublicKeys  struct {
			TotalCount int
		}
		OpenPullRequests struct {
			TotalCount int
		} `graphql:"openPullRequests: pullRequests(states: OPEN)"`
		MergedPullRequests struct {
			TotalCount int
		} `graphql:"mergedPullRequests: pullRequests(states: MERGED)"`
		ClosedPullRequests struct {
			TotalCount int
		} `graphql:"closedPullRequests: pullRequests(states: CLOSED)"`
		SocialAccounts struct {
			Nodes []struct {
				DisplayName string
				Provider    string
			}
		} `graphql:"socialAccounts(first: 10)"`
		Sponsoring struct {
			TotalCount int
		}
		Sponsors struct {
			TotalCount int
		}
		StarredRepositories struct {
			TotalCount int
		}
		Watching struct {
			TotalCount int
		}
	} `graphql:"user(login: $login)"`
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

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	err := client.Query(ctx, &userQuery, variables)
	if err != nil {
		plugin.Logger(ctx).Error("github_user", "api_error", err)
		if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, userQuery.User)

	return nil, nil
}
