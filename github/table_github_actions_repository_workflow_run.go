package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsRepositoryWorkflowRun() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_repository_workflow_run",
		Description: "WorkflowRun represents a repository action workflow run",
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoWorkflowRunList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "event", Require: plugin.Optional},
				{Name: "head_branch", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "conclusion", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoWorkflowRunGet,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that specifies the workflow run."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unque identifier of the workflow run."},
			{Name: "event", Type: proto.ColumnType_STRING, Description: "The event for which workflow triggered off."},
			{Name: "workflow_id", Type: proto.ColumnType_STRING, Description: "The workflow id of the worflow run."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the worflow run."},
			{Name: "conclusion", Type: proto.ColumnType_STRING, Description: "The conclusion for workflow run."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the worflow run."},
			{Name: "run_number", Type: proto.ColumnType_INT, Description: "The number of time workflow has run."},
			{Name: "artifacts_url", Type: proto.ColumnType_STRING, Description: "The address for artifact GitHub web page."},
			{Name: "cancel_url", Type: proto.ColumnType_STRING, Description: "The address for workflow run cancel GitHub web page."},
			{Name: "check_suite_url", Type: proto.ColumnType_STRING, Description: "The address for the workflow check suite GitHub web page."},
			{Name: "head_branch", Type: proto.ColumnType_STRING, Description: "The head branch of the workflow run branch."},
			{Name: "head_sha", Type: proto.ColumnType_STRING, Description: "The head sha of the workflow run.", Transform: transform.FromField("HeadSHA")},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The address for the organization's GitHub web page.", Transform: transform.FromField("HTMLURL")},
			{Name: "jobs_url", Type: proto.ColumnType_STRING, Description: "The address for the workflow job GitHub web page."},
			{Name: "logs_url", Type: proto.ColumnType_STRING, Description: "The address for the workflow logs GitHub web page."},
			{Name: "rerun_url", Type: proto.ColumnType_STRING, Description: "The address for workflow rerun GitHub web page."},
			{Name: "previous_attempt_url", Type: proto.ColumnType_STRING, Description: "The address for the previous attempt GitHub web page."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The address for the workflow run GitHub web page.", Transform: transform.FromField("URL")},
			{Name: "workflow_url", Type: proto.ColumnType_STRING, Description: "The address for workflow GitHub web page."},

			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the workflow run was created."},
			{Name: "head_commit", Type: proto.ColumnType_JSON, Description: "The head commit details for workflow run."},
			{Name: "head_repository", Type: proto.ColumnType_JSON, Description: "The head repository info for the workflow run."},
			{Name: "pull_requests", Type: proto.ColumnType_JSON, Description: "The pull request details for the workflow run."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The repository info for the workflow run."},
			{Name: "run_attempt", Type: proto.ColumnType_INT, Description: "The attempt number of the workflow run."},
			{Name: "run_started_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("RunStartedAt").Transform(convertTimestamp), Description: "Time when the workflow run was started."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "Time when the workflow run was updated."},
			{Name: "actor", Type: proto.ColumnType_JSON, Description: "The user whom initiated the first instance of this workflow run."},
			{Name: "actor_login", Type: proto.ColumnType_STRING, Description: "The login of the user whom initiated the first instance of the workflow run.", Transform: transform.FromField("Actor.Login")},
			{Name: "triggering_actor", Type: proto.ColumnType_JSON, Description: "The user whom initiated the latest instance of this workflow run."},
			{Name: "triggering_actor_login", Type: proto.ColumnType_STRING, Description: "The login of the user whom initiated the latest instance of this workflow run.", Transform: transform.FromField("TriggeringActor.Login")},
		}),
	}
}

func tableGitHubRepoWorkflowRunList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(orgName)
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	equalQuals := d.EqualsQuals
	if equalQuals["event"] != nil {
		if equalQuals["event"].GetStringValue() != "" {
			opts.Event = equalQuals["event"].GetStringValue()
		}
	}
	if equalQuals["head_branch"] != nil {
		if equalQuals["head_branch"].GetStringValue() != "" {
			opts.Branch = equalQuals["head_branch"].GetStringValue()
		}
	}
	if equalQuals["status"] != nil {
		if equalQuals["status"].GetStringValue() != "" {
			opts.Status = equalQuals["status"].GetStringValue()
		}
	}

	// Status param can take the value from both status and conclusion column
	// https://docs.github.com/en/rest/reference/actions#workflow-runs
	if equalQuals["conclusion"] != nil {
		if opts.Status == "" {
			if equalQuals["conclusion"].GetStringValue() != "" {
				opts.Status = equalQuals["conclusion"].GetStringValue()
			}
		}
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		workflowRuns, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range workflowRuns.WorkflowRuns {
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

func tableGitHubRepoWorkflowRunGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	runId := d.EqualsQuals["id"].GetInt64Value()
	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if runId == 0 || orgName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(orgName)
	plugin.Logger(ctx).Trace("tableGitHubRepoWorkflowRunGet", "owner", owner, "repo", repo, "runId", runId)

	client := connect(ctx, d)

	workflowRun, _, err := client.Actions.GetWorkflowRunByID(ctx, owner, repo, runId)
	if err != nil {
		return nil, err
	}

	return workflowRun, nil
}
