package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
		{Name: "hooks", Type: proto.ColumnType_JSON, Description: "The API Hooks URL.", Hydrate: organizationHooksGet, Transform: transform.FromValue()},
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

//// TABLE DEFINITION

func tableGitHubOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization",
		Description: "GitHub Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("login"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           ListOrganizationDetail,
		},
		Columns: gitHubOrganizationColumns(),
	}
}

//// LIST FUNCTION

func ListOrganizationDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item, err := getOrganizationDetail(ctx, d, h)
	if err != nil {
		return nil, err
	}

	if item != nil {
		d.StreamListItem(ctx, item)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

type ListPageResponse struct {
	hooks []*github.Hook
	resp  *github.Response
}

func getOrganizationDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string
	if h.Item != nil {
		org := h.Item.(*github.Organization)
		// Check the null value for hydrated item, while accessing the inner level property of the null value it this throwing panic error
		if org == nil {
			return nil, nil
		}
		login = *org.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

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

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		plugin.Logger(ctx).Error("getOrganizationDetail", err)
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	org := getResp.org

	return org, nil
}

func tableGitHubOrganizationMembersGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("tableGitHubOrganizationMembersGet")

	org := h.Item.(*github.Organization)

	// Check the null value for hydrated item, while accessing the inner level property of the null value it this throwing panic error
	if org == nil {
		return nil, nil
	}
	orgName := *org.Login

	client := connect(ctx, d)

	var repositoryCollaborators []*github.User

	opt := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 100}}

	type ListPageResponse struct {
		users []*github.User
		resp  *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		users, resp, err := client.Organizations.ListMembers(ctx, orgName, opt)
		return ListPageResponse{
			users: users,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
		users := listResponse.users
		resp := listResponse.resp

		repositoryCollaborators = append(repositoryCollaborators, users...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return repositoryCollaborators, nil
}

func organizationHooksGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := h.Item.(*github.Organization)

	// Check the null value for hydrated item, while accessing the inner level property of the null value it this throwing panic error
	if org == nil {
		return nil, nil
	}
	orgName := *org.Login
	client := connect(ctx, d)

	var orgHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		hooks, resp, err := client.Organizations.ListHooks(ctx, orgName, opt)
		return ListPageResponse{
			hooks: hooks,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
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
