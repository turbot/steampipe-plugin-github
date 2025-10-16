package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsRepositoryVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_repository_variable",
		Description: "Variables are unencrypted environment variables created in a repository for use in GitHub Actions workflows.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoVariableList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "name"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoVariableGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the variable."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the variable."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of the variable."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the variable was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "Time when the variable was updated."},
		}),
	}
}

func tableGitHubRepoVariableList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	repoFullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(repoFullName)
	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		variables, resp, err := client.Actions.ListRepoVariables(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		for _, i := range variables.Variables {
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

func tableGitHubRepoVariableGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	repoFullName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if name == "" || repoFullName == "" {
		return nil, nil
	}
	owner, repo := parseRepoFullName(repoFullName)
	plugin.Logger(ctx).Trace("tableGitHubRepoVariableGet", "owner", owner, "repo", repo, "name", name)

	client := connect(ctx, d)

	type GetResponse struct {
		variable *github.ActionsVariable
		resp     *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetRepoVariable(ctx, owner, repo, name)
		return GetResponse{
			variable: detail,
			resp:     resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, retryConfig())
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)

	return getResp.variable, nil
}
