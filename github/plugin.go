package github

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			"github_branch":              tableGitHubBranch(ctx),
			"github_branch_protection":   tableGitHubBranchProtection(ctx),
			"github_commit":              tableGitHubCommit(ctx),
			"github_gist":                tableGitHubGist(),
			"github_gitignore":           tableGitHubGitignore(),
			"github_license":             tableGitHubLicense(),
			"github_issue":               tableGitHubIssue(),
			"github_my_gist":             tableGitHubMyGist(),
			"github_my_issue":            tableGitHubMyIssue(),
			"github_my_organization":     tableGitHubMyOrganization(),
			"github_my_repository":       tableGitHubMyRepository(),
			"github_my_team":             tableGitHubMyTeam(),
			"github_organization":        tableGitHubOrganization(),
			"github_pull_request":        tableGitHubPullRequest(),
			"github_release":             tableGitHubRelease(ctx),
			"github_repository":          tableGitHubRepository(),
			"github_stargazer":           tableGitHubStargazer(ctx),
			"github_tag":                 tableGitHubTag(ctx),
			"github_traffic_view_daily":  tableGitHubTrafficViewDaily(ctx),
			"github_traffic_view_weekly": tableGitHubTrafficViewWeekly(ctx),
			"github_user":                tableGitHubUser(),
			"github_workflow":            tableGitHubWorkflow(ctx),
		},
	}
	return p
}
