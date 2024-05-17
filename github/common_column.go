package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(c []*plugin.Column) []*plugin.Column {
	return append([]*plugin.Column{
		{
			Name:        "login_id",
			Description: "Unique identifier for the user login.",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getLoginId,
			Transform:   transform.FromValue(),
		},
	}, c...)
}

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getLoginIdMemoized = plugin.HydrateFunc(getLoginIdUncached).Memoize(memoize.WithCacheKeyFunction(getLoginIdCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getLoginId(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getLoginIdMemoized(ctx, d, h)
}

// Build a cache key for the call to getLoginIdCacheKey.
func getLoginIdCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getLoginId"
	return key, nil
}

func getLoginIdUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	client := connectV4(ctx, d)

	var query struct {
		Viewer struct {
			Login githubv4.String
			ID    githubv4.ID
		}
	}
	err := client.Query(ctx, &query, nil)
	if err != nil {
		plugin.Logger(ctx).Error("getLoginIdUncached", "api_error", err)
		return nil, err
	}

	return query.Viewer.ID, nil
}
