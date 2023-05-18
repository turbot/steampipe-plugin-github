package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func searchUserColumnsOld() []*plugin.Column {
	return []*plugin.Column{
		{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID of the user or organization."},
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user or organization."},
		{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the the user or organization."},
		{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of the user or organization."},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar."},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user or organization."},
		{Name: "blog", Type: proto.ColumnType_STRING, Description: "The blog address of the user or organization."},
		{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators."},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company the user works for."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was created."},
		{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the user or organization."},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The public email address of the user or organization."},
		{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The event URL of the user or organization."},
		{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the user or organization."},
		{Name: "followers_url", Type: proto.ColumnType_STRING, Description: "The URL to get list of followers."},
		{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the user or organization."},
		{Name: "following_url", Type: proto.ColumnType_STRING, Description: "The URL to get list of users followed by the user or organization."},
		{Name: "gists_url", Type: proto.ColumnType_STRING, Description: "The URL get the gists of the user or organization."},
		{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The gravatar id of the user or organization."},
		{Name: "hireable", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the user or organization is hireable."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user or organization."},
		{Name: "ldap_dn", Type: proto.ColumnType_STRING, Description: "The LDAP distinguished name of the user or organization."},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The URL of the user or organization."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the user or organization."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user or organization."},
		{Name: "organizations_url", Type: proto.ColumnType_STRING, Description: "The URL to get the organization details of the user or organization."},
		{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories by the user or organization."},
		{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists owned by the user or organization."},
		{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists owned by the user or organization."},
		{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories owned by the user or organization."},
		{Name: "received_events_url", Type: proto.ColumnType_STRING, Description: "The URL to get the received events of the user or organization."},
		{Name: "repos_url", Type: proto.ColumnType_STRING, Description: "The URL to get the repositories that the user or organization is part of."},
		{Name: "site_admin", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the user or organization is an administrator."},
		{Name: "starred_url", Type: proto.ColumnType_STRING, Description: "The URL to get the starred details of the user or organization."},
		{Name: "subscriptions_url", Type: proto.ColumnType_STRING, Description: "The URL to get subscription details of the user or organization."},
		{Name: "suspended_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was suspended."},
		{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories of the user or organization."},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The twitter username of the user or organization."},
		{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether two-factor authentication is enabled for the user or organization."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the user or the organization was updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL to get information regarding the user or organization."},
		{Name: "permissions", Type: proto.ColumnType_JSON, Description: "The permission details."},
		{Name: "plan", Type: proto.ColumnType_JSON, Description: "The plan details."},
		{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
	}
}

func gitHubSearchUserColumns() []*plugin.Column {
	userSearchCols := []*plugin.Column{
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user/organization."},
		{Name: "type", Type: proto.ColumnType_STRING, Description: "Indicates if item is User or Organization."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user/organization.", Transform: transform.FromField("Id")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user/organization.", Transform: transform.FromField("NodeId")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the user/organization."},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user/organization."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user/organization's GitHub page.", Transform: transform.FromField("Url")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user/organization was created.", Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user/organization was last updated.", Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user/organization's avatar.", Transform: transform.FromField("AvatarUrl", "Node.AvatarUrl")},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user."},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company on the users profile."},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the organization."},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the user/organization."},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "Twitter username of the user/organization."},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL listing user/organization's projects.", Transform: transform.FromField("ProjectsUrl")},
		{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user/organization."},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user/organization."},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user/organization."},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user/organization."},
		{Name: "is_bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub security bug bounty."},
		{Name: "is_campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program."},
		{Name: "is_developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub developer program member."},
		{Name: "is_employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee."},
		{Name: "is_following_you", Type: proto.ColumnType_BOOL, Description: "If true, user follows you."},
		{Name: "is_github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a member of the GitHub Stars Program.", Transform: transform.FromField("IsGitHubStar")},
		{Name: "is_hireable", Type: proto.ColumnType_BOOL, Description: "If true, user has marked themselves as for hire."},
		{Name: "is_site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is a site administrator."},
		{Name: "is_you", Type: proto.ColumnType_BOOL, Description: "If true, user is you."},
		{Name: "website_url", Type: proto.ColumnType_STRING, Description: "The URL pointing to the user/organization's public website/blog.", Transform: transform.FromField("WebsiteUrl")},
	}
	return append(defaultSearchColumns(), userSearchCols...)
}

func tableGitHubSearchUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_user",
		Description: "Find users via various criteria.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchUserList,
		},
		Columns: gitHubSearchUserColumns(),
	}
}

func tableGitHubSearchUserList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			UserCount int
			PageInfo  models.PageInfo
			Edges     []struct {
				TextMatches []models.TextMatch
				Node        userSearchNode
			}
		} `graphql:"search(type: USER, first: $pageSize, after: $cursor, query: $query)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"query":    githubv4.String(input),
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_user", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_user", "api_error", err)
			return nil, err
		}

		for _, item := range query.Search.Edges {
			d.StreamListItem(ctx, mapToUserSearchRow(&item.Node, &item.TextMatches))

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	return nil, nil
}

func mapToUserSearchRow(node *userSearchNode, matches *[]models.TextMatch) userSearchRow {
	var row userSearchRow

	row.TextMatches = *matches
	row.Type = node.Type

	switch node.Type {
	case "User":
		row.Id = node.User.Id
		row.NodeId = node.User.NodeId
		row.Name = node.User.Name
		row.Login = node.User.Login
		row.Email = node.User.Email
		row.CreatedAt = node.User.CreatedAt
		row.UpdatedAt = node.User.UpdatedAt
		row.Url = node.User.Url
		row.AvatarUrl = node.User.AvatarUrl
		row.Bio = node.User.Bio
		row.Company = node.User.Company
		row.Location = node.User.Location
		row.TwitterUsername = node.User.TwitterUsername
		row.ProjectsUrl = node.User.ProjectsUrl
		row.CanFollow = node.User.CanFollow
		row.CanSponsor = node.User.CanSponsor
		row.IsFollowing = node.User.IsFollowing
		row.IsSponsoring = node.User.IsSponsoring
		row.IsBountyHunter = node.User.IsBountyHunter
		row.IsCampusExpert = node.User.IsCampusExpert
		row.IsDeveloperProgramMember = node.User.IsDeveloperProgramMember
		row.IsYou = node.User.IsYou
		row.IsEmployee = node.User.IsEmployee
		row.IsFollowingYou = node.User.IsFollowingYou
		row.IsGitHubStar = node.User.IsGitHubStar
		row.IsHireable = node.User.IsHireable
		row.IsSiteAdmin = node.User.IsSiteAdmin
		row.WebsiteUrl = node.User.WebsiteUrl
	case "Organization":
		row.Id = node.Organization.Id
		row.NodeId = node.Organization.NodeId
		row.Name = node.Organization.Name
		row.Login = node.Organization.Login
		row.Email = node.Organization.Email
		row.CreatedAt = node.Organization.CreatedAt
		row.UpdatedAt = node.Organization.UpdatedAt
		row.Url = node.Organization.Url
		row.AvatarUrl = node.Organization.AvatarUrl
		row.Description = node.Organization.Description
		row.Location = node.Organization.Location
		row.TwitterUsername = node.Organization.TwitterUsername
		row.ProjectsUrl = node.Organization.ProjectsUrl
		row.CanSponsor = node.Organization.CanSponsor
		row.IsFollowing = node.Organization.IsFollowing
		row.IsSponsoring = node.Organization.IsSponsoring
		row.WebsiteUrl = node.Organization.WebsiteUrl
	}

	return row
}

type userSearchNode struct {
	Type         string              `graphql:"type: __typename"`
	User         models.User         `graphql:"... on User"`
	Organization models.Organization `graphql:"... on Organization"`
}

type userSearchRow struct {
	models.BasicUser
	Type                     string `json:"type"`
	AvatarUrl                string `json:"avatar_url"`
	Bio                      string `json:"bio"`
	Description              string `json:"description"`
	Company                  string `json:"company"`
	Location                 string `json:"location"`
	TwitterUsername          string `json:"twitter_username"`
	ProjectsUrl              string `json:"projects_url"`
	CanFollow                bool   `json:"can_follow"`
	CanSponsor               bool   `json:"can_sponsor"`
	IsFollowing              bool   `json:"is_following"`
	IsSponsoring             bool   `json:"is_sponsoring"`
	IsBountyHunter           bool   `json:"is_bounty_hunter"`
	IsCampusExpert           bool   `json:"is_campus_expert"`
	IsDeveloperProgramMember bool   `json:"is_developer_program_member"`
	IsYou                    bool   `json:"is_you"`
	IsEmployee               bool   `json:"is_employee"`
	IsFollowingYou           bool   `json:"is_following_you"`
	IsGitHubStar             bool   `json:"is_github_star"`
	IsHireable               bool   `json:"is_hireable"`
	IsSiteAdmin              bool   `json:"is_site_admin"`
	WebsiteUrl               string `json:"website_url"`
	TextMatches              []models.TextMatch
}

func tableGitHubSearchUserListOld(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchUserList")

	quals := d.EqualsQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.UsersSearchResult
		resp   *github.Response
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		result, resp, err := client.Search.Users(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchUserList", "error_Search.Users", err)
			return nil, err
		}

		return ListPageResponse{
			result: result,
			resp:   resp,
		}, nil
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			logger.Error("tableGitHubSearchUserList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		users := listResponse.result.Users
		resp := listResponse.resp

		for _, i := range users {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}
