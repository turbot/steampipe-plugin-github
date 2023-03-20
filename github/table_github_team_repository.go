package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func gitHubTeamRepositoryColumns() []*plugin.Column {
	repoColumns := gitHubRepositoryColumns()
	teamColumns := []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Description: "The organization the team is associated with.", Transform: transform.FromQual("organization")},
		{Name: "slug", Type: proto.ColumnType_STRING, Description: "The team slug name.", Transform: transform.FromQual("slug")},
		{Name: "permissions", Type: proto.ColumnType_JSON, Description: "The team's permissions for a repository.", Transform: transform.From(perissionsFromMap)},
	}

	return append(repoColumns, teamColumns...)
}

func tableGitHubTeamRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_team_repository",
		Description: "GitHub Repositories that a given team is associated with. GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
			},
			Hydrate:           tableGitHubTeamRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "slug", Require: plugin.Required},
				{Name: "full_name", Require: plugin.Required},
			},
			Hydrate:           tableGitHubTeamRepositoryGet,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: gitHubTeamRepositoryColumns(),
	}
}

//// LIST FUNCTION

func tableGitHubTeamRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	opt := &github.ListOptions{PerPage: 100}

	org := d.EqualsQuals["organization"].GetStringValue()
	slug := d.EqualsQuals["slug"].GetStringValue()

	type ListPageResponse struct {
		repos []*github.Repository
		resp  *github.Response
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		repos, resp, err := client.Teams.ListTeamReposBySlug(ctx, org, slug, opt)
		return ListPageResponse{
			repos: repos,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		repos := listResponse.repos
		resp := listResponse.resp

		for _, i := range repos {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

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

//// HYDRATE FUNCTIONS

func tableGitHubTeamRepositoryGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var org, slug, owner, repoName string
	if h.Item != nil {
		repo := h.Item.(*github.Repository)
		org = *repo.Organization.Login
		owner = *repo.Owner.Login
		repoName = *repo.Name
		slug = *h.Item.(*github.Team).Slug
	} else {
		org = d.EqualsQuals["organization"].GetStringValue()
		slug = d.EqualsQuals["slug"].GetStringValue()
		fullName := d.EqualsQuals["full_name"].GetStringValue()
		owner, repoName = parseRepoFullName(fullName)
	}

	client := connect(ctx, d)

	type GetResponse struct {
		repo *github.Repository
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Teams.IsTeamRepoBySlug(ctx, org, slug, owner, repoName)
		return GetResponse{
			repo: detail,
			resp: resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)

	if err != nil {
		return nil, err
	}

	repo := getResponse.(GetResponse).repo

	if repo != nil {
		return repo, nil
	}

	return nil, nil
}

func perissionsFromMap(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	permissions := d.HydrateItem.(*github.Repository).Permissions

	var arr []string
	for key, value := range permissions {
		if value {
			arr = append(arr, key)
		}
	}

	return arr, nil
}
