package github

import (
	"context"
	"regexp"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubSearchIssue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_issue",
		Description: "Find issues by state and keyword.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchIssueList,
		},
		Columns: []*plugin.Column{
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the issue."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "The ID of the issue."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the issue."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the issue."},
			{Name: "active_lock_reason", Type: proto.ColumnType_STRING, Description: "The active lock reason of the issue."},
			{Name: "author_association", Type: proto.ColumnType_STRING, Description: "The author association of the issue."},
			{Name: "body", Type: proto.ColumnType_STRING, Description: "The body of the issue."},
			{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the issue closed at."},
			{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments on the issue."},
			{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API URL of the comments for the issue."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the issue created at."},
			{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The API URL of the events for the issue."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The complete URL of the issue."},
			{Name: "labels_url", Type: proto.ColumnType_STRING, Description: "The API URL of the labels for the issue."},
			{Name: "locked", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the issue is locked."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the issue."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "The number of the issue."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.From(extractRepositoryFullName), Description: "The full name of the repository (login/repo-name)."},
			{Name: "repository_url", Type: proto.ColumnType_STRING, Description: "The API URL of the repository for the issue."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the issue updated at."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "The API URL of the issue."},
			{Name: "assignee", Type: proto.ColumnType_JSON, Description: "The assignee details."},
			{Name: "assignees", Type: proto.ColumnType_JSON, Description: "The assignees details."},
			{Name: "closed_by", Type: proto.ColumnType_JSON, Description: "The details of the user that closed the issue."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "The label details."},
			{Name: "milestone", Type: proto.ColumnType_JSON, Description: "The milestone details."},
			{Name: "reactions", Type: proto.ColumnType_JSON, Description: "The reaction details."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The repository details."},
			{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
			{Name: "user", Type: proto.ColumnType_JSON, Description: "The user details."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubSearchIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchIssueList")

	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	query = query + " is:issue"

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.IssuesSearchResult
		resp   *github.Response
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		result, resp, err := client.Search.Issues(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchIssueList", "error_Search.Issues", err)
			return nil, err
		}

		return ListPageResponse{
			result: result,
			resp:   resp,
		}, nil
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

		if err != nil {
			logger.Error("tableGitHubSearchIssueList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		issues := listResponse.result.Issues
		resp := listResponse.resp

		for _, i := range issues {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func extractRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	issue := d.HydrateItem.(*github.Issue)
	if issue.RepositoryURL != nil {
		rx := regexp.MustCompile(`(https?://.+?/repos/)`)
		return rx.ReplaceAllString(*issue.RepositoryURL, ""), nil
	}
	return "", nil
}
