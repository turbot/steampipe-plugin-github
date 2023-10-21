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
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
				Description: "The name of the package.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VersionInfo"),
				Description: "Version info of the package.",
			},
			{
				Name:        "license",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LicenseConcluded"),
				Description: "License of the package.",
			},
			// TODO: there are more fields to be added here!
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

	for _, i := range sbom.SBOM.Packages {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return sbom, nil
}
