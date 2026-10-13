package github

import (
	"context"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func gitHubProjectV2ItemColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "organization", Type: proto.ColumnType_STRING, Transform: transform.FromQual("organization"), Description: "The organization that owns the project. Either this or 'user_login' must be provided."},
		{Name: "user_login", Type: proto.ColumnType_STRING, Transform: transform.FromQual("user_login"), Description: "The user that owns the project. Either this or 'organization' must be provided."},
		{Name: "project_number", Type: proto.ColumnType_INT, Transform: transform.FromQual("project_number"), Description: "The number of the project the item belongs to."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id"), Description: "The ID of the project item."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId"), Description: "The node ID of the project item."},
		{Name: "type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type"), Description: "The type of the item: ISSUE, PULL_REQUEST, DRAFT_ISSUE or REDACTED. REDACTED items are items whose content the token can't access, and have no content columns set."},
		{Name: "is_archived", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsArchived"), Description: "If true, the item is archived."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the item was added to the project."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the item was last updated."},
		{Name: "creator", Type: proto.ColumnType_JSON, Transform: transform.FromField("Creator").NullIfZero(), Description: "The actor who added the item to the project."},
		{Name: "field_values", Type: proto.ColumnType_JSON, Transform: transform.FromField("FieldValues"), Description: "The item's values for the project's custom fields (e.g. Status, Iteration), keyed by field name (first 50 fields). Built-in fields such as Title, Assignees and Labels are available as their own columns."},
		{Name: "content_id", Type: proto.ColumnType_INT, Transform: transform.FromField("ContentId"), Description: "The ID of the issue or pull request. Null for draft issues."},
		{Name: "content_node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ContentNodeId").NullIfZero(), Description: "The node ID of the issue, pull request or draft issue."},
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("RepositoryFullName").NullIfZero(), Description: "The full name of the repository of the issue or pull request."},
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number"), Description: "The number of the issue or pull request."},
		{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Title").NullIfZero(), Description: "The title of the issue, pull request or draft issue."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromField("Body"), Description: "The body of the issue, pull request or draft issue."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url").NullIfZero(), Description: "The URL of the issue or pull request."},
		{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromField("State").NullIfZero(), Description: "The state of the issue (OPEN, CLOSED) or pull request (OPEN, CLOSED, MERGED)."},
		{Name: "closed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Closed"), Description: "If true, the issue or pull request is closed."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ClosedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the issue or pull request was closed."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromField("Author").NullIfZero(), Description: "The actor who authored the issue or pull request, or created the draft issue."},
		{Name: "assignees", Type: proto.ColumnType_JSON, Transform: transform.FromField("Assignees"), Description: "The users assigned to the issue, pull request or draft issue (first 10)."},
		{Name: "assignees_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("AssigneesTotalCount"), Description: "The number of users assigned to the issue, pull request or draft issue."},
		{Name: "labels", Type: proto.ColumnType_JSON, Transform: transform.FromField("Labels"), Description: "A map of the labels on the issue or pull request (first 100)."},
		{Name: "milestone", Type: proto.ColumnType_JSON, Transform: transform.FromField("Milestone").NullIfZero(), Description: "The milestone of the issue or pull request."},
		{Name: "content_created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ContentCreatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the issue, pull request or draft issue was created."},
		{Name: "content_updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ContentUpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "The time when the issue, pull request or draft issue was last updated."},
	}
}

func tableGitHubProjectV2Item() *plugin.Table {
	return &plugin.Table{
		Name:        "github_project_v2_item",
		Description: "Items (issues, pull requests and draft issues) in a GitHub Project (V2).",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "organization", Require: plugin.AnyOf},
				{Name: "user_login", Require: plugin.AnyOf},
				{Name: "project_number", Require: plugin.Required},
				{Name: "type", Require: plugin.Optional},
				{Name: "is_archived", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "repository_full_name", Require: plugin.Optional},
			},
			Hydrate: tableGitHubProjectV2ItemList,
		},
		Columns: commonColumns(gitHubProjectV2ItemColumns()),
	}
}

