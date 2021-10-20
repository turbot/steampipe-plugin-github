package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubWorkflow(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_workflow",
		Description: "GitHub Workflows bundle project files for download by users.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubWorkflowList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			Hydrate:    tableGitHubWorkflowGet,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repositoryFullNameQual, Transform: transform.FromValue(), Description: "Full name of the repository that contains the workflow."},
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
		},
	}
}

func tableGitHubWorkflowList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	type ListPageResponse struct {
		workflows *github.Workflows
		resp      *github.Response
	}

	opts := &github.ListOptions{PerPage: 100}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		workflows, resp, err := client.Actions.ListWorkflows(ctx, owner, repo, opts)
		return ListPageResponse{
			workflows: workflows,
			resp:      resp,
		}, err
	}

	for {

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

		listResponse := listPageResponse.(ListPageResponse)
		workflows := listResponse.workflows
		resp := listResponse.resp

		if err != nil {
			return nil, err
		}

		for _, i := range workflows.Workflows {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

func tableGitHubWorkflowGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var id int64

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	if h.Item != nil {
		workflow := h.Item.(*github.Workflow)
		id = *workflow.ID
	} else {
		id = d.KeyColumnQuals["id"].GetInt64Value()
	}

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubWorkflowGet", "owner", owner, "repo", repo, "id", id)

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

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{shouldRetryError})
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	workflow := getResp.workflow

	return workflow, nil
}
