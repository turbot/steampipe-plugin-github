package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
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
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:    tableGitHubWorkflowGet,
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
		},
	}
}

//// LIST FUNCTION

func tableGitHubWorkflowList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
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
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

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
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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
	id := d.KeyColumnQuals["id"].GetInt64Value()
	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
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

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	workflow := getResp.workflow

	return workflow, nil
}
