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

func tableGitHubSearchCode(ctx context.Context) *plugin.Table {
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

//// LIST FUNCTION

func tableGitHubSearchCodeList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchCodeList")

	quals := d.KeyColumnQuals
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	type ListPageResponse struct {
		result *github.CodeSearchResult
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
		result, resp, err := client.Search.Code(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchCodeList", "error_Search.Code", err)
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
			logger.Error("tableGitHubSearchCodeList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		codeResults := listResponse.result.CodeResults
		resp := listResponse.resp

		for _, i := range codeResults {
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

func extractSearchCodeRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	code := d.HydrateItem.(*github.CodeResult)
	if code.HTMLURL != nil {
		rx := regexp.MustCompile(`https?://.+?/(.+?)/blob`)
		return rx.FindStringSubmatch(*code.HTMLURL)[1], nil
	}
	return "", nil
}