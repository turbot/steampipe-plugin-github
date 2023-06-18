package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubActionsRepositoryRunner() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_repository_runner",
		Description: "The runner is the application that runs a job from a GitHub Actions workflow",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRunnerList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRunnerGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the runners."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromGo(), Description: "The unique identifier of the runner."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the runner."},
			{Name: "os", Type: proto.ColumnType_STRING, Transform: transform.FromField("OS"), Description: "The operating system of the runner."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the runner."},
			{Name: "busy", Type: proto.ColumnType_BOOL, Description: "Indicates wheather the runner is currently in use or not."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels represents a collection of labels attached to each runner."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubRunnerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(orgName)

	type ListPageResponse struct {
		runners *github.Runners
		resp    *github.Response
	}

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		runners, resp, err := client.Actions.ListRunners(ctx, owner, repo, opts)
		return ListPageResponse{
			runners: runners,
			resp:    resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		runners := listResponse.runners
		resp := listResponse.resp

		for _, i := range runners.Runners {
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

func tableGitHubRunnerGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	runnerId := d.EqualsQuals["id"].GetInt64Value()
	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameter
	if runnerId == 0 || orgName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(orgName)
	plugin.Logger(ctx).Trace("tableGitHubRunnerGet", "owner", owner, "repo", repo, "runnerId", runnerId)

	client := connect(ctx, d)

	type GetResponse struct {
		runner *github.Runner
		resp   *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetRunner(ctx, owner, repo, runnerId)
		return GetResponse{
			runner: detail,
			resp:   resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, retryConfig())
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)

	return getResp.runner, nil
}
