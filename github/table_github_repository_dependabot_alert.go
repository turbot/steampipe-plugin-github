package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

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
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
			Hydrate:           tableGitHubRepositoryDependabotAlertList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "alert_number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
			Hydrate:           tableGitHubRepositoryDependabotAlertGet,
		},
		Columns: commonColumns(append(
			gitHubDependabotAlertColumns(),
			[]*plugin.Column{
				{
					Name:        "repository_full_name",
					Type:        proto.ColumnType_STRING,
					Transform:   transform.FromQual("repository_full_name"),
					Description: "The full name of the repository (login/repo-name).",
				},
			}...,
		)),
	}
}

func tableGitHubRepositoryDependabotAlertList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
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
		packageName := quals["dependency_package_name"].GetStringValue()
		opt.Package = &packageName
	}
	if quals["dependency_scope"] != nil {
		scope := quals["dependency_scope"].GetStringValue()
		opt.Scope = &scope
	}

	client := connect(ctx, d)
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListCursorOptions.First) {
			opt.ListCursorOptions.First = int(*limit)
		}
	}

	for {
		alerts, resp, err := client.Dependabot.ListRepoAlerts(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		for _, i := range alerts {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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

func tableGitHubRepositoryDependabotAlertGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var alertNumber int

	logger := plugin.Logger(ctx)
	quals := d.EqualsQuals

	alertNumber = int(d.EqualsQuals["alert_number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubDependabotAlertGet", "owner", owner, "repo", repo, "alertNumber", alertNumber)

	client := connect(ctx, d)
	alert, _, err := client.Dependabot.GetRepoAlert(ctx, owner, repo, alertNumber)
	if err != nil {
		return nil, err
	}

	return alert, nil
}
