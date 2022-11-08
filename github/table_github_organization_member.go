package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func gitHubOrganizationMemberColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the member is associated with.", Transform: transform.FromQual("organization")},
		{Name: "login", Type: proto.ColumnType_STRING, Description: "The login name of the user."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The ID of the user."},

		{Name: "avatar_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's avatar."},
		{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's events."},
		{Name: "followers_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's followers."},
		{Name: "following_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's following."},
		{Name: "gists_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's gists."},
		{Name: "gravatar_id", Type: proto.ColumnType_STRING, Description: "The user's gravatar ID."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub page for the user."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the user."},
		{Name: "organizations_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's organizations."},
		{Name: "received_events_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's received events."},
		{Name: "repos_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's repos."},
		{Name: "site_admin", Type: proto.ColumnType_BOOL, Description: "If true, user is an administrator."},
		{Name: "starred_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's stars."},
		{Name: "subscriptions_url", Type: proto.ColumnType_STRING, Description: "The URL of the user's subscriptions."},
		{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of account."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the user."},

		{Name: "role", Type: proto.ColumnType_STRING, Description: "The organization member's role.", Hydrate: tableGitHubOrganizationMemberGet},
		{Name: "state", Type: proto.ColumnType_STRING, Description: "The membership state.", Hydrate: tableGitHubOrganizationMemberGet},
		{Name: "two_factor_authentication", Type: proto.ColumnType_BOOL, Description: "If true, two-factor authentication is enabled."},
	}
}

func tableGitHubOrganizationMember() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_member",
		Description: "GitHub members for a given organization. GitHub Users are user accounts in GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "role", Require: plugin.Optional},
			},
			Hydrate:           tableGitHubOrganizationMemberList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubOrganizationMemberColumns(),
	}
}

type GithubListMembers struct {
	github.User
	TwoFactorAuthentication *bool
}

//// LIST FUNCTION

func tableGitHubOrganizationMemberList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Role:        "all",
	}
	quals := d.KeyColumnQuals
	org := quals["organization"].GetStringValue()

	// Additional filters
	if quals["role"] != nil {
		opt.Role = quals["role"].GetStringValue()
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

	disabledMember, err := tableGitHubOrganizationMember2FADisabledGet(ctx, d, h)
	if err != nil && !strings.Contains(err.Error(), "Only owners can use this filter") {
		return nil, err
	}

	if err != nil && strings.Contains(err.Error(), "Only owners can use this filter") {
		_, err := tableGitHubOrganizationMemberNo2FAList(ctx, d, h, client, org, opt)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = tableGitHubOrganizationMember2FAList(ctx, d, h, client, org, opt, disabledMember)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func tableGitHubOrganizationMemberGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	org := d.KeyColumnQuals["organization"].GetStringValue()

	user := h.Item.(GithubListMembers)
	username := *user.Login

	client := connect(ctx, d)

	type GetResponse struct {
		membership *github.Membership
		resp       *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Organizations.GetOrgMembership(ctx, username, org)
		return GetResponse{
			membership: detail,
			resp:       resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)

	if err != nil {
		plugin.Logger(ctx).Error("tableGitHubOrganizationMemberGet", "api-err", err)
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	membership := getResp.membership

	return membership, nil
}

func tableGitHubOrganizationMember2FADisabledGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	org := quals["organization"].GetStringValue()

	// Cache 2FA Disabled member list
	cacheKey := org + "2FADisabled"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]string), nil
	}

	client := connect(ctx, d)

	opt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Role:        "all",
		Filter:      "2fa_disabled",
	}

	// Additional filters
	if quals["role"] != nil {
		opt.Role = quals["role"].GetStringValue()
	}

	type ListPageResponse struct {
		members []*github.User
		resp    *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		members, resp, err := client.Organizations.ListMembers(ctx, org, opt)
		return ListPageResponse{
			members: members,
			resp:    resp,
		}, err
	}
	disabledMember := []string{}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubOrganizationMember2FADisabledGet", "api-err", err)
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
		members := listResponse.members
		resp := listResponse.resp

		for _, i := range members {
			if i != nil {
				disabledMember = append(disabledMember, *i.Login)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, disabledMember)

	return disabledMember, nil
}

func tableGitHubOrganizationMember2FAList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client, org string, opt *github.ListMembersOptions, disabledMember interface{}) (interface{}, error) {
	type ListPageResponse struct {
		members []*github.User
		resp    *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		members, resp, err := client.Organizations.ListMembers(ctx, org, opt)
		return ListPageResponse{
			members: members,
			resp:    resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubOrganizationMemberMFAList", "api-err", err)
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
		members := listResponse.members
		resp := listResponse.resp

		for _, i := range members {
			if i != nil {
				if !helpers.StringSliceContains(disabledMember.([]string), *i.Login) {
					d.StreamListItem(ctx, GithubListMembers{*i, types.Bool(true)})
				} else {
					d.StreamListItem(ctx, GithubListMembers{*i, types.Bool(false)})
				}
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

func tableGitHubOrganizationMemberNo2FAList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData, client *github.Client, org string, opt *github.ListMembersOptions) (interface{}, error) {
	type ListPageResponse struct {
		members []*github.User
		resp    *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		members, resp, err := client.Organizations.ListMembers(ctx, org, opt)
		return ListPageResponse{
			members: members,
			resp:    resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubOrganizationMemberNoMFAList", "api-err", err)
			return nil, err
		}
		listResponse := listPageResponse.(ListPageResponse)
		members := listResponse.members
		resp := listResponse.resp

		for _, i := range members {
			if i != nil {
				d.StreamListItem(ctx, GithubListMembers{*i, nil})
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
