package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubLicense() *plugin.Table {
	return &plugin.Table{
		Name:        "github_license",
		Description: "GitHub Licenses are common software licenses that you can associate with your repository.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubLicenseList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("key"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubLicenseGetData,
		},
		Columns: []*plugin.Column{
			{Name: "spdx_id", Description: "The Software Package Data Exchange (SPDX) id of the license.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateSpdxId},
			{Name: "name", Description: "The name of the license.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateName},
			{Name: "url", Description: "The HTML URL of the license.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateUrl},

			// The body is huge and of limited value, exclude it for now
			// {Name: "body", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "conditions", Description: "An array of license conditions (include-copyright,disclose-source, etc).", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: licenseHydrateConditions},
			{Name: "description", Description: "The license description.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateDescription},
			{Name: "featured", Description: "If true, the license is 'featured' in the GitHub UI.", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: licenseHydrateFeatured},
			{Name: "hidden", Description: "Whether the license should be displayed in license pickers.", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: licenseHydrateHidden},
			{Name: "implementation", Description: "Implementation instructions for the license.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateImplementation},
			{Name: "key", Description: "The unique key of the license.", Type: proto.ColumnType_STRING},
			{Name: "limitations", Description: "An array of limitations for the license (trademark-use, liability,warranty, etc).", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: licenseHydrateLimitations},
			{Name: "permissions", Description: "An array of permissions for the license (private-use, commercial-use,modifications, etc).", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Hydrate: licenseHydratePermissions},
			{Name: "nickname", Description: "The customary short name of the license.", Type: proto.ColumnType_STRING, Transform: transform.FromValue(), Hydrate: licenseHydrateNickname},
			{Name: "pseudo_license", Description: "Indicates if the license is a pseudo-license placeholder (e.g. other, no-license).", Type: proto.ColumnType_BOOL, Transform: transform.FromValue(), Hydrate: licenseHydratePseudoLicense},
		},
	}
}

func tableGitHubLicenseList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	var query struct {
		RateLimit models.RateLimit
		Licenses  []models.License `graphql:"licenses"`
	}
	variables := map[string]interface{}{}
	appendLicenseColumnIncludes(&variables, d.QueryContext.Columns)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_license", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_license", "api_error", err)
		return nil, err
	}

	for _, license := range query.Licenses {
		d.StreamListItem(ctx, license)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func tableGitHubLicenseGetData(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	key := d.EqualsQuals["key"].GetStringValue()
	if key == "" {
		return nil, nil
	}

	variables := map[string]interface{}{
		"key": githubv4.String(key),
	}

	client := connectV4(ctx, d)

	var query struct {
		RateLimit models.RateLimit
		License   models.License `graphql:"license(key: $key)"`
	}
	appendLicenseColumnIncludes(&variables, d.QueryContext.Columns)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_license", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_license", "api_error", err)
		return nil, err
	}

	return query.License, nil
}
