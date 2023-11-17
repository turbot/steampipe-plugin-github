package github

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubRepositorySbom() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_sbom",
		Description: "SBOM from a repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
			Hydrate:           tableGitHubRepositorySbomList,
		},
		Columns: []*plugin.Column{
			{
				Name:        "repository_full_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("repository_full_name"),
				Description: "The full name of the repository (login/repo-name).",
			},
			{
				Name:        "spdx_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SPDXID"),
				Description: "The SPDX identifier for the SPDX document.",
			},
			{
				Name:        "spdx_version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SPDXVersion"),
				Description: "The version of the SPDX specification that this document conforms to.",
			},
			{
				Name:        "creation_info",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CreationInfo"),
				Description: "The version of the SPDX specification that this document conforms to.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Description: "The name of the SPDX document.",
			},
			{
				Name:        "data_license",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DataLicense"),
				Description: "The license under which the SPDX document is licensed.",
			},
			{
				Name:        "document_describes",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DocumentDescribes"),
				Description: "The name of the repository that the SPDX document describes.",
			},
			{
				Name:        "document_namespace",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("DocumentNamespace"),
				Description: "The namespace for the SPDX document.",
			},
			{
				Name:        "packages",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Packages"),
				Description: "Array of packages in spdx format.",
			},
		},
	}
}

func tableGitHubRepositorySbomList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string

	logger := plugin.Logger(ctx)
	quals := d.EqualsQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubDependabotSbomGet", "owner", owner, "repo", repo)

	client := connect(ctx, d)
	sbom, _, err := client.DependencyGraph.GetSBOM(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, sbom.SBOM)

	return sbom, nil
}
