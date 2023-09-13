package github

import (
	"context"
	"regexp"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubSearchCode() *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_code",
		Description: "Searches for query terms inside of a file.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchCodeList,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the file where the match has been found."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the code."},
			{Name: "html_url", Type: proto.ColumnType_STRING, Description: "The complete URL of the file where the match has been found."},
			{Name: "sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("SHA"), Description: "The SHA of the file where the match has been found."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path of the file where the match has been found."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.From(extractSearchCodeRepositoryFullName), Description: "The full name of the repository (login/repo-name)."},
			{Name: "repository", Type: proto.ColumnType_JSON, Description: "The repository details of the file where the match has been found."},
			{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
		},
	}
}

func tableGitHubSearchCodeList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchCodeList")

	quals := d.EqualsQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	client := connect(ctx, d)

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	for {
		result, resp, err := client.Search.Code(ctx, query, opt)
		if err != nil {
			logger.Error("tableGitHubSearchCodeList", "error_RetryHydrate", err)
			return nil, err
		}

		codeResults := result.CodeResults
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

func extractSearchCodeRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	code := d.HydrateItem.(*github.CodeResult)
	if code.HTMLURL != nil {
		rx := regexp.MustCompile(`https?://.+?/(.+?)/blob`)
		return rx.FindStringSubmatch(*code.HTMLURL)[1], nil
	}
	return "", nil
}
