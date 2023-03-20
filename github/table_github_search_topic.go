package github

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubSearchTopic(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_search_topic",
		Description: "Find topics via various criteria.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("query"),
			Hydrate:    tableGitHubSearchTopicList,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the topic."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "The query used to match the topic."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt").Transform(convertTimestamp), Description: "The timestamp when the topic was created."},
			{Name: "created_by", Type: proto.ColumnType_STRING, Description: "The creator of the topic."},
			{Name: "curated", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the topic is curated."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The description of the topic."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Description: "The display name of the topic."},
			{Name: "featured", Type: proto.ColumnType_BOOL, Default: false, Description: "Whether the topic is featured."},
			{Name: "score", Type: proto.ColumnType_DOUBLE, Description: "The score of the topic."},
			{Name: "short_description", Type: proto.ColumnType_STRING, Description: "The short description of the topic."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the topic was updated."},
		},
	}
}

//// LIST FUNCTION

func tableGitHubSearchTopicList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("tableGitHubSearchTopicList")

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
		result *github.TopicsSearchResult
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
		result, resp, err := client.Search.Topics(ctx, query, opt)

		if err != nil {
			logger.Error("tableGitHubSearchTopicList", "error_Search.Topics", err)
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
			logger.Error("tableGitHubSearchTopicList", "error_RetryHydrate", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		topics := listResponse.result.Topics
		resp := listResponse.resp

		for _, i := range topics {
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
