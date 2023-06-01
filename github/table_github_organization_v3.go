package github

// import (
// 	"context"
// 	"github.com/google/go-github/v48/github"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
// 	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
// 	"strings"
// )
//
// func gitHubOrganizationV3Columns() []*plugin.Column {
// 	return []*plugin.Column{
// 		// Top columns (login for identity & url for usage in mods as resource link)
// 		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the organization."},
// 		{Name: "url", Type: proto.ColumnType_STRING, Description: "The address for the organization's GitHub web page.", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("HTMLURL")},
// 		// Maybe keep these as v4 had to have an org_owner table to obtain these...?
// 		{Name: "billing_email", Type: proto.ColumnType_STRING, Description: "The email address for billing.", Hydrate: getOrganizationDetailV3},
// 		{Name: "two_factor_requirement_enabled", Type: proto.ColumnType_BOOL, Description: "If true, all members in the organization must have two factor authentication enabled.", Hydrate: getOrganizationDetailV3},
// 		// Fields not in v4
// 		{Name: "default_repo_permission", Type: proto.ColumnType_STRING, Description: "The default repository permissions for the organization.", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_allowed_repository_creation_type", Type: proto.ColumnType_STRING, Description: "Specifies which types of repositories non-admin organization members can create", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_can_create_internal_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create internal repositories.", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_can_create_pages", Type: proto.ColumnType_BOOL, Description: "If true, members can create pages.", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_can_create_private_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create private repositories.", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_can_create_public_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create public repositories.", Hydrate: getOrganizationDetailV3},
// 		{Name: "members_can_create_repos", Type: proto.ColumnType_BOOL, Description: "If true, members can create repositories.", Hydrate: getOrganizationDetailV3},
// 		{Name: "plan_filled_seats", Type: proto.ColumnType_INT, Description: "The number of used seats for the plan.", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("Plan.FilledSeats")},
// 		{Name: "plan_name", Type: proto.ColumnType_STRING, Description: "The name of the GitHub plan.", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("Plan.Name")},
// 		{Name: "plan_private_repos", Type: proto.ColumnType_INT, Description: "The number of private repositories for the plan.", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("Plan.PrivateRepos")},
// 		{Name: "plan_seats", Type: proto.ColumnType_INT, Description: "The number of available seats for the plan", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("Plan.Seats")},
// 		{Name: "plan_space", Type: proto.ColumnType_INT, Description: "The total space allocated for the plan.", Hydrate: getOrganizationDetailV3, Transform: transform.FromField("Plan.Space")},
// 		{Name: "hooks", Type: proto.ColumnType_JSON, Description: "The Hooks of the organization.", Hydrate: organizationHooksGetV3, Transform: transform.FromValue()},
// 		// Below aren't really required, could remove.
// 		{Name: "followers", Type: proto.ColumnType_INT, Description: "The number of users following the organization.", Hydrate: getOrganizationDetailV3},
// 		{Name: "following", Type: proto.ColumnType_INT, Description: "The number of users followed by the organization.", Hydrate: getOrganizationDetailV3},
// 		{Name: "collaborators", Type: proto.ColumnType_INT, Description: "The number of collaborators for the organization.", Hydrate: getOrganizationDetailV3},
// 		{Name: "has_organization_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use organization projects.", Hydrate: getOrganizationDetailV3},
// 		{Name: "has_repository_projects", Type: proto.ColumnType_BOOL, Description: "If true, the organization can use repository projects.", Hydrate: getOrganizationDetailV3},
//
// 		/*
// 		* Note: an alternative approach would be to have some extra tables `org_default_perms` `org_plan` `org_hooks` and populate these from v3 api?
// 		 */
// 	}
// }
//
// func tableGitHubOrganizationV3() *plugin.Table {
// 	return &plugin.Table{
// 		Name:        "github_organization_v3",
// 		Description: "",
// 		List: &plugin.ListConfig{
// 			KeyColumns:        plugin.SingleColumn("login"),
// 			ShouldIgnoreError: isNotFoundError([]string{"404"}),
// 			Hydrate:           tableGitHubOrganizationV3List,
// 		},
// 		Columns: gitHubOrganizationV3Columns(),
// 	}
// }
//
// func tableGitHubOrganizationV3List(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	item, err := getOrganizationDetailV3(ctx, d, h)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	if item != nil {
// 		d.StreamListItem(ctx, item)
// 	}
// 	return nil, nil
// }
//
// type ListHooksResponse struct {
// 	hooks []*github.Hook
// 	resp  *github.Response
// }
//
// func getOrganizationDetailV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	var login string
// 	if h.Item != nil {
// 		org := h.Item.(*github.Organization)
// 		// Check the null value for hydrated item, while accessing the inner level property of the null value it this throwing panic error
// 		if org == nil {
// 			return nil, nil
// 		}
// 		login = *org.Login
// 	} else {
// 		login = d.EqualsQuals["login"].GetStringValue()
// 	}
//
// 	client := connect(ctx, d)
//
// 	type GetResponse struct {
// 		org  *github.Organization
// 		resp *github.Response
// 	}
//
// 	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 		detail, resp, err := client.Organizations.Get(ctx, login)
// 		return GetResponse{
// 			org:  detail,
// 			resp: resp,
// 		}, err
// 	}
//
// 	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, retryConfig())
//
// 	if err != nil {
// 		plugin.Logger(ctx).Error("getOrganizationDetailV3", err)
// 		return nil, err
// 	}
//
// 	getResp := getResponse.(GetResponse)
// 	org := getResp.org
//
// 	return org, nil
// }
//
// func organizationHooksGetV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	org := h.Item.(*github.Organization)
//
// 	// Check the null value for hydrated item, while accessing the inner level property of the null value it this throwing panic error
// 	if org == nil {
// 		return nil, nil
// 	}
// 	orgName := *org.Login
// 	client := connect(ctx, d)
//
// 	var orgHooks []*github.Hook
// 	opt := &github.ListOptions{PerPage: 100}
//
// 	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 		hooks, resp, err := client.Organizations.ListHooks(ctx, orgName, opt)
// 		return ListHooksResponse{
// 			hooks: hooks,
// 			resp:  resp,
// 		}, err
// 	}
//
// 	for {
// 		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
// 		if err != nil && strings.Contains(err.Error(), "Not Found") {
// 			return nil, nil
// 		} else if err != nil {
// 			return nil, err
// 		}
// 		listResponse := listPageResponse.(ListHooksResponse)
// 		hooks := listResponse.hooks
// 		resp := listResponse.resp
// 		orgHooks = append(orgHooks, hooks...)
// 		if resp.NextPage == 0 {
// 			break
// 		}
// 		opt.Page = resp.NextPage
// 	}
// 	return orgHooks, nil
// }
