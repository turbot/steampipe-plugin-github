package github

import (
	"context"
	"net/url"

	"github.com/google/go-github/v55/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubPackage() *plugin.Table {
	return &plugin.Table{
		Name:        "github_package",
		Description: "GitHub Packages allow you to store and manage packages such as container images or other artifacts in your GitHub repositories.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "organization", Require: plugin.Required},
				{Name: "package_type", Require: plugin.Optional},
				{Name: "visibility", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubPackageList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization", "name", "package_type"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubPackageGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique ID of the package."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the package."},
			{Name: "package_type", Type: proto.ColumnType_STRING, Description: "Type of the package. It can be one of 'npm', 'maven', 'rubygems', 'nuget', 'docker', or 'container'."},
			{Name: "organization", Type: proto.ColumnType_STRING, Description: "The name of the GitHub organization.", Transform: transform.FromQual("organization")},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "Visibility of the package (public or private)."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "API URL of the package."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL of the package."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the package was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the package was last updated."},

			// Repository details
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Repository.FullName"), Description: "Full name of the repository associated with the package."},

			// JSON field
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The information about the repository."},
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "The information about the owner."},
		}),
	}
}

func tableGitHubPackageList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	org, packageType := d.EqualsQuals["organization"].GetStringValue(), d.EqualsQuals["package_type"].GetStringValue()
	visibility := d.EqualsQuals["visibility"].GetStringValue()

	if packageType == "" {
		// Default package type
		packageType = "container"
	}

	// Return, if org is not specified
	if org == "" {
		return nil, nil
	}

	opts := &github.PackageListOptions{
		PackageType: &packageType,
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
		packages, resp, err := client.Organizations.ListPackages(ctx, org, opts)
		if err != nil {
			plugin.Logger(ctx).Error("github_package.tableGitHubPackageList", "api_error", err)
			return nil, err
		}

		for _, pkg := range packages {
			d.StreamListItem(ctx, pkg)

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

func tableGitHubPackageGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	org := d.EqualsQuals["organization"].GetStringValue()
	name := d.EqualsQuals["name"].GetStringValue()
	packageType := d.EqualsQuals["package_type"].GetStringValue()

	name = url.QueryEscape(name)

	// Fetch the package
	pkg, _, err := client.Organizations.GetPackage(ctx, org, packageType, name)
	if err != nil {
		plugin.Logger(ctx).Error("github_package.tableGitHubPackageGet", "api_error", err)
		return nil, err
	}

	return pkg, nil
}
