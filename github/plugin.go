package github

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// Plugin returns this plugin
func Plugin(context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-github",
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"github_gist":         tableGitHubGist(),
			"github_license":      tableGitHubLicense(),
			"github_organization": tableGitHubOrganization(),
			"github_repository":   tableGitHubRepository(),
			"github_team":         tableGitHubTeam(),
			"github_user":         tableGitHubUser(),
		},
	}
	return p
}
