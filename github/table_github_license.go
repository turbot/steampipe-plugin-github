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
		Description: "Github License",
		List: &plugin.ListConfig{
			Hydrate: tableGitHubLicenseList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("key"),
			Hydrate:    tableGitHubLicenseGetData,
		},
		Columns: []*plugin.Column{

			// Top columns
			{Name: "spdx_id", Type: pb.ColumnType_STRING, Transform: transform.FromField("SPDXID")},
			{Name: "name", Type: pb.ColumnType_STRING},
			{Name: "html_url", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},

			// The body is huge and of limited value, exclude it for now
			// {Name: "body", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "conditions", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "description", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "featured", Type: pb.ColumnType_BOOL, Hydrate: tableGitHubLicenseGetData},
			{Name: "implementation", Type: pb.ColumnType_STRING, Hydrate: tableGitHubLicenseGetData},
			{Name: "key", Type: pb.ColumnType_STRING},
			{Name: "limitations", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "permissions", Type: pb.ColumnType_JSON, Hydrate: tableGitHubLicenseGetData},
			{Name: "url", Type: pb.ColumnType_STRING},
		},
	}
}

//// list ////

func tableGitHubLicenseList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d.ConnectionManager)

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

	client := connect(ctx, d.ConnectionManager)

	var detail *github.License
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		detail, resp, err = client.Licenses.Get(ctx, key)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}
