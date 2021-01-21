package github

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/connection"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"os"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// create service client
func connect(ctx context.Context, _ *connection.Manager) *github.Client {
	logger := plugin.Logger(ctx)
	logger.Trace("G", os.Getenv("GITHUB_TOKEN"))
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

// transforms

func convertTimestamp(_ context.Context, input *transform.TransformData) (interface{}, error) {
	return input.Value.(*github.Timestamp).Format(time.RFC3339), nil
}
