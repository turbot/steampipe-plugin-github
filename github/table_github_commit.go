package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableGitHubCommit() *plugin.Table {
	return &plugin.Table{
		Name:        "github_commit",
		Description: "GitHub Commits bundle project files for download by users.",
		List: &plugin.ListConfig{
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
				{Name: "authored_date", Require: plugin.Optional, Operators: []string{">", ">=", "=", "<", "<="}},
			},
			Hydrate: tableGitHubCommitList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"repository_full_name", "sha"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubCommitGet,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the commit."},
			{Name: "sha", Type: proto.ColumnType_STRING, Description: "SHA of the commit."},
			{Name: "short_sha", Type: proto.ColumnType_STRING, Description: "Short SHA of the commit."},
			{Name: "message", Type: proto.ColumnType_STRING, Description: "Commit message."},
			{Name: "author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Author.User.Login"), Description: "The login name of the author of the commit."},
			{Name: "authored_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("AuthoredDate").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the author made this commit."},
			{Name: "author", Type: proto.ColumnType_JSON, Transform: transform.FromField("Author").NullIfZero(), Description: "The commit author."},
			{Name: "committer_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Committer.User.Login"), Description: "The login name of the committer."},
			{Name: "committed_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CommittedDate").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when commit was committed."},
			{Name: "committer", Type: proto.ColumnType_JSON, Transform: transform.FromField("Committer").NullIfZero(), Description: "The committer."},
			{Name: "additions", Type: proto.ColumnType_INT, Description: "Number of additions in the commit."},
			{Name: "deletions", Type: proto.ColumnType_INT, Description: "Number of deletions in the commit."},
			{Name: "changed_files", Type: proto.ColumnType_INT, Description: "Count of files changed in the commit."},
			{Name: "committed_via_web", Type: proto.ColumnType_BOOL, Description: "If true, commit was made via GitHub web ui."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("CommitUrl"), Description: "URL of the commit."},
			{Name: "signature", Type: proto.ColumnType_JSON, Transform: transform.FromField("Signature").NullIfZero(), Description: "The signature of commit."},
			{Name: "status", Type: proto.ColumnType_JSON, Transform: transform.FromField("Status").NullIfZero(), Description: "Status of the commit."},
			{Name: "tarball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TarballUrl"), Description: "URL to download a tar of commit."},
			{Name: "zipball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZipballUrl"), Description: "URL to download a zip of commit."},
			{Name: "tree_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("TreeUrl"), Description: "URL to tree of the commit."},
			{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Description: "If true, user can subscribe to this commit."},
			{Name: "subscription", Type: proto.ColumnType_STRING, Description: "Users subscription state of the commit."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url"), Description: "URL of the commit."},
			{Name: "node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("NodeId"), Description: "The node ID of the commit."},
		},
	}
}

func tableGitHubCommitList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)
	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"name":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
		"since":    (*githubv4.GitTimestamp)(nil),
		"until":    (*githubv4.GitTimestamp)(nil),
	}

	if d.Quals["authored_date"] != nil {
		for _, q := range d.Quals["authored_date"].Quals {
			givenTime := q.Value.GetTimestampValue().AsTime()
			beforeTime := givenTime.Add(time.Duration(-1) * time.Second)
			afterTime := givenTime.Add(time.Second * 1)

			switch q.Operator {
			case ">":
				variables["since"] = githubv4.GitTimestamp{Time: afterTime}
			case ">=":
				variables["since"] = githubv4.GitTimestamp{Time: givenTime}
			case "=":
				variables["since"] = githubv4.GitTimestamp{Time: givenTime}
				variables["until"] = githubv4.GitTimestamp{Time: givenTime}
			case "<=":
				variables["until"] = githubv4.GitTimestamp{Time: givenTime}
			case "<":
				variables["until"] = githubv4.GitTimestamp{Time: beforeTime}
			}
		}
	}

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			DefaultBranchRef struct {
				Target struct {
					Commit struct {
						History struct {
							TotalCount int
							PageInfo   models.PageInfo
							Nodes      []models.Commit
						} `graphql:"history(first: $pageSize, after: $cursor, since: $since, until: $until)"`
					} `graphql:"... on Commit"`
				}
			}
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	client := connectV4(ctx, d)
	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	for {
		_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
		plugin.Logger(ctx).Debug(rateLimitLogString("github_commit", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_commit", "api_error", err)
			return nil, err
		}

		for _, commit := range query.Repository.DefaultBranchRef.Target.Commit.History.Nodes {
			d.StreamListItem(ctx, commit)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.DefaultBranchRef.Target.Commit.History.PageInfo.EndCursor)
	}

	return nil, nil
}

func tableGitHubCommitGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	sha := quals["sha"].GetStringValue()

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Object struct {
				Commit models.Commit `graphql:"... on Commit"`
			} `graphql:"object(oid: $sha)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
		"sha":   githubv4.GitObjectID(sha),
	}

	client := connectV4(ctx, d)

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		return nil, client.Query(ctx, &query, variables)
	}

	_, err := plugin.RetryHydrate(ctx, d, h, listPage, retryConfig())
	plugin.Logger(ctx).Debug(rateLimitLogString("github_commit", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_commit", "api_error", err)
		return nil, err
	}

	return query.Repository.Object.Commit, nil
}
