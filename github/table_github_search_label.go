package github

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/go-github/v55/github"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubSearchLabel() *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_label",
		Description: "Find labels in a repository with names or descriptions that match search keywords.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.AllColumns([]string{"query", "repository_id"}),
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubSearchLabelList,
		},
		Columns: []*plugin.Column{
			{Name: "id", Transform: transform.FromField("ID"), Type: proto.ColumnType_INT, Description: "The ID of the label."},
			{Name: "repository_id", Type: proto.ColumnType_INT, Transform: transform.FromQual("repository_id"), Description: "The ID of the repository."},
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.From(extractSearchLabelRepositoryFullName), Description: "The full name of the repository (login/repo-name)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the label."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the label."},
			{Name: "color", Type: proto.ColumnType_STRING, Description: "The color assigned to the label."},
			{Name: "default", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the label is a default one."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the label."},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "The score of the label."},
			{Name: "url", Type: proto.ColumnType_STRING, Description: "The API URL to get the label details."},
			{Name: "text_matches", Type: proto.ColumnType_JSON, Description: "The text match details."},
		},
	}
}

func tableGitHubSearchLabelList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchLabelList")

	quals := d.EqualsQuals
	repoId := d.EqualsQuals["repository_id"].GetInt64Value()
	query := quals["query"].GetStringValue()

	if query == "" {
		return nil, nil
	}

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		TextMatch:   true,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opt.ListOptions.PerPage) {
			opt.ListOptions.PerPage = int(*limit)
		}
	}

	client := connect(ctx, d)
	for {
		result, resp, err := client.Search.Labels(ctx, repoId, query, opt)
		if err != nil {
			logger.Error("tableGitHubSearchLabelList", "error_RetryHydrate", err)
			return nil, err
		}

		labels := result.Labels
		for _, i := range labels {
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

func extractSearchLabelRepositoryFullName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	label := d.HydrateItem.(*github.LabelResult)
	if label.URL != nil {
		rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta("repos/") + `(.*?)` + regexp.QuoteMeta("/labels"))
		replacer := strings.NewReplacer("repos/", "", "/labels", "")
		return replacer.Replace(rx.FindString(*label.URL)), nil
	}
	return "", nil
}
