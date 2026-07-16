package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableGitHubActionsRepositoryWorkflowJob() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_repository_workflow_job",
		Description: "WorkflowJob represents a repository action workflow job",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "run_id", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoWorkflowJobList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required, Operators: []string{"="}},
				{Name: "id", Require: plugin.Required, Operators: []string{"="}},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoWorkflowJobGet,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that specifies the workflow job."},
			{Name: "run_id", Type: proto.ColumnType_INT, Description: "The unique identifier of the workflow run."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique identifier of the workflow job."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the workflow job."},
			{Name: "workflow_name", Type: proto.ColumnType_STRING, Description: "The workflow name of the workflow job."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the workflow job."},
			{Name: "conclusion", Type: proto.ColumnType_STRING, Description: "The conclusion for workflow job."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the workflow job."},
			{Name: "run_url", Type: proto.ColumnType_STRING, Description: "The API address of the workflow run the job belongs to."},
			{Name: "check_run_url", Type: proto.ColumnType_STRING, Description: "The API address of the check run associated with the workflow job."},
			{Name: "head_branch", Type: proto.ColumnType_STRING, Description: "The head branch of the workflow job."},
			{Name: "head_sha", Type: proto.ColumnType_STRING, Description: "The head sha of the workflow job.", Transform: transform.FromField("HeadSHA")},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The address for the workflow job's GitHub web page.", Transform: transform.FromField("HTMLURL")},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The address for the workflow job GitHub web page.", Transform: transform.FromField("URL")},

			// Other columns
			{Name: "steps", Type: proto.ColumnType_JSON, Description: "The list of step details for the workflow job."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "The list of labels for the workflow job."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the workflow job was created."},
			{Name: "started_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("StartedAt").Transform(convertTimestamp), Description: "Time when the workflow job was started."},
			{Name: "completed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CompletedAt").Transform(convertTimestamp), Description: "Time when the workflow job was completed."},
			{Name: "run_attempt", Type: proto.ColumnType_INT, Description: "The attempt number of the workflow run."},
			{Name: "runner_id", Type: proto.ColumnType_INT, Description: "The unique identifier of the workflow job runner."},
			{Name: "runner_name", Type: proto.ColumnType_STRING, Description: "The name of the workflow job runner."},
			{Name: "runner_group_id", Type: proto.ColumnType_INT, Description: "The unique identifier of the job runner group."},
			{Name: "runner_group_name", Type: proto.ColumnType_STRING, Description: "The name of the workflow job runner group."},
		}),
	}
}

func tableGitHubRepoWorkflowJobList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	runId := d.EqualsQuals["run_id"].GetInt64Value()
	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(orgName)
	opts := &github.ListWorkflowJobsOptions{
		// The API defaults to "latest", which only returns jobs from the most recent
		// run attempt. "all" includes jobs from previous attempts.
		Filter:      "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		var (
			workflowJobs *github.Jobs
			resp         *github.Response
			err          error
		)
		workflowJobs, resp, err = client.Actions.ListWorkflowJobs(ctx, owner, repo, runId, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range workflowJobs.Jobs {
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

func tableGitHubRepoWorkflowJobGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	jobId := d.EqualsQuals["id"].GetInt64Value()
	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if jobId == 0 || orgName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(orgName)
	plugin.Logger(ctx).Trace("tableGitHubRepoWorkflowJobGet", "owner", owner, "repo", repo, "jobId", jobId)

	client := connect(ctx, d)

	var (
		workflowJob *github.WorkflowJob
		err         error
	)
	workflowJob, _, err = client.Actions.GetWorkflowJobByID(ctx, owner, repo, jobId)
	if err != nil {
		return nil, err
	}

	return workflowJob, nil
}
