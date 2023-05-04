package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user.", Transform: transform.FromField("Node.Login")},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("Node.DatabaseId")},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar.", Transform: transform.FromField("Node.AvatarUrl")},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user.", Transform: transform.FromField("Node.Url")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("Node.Id")},
		{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator.", Transform: transform.FromField("Node.IsSiteAdmin")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The team member's role.", Transform: transform.FromField("Role").Transform(transform.ToLower)},
		{Name: "status_message", Type: proto.ColumnType_STRING, Description: "The global status message of the team member.", Transform: transform.FromField("Node.Status.Message")},

		// Optional Columns
		// {Name: "followers_count", Type: proto.ColumnType_INT, Description: "Count of users the team member is followed by.", Transform: transform.FromField("Node.Followers.TotalCount")},
		// {Name: "following_count", Type: proto.ColumnType_INT, Description: "Count of users the team member is following.", Transform: transform.FromField("Node.Following.TotalCount")},
		// {Name: "repo_contribution_count", Type: proto.ColumnType_INT, Description: "Count of repositories the team member has contributed to (global).", Transform: transform.FromField("Node.RepositoriesContributedTo.TotalCount")},
		// {Name: "gists_count", Type: proto.ColumnType_INT, Description: "Count of gists the team member has published.", Transform: transform.FromField("Node.Gists.TotalCount")},
		// {Name: "starred_count", Type: proto.ColumnType_INT, Description: "Count of repositories the team member has starred.", Transform: transform.FromField("Node.StarredRepositories.TotalCount")},
	}
}

type teamMemberDetail struct {
	Role            string
	MemberAccessUrl string
	Node            struct {
		Login       string
		DatabaseId  int
		Id          string
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Url         string
		WebsiteUrl  string
		AvatarUrl   string
		ProjectsUrl string
		IsSiteAdmin bool
		Status      struct {
			Message string
		}

		// Optional Counts - decrease speed but give nice insight
		// Followers struct {
		// 	TotalCount int
		// }
		// Following struct {
		// 	TotalCount int
		// }
		// RepositoriesContributedTo struct {
		// 	TotalCount int
		// }
		// Gists struct {
		// 	TotalCount int
		// }
		// StarredRepositories struct {
		// 	TotalCount int
		// }
	}
}

var teamMembersQuery struct {
	Organization struct {
		Team struct {
			Members struct {
				TotalCount int
				PageInfo   struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
				Edges []teamMemberDetail
			} `graphql:"members(first: $pageSize, after: $cursor)"`
		} `graphql:"team(slug: $teamSlug)"`
	} `graphql:"organization(login: $login)"`
}

func tableGitHubTeamMember() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_member",
		Description: "GitHub members for a given team. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
				{Name: "role", Require: plugin.Optional},
			},
			Hydrate:           tableGitHubTeamMemberList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubTeamMemberColumns(),
	}
}

func tableGitHubTeamMemberList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	quals := d.EqualsQuals
	org := quals["organization"].GetStringValue()
	slug := quals["slug"].GetStringValue()

	pageSize := 100
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(pageSize) {
			pageSize = int(*limit)
		}
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"teamSlug": githubv4.String(slug),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &teamMembersQuery, variables)
		if err != nil {
			plugin.Logger(ctx).Error("github_team_member", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, member := range teamMembersQuery.Organization.Team.Members.Edges {
			d.StreamListItem(ctx, member)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !teamMembersQuery.Organization.Team.Members.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(teamMembersQuery.Organization.Team.Members.PageInfo.EndCursor)
	}

	return nil, nil
}
