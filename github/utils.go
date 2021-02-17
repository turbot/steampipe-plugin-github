package github

import (
	"context"
	"os"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/plugin"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"

	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// create service client
func connect(ctx context.Context, d *plugin.QueryData) *github.Client {
	logger := plugin.Logger(ctx)

	// Get connection config for plugin
	githubConfig := GetConfig(d.Connection)
	if &githubConfig != nil {
		if githubConfig.Token != nil {
			os.Setenv("GITHUB_TOKEN", *githubConfig.Token)
		}
	}

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
