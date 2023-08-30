package github

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns this plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-github",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"github_actions_artifact":                tableGitHubActionsArtifact(),
			"github_actions_repository_runner":       tableGitHubActionsRepositoryRunner(),
			"github_actions_repository_secret":       tableGitHubActionsRepositorySecret(),
			"github_actions_repository_workflow_run": tableGitHubActionsRepositoryWorkflowRun(),
			"github_audit_log":                       tableGitHubAuditLog(),
			"github_branch_protection":               tableGitHubBranchProtection(),
			"github_branch":                          tableGitHubBranch(),
			"github_commit":                          tableGitHubCommit(),
			"github_community_profile":               tableGitHubCommunityProfile(),
			"github_code_owner":                      tableGitHubCodeOwner(),
			"github_gist":                            tableGitHubGist(),
			"github_gitignore":                       tableGitHubGitignore(),
			"github_issue":                           tableGitHubIssue(),
			"github_issue_comment":                   tableGitHubIssueComment(),
			"github_license":                         tableGitHubLicense(),
			"github_my_gist":                         tableGitHubMyGist(),
			"github_my_issue":                        tableGitHubMyIssue(),
			"github_my_organization":                 tableGitHubMyOrganization(),
			"github_my_repository":                   tableGitHubMyRepository(),
			"github_my_star":                         tableGitHubMyStar(),
			"github_my_team":                         tableGitHubMyTeam(),
			"github_organization":                    tableGitHubOrganization(),
			"github_organization_member":             tableGitHubOrganizationMember(),
			"github_organization_dependabot_alert":   tableGitHubOrganizationDependabotAlert(),
			"github_organization_external_identity":  tableGitHubOrganizationExternalIdentity(),
			"github_pull_request":                    tableGitHubPullRequest(),
			"github_pull_request_comment":            tableGitHubPullRequestComment(),
			"github_pull_request_review":             tableGitHubPullRequestReview(),
			"github_rate_limit":                      tableGitHubRateLimit(),
			"github_rate_limit_graphql":              tableGitHubRateLimitGraphQL(),
			"github_release":                         tableGitHubRelease(),
			"github_repository":                      tableGitHubRepository(),
			"github_repository_collaborator":         tableGitHubRepositoryCollaborator(),
			"github_repository_dependabot_alert":     tableGitHubRepositoryDependabotAlert(),
			"github_repository_deployment":           tableGitHubRepositoryDeployment(),
			"github_repository_environment":          tableGitHubRepositoryEnvironment(),
			"github_repository_vulnerability_alert":  tableGitHubRepositoryVulnerabilityAlert(),
			"github_search_code":                     tableGitHubSearchCode(),
			"github_search_commit":                   tableGitHubSearchCommit(),
			"github_search_issue":                    tableGitHubSearchIssue(),
			"github_search_label":                    tableGitHubSearchLable(),
			"github_search_pull_request":             tableGitHubSearchPullRequest(),
			"github_search_repository":               tableGitHubSearchRepository(),
			"github_search_topic":                    tableGitHubSearchTopic(),
			"github_search_user":                     tableGitHubSearchUser(),
			"github_stargazer":                       tableGitHubStargazer(),
			"github_tag":                             tableGitHubTag(),
			"github_team_member":                     tableGitHubTeamMember(),
			"github_team_repository":                 tableGitHubTeamRepository(),
			"github_team":                            tableGitHubTeam(),
			"github_traffic_view_daily":              tableGitHubTrafficViewDaily(),
			"github_traffic_view_weekly":             tableGitHubTrafficViewWeekly(),
			"github_tree":                            tableGitHubTree(),
			"github_user":                            tableGitHubUser(),
			"github_workflow":                        tableGitHubWorkflow(),
		},
	}
	return p
}
