package github

import (
	"context"
	"encoding/json"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubSearchPullRequestColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number", "Node.Number"), Description: "The number of the pull request."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID of the pull request."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the pull request."},
		{Name: "active_lock_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("ActiveLockReason", "Node.ActiveLockReason"), Description: "Reason that the conversation was locked."},
		{Name: "additions", Type: proto.ColumnType_INT, Transform: transform.FromField("Additions", "Node.Additions"), Description: "The number of additions in this pull request."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromField("Author", "Node.Author").NullIfZero(), Description: "The author of the pull request."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Transform: transform.FromField("AuthorAssociation", "Node.AuthorAssociation"), Description: "Author's association with the pull request."},
		{Name: "base_ref_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("BaseRefName", "Node.BaseRefName"), Description: "Identifies the name of the base Ref associated with the pull request, even if the ref has been deleted."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromField("Body", "Node.Body"), Description: "The body as Markdown."},
		{Name: "changed_files", Type: proto.ColumnType_INT, Transform: transform.FromField("ChangedFiles", "Node.ChangedFiles"), Description: "The number of files changed in this pull request."},
		{Name: "checks_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ChecksUrl", "Node.ChecksUrl"), Description: "URL for the checks of this pull request."},
		{Name: "closed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Closed", "Node.Closed"), Description: "If true, pull request is closed."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ClosedAt", "Node.ClosedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the pull request was closed."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt", "Node.CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the pull request was created."},
		{Name: "created_via_email", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CreatedViaEmail", "Node.CreatedViaEmail"), Description: "If true, pull request comment was created via email."},
		{Name: "deletions", Type: proto.ColumnType_INT, Transform: transform.FromField("Deletions", "Node.Deletions"), Description: "The number of deletions in this pull request."},
		{Name: "editor", Type: proto.ColumnType_JSON, Transform: transform.FromField("Editor", "Node.Editor").NullIfZero(), Description: "The actor who edited the pull request's body."},
		{Name: "head_ref_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("HeadRefName", "Node.HeadRefName"), Description: "Identifies the name of the head Ref associated with the pull request, even if the ref has been deleted."},
		{Name: "head_ref_oid", Type: proto.ColumnType_STRING, Transform: transform.FromField("HeadRefOid", "Node.HeadRefOid"), Description: "Identifies the oid/sha of the head ref associated with the pull request, even if the ref has been deleted."},
		{Name: "includes_created_edit", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IncludesCreatedEdit", "Node.IncludesCreatedEdit"), Description: "If true, this pull request was edited and includes an edit with the creation data."},
		{Name: "is_cross_repository", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsCrossRepository", "Node.IsCrossRepository"), Description: "If true, head and base repositories are different."},
		{Name: "is_draft", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsDraft", "Node.IsDraft"), Description: "If true, the pull request is a draft."},
		{Name: "is_read_by_user", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsReadByUser", "Node.IsReadByUser"), Description: "If true, this pull request was read by the current user."},
		{Name: "last_edited_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastEditedAt", "Node.LastEditedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp the editor made the last edit."},
		{Name: "locked", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Locked", "Node.Locked"), Description: "If true, the pull request is locked."},
		{Name: "maintainer_can_modify", Type: proto.ColumnType_BOOL, Transform: transform.FromField("MaintainerCanModify", "Node.MaintainerCanModify"), Description: "If true, maintainers can modify the pull request."},
		{Name: "mergeable", Type: proto.ColumnType_STRING, Transform: transform.FromField("Mergeable", "Node.Mergeable"), Description: "Whether or not the pull request can be merged based on the existence of merge conflicts."},
		{Name: "merged", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Merged", "Node.Merged"), Description: "If true, the pull request was merged."},
		{Name: "merged_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("MergedAt", "Node.MergedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when pull request was merged."},
		{Name: "merged_by", Type: proto.ColumnType_JSON, Transform: transform.FromField("MergedBy", "Node.MergedBy").NullIfZero(), Description: "The actor who merged the pull request."},
		{Name: "milestone", Type: proto.ColumnType_JSON, Transform: transform.FromField("Milestone", "Node.Milestone").NullIfZero(), Description: "The milestone associated with the pull request."},
		{Name: "permalink", Type: proto.ColumnType_STRING, Transform: transform.FromField("Permalink", "Node.Permalink"), Description: "Permanent URL for the pull request."},
		{Name: "published_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PublishedAt", "Node.PublishedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp the pull request was published."},
		{Name: "revert_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("RevertUrl", "Node.RevertUrl"), Description: "URL to revert the pull request."},
		{Name: "review_decision", Type: proto.ColumnType_STRING, Transform: transform.FromField("ReviewDecision", "Node.ReviewDecision"), Description: "The current status of this pull request with respect to code review."},
		{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromField("State", "Node.State"), Description: "The current state of the pull request."},
		{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Title", "Node.Title"), Description: "The title of the pull request."},
		{Name: "total_comments_count", Type: proto.ColumnType_INT, Transform: transform.FromField("TotalCommentsCount", "Node.TotalCommentsCount"), Description: "The number of comments on the pull request."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt", "Node.UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the pull request was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url", "Node.Url"), Description: "URL of the pull request."},
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Repo.NameWithOwner", "Node.Repo.NameWithOwner"), Description: "The full name of the repository the pull request belongs to."},
	}

	return append(defaultSearchColumns(), tableCols...)
}

func tableGitHubSearchPullRequest(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_pull_request",
		Description: "Find pull requests by state and keyword.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchPullRequestList,
		},
		Columns: gitHubSearchPullRequestColumns(),
	}
}

func tableGitHubSearchPullRequestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	input := quals["query"].GetStringValue()

	if input == "" {
		return nil, nil
	}

	input += " is:pr"

	var query struct {
		RateLimit models.RateLimit
		Search    struct {
			PageInfo models.PageInfo
			Edges    []struct {
				TextMatches []models.TextMatch
				Node        struct {
					models.BasicPullRequest `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(type: ISSUE, first: $pageSize, after: $cursor, query: $query)"`
	}

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"query":    githubv4.String(input),
	}

	qj, _ := json.Marshal(query)
	plugin.Logger(ctx).Debug(string(qj))

	client := connectV4(ctx, d)
	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_search_pull_request", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_search_pull_request", "api_error", err)
			return nil, err
		}

		for _, pr := range query.Search.Edges {
			d.StreamListItem(ctx, pr)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Search.PageInfo.EndCursor)
	}

	return nil, nil
}
