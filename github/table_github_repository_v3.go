package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

func gitHubRepositoryV3Columns() []*plugin.Column {
	return []*plugin.Column{
		// Fields required for identity / hydration / mod
		{Name: "full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the repo.", Transform: transform.FromField("HTMLURL")},
		{Name: "owner_login", Type: proto.ColumnType_STRING, Description: "The user login name of the repository owner.", Transform: transform.FromField("Owner.Login")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the repository."},
		// Fields not in V4
		{Name: "hooks", Type: proto.ColumnType_JSON, Description: "The API Hooks URL.", Hydrate: repositoryHooksV3Get, Transform: transform.FromValue()},
		{Name: "topics", Type: proto.ColumnType_JSON, Description: "The topics (similar to tags or labels) associated with the repository."},
		// Fields not in V4 but not important
		{Name: "subscribers_count", Type: proto.ColumnType_INT, Description: "The number of users who have subscribed to the repository.", Hydrate: tableGitHubRepositoryV3Get},
		{Name: "has_downloads", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Downloads feature is enabled on the repository."},
		{Name: "has_pages", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Pages feature is enabled on the repository."},
		{Name: "network_count", Type: proto.ColumnType_INT, Description: "The number of member repositories in the network.", Hydrate: tableGitHubRepositoryV3Get},
	}
}

func tableGitHubRepositoryV3() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_v3",
		Description: "",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns:        plugin.SingleColumn("full_name"),
		},
		Columns: gitHubRepositoryV3Columns(),
	}
}

func tableGitHubRepositoryV3Get(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repoName string
	if h.Item != nil {
		repo := h.Item.(*github.Repository)
		owner = *repo.Owner.Login
		repoName = *repo.Name
	} else {
		owner = d.EqualsQuals["owner_login"].GetStringValue()
		repoName = d.EqualsQuals["name"].GetStringValue()
	}

	client := connect(ctx, d)

	type GetResponse struct {
		repo *github.Repository
		resp *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Repositories.Get(ctx, owner, repoName)
		return GetResponse{
			repo: detail,
			resp: resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}
	getResp := getResponse.(GetResponse)
	repo := getResp.repo

	if repo == nil {
		return nil, nil
	}

	return repo, nil
}

func repositoryHooksV3Get(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repo := h.Item.(*github.Repository)
	owner := *repo.Owner.Login
	repoName := *repo.Name

	client := connect(ctx, d)

	var repositoryHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		hooks, resp, err := client.Repositories.ListHooks(ctx, owner, repoName, opt)
		return ListHooksResponse{
			hooks: hooks,
			resp:  resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		listResponse := listPageResponse.(ListHooksResponse)
		hooks := listResponse.hooks
		resp := listResponse.resp
		repositoryHooks = append(repositoryHooks, hooks...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositoryHooks, nil
}
