package github

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubRepositoryDependabotAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_dependabot_alert",
		Description: "Dependabot alerts from a repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:    "state",
					Require: plugin.Optional,
				},
				{
					Name:    "security_advisory_severity",
					Require: plugin.Optional,
				},
				{
					Name:    "dependency_package_ecosystem",
					Require: plugin.Optional,
				},
				{
					Name:    "dependency_package_name",
					Require: plugin.Optional,
				},
				{
					Name:    "dependency_scope",
					Require: plugin.Optional,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryDependabotAlertList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "dependabot_number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryDependabotAlertGet,
		},
		Columns: append(
			gitHubDependabotAlertColumns(),
			[]*plugin.Column{
				{
					Name:        "repository_full_name",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("repository_full_name"),
					Description: "The full name of the repository (login/repo-name).",
				},
			}...,
		),
	}
}

//// LIST FUNCTION

func tableGitHubRepositoryDependabotAlertList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opt := &github.ListAlertsOptions{
		ListCursorOptions: github.ListCursorOptions{First: 100},
	}

	if quals["state"] != nil {
		state := quals["state"].GetStringValue()
		opt.State = &state
	}
	if quals["security_advisory_severity"] != nil {
		severity := quals["security_advisory_severity"].GetStringValue()
		opt.Severity = &severity
	}
	if quals["dependency_package_ecosystem"] != nil {
		ecosystem := quals["dependency_package_ecosystem"].GetStringValue()
		opt.Ecosystem = &ecosystem
	}
	if quals["dependency_package_name"] != nil {
		package_name := quals["dependency_package_name"].GetStringValue()
		opt.Package = &package_name
	}
	if quals["dependency_scope"] != nil {
		scope := quals["dependency_scope"].GetStringValue()
		opt.Scope = &scope
	}

	type ListPageResponse struct {
		alerts []*github.DependabotAlert
		resp   *github.Response
	}

	client := connect(ctx, d)

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListCursorOptions.First) {
			opt.ListCursorOptions.First = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		alerts, resp, err := client.Dependabot.ListRepoAlerts(ctx, owner, repo, opt)
		return ListPageResponse{
			alerts: alerts,
			resp:   resp,
		}, err
	}
	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		alerts := listResponse.alerts
		resp := listResponse.resp

		for _, i := range alerts {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.After == "" {
			break
		}

		opt.ListCursorOptions.After = resp.After
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func tableGitHubRepositoryDependabotAlertGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var alertNumber int

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	alertNumber = int(d.KeyColumnQuals["dependabot_number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubDependabotAlertGet", "owner", owner, "repo", repo, "alertNumber", alertNumber)

	client := connect(ctx, d)

	type GetResponse struct {
		alert *github.DependabotAlert
		resp  *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		alert, resp, err := client.Dependabot.GetRepoAlert(ctx, owner, repo, alertNumber)
		return GetResponse{
			alert: alert,
			resp:  resp,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	alert := getResp.alert

	return alert, nil
}
