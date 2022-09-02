package github

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v47/github"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubCommit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_commit",
		Description: "GitHub Commits bundle project files for download by users.",
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "sha", Require: plugin.Optional},
				{Name: "author_login", Require: plugin.Optional},
				{Name: "author_date", Require: plugin.Optional, Operators: []string{">", ">=", "=", "<", "<="}},
			},
			Hydrate: tableGitHubCommitList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "sha"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubCommitGet,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the commit."},
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA"), Description: "SHA of the commit."},
			// Other columns
			{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.Login"), Description: "The login name of the author of the commit."},
			{Name: "author_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Commit.Author.Date"), Description: "Timestamp when the author made this commit."},
			{Name: "verified", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Commit.Verification.Verified"), Description: "True if the commit was verified with a signature."},
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
		},
	}
}

//// LIST FUNCTION

func tableGitHubCommitList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	opts := &github.CommitsListOptions{ListOptions: github.ListOptions{PerPage: 100}}

	type ListPageResponse struct {
		commits []*github.RepositoryCommit
		resp    *github.Response
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		commits, resp, err := client.Repositories.ListCommits(ctx, owner, repo, opts)
		return ListPageResponse{
			commits: commits,
			resp:    resp,
		}, err
	}

	// Additional filters
	if d.KeyColumnQuals["sha"] != nil {
		opts.SHA = d.KeyColumnQuals["sha"].GetStringValue()
	}

	if d.KeyColumnQuals["author_login"] != nil {
		opts.Author = d.KeyColumnQuals["author_login"].GetStringValue()
	}

	if d.Quals["author_date"] != nil {
		for _, q := range d.Quals["author_date"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			beforeTime := givenTime.Add(time.Duration(-1) * time.Second)
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				opts.Since = afterTime
			case ">=":
				opts.Since = givenTime
			case "=":
				opts.Since = givenTime
				opts.Until = givenTime
			case "<=":
				opts.Until = givenTime
			case "<":
				opts.Until = beforeTime
			}
		}
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.ListOptions.PerPage) {
			opts.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			// Gets this error if repository is not initialized
			if strings.Contains(err.Error(), "409 Git Repository is empty.") {
				return nil, nil
			}
		}

		listResponse := listPageResponse.(ListPageResponse)
		commits := listResponse.commits
		resp := listResponse.resp

		for _, i := range commits {
			if i != nil {
				d.StreamListItem(ctx, i)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

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

	// Return nil, if no input provided
	if fullName == "" || sha == "" {
		return nil, nil
	}

	owner, repo = parseRepoFullName(fullName)
	logger.Trace("tableGitHubCommitGet", "owner", owner, "repo", repo, "sha", sha)

	client := connect(ctx, d)

	opts := &github.ListOptions{
		PerPage: 100,
	}

	type GetResponse struct {
		commit *github.RepositoryCommit
		resp   *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		detail, resp, err := client.Repositories.GetCommit(ctx, owner, repo, sha, opts)
		return GetResponse{
			commit: detail,
			resp:   resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	commit := getResp.commit

	return commit, nil
}
