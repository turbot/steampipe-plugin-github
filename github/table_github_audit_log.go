package github

import (
	"context"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubAuditLog(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_audit_log",
		Description: "Gets the audit logs for an organization.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.Required},
				{Name: "phrase", Require: plugin.Optional},
				{Name: "include", Require: plugin.Optional},
				{Name: "action", Require: plugin.Optional},
				{Name: "actor", Require: plugin.Optional},
				{Name: "created_at", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
			},
			Hydrate: tableGitHubAuditLogList,
		},
		Columns: []*plugin.Column{
			{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The GitHub organization."},
			{Name: "phrase", Type: proto.ColumnType_STRING, Transform: transform.FromQual("phrase"), Description: "The search phrase for your audit events."},
			{Name: "include", Type: proto.ColumnType_STRING, Transform: transform.FromQual("include"), Description: "The event types to include: web, git, all."},

			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The id of the audit event.", Transform: transform.FromField("DocumentID")},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp of the audit event.", Transform: transform.FromField("CreatedAt").Transform(convertTimestamp)},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "The action performed."},
			{Name: "actor", Type: proto.ColumnType_STRING, Description: "The GitHub user who performed the action."},
			{Name: "actor_location", Type: proto.ColumnType_JSON, Description: "The actor's location at the moment of the action."},

			// Optional columns, depending on the audit event
			{Name: "team", Type: proto.ColumnType_STRING, Description: "The GitHub team, when the action relates to a team."},
			{Name: "user", Type: proto.ColumnType_STRING, Description: "The GitHub user, when the action relates to a user."},
			{Name: "repo", Type: proto.ColumnType_STRING, Description: "The GitHub repository, when the action relates to a repository."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Additional data relating to the audit event."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubAuditLogList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	org := quals["organization"].GetStringValue()
	phrase := quals["phrase"].GetStringValue()
	include := quals["include"].GetStringValue()

	opts := &github.GetAuditLogOptions{
		Phrase:            &phrase,
		Include:           &include,
		ListCursorOptions: github.ListCursorOptions{PerPage: 100},
	}

	if d.Quals["created_at"] != nil {
		for _, q := range d.Quals["created_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime().Format(time.RFC3339)

			op := q.Operator
			if op == "=" {
				op = ""
			}

			phrase += " created:" + op + givenTime
			opts.Phrase = &phrase
		}
	}

	if quals["action"] != nil {
		phrase += " action:" + quals["action"].GetStringValue()
		opts.Phrase = &phrase
	}

	if quals["actor"] != nil {
		phrase += " actor:" + quals["actor"].GetStringValue()
		opts.Phrase = &phrase
	}

	type ListPageResponse struct {
		entries []*github.AuditEntry
		resp    *github.Response
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.ListCursorOptions.PerPage) {
			opts.ListCursorOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		plugin.Logger(ctx).Debug("tableGitHubAuditLogs", "org", org, "include", include, "phrase", *opts.Phrase)
		entries, resp, err := client.Organizations.GetAuditLog(ctx, org, opts)

		if err != nil {
			return nil, err
		}

		return ListPageResponse{
			entries: entries,
			resp:    resp,
		}, nil
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		auditResults := listResponse.entries
		resp := listResponse.resp

		for _, i := range auditResults {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.After == "" {
			break
		}

		opts.After = resp.After
	}

	return nil, nil
}
