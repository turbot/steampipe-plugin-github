package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubOrganizationCredentialAuthorization() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_credential_authorization",
		Description: "Classic personal access tokens (and SSH keys) that members have authorized to access an organization's resources through SAML single sign-on (SSO).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404", "403"}),
			Hydrate:           tableGitHubOrganizationCredentialAuthorizationList,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The login name of the organization the credential is authorized against."},
			{Name: "login", Type: proto.ColumnType_STRING, Description: "The login of the user that owns the underlying credential."},
			{Name: "credential_id", Type: proto.ColumnType_INT, Description: "The unique identifier for the credential."},
			{Name: "credential_type", Type: proto.ColumnType_STRING, Description: "A human-readable description of the credential type, for example 'personal access token' or 'SSH key'."},
			{Name: "token_last_eight", Type: proto.ColumnType_STRING, Description: "The last eight characters of the credential. Only included when the credential is a personal access token."},
			{Name: "scopes", Type: proto.ColumnType_JSON, Description: "The list of OAuth scopes the credential has been granted."},
			{Name: "credential_authorized_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CredentialAuthorizedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the credential was authorized for use."},
			{Name: "credential_accessed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CredentialAccessedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the credential was last accessed. May be null if it was never used."},
			{Name: "authorized_credential_id", Type: proto.ColumnType_INT, Description: "The unique identifier for the underlying authorized credential."},
			{Name: "authorized_credential_note", Type: proto.ColumnType_STRING, Description: "The note attached to the token. Only included when the credential is a personal access token."},
			{Name: "authorized_credential_expires_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("AuthorizedCredentialExpiresAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the token expires. Only included when the credential is a personal access token."},
			{Name: "authorized_credential_title", Type: proto.ColumnType_STRING, Description: "The title given to the SSH key. Only included when the credential is an SSH key."},
			{Name: "fingerprint", Type: proto.ColumnType_STRING, Description: "The unique string used to distinguish the credential. Only included when the credential is an SSH key."},
		}),
	}
}

func tableGitHubOrganizationCredentialAuthorizationList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	org := d.EqualsQuals["organization"].GetStringValue()

	// Empty check
	if org == "" {
		return nil, fmt.Errorf("'organization' qual is required for the github_organization_credential_authorization table")
	}

	opts := &github.ListOptions{PerPage: 100}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	for {
		creds, resp, err := client.Organizations.ListCredentialAuthorizations(ctx, org, opts)
		if err != nil {
			plugin.Logger(ctx).Error("github_organization_credential_authorization", "api_error", err)
			return nil, err
		}

		for _, c := range creds {
			if c != nil {
				d.StreamListItem(ctx, c)
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
