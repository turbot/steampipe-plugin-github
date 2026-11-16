package github

import (
	"context"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableGitHubRepositoryCodeScanningAnalysis() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_code_scanning_analysis",
		Description: "Code scanning analyses uploaded to a repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:    "ref",
					Require: plugin.Optional,
				},
			},
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError(codeScanningNotSetUpErrors)},
			Hydrate:      tableGitHubRepositoryCodeScanningAnalysisList,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "Unique identifier for this analysis."},
			{Name: "ref", Type: proto.ColumnType_STRING, Description: "The Git reference of the analysis, for example refs/heads/main."},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("CommitSHA"), Description: "The SHA of the commit to which the analysis you are uploading relates."},
			{Name: "analysis_key", Type: proto.ColumnType_STRING, Description: "Identifies the configuration under which the analysis was executed."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Identifies the configuration under which the analysis was executed, used to distinguish analyses for the same tool and commit."},
			{Name: "environment", Type: proto.ColumnType_STRING, Description: "Identifies the variable values associated with the environment in which this analysis was performed."},
			{Name: "tool_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.Name"), Description: "The name of the tool used to generate the analysis."},
			{Name: "tool_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.Version"), Description: "The version of the tool used to generate the analysis."},
			{Name: "tool_guid", Type: proto.ColumnType_STRING, Transform: transform.FromField("Tool.GUID"), Description: "The GUID of the tool used to generate the analysis, if provided in the uploaded SARIF data."},
			{Name: "results_count", Type: proto.ColumnType_INT, Description: "The total number of results in the analysis."},
			{Name: "rules_count", Type: proto.ColumnType_INT, Description: "The total number of rules used in the analysis."},
			{Name: "sarif_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("SarifID"), Description: "An identifier for the upload."},
			{Name: "deletable", Type: proto.ColumnType_BOOL, Description: "If true, the analysis can be deleted."},
			{Name: "error", Type: proto.ColumnType_STRING, Description: "The error message of the analysis, if any."},
			{Name: "warning", Type: proto.ColumnType_STRING, Description: "Warning generated when processing the analysis."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the analysis was created."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The REST API URL of the analysis resource."},
		}),
	}
}

func tableGitHubRepositoryCodeScanningAnalysisList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	opt := &github.AnalysesListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	if quals["ref"] != nil {
		ref := quals["ref"].GetStringValue()
		opt.Ref = &ref
	}

	client := connect(ctx, d)
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.PerPage) {
			opt.PerPage = int(*limit)
		}
	}

	for {
		analyses, resp, err := client.CodeScanning.ListAnalysesForRepo(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}

		for _, i := range analyses {
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

		opt.Page = resp.NextPage
	}

	return nil, nil
}
