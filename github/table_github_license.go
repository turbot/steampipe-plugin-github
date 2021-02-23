package github

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"

	pb "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubLicense() *plugin.Table {
	return &plugin.Table{
		Name:        "github_license",
		Description: "Github Licenses are common software licenses that you can associate with your repository.",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubLicenseList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("key"),
			Hydrate:    tableGitHubLicenseGetData,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "spdx_id", Description: "The Software Package Data Exchange (SPDX) id of the license.", Type: pb.ColumnType_STRING, Transform: transform.FromField("SPDXID")},
			{Name: "name", Description: "The name of the license.", Type: pb.ColumnType_STRING},
			{Name: "html_url", Description: "The HTML URL of the license.", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},

			// The body is huge and of limited value, exclude it for now
			// {Name: "body", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "conditions", Description: "An array of license conditions (include-copyright,disclose-source, etc).", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "description", Description: "The license description.", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "featured", Description: "If true, the license is 'featured' in the Github UI.", Type: pb.ColumnType_BOOL, Hydrate: tableGitHubLicenseGetData},
			{Name: "implementation", Description: "Implementation instructions for the license.", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "key", Description: "The unique key of the license.", Type: pb.ColumnType_STRING},
			{Name: "limitations", Description: "An array of limitations for the license (trademark-use, liability,warranty, etc).", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "permissions", Description: "An array of permissions for the license (private-use, commercial-use,modifications, etc).", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "url", Description: "The API url of the license.", Type: pb.ColumnType_STRING},
		},
	}
}

//// list ////

func tableGitHubLicenseList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	var items []*github.License
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		items, resp, err = client.Licenses.List(ctx)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, i := range items {
		d.StreamListItem(ctx, i)
	}

	return nil, nil
}

//// hydrate functions ////

func tableGitHubLicenseGetData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var key string

	if h.Item != nil {
		item := h.Item.(*github.License)
		log.Println("[INFO] item:", item.String())
		key = *item.Key
	} else {
		key = d.KeyColumnQuals["key"].GetStringValue()
	}

	client := connect(ctx, d)

	detail, _, err := client.Licenses.Get(ctx, key)

	if err != nil {
		return nil, err
	}
	return detail, nil
}
