package github

import (
	"context"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubOrganization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization",
		Description: "Github Organization",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubOrganizationList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("login"),
			Hydrate:    getOrganizationDetail,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "login", Type: proto.ColumnType_STRING},
			{Name: "type", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "html_url", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail, Transform: transform.FromField("HTMLURL")},
			{Name: "avatar_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("AvatarURL")},
			{Name: "billing_email", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "blog", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "collaborators", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "company", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: getOrganizationDetail},
			{Name: "default_repo_permission", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "description", Type: proto.ColumnType_STRING},
			{Name: "disk_usage", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "email", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "events_url", Type: proto.ColumnType_STRING},
			{Name: "followers", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "following", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "has_organization_projects", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "has_repository_projects", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "hooks_url", Type: proto.ColumnType_STRING},
			{Name: "id", Type: proto.ColumnType_INT},
			{Name: "is_verified", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "issues_url", Type: proto.ColumnType_STRING},
			{Name: "location", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "members_allowed_repository_creation_type", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_internal_repos", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_pages", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_private_repos", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_public_repos", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_repos", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members_can_create_repositories", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "members", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: tableGitHubOrganizationMembersGet},
			{Name: "members_url", Type: proto.ColumnType_STRING},
			{Name: "name", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "node_id", Type: proto.ColumnType_STRING},
			{Name: "owned_private_repos", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "plan_filled_seats", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.FilledSeats")},
			{Name: "plan_name", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Name")},
			{Name: "plan_private_repos", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.PrivateRepos")},
			{Name: "plan_seats", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Seats")},
			{Name: "plan_space", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail, Transform: transform.FromField("Plan.Space")},
			{Name: "private_gists", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "public_gists", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "public_members_url", Type: proto.ColumnType_STRING},
			{Name: "public_repos", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "repos_url", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail, Transform: transform.FromField("ReposURL")},
			{Name: "total_private_repos", Type: proto.ColumnType_INT, Hydrate: getOrganizationDetail},
			{Name: "twitter_username", Type: proto.ColumnType_STRING, Hydrate: getOrganizationDetail},
			{Name: "two_factor_requirement_enabled", Type: proto.ColumnType_BOOL, Hydrate: getOrganizationDetail},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: getOrganizationDetail},
			{Name: "url", Type: proto.ColumnType_STRING},
		},
	}
}

//// list ////

func tableGitHubOrganizationList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d.ConnectionManager)

	opt := &github.ListOptions{PerPage: 100}

	for {

		var orgs []*github.Organization
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			orgs, resp, err = client.Organizations.List(ctx, "", opt)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range orgs {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

//// hydrate functions ////

func getOrganizationDetail(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var login string

	if h.Item != nil {
		org := h.Item.(*github.Organization)
		login = *org.Login
	} else {
		login = d.KeyColumnQuals["login"].GetStringValue()
	}

	client := connect(ctx, d.ConnectionManager)

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

	client := connect(ctx, d.ConnectionManager)

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
