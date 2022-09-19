package github

import (
	"context"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

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

			// Top columns
			{Name: "spdx_id", Description: "The Software Package Data Exchange (SPDX) id of the license.", Type: proto.ColumnType_STRING, Transform: transform.FromField("SPDXID")},
			{Name: "name", Description: "The name of the license.", Type: proto.ColumnType_STRING},
			{Name: "html_url", Description: "The HTML URL of the license.", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},

			// The body is huge and of limited value, exclude it for now
			// {Name: "body", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "conditions", Description: "An array of license conditions (include-copyright,disclose-source, etc).", Type: proto.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "description", Description: "The license description.", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "featured", Description: "If true, the license is 'featured' in the GitHub UI.", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubLicenseGetData},
			{Name: "implementation", Description: "Implementation instructions for the license.", Type: proto.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "key", Description: "The unique key of the license.", Type: proto.ColumnType_STRING},
			{Name: "limitations", Description: "An array of limitations for the license (trademark-use, liability,warranty, etc).", Type: proto.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "permissions", Description: "An array of permissions for the license (private-use, commercial-use,modifications, etc).", Type: proto.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "url", Description: "The API url of the license.", Type: proto.ColumnType_STRING},
		},
	}
}

//// LIST FUNCTION

func tableGitHubLicenseList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	type ListPageResponse struct {
		licenses []*github.License
		resp     *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		licenses, resp, err := client.Licenses.List(ctx)
		return ListPageResponse{
			licenses: licenses,
			resp:     resp,
		}, err
	}

	listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	listResponse := listPageResponse.(ListPageResponse)
	licenses := listResponse.licenses

	if err != nil {
		return nil, err
	}

	for _, i := range licenses {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func tableGitHubLicenseGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var key string
	if h.Item != nil {
		item := h.Item.(*github.License)
		key = *item.Key
	} else {
		key = d.KeyColumnQuals["key"].GetStringValue()
	}

	// Return nil, if no input provided
	if key == "" {
		return nil, nil
	}

	client := connect(ctx, d)

	type GetResponse struct {
		license *github.License
		resp    *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		license, resp, err := client.Licenses.Get(ctx, key)
		return GetResponse{
			license: license,
			resp:    resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}
	getResp := getResponse.(GetResponse)
	license := getResp.license

	return license, nil
}
