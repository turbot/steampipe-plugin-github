package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v45/github"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubCommunityProfile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_community_profile",
		Description: "Community profile information for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			Hydrate:           tableGitHubCommunityProfileList,
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
			{Name: "security", Type: proto.ColumnType_JSON, Transform: transform.FromValue(), Description: "Security for the repository.", Hydrate: securityFileGet},
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

func securityFileGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connect(ctx, d)

	fullName := d.KeyColumnQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	optionalSecurityMDNames := []string{"SECURITY.md", "security.md", "Security.md"}

	for _, filePath := range optionalSecurityMDNames {
		fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, &github.RepositoryContentGetOptions{})

		if err != nil {
			// Gets this error of repository is not initialized
			if strings.Contains(err.Error(), "404 This repository is empty.") {
				return nil, nil
			}
		}

		if fileContent != nil {
			return fileContent, nil
		} else if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			plugin.Logger(ctx).Error("github_community_profile.securityFileGet", "api_error", err)
			return nil, err
		}
	}
	return nil, nil
}
