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
			{Name: "code_of_conduct", Type: proto.ColumnType_JSON, Transform: transform.FromField("CodeOfConduct").NullIfZero(), Description: "Code of conduct for the repository."},
			{Name: "contributing", Type: proto.ColumnType_JSON, Transform: transform.FromField("ContributingUpper.Blob", "ContributingLower.Blob", "ContributingTitle.Blob").NullIfZero(), Description: "Contributing guidelines for the repository."},
			{Name: "issue_templates", Type: proto.ColumnType_JSON, Transform: transform.FromField("IssueTemplates").NullIfZero(), Description: "Issue template for the repository."},
			{Name: "pull_request_templates", Type: proto.ColumnType_JSON, Transform: transform.FromField("PullRequestTemplates").NullIfZero(), Description: "Pull request template for the repository."},
			{Name: "license_info", Type: proto.ColumnType_JSON, Transform: transform.FromField("LicenseInfo").NullIfZero(), Description: "License for the repository."},
			{Name: "readme", Type: proto.ColumnType_JSON, Transform: transform.FromField("ReadMeUpper.Blob", "ReadMeLower.Blob", "ReadMeTitle.Blob"), Description: "README for the repository."},
			{Name: "security", Type: proto.ColumnType_JSON, Transform: transform.FromField("SecurityUpper.Blob", "SecurityLower.Blob", "SecurityTitle.Blob"), Description: "Security for the repository."},
		},
	}
}

func tableGitHubCommunityProfileList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			LicenseInfo          models.License
			CodeOfConduct        models.RepositoryCodeOfConduct
			IssueTemplates       []models.IssueTemplate
			PullRequestTemplates []models.PullRequestTemplate
			// readme
			ReadMeLower struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"readMeLower: object(expression: \"HEAD:readme.md\")"`
			ReadMeUpper struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"readMeUpper: object(expression: \"HEAD:README.md\")"`
			// contributing
			ContributingLower struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"contributingLower: object(expression: \"HEAD:contributing.md\")"`
			ContributingTitle struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"contributingTitle: object(expression: \"HEAD:Contributing.md\")"`
			ContributingUpper struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"contributingUpper: object(expression: \"HEAD:CONTRIBUTING.md\")"`
			// security
			SecurityLower struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"securityLower: object(expression: \"HEAD:security.md\")"`
			SecurityTitle struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"securityTitle: object(expression: \"HEAD:Security.md\")"`
			SecurityUpper struct {
				Blob models.Blob `graphql:"... on Blob"`
			} `graphql:"securityUpper: object(expression: \"HEAD:SECURITY.md\")"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	}

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_community_profile", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_community_profile", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, query.Repository)

	return nil, nil
}
