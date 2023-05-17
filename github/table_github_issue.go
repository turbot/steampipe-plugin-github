package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubIssueColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
	}

	return append(tableCols, sharedIssueColumns()...)
}

func sharedIssueColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number", "Node.Number"), Description: "The issue number."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID of the issue."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the issue."},
		{Name: "active_lock_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("ActiveLockReason", "Node.ActiveLockReason"), Description: "Reason that the conversation was locked."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromField("Author", "Node.Author").NullIfZero(), Description: "The actor who authored the issue."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.Login", "Node.Author.Login"), Description: "The login of the issue author."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Transform: transform.FromField("AuthorAssociation", "Node.AuthorAssociation"), Description: "Author's association with the subject of the issue."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromField("Body", "Node.Body"), Description: "Identifies the body of the issue."},
		{Name: "body_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("BodyUrl", "Node.BodyUrl"), Description: "URL for this issue body."},
		{Name: "closed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Closed", "Node.Closed"), Description: "If true, issue is closed."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ClosedAt", "Node.ClosedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when issue was closed."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt", "Node.CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when issue was created."},
		{Name: "created_via_email", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CreatedViaEmail", "Node.CreatedViaEmail"), Description: "If true, issue was created via email."},
		{Name: "editor", Type: proto.ColumnType_JSON, Transform: transform.FromField("Editor", "Node.Editor").NullIfZero(), Description: "The actor who edited the issue."},
		{Name: "full_database_id", Type: proto.ColumnType_INT, Transform: transform.FromField("FullDatabaseId", "Node.FullDatabaseId"), Description: "Identifies the primary key from the database as a BigInt."},
		{Name: "includes_created_edit", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IncludesCreatedEdit", "Node.IncludesCreatedEdit"), Description: "If true, issue was edited and includes an edit with the creation data."},
		{Name: "is_pinned", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsPinned", "Node.IsPinned"), Description: "if true, this issue is currently pinned to the repository issues list."},
		{Name: "is_read_by_user", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsReadByUser", "Node.IsReadByUser"), Description: "if true, this issue has been read by the user."},
		{Name: "last_edited_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastEditedAt", "Node.LastEditedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when issue was last edited."},
		{Name: "locked", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Locked", "Node.Locked"), Description: "If true, issue is locked."},
		{Name: "milestone", Type: proto.ColumnType_JSON, Transform: transform.FromField("Milestone", "Node.Milestone").NullIfZero(), Description: "The milestone associated with the issue."},
		{Name: "published_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PublishedAt", "Node.PublishedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when issue was published."},
		{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromField("State", "Node.State"), Description: "The state of the issue."},
		{Name: "state_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("StateReason", "Node.StateReason"), Description: "The reason for the issue state."},
		{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Title", "Node.Title"), Description: "The title of the issue."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt", "Node.UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when issue was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url", "Node.Url"), Description: "URL for the issue."},
		{Name: "user_can_close", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserCanClose", "Node.UserCanClose"), Description: "If true, user can close the issue."},
		{Name: "user_can_react", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserCanReact", "Node.UserCanReact"), Description: "If true, user can react on the issue."},
		{Name: "user_can_reopen", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserCanReopen", "Node.UserCanReopen"), Description: "If true, user can reopen the issue."},
		{Name: "user_can_subscribe", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserCanSubscribe", "Node.UserCanSubscribe"), Description: "If true, user can subscribe to the issue."},
		{Name: "user_can_update", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserCanUpdate", "Node.UserCanUpdate"), Description: "If true, user can update the issue,"},
		{Name: "user_cannot_update_reasons", Type: proto.ColumnType_JSON, Transform: transform.FromField("UserCannotUpdateReasons", "Node.UserCannotUpdateReasons").NullIfZero(), Description: "A list of reason why user cannot update the issue."},
		{Name: "user_did_author", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UserDidAuthor", "Node.UserDidAuthor"), Description: "If true, user authored the issue."},
		{Name: "user_subscription", Type: proto.ColumnType_STRING, Transform: transform.FromField("UserSubscription", "Node.UserSubscription"), Description: "Subscription state of the user to the issue."},
	}
}

func tableGitHubIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_issue",
		Description: "GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "repository_full_name",
					Require: plugin.Required,
				},
				{
					Name:    "author_login",
					Require: plugin.Optional,
				},
				{
					Name:    "state",
					Require: plugin.Optional,
				},
				{
					Name:      "updated_at",
					Require:   plugin.Optional,
					Operators: []string{">", ">="},
				},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryIssueList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubRepositoryIssueGet,
		},
		Columns: gitHubIssueColumns(),
	}
}

func tableGitHubRepositoryIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var filters githubv4.IssueFilters

	if quals["state"] != nil {
		state := quals["state"].GetStringValue()
		switch state {
		case "OPEN":
			filters.States = &[]githubv4.IssueState{githubv4.IssueStateOpen}
		case "CLOSED":
			filters.States = &[]githubv4.IssueState{githubv4.IssueStateClosed}
		default:
			plugin.Logger(ctx).Error("github_issue", "invalid filter", "state", state)
			return nil, fmt.Errorf("invalid value for 'state' can only filter for 'OPEN' or 'CLOSED' - you attempted to filter for '%s'", state)
		}
	} else {
		filters.States = &[]githubv4.IssueState{githubv4.IssueStateOpen, githubv4.IssueStateClosed}
	}

	if quals["author_login"] != nil {
		author := quals["author_login"].GetStringValue()
		filters.CreatedBy = githubv4.NewString(githubv4.String(author))
	}

	if d.Quals["updated_at"] != nil {
		for _, q := range d.Quals["updated_at"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				filters.Since = githubv4.NewDateTime(githubv4.DateTime{Time: afterTime})
			case ">=":
				filters.Since = githubv4.NewDateTime(githubv4.DateTime{Time: givenTime})
			}
		}
	}

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Issues struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.Issue
			} `graphql:"issues(first: $pageSize, after: $cursor, filterBy: $filters)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repoName),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"filters":  filters,
	}

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_issue", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_issue", "api_error", err)
			return nil, err
		}

		for _, issue := range query.Repository.Issues.Nodes {
			d.StreamListItem(ctx, issue)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Issues.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubRepositoryIssueGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	issueNumber := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	client := connectV4(ctx, d)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Issue models.Issue `graphql:"issue(number: $issueNumber)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"repo":        githubv4.String(repo),
		"issueNumber": githubv4.Int(issueNumber),
	}

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_issue", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_issue", "api_error", err)
		return nil, err
	}

	return query.Repository.Issue, nil
}
