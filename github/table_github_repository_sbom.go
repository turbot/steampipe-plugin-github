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
		Description: "Get the software bill of materials (SBOM) for a repository.",
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
				Description: "It represents when the SBOM was created and who created it.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the SPDX document.",
			},
			{
				Name:        "data_license",
				Type:        proto.ColumnType_STRING,
				Description: "The license under which the SPDX document is licensed.",
			},
			{
				Name:        "document_describes",
				Type:        proto.ColumnType_JSON,
				Description: "The name of the repository that the SPDX document describes.",
			},
			{
				Name:        "document_namespace",
				Type:        proto.ColumnType_STRING,
				Description: "The namespace for the SPDX document.",
			},
			{
				Name:        "packages",
				Type:        proto.ColumnType_JSON,
				Description: "Array of packages in SPDX format.",
			},
		},
	}
}

//// LIST FUNCTION
func listRepositorySboms(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string

	logger := plugin.Logger(ctx)

	fullName := d.EqualsQualString("repository_full_name")
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubRepositorySbomList", "owner", owner, "repo", repo)

	client := connect(ctx, d)
	sbom, _, err := client.DependencyGraph.GetSBOM(ctx, owner, repo)
	if err != nil {
		logger.Error("github_repository_sbom.listRepositorySboms", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, sbom.SBOM)

	return sbom, nil
}