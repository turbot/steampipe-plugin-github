package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/shurcooL/githubv4"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

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
			{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of account."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user."},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar"},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user."},
			{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The user's gravatar ID"},
			{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
			{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user."},
			{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographic location of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user."},
			{Name: "hireable", Type: proto.ColumnType_BOOL, Description: "Whether the user currently hireable."},
			{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user."},
			{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user."},
			{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user."},
			{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user."},
			{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was created.", Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user was last updated.", Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp)},
			//{Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP},
			{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator."},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories."},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories."},
			{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user."},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user."},
			{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
			{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
			{Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user."},

			{Name: "starred_repositories_count", Type: proto.ColumnType_INT, Description: "Total repositories the user has starred.", Hydrate: getGitHubUserContributions, Transform: transform.FromField("StarredRepositories.TotalCount")},
			{Name: "repositories_contributed_to_count", Type: proto.ColumnType_INT, Description: "Total repositories that the user recently contributed to.", Hydrate: getGitHubUserContributions, Transform: transform.FromField("RepositoriesContributedTo.TotalCount")},
			{Name: "contributions_collection", Type: proto.ColumnType_JSON, Description: "The collection of contributions this user has made to different repositories.", Hydrate: getGitHubUserContributions},
			{Name: "sponsoring", Type: proto.ColumnType_INT, Description: "Total users and organizations this entity is sponsoring.", Hydrate: getGitHubUserContributions, Transform: transform.FromField("Sponsoring.TotalCount")},
			{Name: "sponsors", Type: proto.ColumnType_INT, Description: "Total sponsors for this user or organization.", Hydrate: getGitHubUserContributions, Transform: transform.FromField("Sponsors.TotalCount")},
			{Name: "keys", Type: proto.ColumnType_JSON, Description: "The verified public keys for a user.", Transform: transform.FromValue(), Hydrate: getGitHubUserKeys},
			{Name: "gpg_keys", Type: proto.ColumnType_JSON, Description: "Repositories the user has contributed to, ordered by contribution rank, plus repositories the user has created.", Transform: transform.FromValue(), Hydrate: getGitHubUserGPGKeys},
		},
	}
}

//// HYDRATE FUNCTRIONS

var userQuery struct {
	User struct {
		Login               string
		StarredRepositories struct {
			TotalCount githubv4.Int
		}
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
		}
		ContributionsCollection struct {
			TotalIssueContributions                      githubv4.Int
			TotalCommitContributions                     githubv4.Int
			TotalPullRequestContributions                githubv4.Int
			TotalPullRequestReviewContributions          githubv4.Int
			TotalRepositoriesWithContributedCommits      githubv4.Int
			TotalRepositoriesWithContributedIssues       githubv4.Int
			TotalRepositoriesWithContributedPullRequests githubv4.Int
			TotalRepositoryContributions                 githubv4.Int
			RestrictedContributionsCount                 githubv4.Int
		}
		Sponsoring struct {
			TotalCount githubv4.Int
		}
		Sponsors struct {
			TotalCount githubv4.Int
		}
	} `graphql:"user(login: $login)"`
}

// Listing all users is not terribly useful, so we require a 'login' qual and essentially always
// do a 'get':  from GitHub API docs: https://developer.github.com/v3/users/#list-users:
//
//	    	Lists all users, in the order that they signed up on GitHub. This list includes personal user
//			accounts and organization accounts.
//	    	Note: Pagination is powered exclusively by the since parameter. Use the Link header to get
//			the URL for the next page of users.
func tableGitHubUserGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var login string

	if h.Item != nil {
		item := h.Item.(*github.User)
		logger.Trace("tableGitHubUserGet", item.String())
		login = *item.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	client := connect(ctx, d)

	type GetResponse struct {
		user *github.User
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Users.Get(ctx, login)
		return GetResponse{
			user: detail,
			resp: resp,
		}, err
	}
	getResponse, err := retryHydrate(ctx, d, h, getDetails)

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	user := getResp.user

	if user != nil {
		d.StreamListItem(ctx, user)
	}

	return nil, nil
}

func getGitHubUserContributions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	var login string

	if h.Item != nil {
		item := h.Item.(*github.User)
		login = *item.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	err := client.Query(ctx, &userQuery, variables)
	if err != nil {
		plugin.Logger(ctx).Error("github_user", "api_error", err)
		// if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
		// 	return nil, nil
		// }
		return nil, err
	}

	return userQuery.User, nil
}

func getGitHubUserGPGKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string

	if h.Item != nil {
		item := h.Item.(*github.User)
		login = *item.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	opts := github.ListOptions{PerPage: 100}

	type ListPageResponse struct {
		keys []*github.GPGKey
		resp *github.Response
	}

	client := connect(ctx, d)

	var keyList []*github.GPGKey

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		keys, resp, err := client.Users.ListGPGKeys(ctx, login, &opts)
		return ListPageResponse{
			keys: keys,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		keys := listResponse.keys
		resp := listResponse.resp

		keyList = append(keyList, keys...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return keyList, nil
}

func getGitHubUserKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string

	if h.Item != nil {
		item := h.Item.(*github.User)
		login = *item.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	opts := github.ListOptions{PerPage: 100}

	type ListPageResponse struct {
		keys []*github.Key
		resp *github.Response
	}

	client := connect(ctx, d)

	var keyList []*github.Key

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		keys, resp, err := client.Users.ListKeys(ctx, login, &opts)
		return ListPageResponse{
			keys: keys,
			resp: resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		keys := listResponse.keys
		resp := listResponse.resp

		keyList = append(keyList, keys...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return keyList, nil
}
