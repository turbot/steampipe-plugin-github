package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func gitHubOrganizationColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the organization."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The organization name.", Hydrate: getOrganizationDetail},
		{Name: "type", Type: proto.ColumnType_STRING, Description: "The user type of the organization.", Hydrate: getOrganizationDetail},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The address for the organization's GitHub web page.", Hydrate: getOrganizationDetail, Transform: transform.FromField("HTMLURL")},

		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the organization's avatar.", Transform: transform.FromField("AvatarURL")},
		{Name: "billing_email", Type: proto.ColumnType_STRING, Description: "The email address for billing.", Hydrate: getOrganizationDetail},
		{Name: "blog", Type: proto.ColumnType_STRING, Description: "The URL of the organizations blog.", Hydrate: getOrganizationDetail},
		{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators for the organization.", Hydrate: getOrganizationDetail},
		{Name: "company", Type: proto.ColumnType_STRING, Description: "The name of the associated company.", Hydrate: getOrganizationDetail},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the organization was created.", Hydrate: getOrganizationDetail},
		{Name: "default_repo_permission", Type: proto.ColumnType_STRING, Description: "The default repository permissions for the organization.", Hydrate: getOrganizationDetail},
		{Name: "description", Type: proto.ColumnType_STRING, Description: "The organization description."},
		{Name: "disk_usage", Type: proto.ColumnType_INT, Description: "The total disk usage for the organization.", Hydrate: getOrganizationDetail},
		{Name: "email", Type: proto.ColumnType_STRING, Description: "The email address associated with the organization.", Hydrate: getOrganizationDetail},
		{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The API Events URL."},
		{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the organization.", Hydrate: getOrganizationDetail},
		{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the organization.", Hydrate: getOrganizationDetail},
		{Name: "has_organization_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use organization projects.", Hydrate: getOrganizationDetail},
		{Name: "has_repository_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use repository projects.", Hydrate: getOrganizationDetail},
		{Name: "hooks_url", Type: proto.ColumnType_STRING, Description: "The API Hooks URL."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the organization."},
		{Name: "is_verified", Type: proto.ColumnType_BOOL, Description: "If true, the organization has been verified with domain verification.", Hydrate: getOrganizationDetail},
		{Name: "issues_url", Type: proto.ColumnType_STRING, Description: "The API Issues URL."},
		{Name: "location", Type: proto.ColumnType_STRING, Description: "The geographical location.", Hydrate: getOrganizationDetail},
		{Name: "members_allowed_repository_creation_type", Type: proto.ColumnType_STRING, Description: "Specifies which types of repositories non-admin organization members can create", Hydrate: getOrganizationDetail},
		{Name: "members_can_create_internal_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create internal repositories.", Hydrate: getOrganizationDetail},
		{Name: "members_can_create_pages", Type: proto.ColumnType_BOOL, Description: "If true, members can create pages.", Hydrate: getOrganizationDetail},
		{Name: "members_can_create_private_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create private repositories.", Hydrate: getOrganizationDetail},
		{Name: "members_can_create_public_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create public repositories.", Hydrate: getOrganizationDetail},
		{Name: "members_can_create_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create repositories.", Hydrate: getOrganizationDetail},
		{Name: "members", Type: proto.ColumnType_JSON, Description: "An array of users that are members of the organization.", Transform: transform.FromValue(), Hydrate: tableGitHubOrganizationMembersGet},
		{Name: "member_logins", Type: proto.ColumnType_JSON, Description: "An array of user logins that are members of the organization.", Transform: transform.FromValue().Transform(filterUserLogins), Hydrate: tableGitHubOrganizationMembersGet},
		{Name: "members_url", Type: proto.ColumnType_STRING, Description: "The API Members URL."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the organization."},
		{Name: "owned_private_repos", Type: proto.ColumnType_INT, Description: "The number of owned private repositories.", Hydrate: getOrganizationDetail},
		{Name: "plan_filled_seats", Type: proto.ColumnType_INT, Description: "The number of used seats for the plan.", Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.FilledSeats")},
		{Name: "plan_name", Type: proto.ColumnType_STRING, Description: "The name of the GitHub plan.", Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Name")},
		{Name: "plan_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories for the plan.", Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.PrivateRepos")},
		{Name: "plan_seats", Type: proto.ColumnType_INT, Description: "The number of available seats for the plan", Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Seats")},
		{Name: "plan_space", Type: proto.ColumnType_INT, Description: "The total space allocated for the plan.", Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Space")},
		{Name: "private_gists", Type: proto.ColumnType_INT, Description: "The number of private gists.", Hydrate: getOrganizationDetail},
		{Name: "public_gists", Type: proto.ColumnType_INT, Description: "The number of public gists.", Hydrate: getOrganizationDetail},
		{Name: "public_members_url", Type: proto.ColumnType_STRING, Description: "The API Public Members URL."},
		{Name: "public_repos", Type: proto.ColumnType_INT, Description: "The number of public repositories.", Hydrate: getOrganizationDetail},
		{Name: "repos_url", Type: proto.ColumnType_STRING, Description: "The API Repos URL.", Hydrate: getOrganizationDetail, Transform: transform.FromField("ReposURL")},
		{Name: "total_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories.", Hydrate: getOrganizationDetail},
		{Name: "twitter_username", Type: proto.ColumnType_STRING, Description: "The organizations twitter handle.", Hydrate: getOrganizationDetail},
		{Name: "two_factor_requirement_enabled", Type: proto.ColumnType_BOOL, Description: "If true, all members in the organization must have two factor authentication enabled.", Hydrate: getOrganizationDetail},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the organization was last updated.", Hydrate: getOrganizationDetail},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL of the organization."},
	}
}

func tableGitHubOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization",
		Description: "GitHub Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("login"),
			Hydrate:    ListOrganizationDetail,
		},
		Columns: gitHubOrganizationColumns(),
	}
}

//// hydrate functions ////

func ListOrganizationDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item, err := getOrganizationDetail(ctx, d, h)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, item)
	return nil, nil
}

func getOrganizationDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string

	if h.Item != nil {
		org := h.Item.(*github.Organization)
		login = *org.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	client := connect(ctx, d)

	var detail *github.Organization
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		detail, resp, err = client.Organizations.Get(ctx, login)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return detail, nil
}

func tableGitHubOrganizationMembersGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	org := h.Item.(*github.Organization)
	orgName := *org.Login

	client := connect(ctx, d)

	var repositoryCollaborators []*github.User

	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {

		var users []*github.User
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			users, resp, err = client.Organizations.ListMembers(ctx, orgName, opt)
			logger.Info("Users", users)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range users {
			repositoryCollaborators = append(repositoryCollaborators, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	logger.Trace("OrganizationMembers", repositoryCollaborators)

	return repositoryCollaborators, nil
}
