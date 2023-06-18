package github

import (
	"context"
	pipelineConsts "github.com/argonsecurity/pipeline-parser/pkg/consts"
	pipelineHandler "github.com/argonsecurity/pipeline-parser/pkg/handler"
	pipelineModels "github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/ghodss/yaml"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubWorkflowColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the workflow."},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the workflow."},
		{Name: "path", Type: proto.ColumnType_STRING, Description: "Path of the workflow."},
		{Name: "line_count", Type: proto.ColumnType_INT, Description: "The line count of the workflow file."},
		{Name: "size", Type: proto.ColumnType_INT, Description: "Size in bytes of the workflow file."},
		{Name: "language", Type: proto.ColumnType_STRING, Description: "Language of the workflow file.", Transform: transform.FromField("Language.Name")},
		{Name: "text", Type: proto.ColumnType_STRING, Description: "Contents of the workflow file.", Transform: transform.FromField("Object.Blob.Text")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the workflow.", Transform: transform.FromField("Object.Blob.NodeId")},
		{Name: "is_truncated", Type: proto.ColumnType_BOOL, Description: "If true, the text has been truncated due to length exceeding limits.", Transform: transform.FromField("Object.Blob.IsTruncated")},
		{Name: "is_binary", Type: proto.ColumnType_BOOL, Description: "If true, file is binary and therefore contents will be displayed as null.", Transform: transform.FromField("Object.Blob.IsBinary")},
		{Name: "is_generated", Type: proto.ColumnType_BOOL, Description: "If true, this workflow file was generated."},
		{Name: "commit_sha", Type: proto.ColumnType_STRING, Description: "Commit SHA associated with this file.", Transform: transform.FromField("Object.Blob.CommitSha")},
		{Name: "commit_url", Type: proto.ColumnType_STRING, Description: "URL of the commit associated with this file.", Transform: transform.FromField("Object.Blob.CommitUrl")},
		{Name: "text_json", Type: proto.ColumnType_JSON, Description: "Contents of workflow file in JSON format.", Transform: transform.FromField("Object.Blob.Text").Transform(unmarshalYAML)},
		{Name: "pipeline", Type: proto.ColumnType_JSON, Description: "GitHub workflow in the generic pipeline entity format to be used across CI/CD platforms.", Transform: transform.FromField("Object.Blob.Text").Transform(decodeFileContentToPipeline)},
	}
}

func tableGitHubWorkflow() *plugin.Table {
	return &plugin.Table{
		Name:        "github_workflow",
		Description: "GitHub Workflows bundle project files for download by users.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubWorkflowList,
		},
		Columns: gitHubWorkflowColumns(),
	}
}

func tableGitHubWorkflowList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Tree struct {
					Entries []models.TreeEntry
				} `graphql:"... on Tree"`
			} `graphql:"object(expression: \"HEAD:.github/workflows\")"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}

	client := connectV4(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
	plugin.Logger(ctx).Debug(rateLimitLogString("github_workflow", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_workflow", "api_error", err)
		return nil, err
	}

	for _, workflow := range query.Repository.Object.Tree.Entries {
		if workflow.Extension == ".yml" || workflow.Extension == ".yaml" {
			d.StreamListItem(ctx, workflow)
		}
		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// toPipeline:: Converts the github workflow buffer to generic CI pipeline format
func toPipeline(buf []byte) (*pipelineModels.Pipeline, error) {
	return pipelineHandler.Handle(buf, pipelineConsts.GitHubPlatform)
}

// decodeFileContentToPipeline:: Converts the workflow decoded text to generic CI pipeline.
func decodeFileContentToPipeline(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	repContent, ok := d.Value.(string)
	if !ok {
		return nil, nil
	}

	pipeline, err := toPipeline([]byte(repContent))
	if err != nil {
		plugin.Logger(ctx).Error("github_workflow.decodeFileContentToPipeline", "Pipeline conversion error", err)
		return nil, err
	}

	return pipeline, nil
}

// UnmarshalYAML parse the yaml-encoded data and return the result
func unmarshalYAML(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	inputStr := types.SafeString(d.Value)
	var result interface{}
	if inputStr != "" {
		err := yaml.Unmarshal([]byte(inputStr), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
