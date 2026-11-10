package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v93/github"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
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
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"})},
			Hydrate:      listRepositorySboms,
		},
		Columns: commonColumns([]*plugin.Column{
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
		}),
	}
}

// // LIST FUNCTION
func listRepositorySboms(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string

	logger := plugin.Logger(ctx)

	fullName := d.EqualsQualString("repository_full_name")
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubRepositorySbomList", "owner", owner, "repo", repo)

	client := connect(ctx, d)
	sbom, err := getRepositorySBOM(ctx, client, owner, repo, 2*time.Second)
	if err != nil {
		logger.Error("github_repository_sbom.listRepositorySboms", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, sbom.SBOM)

	return sbom, nil
}

// getRepositorySBOM waits for one generated report within the query deadline or five minutes.
func getRepositorySBOM(ctx context.Context, client *github.Client, owner, repo string, pollInterval time.Duration) (*github.SBOM, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	generation, _, err := client.DependencyGraph.GenerateSBOM(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	reportURL, err := url.Parse(generation.GetSBOMURL())
	if err != nil {
		return nil, fmt.Errorf("invalid SBOM report URL: %w", err)
	}
	_, uuid, found := strings.Cut(reportURL.Path, "/fetch-report/")
	if !found || uuid == "" || strings.Contains(uuid, "/") {
		return nil, fmt.Errorf("SBOM generation returned no valid report identifier")
	}

	// Signed download URLs require a client without GitHub authentication.
	downloadClient := &http.Client{Timeout: time.Minute}
	for {
		sbom, _, _, err := client.DependencyGraph.FetchSBOM(ctx, owner, repo, uuid, downloadClient)
		if err == nil {
			if sbom == nil || sbom.SBOM == nil {
				return nil, fmt.Errorf("SBOM download returned an empty document")
			}
			return sbom, nil
		}
		var accepted *github.AcceptedError
		if !errors.As(err, &accepted) {
			return nil, err
		}

		timer := time.NewTimer(pollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}
