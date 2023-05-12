package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubTeamMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user.", Transform: transform.FromField("Node.Login")},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user.", Transform: transform.FromField("Node.Id")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user.", Transform: transform.FromField("Node.NodeId")},
		{Name: "role", Type: proto.ColumnType_STRING, Description: "The team member's role (MEMBER, MAINTAINER)."},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The email of the user.", Transform: transform.FromField("Node.Email")},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user's GitHub page.", Transform: transform.FromField("Node.Url")},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was created.", Transform: transform.FromField("Node.CreatedAt")},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when user was last updated.", Transform: transform.FromField("Node.UpdatedAt")},
		{Name: "any_pinnable_items", Type: proto.ColumnType_BOOL, Description: "If true, user has pinnable items.", Transform: transform.FromField("Node.AnyPinnableItems")},
		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar.", Transform: transform.FromField("Node.AvatarUrl")},
		{Name: "bio", Type: proto.ColumnType_STRING, Description: "The biography of the user.", Transform: transform.FromField("Node.Bio")},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The company on the users profile.", Transform: transform.FromField("Node.Company")},
		{Name: "estimated_next_sponsors_payout_in_cents", Type: proto.ColumnType_INT, Description: "The estimated next GitHub sponsors payout for this user in cents (USD).", Transform: transform.FromField("Node.EstimatedNextSponsorsPayoutInCents")},
		{Name: "has_sponsors_listing", Type: proto.ColumnType_BOOL, Description: "If true, user has a GitHub sponsors listing.", Transform: transform.FromField("Node.HasSponsorsListing")},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Description: "The interaction ability settings for this user.", Transform: transform.FromField("Node.InteractionAbility").NullIfZero()},
		{Name: "is_bounty_hunter", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub security bug bounty.", Transform: transform.FromField("Node.IsBountyHunter")},
		{Name: "is_campus_expert", Type: proto.ColumnType_BOOL, Description: "If true, user is a participant in the GitHub campus experts program.", Transform: transform.FromField("Node.IsCampusExpert")},
		{Name: "is_developer_program_member", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub developer program member.", Transform: transform.FromField("Node.IsDeveloperProgramMember")},
		{Name: "is_employee", Type: proto.ColumnType_BOOL, Description: "If true, user is a GitHub employee.", Transform: transform.FromField("Node.IsEmployee")},
		{Name: "is_following_you", Type: proto.ColumnType_BOOL, Description: "If true, user follows you.", Transform: transform.FromField("Node.IsFollowingYou")},
		{Name: "is_github_star", Type: proto.ColumnType_BOOL, Description: "If true, user is a member of the GitHub Stars Program.", Transform: transform.FromField("Node.IsGitHubStar")},
		{Name: "is_hireable", Type: proto.ColumnType_BOOL, Description: "If true, user has marked themselves as for hire.", Transform: transform.FromField("Node.IsHireable")},
		{Name: "is_site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is a site administrator.", Transform: transform.FromField("Node.IsSiteAdmin")},
		{Name: "is_sponsoring_you", Type: proto.ColumnType_BOOL, Description: "If true, this user is sponsoring you.", Transform: transform.FromField("Node.IsSponsoringYou")},
		{Name: "is_you", Type: proto.ColumnType_BOOL, Description: "If true, user is you.", Transform: transform.FromField("Node.IsYou")},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The location of the user.", Transform: transform.FromField("Node.Location")},
		{Name: "monthly_estimated_sponsors_income_in_cents", Type: proto.ColumnType_INT, Description: "The estimated monthly GitHub sponsors income for this user in cents (USD).", Transform: transform.FromField("Node.MonthlyEstimatedSponsorsIncomeInCents")},
		{Name: "pinned_items_remaining", Type: proto.ColumnType_INT, Description: "How many more items this user can pin to their profile.", Transform: transform.FromField("Node.PinnedItemsRemaining")},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Description: "The URL listing user's projects.", Transform: transform.FromField("Node.ProjectsUrl")},
		{Name: "pronouns", Type: proto.ColumnType_STRING, Description: "The user's pronouns.", Transform: transform.FromField("Node.Pronouns")},
		{Name: "sponsors_listing", Type: proto.ColumnType_JSON, Description: "The user's sponsors listing.", Transform: transform.FromField("Node.SponsorsListing").NullIfZero()},
		{Name: "status", Type: proto.ColumnType_JSON, Description: "The user's status.", Transform: transform.FromField("Node.Status").NullIfZero()},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "Twitter username of the user.", Transform: transform.FromField("Node.TwitterUsername")},
		{Name: "can_changed_pinned_items", Type: proto.ColumnType_BOOL, Description: "If true, you can change the pinned items for this user.", Transform: transform.FromField("Node.CanChangedPinnedItems")},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Description: "If true, you can create projects for this user.", Transform: transform.FromField("Node.CanCreateProjects")},
		{Name: "can_follow", Type: proto.ColumnType_BOOL, Description: "If true, you can follow this user.", Transform: transform.FromField("Node.CanFollow")},
		{Name: "can_sponsor", Type: proto.ColumnType_BOOL, Description: "If true, you can sponsor this user.", Transform: transform.FromField("Node.CanSponsor")},
		{Name: "is_following", Type: proto.ColumnType_BOOL, Description: "If true, you are following this user.", Transform: transform.FromField("Node.IsFollowing")},
		{Name: "is_sponsoring", Type: proto.ColumnType_BOOL, Description: "If true, you are sponsoring this user.", Transform: transform.FromField("Node.IsSponsoring")},
		{Name: "website_url", Type: proto.ColumnType_STRING, Description: "The URL pointing to the user's public website/blog.", Transform: transform.FromField("Node.WebsiteUrl")},
	}
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

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit    models.RateLimit
		Organization struct {
			Team struct {
				Members struct {
					TotalCount int
					PageInfo   struct {
						EndCursor   githubv4.String
						HasNextPage bool
					}
					Edges []models.TeamMemberWithRole
				} `graphql:"members(first: $pageSize, after: $cursor)"`
			} `graphql:"team(slug: $slug)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":    githubv4.String(org),
		"slug":     githubv4.String(slug),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_team_member", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_team_member", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, member := range query.Organization.Team.Members.Edges {
			d.StreamListItem(ctx, member)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Organization.Team.Members.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Organization.Team.Members.PageInfo.EndCursor)
	}

	return nil, nil
}
