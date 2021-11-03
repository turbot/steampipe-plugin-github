package github

import (
	"context"

	"github.com/google/go-github/v33/github"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubCommunityProfile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_community_profile",
		Description: "Community profile information for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("repository_full_name"),
			Hydrate:    tableGitHubCommunityProfileList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tag."},
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

//// LIST FUNCTION

func tableGitHubCommunityProfileList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	type GetResponse struct {
		result *github.CommunityHealthMetrics
		resp   *github.Response
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		details, resp, err := client.Repositories.GetCommunityHealthMetrics(ctx, owner, repo)
		return GetResponse{
			result: details,
			resp:   resp,
		}, err
	}

	getResponse, err := plugin.RetryHydrate(ctx, d, h, getDetails, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})

	if err != nil {
		return nil, err
	}

	getResp := getResponse.(GetResponse)
	result := getResp.result

	d.StreamListItem(ctx, result)
	return nil, nil
}
