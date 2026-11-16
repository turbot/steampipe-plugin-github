package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// codeScanningNotSetUpErrors are the responses for a repository whose code scanning has never run:
// 404 "no analysis found", or 403 when GitHub Advanced Security / Code Security is off for it.
// A 403 for a missing token permission ("Resource not accessible") is deliberately not matched.
var codeScanningNotSetUpErrors = []string{"404", "must be enabled for this repository to use code scanning"}

func tableGitHubRepositoryCodeScanningAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_code_scanning_alert",
		Description: "Code scanning alerts from a repository.",
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
			},
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError(codeScanningNotSetUpErrors)},
			Hydrate:      tableGitHubRepositoryCodeScanningAlertList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:   plugin.AllColumns([]string{"repository_full_name", "number"}),
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError(codeScanningNotSetUpErrors)},
			Hydrate:      tableGitHubRepositoryCodeScanningAlertGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "The security alert number."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the alert: open, dismissed or fixed."},
			{Name: "rule_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Rule.ID"), Description: "A unique identifier for the rule used to detect the alert."},
			{Name: "rule_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Rule.Name"), Description: "The name of the rule used to detect the alert."},
			{Name: "rule_severity", Type: proto.ColumnType_STRING, Transform: transform.FromField("Rule.Severity"), Description: "The severity of the alert: none, note, warning or error."},
			{Name: "rule_security_severity_level", Type: proto.ColumnType_STRING, Transform: transform.FromField("Rule.SecuritySeverityLevel"), Description: "The security severity of the alert: low, medium, high or critical. Null for rules that are not security rules."},
			{Name: "rule_description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Rule.Description"), Description: "A short description of the rule used to detect the alert."},
			{Name: "rule_tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Rule.Tags"), Description: "A set of tags applicable for the rule."},
			{Name: "tool_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.Name"), Description: "The name of the tool used to generate the code scanning analysis."},
			{Name: "tool_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.Version"), Description: "The version of the tool used to generate the code scanning analysis."},
			{Name: "tool_guid", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.GUID"), Description: "The GUID of the tool used to generate the code scanning analysis, if provided in the uploaded SARIF data."},
			{Name: "most_recent_instance_ref", Type: proto.ColumnType_STRING, Transform: transform.FromField("MostRecentInstance.Ref"), Description: "The Git reference of the most recent instance of the alert."},
			{Name: "most_recent_instance_state", Type: proto.ColumnType_STRING, Transform: transform.FromField("MostRecentInstance.State"), Description: "The state of the most recent instance of the alert."},
			{Name: "most_recent_instance_commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("MostRecentInstance.CommitSHA"), Description: "The commit SHA of the most recent instance of the alert."},
			{Name: "most_recent_instance_location_path", Type: proto.ColumnType_STRING, Transform: transform.FromField("MostRecentInstance.Location.Path"), Description: "The file path of the most recent instance of the alert."},
			{Name: "most_recent_instance", Type: proto.ColumnType_JSON, Description: "The most recent instance of the alert on the default branch or the specified Git reference."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was last updated."},
			{Name: "fixed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("FixedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was no longer detected and was considered fixed."},
			{Name: "dismissed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("DismissedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was dismissed."},
			{Name: "dismissed_by", Type: proto.ColumnType_STRING, Transform: transform.FromField("DismissedBy.Login"), Description: "The login of the user who dismissed the alert."},
			{Name: "dismissed_reason", Type: proto.ColumnType_STRING, Description: "The reason for dismissing or closing the alert: false positive, won't fix or used in tests."},
			{Name: "dismissed_comment", Type: proto.ColumnType_STRING, Description: "The dismissal comment associated with the dismissal of the alert."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The REST API URL of the alert resource."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub URL of the alert resource."},
		}),
	}
}

func tableGitHubRepositoryCodeScanningAlertList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opt := &github.AlertListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	if quals["state"] != nil {
		opt.State = quals["state"].GetStringValue()
	}

	client := connect(ctx, d)
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		alerts, resp, err := client.CodeScanning.ListAlertsForRepo(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		for _, i := range alerts {
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

		opt.ListOptions.Page = resp.NextPage
	}

	return nil, nil
}

func tableGitHubRepositoryCodeScanningAlertGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	number := quals["number"].GetInt64Value()
	fullName := quals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if number == 0 || fullName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubRepositoryCodeScanningAlertGet", "owner", owner, "repo", repo, "number", number)

	client := connect(ctx, d)
	alert, _, err := client.CodeScanning.GetAlert(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	return alert, nil
}
