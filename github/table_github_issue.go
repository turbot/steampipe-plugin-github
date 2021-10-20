package github

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func gitHubIssueColumns() []*plugin.Column {

	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repoNameQual, Transform: transform.FromValue(), Description: "The full name of the repository (login/repo-name)."},
		{Name: "issue_number", Type: proto.ColumnType_INT, Description: "The issue number.", Transform: transform.FromField("Number")},
		{Name: "title", Type: proto.ColumnType_STRING, Description: "The issue title."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Description: "The login name of the user that submitted the PR.", Transform: transform.FromField("User.Login")},
		{Name: "author_association", Type: proto.ColumnType_STRING, Description: "The association of the issue author to the repository (COLLABORATOR,CONTRIBUTOR, etc)."},
		{Name: "assignee_logins", Type: proto.ColumnType_JSON, Description: "An array of user login names that are assigned to the issue.", Transform: transform.FromField("Assignees").Transform(filterUserLogins)},

		{Name: "body", Type: proto.ColumnType_STRING, Description: "The body of the issue text."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the issue was closed."},
		{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments on the issue."},
		{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API Comments URL."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the issue was created."},
		{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The API Events URL."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the issue page in GitHub."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the issue."},
		{Name: "labels", Type: proto.ColumnType_JSON, Description: "An array of labels associated with this issue."},
		{Name: "labels_url", Type: proto.ColumnType_STRING, Description: "The API Labels URL."},
		{Name: "locked", Type: proto.ColumnType_BOOL, Description: "If true, the issue is locked."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the issue."},
		{Name: "milestone_id", Type: proto.ColumnType_INT, Description: "The milestone id this issue is associated with.", Transform: transform.FromField("Milestone.ID")},
		{Name: "milestone_title", Type: proto.ColumnType_STRING, Description: "The title of the milestone this issue is associated with.", Transform: transform.FromField("Milestone.Title")},
		{Name: "repository_url", Type: proto.ColumnType_STRING, Description: "The API Repository URL."},
		{Name: "reactions", Type: proto.ColumnType_JSON, Description: "An object containing the count of various reactions on the issue."},
		{Name: "state", Type: proto.ColumnType_STRING, Description: "The state or the issue (open, closed)."},
		{Name: "tags", Type: proto.ColumnType_JSON, Description: "A map of label names associated with this issue, in Steampipe standard format.", Transform: transform.From(getIssueTags)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the issue was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL of the issue."},
	}
}

func tableGitHubIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_issue",
		Description: "GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubRepositoryIssueList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "issue_number"}),
			Hydrate:    tableGitHubRepositoryIssueGet,
		},
		Columns: gitHubIssueColumns(),
	}
}

//// HYDRATE FUNCTIONS

func tableGitHubRepositoryIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	logger.Trace("tableGitHubRepositoryIssueList", "owner", owner, "repo", repo)

	type ListPageResponse struct {
		issues []*github.Issue
		resp   *github.Response
	}

	// TO DO - get state and other filters from the quals
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	client := connect(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
		return ListPageResponse{
			issues: issues,
			resp:   resp,
		}, err
	}

	for {

		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{shouldRetryError})

		if err != nil {
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		issues := listResponse.issues
		resp := listResponse.resp

		for _, i := range issues {
			// Only issues, not PRs (those are in the pull_request table...)
			if !i.IsPullRequest() {
				d.StreamListItem(ctx, i)
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return nil, nil
}

func tableGitHubRepositoryIssueGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var issueNumber int

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	if h.Item != nil {
		issue := h.Item.(*github.Issue)
		issueNumber = *issue.Number

	} else {
		issueNumber = int(d.KeyColumnQuals["issue_number"].GetInt64Value())
	}

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubRepositoryIssueGet", "owner", owner, "repo", repo, "issueNumber", issueNumber)

	client := connect(ctx, d)

	type GetResponse struct {
		issue *github.Issue
		resp  *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Issues.Get(ctx, owner, repo, issueNumber)
		return GetResponse{
			issue: detail,
			resp:  resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	issue := getResp.issue

	return issue, nil
}

func repoNameQual(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	item := h.Item.(*github.Issue)
	if item.Repository != nil && item.Repository.FullName != nil {
		return item.Repository.FullName, nil
	}

	return d.KeyColumnQuals["repository_full_name"].GetStringValue(), nil
}

//// TRANSFORM FUNCTIONS

func getIssueTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	issue := d.HydrateItem.(*github.Issue)

	tags := make(map[string]bool)
	if issue.Labels != nil {
		for _, i := range issue.Labels {
			tags[*i.Name] = true
		}
	}
	return tags, nil
}

func isPullRequest(_ context.Context, d *transform.TransformData) (interface{}, error) {
	issue := d.HydrateItem.(*github.Issue)
	return issue.IsPullRequest(), nil
}
