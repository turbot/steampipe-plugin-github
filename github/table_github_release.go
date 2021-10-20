package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubRelease(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_release",
		Description: "GitHub Releases bundle project files for download by users.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubReleaseList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    tableGitHubReleaseGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repositoryFullNameQual, Transform: transform.FromValue(), Description: "Full name of the repository that contains the release."},

			// Other columns
			{Name: "assets", Type: proto.ColumnType_JSON, Description: "List of assets contained in the release."},
			{Name: "assets_url", Type: proto.ColumnType_STRING, Description: "Assets URL for the release."},
			{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.Login"), Description: "The login name of the user that created the release."},
			{Name: "body", Type: proto.ColumnType_STRING, Description: "Text describing the contents of the tag."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the release was created."},
			{Name: "draft", Type: proto.ColumnType_BOOL, Description: "True if this is a draft (unpublished) release."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL for the release."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the release."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the release."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "Node where GitHub stores this data internally."},
			{Name: "prerelease", Type: proto.ColumnType_BOOL, Description: "True if this is a prerelease version."},
			{Name: "published_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PublishedAt").Transform(convertTimestamp), Description: "Time when the release was published."},
			{Name: "tag_name", Type: proto.ColumnType_STRING, Description: "The name of the tag the release is associated with."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING, Description: "Tarball URL for the release."},
			{Name: "target_commitish", Type: proto.ColumnType_STRING, Description: "Specifies the commitish value that determines where the Git tag is created from. Can be any branch or commit SHA."},
			{Name: "upload_url", Type: proto.ColumnType_STRING, Description: "Upload URL for the release."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL of the release."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING, Description: "Zipball URL for the release."},
		},
	}
}

func tableGitHubReleaseList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opts := &github.ListOptions{PerPage: 100}

	type ListPageResponse struct {
		releases []*github.RepositoryRelease
		resp     *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		releases, resp, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
		return ListPageResponse{
			releases: releases,
			resp:     resp,
		}, err
	}

	for {

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		releases := listResponse.releases
		resp := listResponse.resp

		for _, i := range releases {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

func tableGitHubReleaseGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var id int64

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	if h.Item != nil {
		release := h.Item.(*github.RepositoryRelease)
		id = *release.ID
	} else {
		id = d.KeyColumnQuals["id"].GetInt64Value()
	}

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubReleaseGet", "owner", owner, "repo", repo, "id", id)

	client := connect(ctx, d)

	type GetResponse struct {
		release *github.RepositoryRelease
		resp    *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Repositories.GetRelease(ctx, owner, repo, id)
		return GetResponse{
			release: detail,
			resp:    resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	release := getResp.release

	return release, nil
}
