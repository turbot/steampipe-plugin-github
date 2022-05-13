package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubActionsOrganozationSecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_organization_secret",
		Description: "Secrets are encrypted environment variables that you create in an organization",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("organization_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubOrgSecretList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"organization_name", "name"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubOrgSecretGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "organization_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization_name"), Description: "Full name of the orgainazation that contains the secrets."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the secret."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The vicibility of the secret"},
			{Name: "selected_repositories_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("SelectedRepositoriesURL"), Description: "Size of the artifact in bytes."},

			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the secret was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ExpiresAt").Transform(convertTimestamp), Description: "Time when the secret was updated."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubOrgSecretList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	orgName := d.KeyColumnQuals["organization_name"].GetStringValue()

	type ListPageResponse struct {
		orgSecrets *github.Secrets
		resp       *github.Response
	}

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		orgSecrets, resp, err := client.Actions.ListOrgSecrets(ctx, orgName, opts)
		return ListPageResponse{
			orgSecrets: orgSecrets,
			resp:       resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		orgSecrets := listResponse.orgSecrets
		resp := listResponse.resp

		for _, i := range orgSecrets.Secrets {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
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

//// HYDRATE FUNCTIONS

func tableGitHubOrgSecretGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()
	orgName := d.KeyColumnQuals["organization_name"].GetStringValue()
	plugin.Logger(ctx).Trace("tableGitHubOrgSecretGet", "owner", orgName, "name", name)

	client := connect(ctx, d)

	type GetResponse struct {
		orgSecret *github.Secret
		resp      *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetOrgSecret(ctx, orgName, name)
		return GetResponse{
			orgSecret: detail,
			resp:      resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	artifact := getResp.orgSecret

	return artifact, nil
}
