package github

import (
	"context"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubActionsRepositorySecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_actions_repository_secret",
		Description: "Secrets are encrypted environment variables that you create in a repository",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoSecretList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "name"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepoSecretGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the secrets."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the secret."},
			{Name: "visibility", Type: proto.ColumnType_STRING, Description: "The visibility of the secret."},
			{Name: "selected_repositories_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("SelectedRepositoriesURL"), Description: "The GitHub URL of the repository."},

			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "Time when the secret was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").Transform(convertTimestamp), Description: "Time when the secret was updated."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubRepoSecretList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(orgName)

	type ListPageResponse struct {
		secrets *github.Secrets
		resp    *github.Response
	}

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		secrets, resp, err := client.Actions.ListRepoSecrets(ctx, owner, repo, opts)
		return ListPageResponse{
			secrets: secrets,
			resp:    resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		secrets := listResponse.secrets
		resp := listResponse.resp

		for _, i := range secrets.Secrets {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
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

//// HYDRATE FUNCTIONS

func tableGitHubRepoSecretGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	orgName := d.EqualsQuals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if name == "" || orgName == "" {
		return nil, nil
	}
	owner, repo := parseRepoFullName(orgName)
	plugin.Logger(ctx).Trace("tableGitHubRepoSecretGet", "owner", owner, "repo", repo, "name", name)

	client := connect(ctx, d)

	type GetResponse struct {
		secret *github.Secret
		resp   *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Actions.GetRepoSecret(ctx, owner, repo, name)
		return GetResponse{
			secret: detail,
			resp:   resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, retryConfig())
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)

	return getResp.secret, nil
}
