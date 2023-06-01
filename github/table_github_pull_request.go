package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func sharedPullRequestColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Repo.NameWithOwner", "Node.Repo.NameWithOwner"), Description: "The full name of the repository the pull request belongs to."},
		{Name: "number", Type: proto.ColumnType_INT, Transform: transform.FromField("Number", "Node.Number"), Description: "The number of the pull request."},
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Id", "Node.Id"), Description: "The ID of the pull request."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId", "Node.NodeId"), Description: "The node ID of the pull request."},
		{Name: "active_lock_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("ActiveLockReason", "Node.ActiveLockReason"), Description: "Reason that the conversation was locked."},
		{Name: "additions", Type: proto.ColumnType_INT, Transform: transform.FromField("Additions", "Node.Additions"), Description: "The number of additions in this pull request."},
		{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromField("Author", "Node.Author").NullIfZero(), Description: "The author of the pull request."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Transform: transform.FromField("AuthorAssociation", "Node.AuthorAssociation"), Description: "Author's association with the pull request."},
		{Name: "base_ref_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("BaseRefName", "Node.BaseRefName"), Description: "Identifies the name of the base Ref associated with the pull request, even if the ref has been deleted."},
		{Name: "body", Type: proto.ColumnType_STRING, Transform: transform.FromField("Body", "Node.Body"), Description: "The body as Markdown."},
		// {Name: "can_be_rebased", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanBeRebased", "Node.CanBeRebased"), Description: "If true, the pull request is rebaseable."},
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
	}
}

