package github

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

// secretScanningAlert mirrors the REST alert without its `secret` field, so the leaked value is never decoded.
// go-github v55's SecretScanningAlert lacks updated_at, validity and the push protection fields.
type secretScanningAlert struct {
	Number                   *int              `json:"number,omitempty"`
	CreatedAt                *github.Timestamp `json:"created_at,omitempty"`
	UpdatedAt                *github.Timestamp `json:"updated_at,omitempty"`
	URL                      *string           `json:"url,omitempty"`
	HTMLURL                  *string           `json:"html_url,omitempty"`
	LocationsURL             *string           `json:"locations_url,omitempty"`
	State                    *string           `json:"state,omitempty"`
	SecretType               *string           `json:"secret_type,omitempty"`
	SecretTypeDisplayName    *string           `json:"secret_type_display_name,omitempty"`
	Validity                 *string           `json:"validity,omitempty"`
	Resolution               *string           `json:"resolution,omitempty"`
	ResolutionComment        *string           `json:"resolution_comment,omitempty"`
	ResolvedAt               *github.Timestamp `json:"resolved_at,omitempty"`
	ResolvedBy               *github.User      `json:"resolved_by,omitempty"`
	PushProtectionBypassed   *bool             `json:"push_protection_bypassed,omitempty"`
	PushProtectionBypassedAt *github.Timestamp `json:"push_protection_bypassed_at,omitempty"`
	PushProtectionBypassedBy *github.User      `json:"push_protection_bypassed_by,omitempty"`
	PubliclyLeaked           *bool             `json:"publicly_leaked,omitempty"`
	MultiRepo                *bool             `json:"multi_repo,omitempty"`
}

func tableGitHubRepositorySecretScanningAlert() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_secret_scanning_alert",
		Description: "Secret scanning alerts from a repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:    "state",
					Require: plugin.Optional,
				},
			},
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"})},
			Hydrate:      tableGitHubRepositorySecretScanningAlertList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:   plugin.AllColumns([]string{"repository_full_name", "number"}),
			IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: isNotFoundError([]string{"404"})},
			Hydrate:      tableGitHubRepositorySecretScanningAlertGet,
		},
		Columns: commonColumns([]*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "The security alert number."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the alert, either open or resolved."},
			{Name: "secret_type", Type: proto.ColumnType_STRING, Description: "The type of secret that secret scanning detected."},
			{Name: "secret_type_display_name", Type: proto.ColumnType_STRING, Description: "The display name of the type of secret that secret scanning detected."},
			{Name: "validity", Type: proto.ColumnType_STRING, Description: "The token status as of the latest validity check: active, inactive or unknown."},
			{Name: "resolution", Type: proto.ColumnType_STRING, Description: "The reason for resolving the alert: false_positive, wont_fix, revoked, used_in_tests, pattern_edited or pattern_deleted."},
			{Name: "resolution_comment", Type: proto.ColumnType_STRING, Description: "The comment that was optionally added when the alert was resolved."},
			{Name: "resolved_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ResolvedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was resolved."},
			{Name: "resolved_by", Type: proto.ColumnType_STRING, Transform: transform.FromField("ResolvedBy.Login"), Description: "The login of the user who resolved the alert."},
			{Name: "push_protection_bypassed", Type: proto.ColumnType_BOOL, Description: "If true, push protection was bypassed for the detected secret."},
			{Name: "push_protection_bypassed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PushProtectionBypassedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that push protection was bypassed."},
			{Name: "push_protection_bypassed_by", Type: proto.ColumnType_STRING, Transform: transform.FromField("PushProtectionBypassedBy.Login"), Description: "The login of the user who bypassed push protection."},
			{Name: "publicly_leaked", Type: proto.ColumnType_BOOL, Description: "If true, the detected secret was publicly leaked."},
			{Name: "multi_repo", Type: proto.ColumnType_BOOL, Description: "If true, the detected secret was found in multiple repositories in the same organization or enterprise."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time that the alert was last updated."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The REST API URL of the alert resource."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The GitHub URL of the alert resource."},
			{Name: "locations_url", Type: proto.ColumnType_STRING, Description: "The REST API URL of the code locations for this alert."},
		}),
	}
}

func tableGitHubRepositorySecretScanningAlertList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	perPage := 100
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(perPage) {
			perPage = int(*limit)
		}
	}

	params := url.Values{}
	params.Set("hide_secret", "true")
	params.Set("per_page", strconv.Itoa(perPage))
	if quals["state"] != nil {
		params.Set("state", quals["state"].GetStringValue())
	}

	client := connect(ctx, d)

	for {
		u := fmt.Sprintf("repos/%v/%v/secret-scanning/alerts?%s", owner, repo, params.Encode())
		req, err := client.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}

		var alerts []*secretScanningAlert
		resp, err := client.Do(ctx, req, &alerts)
		if err != nil {
			return nil, err
		}

		for _, i := range alerts {
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

		params.Set("page", strconv.Itoa(resp.NextPage))
	}

	return nil, nil
}

func tableGitHubRepositorySecretScanningAlertGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	number := quals["number"].GetInt64Value()
	fullName := quals["repository_full_name"].GetStringValue()

	// Empty check for the parameters
	if number == 0 || fullName == "" {
		return nil, nil
	}

	owner, repo := parseRepoFullName(fullName)
	plugin.Logger(ctx).Trace("tableGitHubRepositorySecretScanningAlertGet", "owner", owner, "repo", repo, "number", number)

	client := connect(ctx, d)
	u := fmt.Sprintf("repos/%v/%v/secret-scanning/alerts/%d?hide_secret=true", owner, repo, number)
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var alert secretScanningAlert
	_, err = client.Do(ctx, req, &alert)
	if err != nil {
		return nil, err
	}

	return &alert, nil
}
