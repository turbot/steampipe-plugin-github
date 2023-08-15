package github

import (
	"context"
	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubRepositoryContent() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_content",
		Description: "List the content in a repository (list directory, or get file content",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubRepositoryContentList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "repository_content_path", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Description: "The full name of the repository (login/repo-name).", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name")},
			{Name: "type", Description: "The file type (directory or file).", Type: proto.ColumnType_STRING},
			{Name: "name", Description: "The file name.", Type: proto.ColumnType_STRING},
			{Name: "repository_content_path", Description: "The requested path in repository search.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_content_path")},
			{Name: "path", Description: "The path of the file.", Type: proto.ColumnType_STRING},
			{Name: "size", Description: "The size of the file (in MB).", Type: proto.ColumnType_INT},
			{Name: "content", Description: "The decoded file content (if the element is a file).", Type: proto.ColumnType_STRING, Transform: transform.From(transformFileContent), Hydrate: tableGitHubRepositoryContentGet},
			{Name: "target", Description: "Target is only set if the type is \"symlink\" and the target is not a normal file. If Target is set, Path will be the symlink path.", Type: proto.ColumnType_STRING},
			{Name: "sha", Description: "The sha of the file.", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA")},
			{Name: "url", Description: "URL of file's metadata.", Type: proto.ColumnType_STRING},
			{Name: "git_url", Description: "Git URL (with SHA) of the file.", Type: proto.ColumnType_STRING},
			{Name: "html_url", Description: "Raw file URL in GitHub.", Type: proto.ColumnType_STRING},
			{Name: "download_url", Description: "Download URL : it expires and can be be used just once.", Type: proto.ColumnType_STRING},
		},
	}
}

//// LIST FUNCTION

func tableGitHubRepositoryContentList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner, repo := parseRepoFullName(d.KeyColumnQuals["repository_full_name"].GetStringValue())
	var filterPath string
	if d.KeyColumnQuals["repository_content_path"] != nil {
		filterPath = d.KeyColumnQuals["repository_content_path"].GetStringValue()
	}
	plugin.Logger(ctx).Trace("tableGitHubRepositoryContentList", "owner", owner, "repo", repo, "path", filterPath)

	type ListPageResponse struct {
		repositoryContent []*github.RepositoryContent
		resp              *github.Response
	}
	client := connect(ctx, d)
	opt := &github.RepositoryContentGetOptions{}
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		fileContent, directoryContent, resp, err := client.Repositories.GetContents(ctx, owner, repo, filterPath, opt)

		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "api_error", err, "path", filterPath)
			return nil, err
		}

		if fileContent != nil {
			directoryContent = []*github.RepositoryContent{fileContent}
		}

		return ListPageResponse{
			repositoryContent: directoryContent,
			resp:              resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)
		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubRepositoryContentList", "retry_hydrate_error", err)
			return nil, err
		}

		for _, i := range listPageResponse.(ListPageResponse).repositoryContent {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if listPageResponse.(ListPageResponse).resp.NextPage == 0 {
			break
		}
	}
	return nil, nil
}

//// GET FUNCTION

func tableGitHubRepositoryContentGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner, repo := parseRepoFullName(d.KeyColumnQuals["repository_full_name"].GetStringValue())
	filterPath := *h.Item.(*github.RepositoryContent).Path

	plugin.Logger(ctx).Trace("tableGitHubRepositoryContentGet", "owner", owner, "repo", repo, "path", filterPath)

	type GetResponse struct {
		repositoryContent *github.RepositoryContent
		resp              *github.Response
	}

	client := connect(ctx, d)
	getFileContent := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		fileContent, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, filterPath, &github.RepositoryContentGetOptions{})

		if err != nil {
			plugin.Logger(ctx).Error("tableGitHubRepositoryContentGet", "api_error", err, "path", filterPath)
			return nil, err
		}

		return GetResponse{
			repositoryContent: fileContent,
			resp:              resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getFileContent)
	if err != nil {
		return nil, err
	}

	return getResponse.(GetResponse).repositoryContent, nil
}

func transformFileContent(_ context.Context, d *transform.TransformData) (interface{}, error) {
	content := d.HydrateItem.(*github.RepositoryContent)
	// directory use case. By definition, a directory doesn't have a raw content
	if content.Content == nil {
		return nil, nil
	}
	// empty file with "none" encoding,
	// or too big file (greater than 100MB, the RepositoryContent endpoint is not supported)
	if *content.Content == "" {
		return "", nil
	}
	return content.GetContent()
}
