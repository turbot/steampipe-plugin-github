package github

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubSearchCommit(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_commit",
		Description: "Find commits via various criteria on the default branch (usually master).",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchCommitList,
		},
		Columns: []*plugin.Column{
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA"), Description: "The SHA of the commit."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the commit."},
			{Name: "comments_url", Type: proto.ColumnType_STRING, Description: "The API URL of the comments made on the commit."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The Github URL of the commit."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("URL"), Description: "The API URL of the commit."},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "The score of the commit."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.From(extractSearchCommitRepositoryFullName), Description: "The full name of the repository (login/repo-name)."},
			{Name: "author", Type: proto.ColumnType_JSON, Description: "The author details."},
			{Name: "commit", Type: proto.ColumnType_JSON, Description: "The commit details."},
			{Name: "committer", Type: proto.ColumnType_JSON, Description: "The committer details."},
			{Name: "parents", Type: proto.ColumnType_JSON, Description: "The parent details."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The repository details of the commit."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubSearchCommitList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchCommitList")

	quals := d.EqualsQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.CommitsSearchResult
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
		result, resp, err := client.Search.Commits(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchCommitList", "error_Search.Commits", err)
			return nil, err
		}

		return ListPageResponse{
			result: result,
			resp:   resp,
		}, nil
	}

	for {
		listPageResponse, err := retryHydrate(ctx, d, h, listPage)

		if err != nil {
			logger.Error("tableGitHubSearchCommitList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		codeResults := listResponse.result.Commits
		resp := listResponse.resp

		for _, i := range codeResults {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
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

func extractSearchCommitRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	commit := d.HydrateItem.(*github.CommitResult)
	if commit.URL != nil {
		rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta("repos/") + `(.*?)` + regexp.QuoteMeta("/commits"))
		replacer := strings.NewReplacer("repos/", "", "/commits", "")
		return replacer.Replace(rx.FindString(*commit.URL)), nil
	}
	return "", nil
}