func gitHubPullRequestColumns() []*plugin.Column {
	tableCols := []*plugin.Column{
		{Name: "base_ref", Type: proto.ColumnType_JSON, Transform: transform.FromField("BaseRef", "Node.BaseRef").NullIfZero(), Description: "The base ref associated with the pull request."},
		{Name: "head_ref", Type: proto.ColumnType_JSON, Transform: transform.FromField("HeadRef", "Node.HeadRef").NullIfZero(), Description: "The head ref associated with the pull request."},
		{Name: "merge_commit", Type: proto.ColumnType_JSON, Transform: transform.FromField("MergeCommit", "Node.MergeCommit").NullIfZero(), Description: "The merge commit associated the pull request, null if not merged."},
		{Name: "suggested_reviewers", Type: proto.ColumnType_JSON, Transform: transform.FromField("SuggestedReviewers", "Node.SuggestedReviewers").NullIfZero(), Description: "Suggested reviewers for the pull request."},
		{Name: "can_apply_suggestion", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanApplySuggestion", "Node.CanApplySuggestion"), Description: "If true, current user can apply suggestions."},
		{Name: "can_close", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanClose", "Node.CanClose"), Description: "If true, current user can close the pull request."},
		{Name: "can_delete_head_ref", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanDeleteHeadRef", "Node.CanDeleteHeadRef"), Description: "If true, current user can delete/restore head ref."},
		{Name: "can_disable_auto_merge", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanDisableAutoMerge", "Node.CanDisableAutoMerge"), Description: "If true, current user can disable auto-merge."},
		{Name: "can_edit_files", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanEditFiles", "Node.CanEditFiles"), Description: "If true, current user can edit files within this pull request."},
		{Name: "can_enable_auto_merge", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanEnableAutoMerge", "Node.CanEnableAutoMerge"), Description: "If true, current user can enable auto-merge."},
		{Name: "can_merge_as_admin", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanMergeAsAdmin", "Node.CanMergeAsAdmin"), Description: "If true, current user can bypass branch protections and merge the pull request immediately."},
		{Name: "can_react", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanReact", "Node.CanReact"), Description: "If true, current user can react to the pull request."},
		{Name: "can_reopen", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanReopen", "Node.CanReopen"), Description: "If true, current user can reopen the pull request."},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanSubscribe", "Node.CanSubscribe"), Description: "If true, current user can subscribe to the pull request."},
		{Name: "can_update", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanUpdate", "Node.CanUpdate"), Description: "If true, current user can update the pull request."},
		{Name: "can_update_branch", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanUpdateBranch", "Node.CanUpdateBranch"), Description: "If true, current user can update the head ref of the pull request by merging or rebasing the base ref."},
		{Name: "did_author", Type: proto.ColumnType_BOOL, Transform: transform.FromField("DidAuthor", "Node.DidAuthor"), Description: "If true, current user authored the pull request."},
		{Name: "cannot_update_reasons", Type: proto.ColumnType_JSON, Transform: transform.FromField("CannotUpdateReasons", "Node.CannotUpdateReasons").NullIfZero(), Description: "Reasons why the current user cannot update the pull request, if applicable."},
		{Name: "subscription", Type: proto.ColumnType_STRING, Transform: transform.FromField("Subscription", "Node.Subscription"), Description: "Status of current users subscription to the pull request."},
		{Name: "assignees_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Assignees.TotalCount", "Node.Assignees.TotalCount"), Description: "A count of users assigned to the pull request."},
		{Name: "labels_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Labels.TotalCount", "Node.Labels.TotalCount"), Description: "A count of labels applied to the pull request."},
		{Name: "commits_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Commits.TotalCount", "Node.Commits.TotalCount"), Description: "A count of commits in the pull request."},
		{Name: "review_requests_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("ReviewRequests.TotalCount", "Node.ReviewRequests.TotalCount"), Description: "A count of reviews requested on the pull request."},
		{Name: "reviews_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Reviews.TotalCount", "Node.Reviews.TotalCount"), Description: "A count of completed reviews on the pull request."},
	}

	return append(sharedPullRequestColumns(), tableCols...)
}

func tableGitHubPullRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "github_pull_request",
		Description: "GitHub Pull requests let you tell others about changes you've pushed to a branch in a repository on GitHub. Once a pull request is opened, you can discuss and review the potential changes with collaborators and add follow-up commits before your changes are merged into the base branch.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "state", Require: plugin.Optional},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubPullRequestList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "number"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubPullRequestGet,
		},
		Columns: gitHubPullRequestColumns(),
	}
}

func tableGitHubPullRequestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := adjustPageSize(75, d.QueryContext.Limit)

	states := []githubv4.PullRequestState{githubv4.PullRequestStateOpen, githubv4.PullRequestStateClosed, githubv4.PullRequestStateMerged}
	if quals["state"] != nil {
		state := quals["state"].GetStringValue()
		switch state {
		case "OPEN":
			states = []githubv4.PullRequestState{githubv4.PullRequestStateOpen}
		case "CLOSED":
			states = []githubv4.PullRequestState{githubv4.PullRequestStateClosed}
		case "MERGED":
			states = []githubv4.PullRequestState{githubv4.PullRequestStateMerged}
		default:
			plugin.Logger(ctx).Error("github_pull_request", "invalid filter", "state", state)
			return nil, fmt.Errorf("invalid value for 'state' can only filter for 'OPEN', 'CLOSED' or 'MERGED', value provided was '%s'", state)
		}
	}

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			PullRequests struct {
				PageInfo   models.PageInfo
				TotalCount int
				Nodes      []models.PullRequest
			} `graphql:"pullRequests(first: $pageSize, after: $cursor, states: $states)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"states":   states,
	}

	client := connectV4(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := retryHydrate(ctx, d, h, listPage)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_pull_request", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_pull_request", "api_error", err)
			return nil, err
		}

		for _, issue := range query.Repository.PullRequests.Nodes {
			d.StreamListItem(ctx, issue)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.PullRequests.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubPullRequestGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	number := int(quals["number"].GetInt64Value())
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	client := connectV4(ctx, d)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			PullRequest models.PullRequest `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"number": githubv4.Int(number),
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	_, err := retryHydrate(ctx, d, h, listPage)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_pull_request", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_pull_request", "api_error", err)
		return nil, err
	}

	return query.Repository.PullRequest, nil
}
