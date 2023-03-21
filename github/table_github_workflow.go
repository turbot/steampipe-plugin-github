package github

import (
	"context"
	"encoding/base64"
	"strings"

	pipelineConsts "github.com/argonsecurity/pipeline-parser/pkg/consts"
	pipelineHandler "github.com/argonsecurity/pipeline-parser/pkg/handler"
	pipelineModels "github.com/argonsecurity/pipeline-parser/pkg/models"

	"github.com/ghodss/yaml"
	"github.com/google/go-github/v48/github"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubWorkflow(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_workflow",
		Description: "GitHub Workflows bundle project files for download by users.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubWorkflowList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubWorkflowGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the workflow."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the workflow."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the workflow."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "Path of the workflow."},

			// Other columns
			{Name: "badge_url", Type: proto.ColumnType_STRING, Description: "Badge URL for the workflow."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the workflow was created."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL for the workflow."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "Node where GitHub stores this data internally."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the workflow."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "Time when the workflow was updated."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL of the workflow."},
			{Name: "workflow_file_content", Type: proto.ColumnType_STRING, Hydrate: GitHubWorkflowFileContent, Transform: transform.FromValue().Transform(decodeFileContentBase64), Description: "Content of github workflow file in text format."},
			{Name: "workflow_file_content_json", Type: proto.ColumnType_JSON, Hydrate: GitHubWorkflowFileContent, Transform: transform.FromValue().Transform(decodeFileContentBase64).Transform(unmarshalYAML), Description: "Content of github workflow file in the JSON format."},
			{Name: "pipeline", Type: proto.ColumnType_JSON, Hydrate: GitHubWorkflowFileContent, Transform: transform.FromValue().Transform(decodeFileContentBase64).Transform(decodeFileContentToPipeline), Description: "Github workflow in the generic pipeline entity format to be used across CI/CD platforms."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubWorkflowList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	type ListPageResponse struct {
		workflows *github.Workflows
		resp      *github.Response
	}

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		workflows, resp, err := client.Actions.ListWorkflows(ctx, owner, repo, opts)
		return ListPageResponse{
			workflows: workflows,
			resp:      resp,
		}, err
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		workflows := listResponse.workflows
		resp := listResponse.resp

		for _, i := range workflows.Workflows {
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

		opts.Page = resp.NextPage
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func tableGitHubWorkflowGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := d.EqualsQuals["id"].GetInt64Value()
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubWorkflowGet", "owner", owner, "repo", repo, "id", id)

	client := connect(ctx, d)

	type GetResponse struct {
		workflow *github.Workflow
		resp     *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetWorkflowByID(ctx, owner, repo, id)
		return GetResponse{
			workflow: detail,
			resp:     resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	workflow := getResp.workflow

	return workflow, nil
}

func GitHubWorkflowFileContent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	workflow := h.Item.(*github.Workflow)
	if workflow.Path == nil {
		return nil, nil
	}

	id := d.EqualsQuals["id"].GetInt64Value()
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubWorkflowGet", "owner", owner, "repo", repo, "id", id)

	client := connect(ctx, d)

	type GetFileContentResponse struct {
		content *github.RepositoryContent
		resp    *github.Response
	}

	// Get the name of the default branch for the repository
	workflowUrlParts := strings.Split(*workflow.HTMLURL, "/")
	defaultBranch := "main"
	if len(workflowUrlParts) > 6 {
		defaultBranch = workflowUrlParts[6]
	}

	// Get workflow file content
	getFileContent := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		content, _, resp, err := client.Repositories.GetContents(ctx, owner, repo, *workflow.Path, &github.RepositoryContentGetOptions{Ref: defaultBranch})
		return GetFileContentResponse{
			content: content,
			resp:    resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getFileContent)
	if err != nil {
		// the workflow object exists, but the file is deleted
		if strings.Contains(err.Error(), "404 Not Found") {
			return nil, nil
		}
		return nil, err
	}

	getResp := getResponse.(GetFileContentResponse)
	content := getResp.content

	return content, nil
}

//// TRANSFORM FUNCTIONS

// decodeFileContentBase64:: Decode the workflow file content from Base64 encoded string to simple text
func decodeFileContentBase64(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	repContent, ok := d.Value.(*github.RepositoryContent)
	if !ok {
		return nil, nil
	}

	decodedText, err := base64.StdEncoding.DecodeString(*repContent.Content)
	if err != nil {
		plugin.Logger(ctx).Error("github_workflow.decodeFileContentBase64", "Decoding file content error", err)
		return nil, err
	}

	return string(decodedText), nil
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