// GitHub sometimes fails a page of items with a generic "Something went wrong"
// error when content assignees and field values are fetched together, while
// pages of 10 of the same items succeed.
const projectV2ItemFallbackPageSize = 10

type projectV2ItemConnection struct {
	PageInfo models.PageInfo
	Nodes    []models.ProjectV2Item
}

func tableGitHubProjectV2ItemList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	organization := d.EqualsQualString("organization")
	user := d.EqualsQualString("user_login")
	projectNumber := int(d.EqualsQuals["project_number"].GetInt64Value())

	maxPageSize := adjustPageSize(100, d.QueryContext.Limit)
	pageSize := maxPageSize
	// The number of fallback-sized pages left before going back to maxPageSize
	fallbackPagesLeft := 0
	variables := map[string]interface{}{
		"owner":          githubv4.String(organization),
		"number":         githubv4.Int(projectNumber),
		"pageSize":       githubv4.Int(pageSize),
		"cursor":         (*githubv4.String)(nil),
		"archivedStates": projectV2ItemArchivedStates(d),
		"query":          githubv4.String(projectV2ItemQuery(d)),
	}
	if organization == "" {
		variables["owner"] = githubv4.String(user)
	}
	appendProjectV2ItemColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	for {
		var items projectV2ItemConnection
		var rateLimit models.RateLimit
		var err error

		// The items connection is the same for both owner types, but the root
		// query field differs.
		if organization != "" {
			var query struct {
				RateLimit    models.RateLimit
				Organization struct {
					ProjectV2 struct {
						Items projectV2ItemConnection `graphql:"items(first: $pageSize, after: $cursor, archivedStates: $archivedStates, query: $query)"`
					} `graphql:"projectV2(number: $number)"`
				} `graphql:"organization(login: $owner)"`
			}
			err = client.Query(ctx, &query, variables)
			items, rateLimit = query.Organization.ProjectV2.Items, query.RateLimit
		} else {
			var query struct {
				RateLimit models.RateLimit
				User      struct {
					ProjectV2 struct {
						Items projectV2ItemConnection `graphql:"items(first: $pageSize, after: $cursor, archivedStates: $archivedStates, query: $query)"`
					} `graphql:"projectV2(number: $number)"`
				} `graphql:"user(login: $owner)"`
			}
			err = client.Query(ctx, &query, variables)
			items, rateLimit = query.User.ProjectV2.Items, query.RateLimit
		}
		plugin.Logger(ctx).Debug(rateLimitLogString("github_project_v2_item", &rateLimit))
		if err != nil {
			// Fetch the items from the failed page in fallback-sized pages, then go
			// back to the full page size.
			if strings.Contains(err.Error(), "Something went wrong while executing your query") && pageSize > projectV2ItemFallbackPageSize {
				fallbackPagesLeft = (pageSize + projectV2ItemFallbackPageSize - 1) / projectV2ItemFallbackPageSize
				pageSize = projectV2ItemFallbackPageSize
				variables["pageSize"] = githubv4.Int(pageSize)
				plugin.Logger(ctx).Warn("github_project_v2_item", "retrying_with_page_size", pageSize, "api_error", err)
				continue
			}
			plugin.Logger(ctx).Error("github_project_v2_item", "api_error", err)
			if strings.Contains(err.Error(), "Could not resolve to a ProjectV2") ||
				strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") ||
				strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
				return nil, nil
			}
			return nil, err
		}

		for _, item := range items.Nodes {
			d.StreamListItem(ctx, newProjectV2ItemRow(item))

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !items.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(items.PageInfo.EndCursor)
		if fallbackPagesLeft > 0 {
			fallbackPagesLeft--
			if fallbackPagesLeft == 0 {
				pageSize = maxPageSize
				variables["pageSize"] = githubv4.Int(pageSize)
			}
		}
	}

	return nil, nil
}
