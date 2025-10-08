package github

import (
	"context"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubActionsOrganizationVariable() *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_organization_variable",
		Description: "Variables are unencrypted environment variables that you create in an organization for use in GitHub Actions workflows.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("org"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubOrgVariableList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"org", "name"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubOrgVariableGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "org", Type: proto.ColumnType_STRING, Transform: transform.FromQual("org"), Description: "The organization name."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the variable."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The value of the variable."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the variable (all, private, or selected)."},
			{Name: "selected_repositories_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("SelectedRepositoriesURL"), Description: "The API URL for the selected repositories (if visibility is 'selected')."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the variable was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "Time when the variable was updated."},
		}),
	}
}

func tableGitHubOrgVariableList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	org := d.EqualsQuals["org"].GetStringValue()
	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		variables, resp, err := client.Actions.ListOrgVariables(ctx, org, opts)
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

func tableGitHubOrgVariableGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	org := d.EqualsQuals["org"].GetStringValue()

	// Empty check for the parameters
	if name == "" || org == "" {
		return nil, nil
	}

	plugin.Logger(ctx).Trace("tableGitHubOrgVariableGet", "org", org, "name", name)

	client := connect(ctx, d)

	type GetResponse struct {
		variable *github.ActionsVariable
		resp     *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetOrgVariable(ctx, org, name)
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
