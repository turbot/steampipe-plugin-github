package github

import (
	"context"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/sethvargo/go-retry"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableGitHubCommunityProfile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_community_profile",
		Description: "Community profile information for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubCommunityProfileList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Hydrate: repositoryFullNameQual, Transform: transform.FromValue(), Description: "Full name of the repository that contains the tag."},
			{Name: "health_percentage", Type: proto.ColumnType_INT, Description: "Community profile health as a percentage metric."},
			{Name: "code_of_conduct", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.CodeOfConduct"), Description: "Code of conduct for the repository."},
			{Name: "contributing", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.Contributing"), Description: "Contributing guidelines for the repository."},
			{Name: "issue_template", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.IssueTemplate"), Description: "Issue template for the repository."},
			{Name: "pull_request_template", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.PullRequestTemplate"), Description: "Pull request template for the repository."},
			{Name: "license", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.License"), Description: "License for the repository."},
			{Name: "readme", Type: proto.ColumnType_JSON, Transform: transform.FromField("Files.Readme"), Description: "README for the repository."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Time when the community profile was last updated."},
		},
	}
}

func tableGitHubCommunityProfileList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)
	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	var result *github.CommunityHealthMetrics
	b, err := retry.NewFibonacci(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}
	err = retry.Do(ctx, retry.WithMaxRetries(10, b), func(ctx context.Context) error {
		var err error
		result, _, err = client.Repositories.GetCommunityHealthMetrics(ctx, owner, repo)
		if _, ok := err.(*github.RateLimitError); ok {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, result)
	return nil, nil
}
