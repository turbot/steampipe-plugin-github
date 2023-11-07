package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			{Name: "oid", Description: "The Git object ID.", Type: proto.ColumnType_STRING},
			{Name: "abbreviated_oid", Description: "An abbreviated version of the Git object ID.", Type: proto.ColumnType_STRING},
			{Name: "repository_content_path", Description: "The requested path in repository search.", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_content_path")},
			{Name: "path", Description: "The path of the file.", Type: proto.ColumnType_STRING},
			{Name: "path_raw", Description: "A Base64-encoded representation of the file's path.", Type: proto.ColumnType_STRING},
			{Name: "mode", Description: "The mode of the file.", Type: proto.ColumnType_INT},
			{Name: "size", Description: "The size of the file (in KB).", Type: proto.ColumnType_INT},
			{Name: "line_count", Description: "The number of lines available in the file.", Type: proto.ColumnType_INT},
			{Name: "content", Description: "The decoded file content (if the element is a file).", Type: proto.ColumnType_STRING},
			{Name: "is_generated", Description: "Whether or not this tree entry is generated.", Type: proto.ColumnType_BOOL},
			{Name: "is_binary", Description: "Indicates whether the Blob is binary or text.", Type: proto.ColumnType_BOOL},
			{Name: "commit_url", Description: "Git URL (with SHA) of the file.", Type: proto.ColumnType_STRING},
		},
	}
}

type ContentInfo struct {
	Oid            string
	AbbreviatedOid string
	Name           string
	Mode           int
	PathRaw        string
	IsGenerated    bool
	Path           string
	Size           int
	LineCount      int
	Type           string
	Content        string
	CommitUrl      string
	IsBinary       bool
}

//// LIST FUNCTION

func tableGitHubRepositoryContentList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	owner, repo := parseRepoFullName(d.EqualsQualString("repository_full_name"))
	var filterPath string
	if d.EqualsQualString("repository_content_path") != "" {
		filterPath = d.EqualsQualString("repository_content_path")
	}
	plugin.Logger(ctx).Trace("tableGitHubRepositoryContentList", "owner", owner, "repo", repo, "path", filterPath)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Tree struct {
					Oid            githubv4.String
					AbbreviatedOid githubv4.String
					Entries        []struct {
						Name        githubv4.String
						Path        githubv4.String
						Size        githubv4.Int
						LineCount   githubv4.Int
						Mode        githubv4.Int
						PathRaw     githubv4.String
						IsGenerated githubv4.Boolean
						Type        githubv4.String
						Object      struct {
							Blob struct {
								Oid            githubv4.String
								AbbreviatedOid githubv4.String
								Text           githubv4.String
								IsBinary       githubv4.Boolean
								CommitUrl      githubv4.String
							} `graphql:"... on Blob"`
							SubTree struct {
								Entries []struct {
									Name        githubv4.String
									Path        githubv4.String
									Size        githubv4.Int
									LineCount   githubv4.Int
									Mode        githubv4.Int
									PathRaw     githubv4.String
									IsGenerated githubv4.Boolean
									Type        githubv4.String
									Object      struct {
										Blob struct {
											Oid            githubv4.String
											AbbreviatedOid githubv4.String
											Text           githubv4.String
											IsBinary       githubv4.Boolean
											CommitUrl      githubv4.String
										} `graphql:"... on Blob"`
									}
								}
							} `graphql:"... on Tree"`
						}
					}
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: $expression)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repo":       githubv4.String(repo),
		"expression": githubv4.String("HEAD:" + filterPath),
	}

	client := connectV4(ctx, d)
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	// for {
	_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
	if err != nil {
		plugin.Logger(ctx).Error("github_repository_content", "api_error", err, "repository", repo)
		return nil, err
	}

	var contents []ContentInfo

	for _, data := range query.Repository.Object.Tree.Entries {
		contents = append(contents, ContentInfo{
			Oid:            string(data.Object.Blob.Oid),
			AbbreviatedOid: string(data.Object.Blob.AbbreviatedOid),
			Name:           string(data.Name),
			Mode:           int(data.Mode),
			PathRaw:        string(data.PathRaw),
			IsGenerated:    bool(data.IsGenerated),
			Path:           string(data.Path),
			Size:           int(data.Size),
			LineCount:      int(data.LineCount),
			Type:           string(data.Type),
			Content:        string(data.Object.Blob.Text),
			IsBinary:       bool(data.Object.Blob.IsBinary),
			CommitUrl:      string(data.Object.Blob.CommitUrl),
		})
		if len(data.Object.SubTree.Entries) > 0 {
			for _, subData := range data.Object.SubTree.Entries {
				contents = append(contents, ContentInfo{
					Oid:            string(subData.Object.Blob.Oid),
					AbbreviatedOid: string(subData.Object.Blob.AbbreviatedOid),
					Name:           string(subData.Name),
					Mode:           int(subData.Mode),
					PathRaw:        string(subData.PathRaw),
					IsGenerated:    bool(subData.IsGenerated),
					Path:           string(subData.Path),
					Size:           int(subData.Size),
					LineCount:      int(subData.LineCount),
					Type:           string(subData.Type),
					Content:        string(subData.Object.Blob.Text),
					IsBinary:       bool(subData.Object.Blob.IsBinary),
					CommitUrl:      string(subData.Object.Blob.CommitUrl),
				})
			}
		}
	}

	for _, c := range contents {
		d.StreamListItem(ctx, c)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
