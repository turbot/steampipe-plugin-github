package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubCommit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_commit",
		Description: "Github Commits bundle project files for download by users.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubCommitList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "sha"}),
			Hydrate:    tableGitHubCommitGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repositoryFullNameQual, Transform: transform.FromValue(), Description: "Full name of the repository that contains the commit."},
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA"), Description: "SHA of the commit."},
			// Other columns
			{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.Login"), Description: "The login name of the author of the commit."},
			{Name: "author_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Commit.Author.Date"), Description: "Timestamp when the author made this commit."},
			{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "Comments URL of the commit."},
			{Name: "commit", Type: proto.ColumnType_JSON, Description: "Commit details."},
			{Name: "committer_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Committer.Login"), Description: "The login name of committer of the commit."},
			{Name: "committer_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Commit.Committer.Date"), Description: "Timestamp when the committer made this commit."},
			{Name: "files", Type: proto.ColumnType_JSON, Hydrate: tableGitHubCommitGet, Description: "Files of the commit."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "HTML URL of the commit."},
			{Name: "message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Commit.Message"), Description: "Commit message."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Description: "Node where GitHub stores this data internally."},
			{Name: "parents", Type: proto.ColumnType_JSON, Description: "Parent commits of the commit."},
			{Name: "stats", Type: proto.ColumnType_JSON, Hydrate: tableGitHubCommitGet, Description: "Statistics of the commit."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "URL of the commit."},
			{Name: "verified", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Commit.Verification.Verified"), Description: "True if the commit was verified with a signature."},
		},
	}
}

func tableGitHubCommitList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	s := strings.Split(fullName, "/")
	owner := s[0]
	repo := s[1]

	opts := &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {

		var commits []*github.RepositoryCommit
		var resp *github.Response

		b, err := retry.NewFibonacci(100 * time.Millisecond)
		if err != nil {
			return nil, err
		}

		err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
			var err error
			commits, resp, err = client.Repositories.ListCommits(ctx, owner, repo, opts)
			if _, ok := err.(*github.RateLimitError); ok {
				return retry.RetryableError(err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		for _, i := range commits {
			d.StreamListItem(ctx, i)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

func tableGitHubCommitGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var owner, repo string
	var sha string

	logger := plugin.Logger(ctx)
	quals := d.KeyColumnQuals

	if h.Item != nil {
		commit := h.Item.(*github.RepositoryCommit)
		sha = *commit.SHA
	} else {
		sha = d.KeyColumnQuals["sha"].GetStringValue()
	}

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubCommitGet", "owner", owner, "repo", repo, "sha", sha)

	client := connect(ctx, d)

	var detail *github.RepositoryCommit
	var resp *github.Response

	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return detail, err
	}

	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		detail, resp, err = client.Repositories.GetCommit(ctx, owner, repo, sha)
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
