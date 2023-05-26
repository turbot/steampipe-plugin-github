package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func gitHubRepositoryColumns() []*plugin.Column {
	repoColumns := []*plugin.Column{
		{Name: "full_name", Type: proto.ColumnType_STRING, Description: "The full name of the repository, including the owner and repo name.", Transform: transform.FromQual("full_name")},
	}

	return append(repoColumns, sharedRepositoryColumns()...)
}

func sharedRepositoryColumns() []*plugin.Column {
	return []*plugin.Column{
		{Name: "id", Type: proto.ColumnType_INT, Description: "The numeric ID of the repository.", Transform: transform.FromField("Id", "Node.Id")},
		{Name: "node_id", Type: proto.ColumnType_STRING, Description: "The node ID of the repository.", Transform: transform.FromField("NodeId", "Node.NodeId")},
		{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the repository.", Transform: transform.FromField("Name", "Node.Name")},
		{Name: "allow_update_branch", Type: proto.ColumnType_BOOL, Description: "If true, a pull request head branch that is behind its base branch can always be updated even if it is not required to be up to date before merging.", Transform: transform.FromField("AllowUpdateBranch", "Node.AllowUpdateBranch")},
		{Name: "archived_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when repository was archived.", Transform: transform.FromField("ArchivedAt", "Node.ArchivedAt").NullIfZero().Transform(convertTimestamp)},
		{Name: "auto_merge_allowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("AutoMergeAllowed", "Node.AutoMergeAllowed"), Description: "If true, auto-merge can be enabled on pull requests in this repository."},
		{Name: "code_of_conduct", Type: proto.ColumnType_JSON, Transform: transform.FromField("CodeOfConduct", "Node.CodeOfConduct").NullIfZero(), Description: "The code of conduct for this repository."},
		{Name: "contact_links", Type: proto.ColumnType_JSON, Transform: transform.FromField("ContactLinks", "Node.ContactLinks").NullIfZero(), Description: "List of contact links associated to the repository."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("CreatedAt", "Node.CreatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the repository was created."},
		{Name: "default_branch_ref", Type: proto.ColumnType_JSON, Transform: transform.FromField("DefaultBranchRef", "Node.DefaultBranchRef").NullIfZero(), Description: "Default ref information."},
		{Name: "delete_branch_on_merge", Type: proto.ColumnType_BOOL, Transform: transform.FromField("DeleteBranchOnMerge", "Node.DeleteBranchOnMerge"), Description: "If true, branches are automatically deleted when merged in this repository."},
		{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Description", "Node.Description"), Description: "The description of the repository."},
		{Name: "disk_usage", Type: proto.ColumnType_INT, Transform: transform.FromField("DiskUsage", "Node.DiskUsage"), Description: "Number of kilobytes this repository occupies on disk."},
		{Name: "fork_count", Type: proto.ColumnType_INT, Transform: transform.FromField("ForkCount", "Node.ForkCount"), Description: "Number of forks there are of this repository in the whole network."},
		{Name: "forking_allowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ForkingAllowed", "Node.ForkingAllowed"), Description: "If true, repository allows forks."},
		{Name: "funding_links", Type: proto.ColumnType_JSON, Transform: transform.FromField("FundingLinks", "Node.FundingLinks").NullIfZero(), Description: "The funding links for this repository."},
		{Name: "has_discussions_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasDiscussionsEnabled", "Node.HasDiscussionsEnabled"), Description: "If true, the repository has the Discussions feature enabled."},
		{Name: "has_issues_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasIssuesEnabled", "Node.HasIssuesEnabled"), Description: "If true, the repository has issues feature enabled."},
		{Name: "has_projects_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasProjectsEnabled", "Node.HasProjectsEnabled"), Description: "If true, the repository has the Projects feature enabled."},
		{Name: "has_vulnerability_alerts_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasVulnerabilityAlertsEnabled", "Node.HasVulnerabilityAlertsEnabled"), Description: "If true, vulnerability alerts are enabled for the repository."},
		{Name: "has_wiki_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasWikiEnabled", "Node.HasWikiEnabled"), Description: "If true, the repository has wiki feature enabled."},
		{Name: "homepage_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("HomepageUrl", "Node.HomepageUrl"), Description: "The external URL of the repository if set."},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Transform: transform.FromField("InteractionAbility", "Node.InteractionAbility").NullIfZero(), Description: "The interaction ability settings for this repository."},
		{Name: "is_archived", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsArchived", "Node.IsArchived"), Description: "If true, the repository is unmaintained (archived)."},
		{Name: "is_blank_issues_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsBlankIssuesEnabled", "Node.IsBlankIssuesEnabled"), Description: "If true, blank issue creation is allowed."},
		{Name: "is_disabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsDisabled", "Node.IsDisabled"), Description: "If true, this repository disabled."},
		{Name: "is_empty", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsEmpty", "Node.IsEmpty"), Description: "If true, this repository is empty."},
		{Name: "is_fork", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsFork", "Node.IsFork"), Description: "If true, the repository is a fork."},
		{Name: "is_in_organization", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsInOrganization", "Node.IsInOrganization"), Description: "If true, repository is either owned by an organization, or is a private fork of an organization repository."},
		{Name: "is_locked", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsLocked", "Node.IsLocked"), Description: "If true, repository is locked."},
		{Name: "is_mirror", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsMirror", "Node.IsMirror"), Description: "If true, the repository is a mirror."},
		{Name: "is_private", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsPrivate", "Node.IsPrivate"), Description: "If true, the repository is private or internal."},
		{Name: "is_security_policy_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsSecurityPolicyEnabled", "Node.IsSecurityPolicyEnabled"), Description: "If true, repository has a security policy."},
		{Name: "is_template", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsTemplate", "Node.IsTemplate"), Description: "If true, the repository is a template that can be used to generate new repositories."},
		{Name: "is_user_configuration_repository", Type: proto.ColumnType_BOOL, Transform: transform.FromField("IsUserConfigurationRepository", "Node.IsUserConfigurationRepository"), Description: "If true, this is a user configuration repository."},
		{Name: "issue_templates", Type: proto.ColumnType_JSON, Transform: transform.FromField("IssueTemplates", "Node.IssueTemplates").NullIfZero(), Description: "A list of issue templates associated to the repository."},
		{Name: "license_info", Type: proto.ColumnType_JSON, Transform: transform.FromField("LicenseInfo", "Node.LicenseInfo").NullIfZero(), Description: "The license associated with the repository."},
		{Name: "lock_reason", Type: proto.ColumnType_STRING, Transform: transform.FromField("LockReason", "Node.LockReason"), Description: "The reason the repository has been locked."},
		{Name: "merge_commit_allowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("MergeCommitAllowed", "Node.MergeCommitAllowed"), Description: "If true, PRs are merged with a merge commit on this repository."},
		{Name: "merge_commit_message", Type: proto.ColumnType_STRING, Transform: transform.FromField("MergeCommitMessage", "Node.MergeCommitMessage"), Description: "How the default commit message will be generated when merging a pull request."},
		{Name: "merge_commit_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("MergeCommitTitle", "Node.MergeCommitTitle"), Description: "How the default commit title will be generated when merging a pull request."},
		{Name: "mirror_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("MirrorUrl", "Node.MirrorUrl"), Description: "The repository's original mirror URL."},
		{Name: "name_with_owner", Type: proto.ColumnType_STRING, Transform: transform.FromField("NameWithOwner", "Node.NameWithOwner"), Description: "The repository's name with owner."},
		{Name: "open_graph_image_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("MirrorUrl", "Node.MirrorUrl"), Description: "The image used to represent this repository in Open Graph data."},
		{Name: "owner_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Owner.Login", "Node.Owner.Login"), Description: "Login of the repository owner."},
		{Name: "primary_language", Type: proto.ColumnType_JSON, Transform: transform.FromField("PrimaryLanguage", "Node.PrimaryLanguage").NullIfZero(), Description: "The primary language of the repository's code."},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("ProjectsUrl", "Node.ProjectsUrl"), Description: "The URL listing the repository's projects."},
		{Name: "pull_request_templates", Type: proto.ColumnType_JSON, Transform: transform.FromField("PullRequestTemplates", "Node.PullRequestTemplates"), Description: "Returns a list of pull request templates associated to the repository."},
		{Name: "pushed_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("PushedAt", "Node.PushedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the repository was last pushed to."},
		{Name: "rebase_merge_allowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("RebaseMergeAllowed", "Node.RebaseMergeAllowed"), Description: "If true, rebase-merging is enabled on this repository."},
		{Name: "security_policy_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("SecurityPolicyUrl", "Node.SecurityPolicyUrl"), Description: "The security policy URL."},
		{Name: "squash_merge_allowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("SquashMergeAllowed", "Node.SquashMergeAllowed"), Description: "If true, squash-merging is enabled on this repository."},
		{Name: "squash_merge_commit_message", Type: proto.ColumnType_STRING, Transform: transform.FromField("SquashMergeCommitMessage", "Node.SquashMergeCommitMessage"), Description: "How the default commit message will be generated when squash merging a pull request."},
		{Name: "squash_merge_commit_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("SquashMergeCommitTitle", "Node.SquashMergeCommitTitle"), Description: "How the default commit title will be generated when squash merging a pull request."},
		{Name: "ssh_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("SshUrl", "Node.SshUrl"), Description: "The SSH URL to clone this repository."},
		{Name: "stargazer_count", Type: proto.ColumnType_INT, Transform: transform.FromField("StargazerCount", "Node.StargazerCount"), Description: "Returns a count of how many stargazers there are on this repository."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("UpdatedAt", "Node.UpdatedAt").NullIfZero().Transform(convertTimestamp), Description: "Timestamp when repository was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Url", "Node.Url"), Description: "The URL of the repository."},
		{Name: "uses_custom_open_graph_image", Type: proto.ColumnType_BOOL, Transform: transform.FromField("UsesCustomOpenGraphImage", "Node.UsesCustomOpenGraphImage"), Description: "if true, this repository has a custom image to use with Open Graph as opposed to being represented by the owner's avatar."},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanAdminister", "Node.CanAdminister"), Description: "If true, you can administer this repository."},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanCreateProjects", "Node.CanCreateProjects"), Description: "If true, you can create projects in this repository."},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanSubscribe", "Node.CanSubscribe"), Description: "If true, you can subscribe to this repository."},
		{Name: "can_update_topics", Type: proto.ColumnType_BOOL, Transform: transform.FromField("CanUpdateTopics", "Node.CanUpdateTopics"), Description: "If true, you can update topics on this repository."},
		{Name: "has_starred", Type: proto.ColumnType_BOOL, Transform: transform.FromField("HasStarred", "Node.HasStarred"), Description: "If true, you have starred this repository."},
		{Name: "temp_clone_token", Type: proto.ColumnType_STRING, Transform: transform.FromField("TempCloneToken", "Node.TempCloneToken"), Description: "Temporary authentication token for cloning this repository."},
		{Name: "possible_commit_emails", Type: proto.ColumnType_JSON, Transform: transform.FromField("PossibleCommitEmails", "Node.PossibleCommitEmails").NullIfZero(), Description: "A list of emails you can commit to this repository with."},
		{Name: "subscription", Type: proto.ColumnType_STRING, Transform: transform.FromField("Subscription", "Node.Subscription"), Description: "Identifies if the current user is watching, not watching, or ignoring the repository."},
		{Name: "visibility", Type: proto.ColumnType_STRING, Transform: transform.FromField("Visibility", "Node.Visibility"), Description: "Indicates the repository's visibility level."},
		{Name: "your_permission", Type: proto.ColumnType_STRING, Transform: transform.FromField("YourPermission", "Node.YourPermission"), Description: "Your permission level on the repository. Will return null if authenticated as an GitHub App."},
		{Name: "web_commit_signoff_required", Type: proto.ColumnType_BOOL, Transform: transform.FromField("WebCommitSignoffRequired", "Node.WebCommitSignoffRequired"), Description: "If true, contributors are required to sign off on web-based commits in this repository."},
		{Name: "outside_collaborators_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("OutsideCollaborators.TotalCount", "Node.OutsideCollaborators.TotalCount"), Description: "Count of outside collaborators."},
		{Name: "repository_topics_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("RepositoryTopics.TotalCount", "Node.RepositoryTopics.TotalCount"), Description: "Count of topics associated with the repository."},
		{Name: "open_issues_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("OpenIssues.TotalCount", "Node.OpenIssues.TotalCount"), Description: "Count of issues open on the repository."},
		{Name: "watchers_total_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Watchers.TotalCount", "Node.Watchers.TotalCount"), Description: "Count of watchers on the repository."},
	}
}

func tableGitHubRepository() *plugin.Table {
	return &plugin.Table{
		Name:        "github_repository",
		Description: "GitHub Repositories contain all of your project's files and each file's revision history.",
		List: &plugin.ListConfig{
			Hydrate:           tableGitHubRepositoryList,
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			KeyColumns:        plugin.SingleColumn("full_name"),
		},
		Columns: gitHubRepositoryColumns(),
	}
}

func tableGitHubRepositoryList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	repoFullName := d.EqualsQuals["full_name"].GetStringValue()
	owner, repoName := parseRepoFullName(repoFullName)

	var query struct {
		RateLimit  models.RateLimit
		Repository models.Repository `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repoName),
	}

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_repository", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_repository", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, query.Repository)

	return nil, nil
}
