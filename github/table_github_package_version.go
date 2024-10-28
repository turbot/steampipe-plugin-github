package github

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubPackageVersion() *plugin.Table {
	return &plugin.Table{
		Name:        "github_package_version",
		Description: "GitHub Packages allow you to store and manage packages such as container images or other artifacts in your GitHub repositories.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "organization", Require: plugin.Required},
				{Name: "package_type", Require: plugin.Optional},
				{Name: "package_name", Require: plugin.Optional},
				{Name: "visibility", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			ParentHydrate:     tableGitHubPackageList,
			Hydrate:           tableGitHubPackageVersionList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "package_name", "package_type", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubPackageVersionGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "Name of the package version."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the package version.", Transform: transform.FromField("PackageVersion.ID")},
			{Name: "digest", Type: proto.ColumnType_STRING, Description: "The digest (shasum) of the package version.", Transform: transform.FromField("PackageVersion.Name")},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL of the package version.", Transform: transform.FromField("PackageVersion.HTMLURL")},

			{Name: "organization", Type: proto.ColumnType_STRING, Description: "The name of the GitHub organization.", Transform: transform.FromQual("organization")},
			{Name: "package_type", Type: proto.ColumnType_STRING, Description: "Type of the package (e.g., container, npm, etc.)."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Visibility of the package (public or private)."},
			{Name: "prerelease", Type: proto.ColumnType_BOOL, Description: "Indicates if the package version is a pre-release.", Transform: transform.FromField("PackageVersion.Prerelease")},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The URL of the package version.", Transform: transform.FromField("PackageVersion.URL")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the package version was created.", Transform: transform.FromField("PackageVersion.CreatedAt").Transform(convertTimestamp)},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the package version was last updated.", Transform: transform.FromField("PackageVersion.UpdatedAt").Transform(convertTimestamp)},

			// JSON fields
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "Tags associated with the package version.", Transform: transform.FromField("PackageVersion.Metadata.Container.Tags")},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Metadata of the package version.", Transform: transform.FromField("PackageVersion.Metadata")},
			{Name: "author", Type: proto.ColumnType_JSON, Description: "Author of the package version.", Transform: transform.FromField("PackageVersion.Author")},
			{Name: "release", Type: proto.ColumnType_JSON, Description: "Release information of the package version.", Transform: transform.FromField("PackageVersion.Release")},
		}),
	}
}

type PackageVersionInfo struct {
	PackageName    string
	PackageType    string
	Visibility     string
	PackageVersion *github.PackageVersion
}

func tableGitHubPackageVersionList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	pkg := h.Item.(*github.Package)
	org, packageType, packageName := d.EqualsQuals["organization"].GetStringValue(), d.EqualsQuals["package_type"].GetStringValue(), d.EqualsQuals["package_name"].GetStringValue()
	visibility := d.EqualsQuals["visibility"].GetStringValue()

	if packageName != "" && packageName != *pkg.Name {
		return nil, nil
	}
	if packageType != "" && packageType != *pkg.PackageType {
		return nil, nil
	}
	if visibility != "" && visibility != *pkg.Visibility {
		return nil, nil
	}

	// Return, if org is not specified
	if org == "" {
		return nil, nil
	}

	// Encode the package name, otherwise we will get 404 err if the package name looks like "steampipe/plugins/turbot/jira"
	packageName = url.QueryEscape(*pkg.Name)

	opts := &github.PackageListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	if visibility != "" {
		opts.Visibility = &visibility
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		packageVersions, resp, err := client.Organizations.PackageGetAllVersions(ctx, org, *pkg.PackageType, packageName, opts)
		if err != nil {
			// In the case of parent hydrate the ignore config seems to not work for the child table. So we need to handle it manually.
			// Steampipe SDK issue ref: https://github.com/turbot/steampipe-plugin-sdk/issues/544
			if strings.Contains(err.Error(), "404") {
				return nil, nil
			}

			plugin.Logger(ctx).Error("github_package_version.tableGitHubPackageVersionList", "api_error", err)
			return nil, err
		}

		// Find latest tag version
		latestTagVersion := findLatestPackage(ctx, packageVersions)

		for _, pkgVersion := range packageVersions {

			if helpers.StringSliceContains(pkgVersion.Metadata.Container.Tags, latestTagVersion) {
				pkgVersion.Metadata.Container.Tags = append(pkgVersion.Metadata.Container.Tags, "latest")
			}
			d.StreamListItem(ctx, PackageVersionInfo{*pkg.Name, *pkg.PackageType, *pkg.Visibility, pkgVersion})

			// Stop if we've hit the limit set in the query context
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

//// HYDRATE FUNCTION

func tableGitHubPackageVersionGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	org := d.EqualsQuals["organization"].GetStringValue()
	name := d.EqualsQuals["package_name"].GetStringValue()
	packageType := d.EqualsQuals["package_type"].GetStringValue()
	packageVersionID := d.EqualsQuals["id"].GetInt64Value()

	// Encode the package name, otherwise we will get 404 err if the package name looks like
	name = url.QueryEscape(name)

	// Fetch the package
	pkgVersion, _, err := client.Organizations.PackageGetVersion(ctx, org, packageType, name, packageVersionID)
	if err != nil {
		plugin.Logger(ctx).Error("github_package_version.tableGitHubPackageVersionGet", "api_error", err)
		return nil, err
	}

	return PackageVersionInfo{name, packageType, "", pkgVersion}, nil
}

//// HELPER FUNCTION

// Find the latest package based on CreatedAt field
func findLatestPackage(_ context.Context, packages []*github.PackageVersion) string {
	var latest *github.PackageVersion
	if len(packages) > 1 {
		latest = packages[0]
	}

	for _, pkg := range packages {
		if pkg.CreatedAt != nil {
			t := pkg.CreatedAt.UTC()
			if pkg.CreatedAt.After(t) {
				latest = pkg
			}
		}
	}

	container := latest.GetMetadata().GetContainer()
	if container != nil && len(container.Tags) > 1 {
		return container.Tags[1] // return the version tag
	}
	return ""
}
