package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubDependabotAlertColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "alert_number",
			Type:        proto.ColumnType_INT,
			Description: "The security alert number.",
			Transform:   transform.FromField("Number"),
		},
		{
			Name:        "state",
			Type:        proto.ColumnType_STRING,
			Description: "The state of the Dependabot alert.",
		},
		{
			Name:        "dependency_package_ecosystem",
			Type:        proto.ColumnType_STRING,
			Description: "The package's language or package management ecosystem.",
			Transform:   transform.FromField("Dependency.Package.Ecosystem"),
		},
		{
			Name:        "dependency_package_name",
			Type:        proto.ColumnType_STRING,
			Description: "The unique package name within its ecosystem.",
			Transform:   transform.FromField("Dependency.Package.Name"),
		},
		{
			Name:        "dependency_manifest_path",
			Type:        proto.ColumnType_STRING,
			Description: "The unique manifestation path within the ecosystem.",
			Transform:   transform.FromField("Dependency.ManifestPath"),
		},
		{
			Name:        "dependency_scope",
			Type:        proto.ColumnType_STRING,
			Description: "The execution scope of the vulnerable dependency.",
			Transform:   transform.FromField("Dependency.Scope"),
		},
		{
			Name:        "security_advisory_ghsa_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique GitHub Security Advisory ID assigned to the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.GHSAID"),
		},
		{
			Name:        "security_advisory_cve_id",
			Type:        proto.ColumnType_STRING,
			Description: "The unique CVE ID assigned to the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.CVEID"),
		},
		{
			Name:        "security_advisory_summary",
			Type:        proto.ColumnType_STRING,
			Description: "A short, plain text summary of the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.Summary"),
		},
		{
			Name:        "security_advisory_description",
			Type:        proto.ColumnType_STRING,
			Description: "A long-form Markdown-supported description of the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.Description"),
		},
		{
			Name:        "security_advisory_severity",
			Type:        proto.ColumnType_STRING,
			Description: "The severity of the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.Severity"),
		},
		{
			Name:        "security_advisory_cvss_score",
			Type:        proto.ColumnType_DOUBLE,
			Description: "The overall CVSS score of the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.CVSS.Score"),
		},
		{
			Name:        "security_advisory_cvss_vector_string",
			Type:        proto.ColumnType_STRING,
			Description: "The full CVSS vector string for the advisory.",
			Transform:   transform.FromField("SecurityAdvisory.CVSS.VectorString"),
		},
		// {
		// 	Name:        "security_advisory_cwes_cweid",
		// 	Type:        proto.ColumnType_STRING,
		// 	Description: "The unique CWE ID.",
		// 	Transform:   transform.FromField("SecurityAdvisory.CWEs[0].CWEID"),
		// },
		// {
		// 	Name:        "security_advisory_cwes_name",
		// 	Type:        proto.ColumnType_STRING,
		// 	Description: "The short, plain text name of the CWE.",
		// 	Transform:   transform.FromField("SecurityAdvisory.CWEs[0].Name"),
		// },
		{
			Name:        "security_advisory_cwes",
			Type:        proto.ColumnType_JSON,
			Description: "The associated CWEs",
			Transform:   transform.FromField("SecurityAdvisory.CWEs"),
		},
		{
			Name:        "security_advisory_published_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was published.",
			Transform:   transform.FromField("SecurityAdvisory.PublishedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "security_advisory_updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was last modified.",
			Transform:   transform.FromField("SecurityAdvisory.UpdatedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "security_advisory_withdrawn_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the advisory was withdrawn.",
			Transform:   transform.FromField("SecurityAdvisory.WithdrawnAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "url",
			Type:        proto.ColumnType_STRING,
			Description: "The REST API URL of the alert resource.",
		},
		{
			Name:        "html_url",
			Type:        proto.ColumnType_STRING,
			Description: "The GitHub URL of the alert resource.",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was created.",
			Transform:   transform.FromField("CreatedAt").Transform(convertTimestamp),
		},
		{
			Name:        "updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was last updated.",
			Transform:   transform.FromField("UpdatedAt").Transform(convertTimestamp),
		},
		{
			Name:        "dismissed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was dismissed.",
			Transform:   transform.FromField("DismissedAt").NullIfZero().Transform(convertTimestamp),
		},
		{
			Name:        "dismissed_reason",
			Type:        proto.ColumnType_STRING,
			Description: "The reason that the alert was dismissed.",
		},
		{
			Name:        "dismissed_comment",
			Type:        proto.ColumnType_STRING,
			Description: "An optional comment associated with the alert's dismissal.",
		},
		{
			Name:        "fixed_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "The time that the alert was no longer detected and was considered fixed.",
			Transform:   transform.FromField("FixedAt").NullIfZero().Transform(convertTimestamp),
		},
	}
}

func tableGitHubOrganizationDependabotAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_dependabot_alert",
		Description: "Dependabot alerts from an organization.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "organization",
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
			Hydrate:           tableGitHubOrganizationDependabotAlertList,
		},
		Columns: append(
			gitHubDependabotAlertColumns(),
			[]*plugin.Column{
				{
					Name:        "organization",
					Type:        proto.ColumnType_STRING,
					Description: "The login name of the organization.",
					Transform:   transform.FromQual("organization"),
				},
			}...,
		),
	}
}

func tableGitHubOrganizationDependabotAlertList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals

	org := quals["organization"].GetStringValue()

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
		alerts, resp, err := client.Dependabot.ListOrgAlerts(ctx, org, opt)
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
