package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubRepositoryIssue() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository_issue",
		Description: "Github Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubRepositoryIssueList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "issue_number"}),
			Hydrate:    tableGitHubRepositoryIssueGet,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repoNameQual, Transform: transform.FromValue(), Description: "The full name of the repository (login/repo-name)."},
			{Name: "issue_number", Type: proto.ColumnType_INT, Description: "The issue number.", Transform: transform.FromField("Number")},

			{Name: "assignees", Type: proto.ColumnType_JSON, Description: "An array of users that are assigned to the issue."},
			{Name: "author", Type: proto.ColumnType_JSON, Description: "The user that submitted the issue.", Transform: transform.FromField("User")},
			{Name: "author_association", Type: proto.ColumnType_STRING, Description: "The association of the issue author to the repository (COLLABORATOR,CONTRIBUTOR, etc)."},
			{Name: "body", Type: proto.ColumnType_STRING, Description: "The body of the issue text."},
			{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the issue was closed."},
			{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments on the issue."},
			{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API Comments URL."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "he timestamp when the issue was created."},
			{Name: "events_url", Type: proto.ColumnType_STRING, Description: "The API Events URL."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the issue page in GitHub."},
			{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the issue."},
			{Name: "is_pull_request", Type: proto.ColumnType_BOOL, Description: "It true, the issue is a pull request.", Transform: transform.From(isPullRequest)},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "An array of labels associated with this issue."},
			{Name: "labels_url", Type: proto.ColumnType_STRING, Description: "The API Labels URL."},
			{Name: "locked", Type: proto.ColumnType_BOOL, Description: "If true, the issue is locked."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the issue."},
			{Name: "milestone", Type: proto.ColumnType_JSON, Description: "The milestone this issue is associated with."},
			{Name: "pull_request_links", Type: proto.ColumnType_JSON, Description: "Links to pull request details (If this issue is a pull request)."},
			{Name: "repository_url", Type: proto.ColumnType_STRING, Description: "The API Repository URL."},
			{Name: "reactions", Type: proto.ColumnType_JSON, Description: "An object containing the count of various reactions on the issue."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "The state or the issue (open, closed)."},
			{Name: "tags", Type: proto.ColumnType_JSON, Description: "A map of label names associated with this issue, in Steampipe standard format.", Transform: transform.From(getIssueTags)},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "The issue title."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the issue was last updated."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL of the issue."},
		},
	}
}

//// HYDRATE FUNCTIONS

func tableGitHubRepositoryIssueList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	logger.Trace("tableGitHubRepositoryIssueList", "owner", owner, "repo", repo)

	// TO DO - get state and other filters from the quals
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	client := connect(ctx, d.ConnectionManager)

	for {
		var issues []*github.Issue
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			issues, resp, err = client.Issues.ListByRepo(ctx, owner, repo, opt)

			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range issues {
			d.StreamListItem(ctx, i)
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

	client := connect(ctx, d.ConnectionManager)

	var detail *github.Issue
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		detail, resp, err = client.Issues.Get(ctx, owner, repo, issueNumber)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return detail, nil
}

func repoNameQual(_ context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

//// HELPER FUNCTIONS

func parseRepoFullName(fullName string) (string, string) {
	owner := ""
	repo := ""
	s := strings.Split(fullName, "/")
	owner = s[0]
	if len(s) > 1 {
		repo = s[1]
	}
	return owner, repo

}
