package github

import (
	"context"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubCommunityProfile() *plugin.Table {
	return &plugin.Table{
		Name:        "github_community_profile",
		Description: "Community profile information for the given repository.",
		List: &plugin.ListConfig{
			KeyColumns:        plugin.SingleColumn("repository_full_name"),
			Hydrate:           tableGitHubCommunityProfileList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the tag."},
			{Name: "code_of_conduct", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateCodeOfConduct, Description: "Code of conduct for the repository."},
			{Name: "contributing", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateContributing, Description: "Contributing guidelines for the repository."},
			{Name: "issue_templates", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateIssueTemplates, Description: "Issue template for the repository."},
			{Name: "pull_request_templates", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydratePullRequestTemplates, Description: "Pull request template for the repository."},
			{Name: "license_info", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateLicense, Description: "License for the repository."},
			{Name: "readme", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateReadme, Description: "README for the repository."},
			{Name: "security", Type: proto.ColumnType_JSON, Transform: transform.FromValue().NullIfZero(), Hydrate: cpHydrateSecurity, Description: "Security for the repository."},
		},
	}
}

func tableGitHubCommunityProfileList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			models.CommunityProfile
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	}
	appendCommunityProfileColumnIncludes(&variables, d.QueryContext.Columns)

	client := connectV4(ctx, d)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_community_profile", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_community_profile", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, query.Repository.CommunityProfile)

	return nil, nil
}
