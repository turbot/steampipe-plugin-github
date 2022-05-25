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

func tableGitHubSearchPullRequest(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_pull_request",
		Description: "Find pull requests by state and keyword.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchPullRequestList,
		},
		Columns: []*plugin.Column{
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The title of the pull request."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "The ID of the pull request."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the pull request."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state of the pull request."},
			{Name: "active_lock_reason", Type: proto.ColumnType_STRING, Description: "The active lock reason of the pull request."},
			{Name: "author_association", Type: proto.ColumnType_STRING, Description: "The author association of the pull request."},
			{Name: "body", Type: proto.ColumnType_STRING, Description: "The body of the pull request."},
			{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the pull request closed at."},
			{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments on the pull request."},
			{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API URL of the comments for the pull request."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the pull request created at."},
			{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The API URL of the events for the pull request."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The complete URL of the pull request."},
			{Name: "labels_url", Type: proto.ColumnType_STRING, Description: "The API URL of the labels for the pull request."},
			{Name: "locked", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the pull request is locked."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the pull request."},
			{Name: "number", Type: proto.ColumnType_INT, Description: "The number of the pull request."},
			{Name: "repository_url", Type: proto.ColumnType_STRING, Description: "The API URL of the repository for the pull request."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.From(extractSearchPullReqRepositoryFullName), Description: "The full name of the repository (login/repo-name)."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp the pull request updated at."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "The API URL of the pull request."},
			{Name: "assignee", Type: proto.ColumnType_JSON, Description: "The assignee details."},
			{Name: "assignees", Type: proto.ColumnType_JSON, Description: "The assignees details."},
			{Name: "closed_by", Type: proto.ColumnType_JSON, Description: "The details of the user that closed the pull request."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "The label details."},
			{Name: "milestone", Type: proto.ColumnType_JSON, Description: "The milestone details."},
			{Name: "pull_request_links", Type: proto.ColumnType_JSON, Description: "The pull request link details."},
			{Name: "reactions", Type: proto.ColumnType_JSON, Description: "The reaction details."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The repository details."},
			{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
			{Name: "user", Type: proto.ColumnType_JSON, Description: "The user details."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubSearchPullRequestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchPullRequestList")

	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	query = query + " is:pr"

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
			logger.Error("tableGitHubSearchPullRequestList", "error_Search.Issues", err)
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
			logger.Error("tableGitHubSearchPullRequestList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		prs := listResponse.result.Issues
		resp := listResponse.resp

		for _, i := range prs {
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

func extractSearchPullReqRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	pr := d.HydrateItem.(*github.Issue)
	if pr.RepositoryURL != nil {
		rx := regexp.MustCompile(`(https?://.+?/repos/)`)
		return rx.ReplaceAllString(*pr.RepositoryURL, ""), nil
	}
	return "", nil
}
