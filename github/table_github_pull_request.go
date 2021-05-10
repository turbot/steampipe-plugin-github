package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func gitHubPullRequestColumns() []*plugin.Column {

	return []*plugin.Column{
		{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: RepoNameFromQuals, Transform: transform.FromValue(), Description: "The full name of the repository (login/repo-name)."},
		{Name: "issue_number", Type: proto.ColumnType_INT, Description: "The PR issue number.", Transform: transform.FromField("Number")},
		{Name: "title", Type: proto.ColumnType_STRING, Description: "The PR issue title."},
		{Name: "author_login", Type: proto.ColumnType_STRING, Description: "The login name of the user that submitted the PR.", Transform: transform.FromField("User.Login")},
		{Name: "assignee_logins", Type: proto.ColumnType_JSON, Description: "An array of user login names that are assigned to the issue.", Transform: transform.FromField("Assignees").Transform(filterUserLogins)},

		{Name: "additions", Type: proto.ColumnType_INT, Hydrate: tableGitHubPullRequestGet, Description: "The number of additions in this PR."},
		{Name: "author_association", Type: proto.ColumnType_STRING, Description: "The association of the PR issue author to the repository (COLLABORATOR,CONTRIBUTOR, etc)."},
		{Name: "body", Type: proto.ColumnType_STRING, Description: "The body of the PR issue text."},
		{Name: "changed_files", Type: proto.ColumnType_INT, Hydrate: tableGitHubPullRequestGet, Description: "The number of changed files in this PR."},
		{Name: "closed_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the PR was closed."},
		{Name: "comments", Type: proto.ColumnType_INT, Description: "The number of comments on the PR."},
		{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API Comments URL."},
		{Name: "commits", Type: proto.ColumnType_INT, Hydrate: tableGitHubPullRequestGet, Description: "The number of commits in this PR."},
		{Name: "commits_url", Type: proto.ColumnType_STRING, Description: "The URL of the Commits page in GitHub."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the PR was created."},
		{Name: "deletions", Type: proto.ColumnType_INT, Hydrate: tableGitHubPullRequestGet, Description: "The number of deletions in this PR."},
		{Name: "diff_url", Type: proto.ColumnType_STRING, Description: "The URL of the Diff page in GitHub."},
		{Name: "draft", Type: proto.ColumnType_BOOL, Description: "If true, the PR is in draft."},
		{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The URL of the PR page in GitHub."},
		{Name: "id", Type: proto.ColumnType_INT, Description: "The unique ID number of the PR issue."},
		{Name: "issue_url", Type: proto.ColumnType_STRING, Description: "The URL of the Issue page in GitHub."},
		{Name: "labels", Type: proto.ColumnType_JSON, Description: "An array of labels associated with this PR."},
		{Name: "locked", Type: proto.ColumnType_BOOL, Description: "If true, the PR is locked."},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node id of the PR."},
		{Name: "maintainer_can_modify", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubPullRequestGet, Description: "If true, people with push access to the upstream repository of a fork owned by a user account can commit to the forked branches."},
		{Name: "mergeable", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubPullRequestGet, Description: "If true, the PR can be merged."},
		{Name: "mergeable_state", Type: proto.ColumnType_STRING, Hydrate: tableGitHubPullRequestGet, Description: "The mergeability state of the PR."},
		{Name: "merged", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubPullRequestGet, Description: "If true, the PR has been merged."},
		{Name: "merged_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the PR was merged."},
		{Name: "merged_by_login", Type: proto.ColumnType_STRING, Hydrate: tableGitHubPullRequestGet, Description: "The login name of the user that submitted the PR.", Transform: transform.FromField("MergedBy.Login")},
		{Name: "merge_commit_sha", Type: proto.ColumnType_STRING, Description: "The commit sha of the merged PR.", Transform: transform.FromField("MergeCommitSHA")},
		{Name: "milestone_id", Type: proto.ColumnType_INT, Description: "The milestone id this issue is associated with.", Transform: transform.FromField("Milestone.ID")},
		{Name: "milestone_title", Type: proto.ColumnType_STRING, Description: "The title of the milestone this issue is associated with.", Transform: transform.FromField("Milestone.Title")},
		{Name: "patch_url", Type: proto.ColumnType_STRING, Description: "The URL of the Patch page in GitHub."},
		{Name: "rebaseable", Type: proto.ColumnType_BOOL, Hydrate: tableGitHubPullRequestGet, Description: "If true, the PR can be rebased."},
		{Name: "requested_reviewer_logins", Type: proto.ColumnType_JSON, Description: "An array of user login names that are requested reviewers of the PR.", Transform: transform.FromField("RequestedReviewers").Transform(filterUserLogins)},
		{Name: "review_comments", Type: proto.ColumnType_INT, Hydrate: tableGitHubPullRequestGet, Description: "The number of review comments in this PR."},
		{Name: "review_comments_url", Type: proto.ColumnType_STRING, Description: "The URL of the Review Comments page in GitHub."},
		{Name: "review_comment_url", Type: proto.ColumnType_STRING, Description: "The URL of the Review Comment page in GitHub."},
		{Name: "state", Type: proto.ColumnType_STRING, Description: "The state or the PR (open, closed)."},
		{Name: "statuses_url", Type: proto.ColumnType_STRING, Description: "The URL of the Statuses page in GitHub."},
		{Name: "tags", Type: proto.ColumnType_JSON, Description: "A map of label names associated with this PR, in Steampipe standard format.", Transform: transform.From(getPullRequestTags)},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the PR was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL of the PR."},
	}
}

func tableGitHubPullRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "github_pull_request",
		Description: "GitHub Pull requests let you tell others about changes you've pushed to a branch in a repository on GitHub. Once a pull request is opened, you can discuss and review the potential changes with collaborators and add follow-up commits before your changes are merged into the base branch.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubPullRequestList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "issue_number"}),
			Hydrate:    tableGitHubPullRequestGet,
		},
		Columns: gitHubPullRequestColumns(),
	}
}

//// HYDRATE FUNCTIONS

func tableGitHubPullRequestList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	logger.Trace("tableGitHubPullRequestList", "owner", owner, "repo", repo)

	// TO DO - get state and other filters from the quals
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}

	client := connect(ctx, d)

	for {
		var issues []*github.PullRequest
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			issues, resp, err = client.PullRequests.List(ctx, owner, repo, opt)

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

func tableGitHubPullRequestGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var issueNumber int

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	if h.Item != nil {
		issue := h.Item.(*github.PullRequest)
		issueNumber = *issue.Number

	} else {
		issueNumber = int(d.KeyColumnQuals["issue_number"].GetInt64Value())
	}

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubPullRequestGet", "owner", owner, "repo", repo, "issueNumber", issueNumber)

	client := connect(ctx, d)

	var detail *github.PullRequest
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error

		detail, resp, err = client.PullRequests.Get(ctx, owner, repo, issueNumber)
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

func RepoNameFromQuals(_ context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return d.KeyColumnQuals["repository_full_name"].GetStringValue(), nil
}

//// TRANSFORM FUNCTIONS

func getPullRequestTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	issue := d.HydrateItem.(*github.PullRequest)

	tags := make(map[string]bool)
	if issue.Labels != nil {
		for _, i := range issue.Labels {
			tags[*i.Name] = true
		}
	}
	return tags, nil
}
