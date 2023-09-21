package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v55/github"
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
		{Name: "allow_update_branch", Type: proto.ColumnType_BOOL, Description: "If true, a pull request head branch that is behind its base branch can always be updated even if it is not required to be up to date before merging.", Hydrate: repoHydrateAllowUpdateBranch, Transform: transform.FromValue()},
		{Name: "archived_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when repository was archived.", Hydrate: repoHydrateArchivedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp)},
		{Name: "auto_merge_allowed", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateAutoMergeAllowed, Transform: transform.FromValue(), Description: "If true, auto-merge can be enabled on pull requests in this repository."},
		{Name: "code_of_conduct", Type: proto.ColumnType_JSON, Hydrate: repoHydrateCodeOfConduct, Transform: transform.FromValue().NullIfZero(), Description: "The code of conduct for this repository."},
		{Name: "contact_links", Type: proto.ColumnType_JSON, Hydrate: repoHydrateContactLinks, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "List of contact links associated to the repository."},
		{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: repoHydrateCreatedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the repository was created."},
		{Name: "default_branch_ref", Type: proto.ColumnType_JSON, Hydrate: repoHydrateDefaultBranchRef, Transform: transform.FromValue().NullIfZero(), Description: "Default ref information."},
		{Name: "delete_branch_on_merge", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateDeleteBranchOnMerge, Transform: transform.FromValue(), Description: "If true, branches are automatically deleted when merged in this repository."},
		{Name: "description", Type: proto.ColumnType_STRING, Hydrate: repoHydrateDescription, Transform: transform.FromValue(), Description: "The description of the repository."},
		{Name: "disk_usage", Type: proto.ColumnType_INT, Hydrate: repoHydrateDiskUsage, Transform: transform.FromValue(), Description: "Number of kilobytes this repository occupies on disk."},
		{Name: "fork_count", Type: proto.ColumnType_INT, Hydrate: repoHydrateForkCount, Transform: transform.FromValue(), Description: "Number of forks there are of this repository in the whole network."},
		{Name: "forking_allowed", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateForkingAllowed, Transform: transform.FromValue(), Description: "If true, repository allows forks."},
		{Name: "funding_links", Type: proto.ColumnType_JSON, Hydrate: repoHydrateFundingLinks, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "The funding links for this repository."},
		{Name: "has_discussions_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasDiscussionsEnabled, Transform: transform.FromValue(), Description: "If true, the repository has the Discussions feature enabled."},
		{Name: "has_issues_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasIssuesEnabled, Transform: transform.FromValue(), Description: "If true, the repository has issues feature enabled."},
		{Name: "has_projects_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasProjectsEnabled, Transform: transform.FromValue(), Description: "If true, the repository has the Projects feature enabled."},
		{Name: "has_vulnerability_alerts_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasVulnerabilityAlertsEnabled, Transform: transform.FromValue(), Description: "If true, vulnerability alerts are enabled for the repository."},
		{Name: "has_wiki_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasWikiEnabled, Transform: transform.FromValue(), Description: "If true, the repository has wiki feature enabled."},
		{Name: "homepage_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateHomepageUrl, Transform: transform.FromValue(), Description: "The external URL of the repository if set."},
		{Name: "interaction_ability", Type: proto.ColumnType_JSON, Hydrate: repoHydrateInteractionAbility, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "The interaction ability settings for this repository."},
		{Name: "is_archived", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsArchived, Transform: transform.FromValue(), Description: "If true, the repository is unmaintained (archived)."},
		{Name: "is_blank_issues_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsBlankIssuesEnabled, Transform: transform.FromValue(), Description: "If true, blank issue creation is allowed."},
		{Name: "is_disabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsDisabled, Transform: transform.FromValue(), Description: "If true, this repository disabled."},
		{Name: "is_empty", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsEmpty, Transform: transform.FromValue(), Description: "If true, this repository is empty."},
		{Name: "is_fork", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsFork, Transform: transform.FromValue(), Description: "If true, the repository is a fork."},
		{Name: "is_in_organization", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsInOrganization, Transform: transform.FromValue(), Description: "If true, repository is either owned by an organization, or is a private fork of an organization repository."},
		{Name: "is_locked", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsLocked, Transform: transform.FromValue(), Description: "If true, repository is locked."},
		{Name: "is_mirror", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsMirror, Transform: transform.FromValue(), Description: "If true, the repository is a mirror."},
		{Name: "is_private", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsPrivate, Transform: transform.FromValue(), Description: "If true, the repository is private or internal."},
		{Name: "is_security_policy_enabled", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsSecurityPolicyEnabled, Transform: transform.FromValue(), Description: "If true, repository has a security policy."},
		{Name: "is_template", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsTemplate, Transform: transform.FromValue(), Description: "If true, the repository is a template that can be used to generate new repositories."},
		{Name: "is_user_configuration_repository", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateIsUserConfigurationRepository, Transform: transform.FromValue(), Description: "If true, this is a user configuration repository."},
		{Name: "issue_templates", Type: proto.ColumnType_JSON, Hydrate: repoHydrateIssueTemplates, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "A list of issue templates associated to the repository."},
		{Name: "license_info", Type: proto.ColumnType_JSON, Hydrate: repoHydrateLicenseInfo, Transform: transform.FromValue().NullIfZero(), Description: "The license associated with the repository."},
		{Name: "lock_reason", Type: proto.ColumnType_STRING, Hydrate: repoHydrateLockReason, Transform: transform.FromValue(), Description: "The reason the repository has been locked."},
		{Name: "merge_commit_allowed", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateMergeCommitAllowed, Transform: transform.FromValue(), Description: "If true, PRs are merged with a merge commit on this repository."},
		{Name: "merge_commit_message", Type: proto.ColumnType_STRING, Hydrate: repoHydrateMergeCommitMessage, Transform: transform.FromValue(), Description: "How the default commit message will be generated when merging a pull request."},
		{Name: "merge_commit_title", Type: proto.ColumnType_STRING, Hydrate: repoHydrateMergeCommitTitle, Transform: transform.FromValue(), Description: "How the default commit title will be generated when merging a pull request."},
		{Name: "mirror_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateMirrorUrl, Transform: transform.FromValue(), Description: "The repository's original mirror URL."},
		{Name: "name_with_owner", Type: proto.ColumnType_STRING, Transform: transform.FromField("NameWithOwner", "Node.NameWithOwner"), Description: "The repository's name with owner."},
		{Name: "open_graph_image_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateOpenGraphImageUrl, Transform: transform.FromValue(), Description: "The image used to represent this repository in Open Graph data."},
		{Name: "owner_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Owner.Login", "Node.Owner.Login"), Description: "Login of the repository owner."},
		{Name: "primary_language", Type: proto.ColumnType_JSON, Hydrate: repoHydratePrimaryLanguage, Transform: transform.FromValue().NullIfZero(), Description: "The primary language of the repository's code."},
		{Name: "projects_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateProjectsUrl, Transform: transform.FromValue(), Description: "The URL listing the repository's projects."},
		{Name: "pull_request_templates", Type: proto.ColumnType_JSON, Hydrate: repoHydratePullRequestTemplates, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "Returns a list of pull request templates associated to the repository."},
		{Name: "pushed_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: repoHydratePushedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "Timestamp when the repository was last pushed to."},
		{Name: "rebase_merge_allowed", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateRebaseMergeAllowed, Transform: transform.FromValue(), Description: "If true, rebase-merging is enabled on this repository."},
		{Name: "security_policy_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateSecurityPolicyUrl, Transform: transform.FromValue(), Description: "The security policy URL."},
		{Name: "squash_merge_allowed", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateSquashMergeAllowed, Transform: transform.FromValue(), Description: "If true, squash-merging is enabled on this repository."},
		{Name: "squash_merge_commit_message", Type: proto.ColumnType_STRING, Hydrate: repoHydrateSquashMergeCommitMessage, Transform: transform.FromValue(), Description: "How the default commit message will be generated when squash merging a pull request."},
		{Name: "squash_merge_commit_title", Type: proto.ColumnType_STRING, Hydrate: repoHydrateSquashMergeCommitTitle, Transform: transform.FromValue(), Description: "How the default commit title will be generated when squash merging a pull request."},
		{Name: "ssh_url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateSshUrl, Transform: transform.FromValue(), Description: "The SSH URL to clone this repository."},
		{Name: "stargazer_count", Type: proto.ColumnType_INT, Hydrate: repoHydrateStargazerCount, Transform: transform.FromValue(), Description: "Returns a count of how many stargazers there are on this repository."},
		{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Hydrate: repoHydrateUpdatedAt, Transform: transform.FromValue().NullIfZero().Transform(convertTimestamp), Description: "Timestamp when repository was last updated."},
		{Name: "url", Type: proto.ColumnType_STRING, Hydrate: repoHydrateUrl, Transform: transform.FromValue(), Description: "The URL of the repository."},
		{Name: "uses_custom_open_graph_image", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateUsesCustomOpenGraphImage, Transform: transform.FromValue(), Description: "if true, this repository has a custom image to use with Open Graph as opposed to being represented by the owner's avatar."},
		{Name: "can_administer", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateCanAdminister, Transform: transform.FromValue(), Description: "If true, you can administer this repository."},
		{Name: "can_create_projects", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateCanCreateProjects, Transform: transform.FromValue(), Description: "If true, you can create projects in this repository."},
		{Name: "can_subscribe", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateCanSubscribe, Transform: transform.FromValue(), Description: "If true, you can subscribe to this repository."},
		{Name: "can_update_topics", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateCanUpdateTopics, Transform: transform.FromValue(), Description: "If true, you can update topics on this repository."},
		{Name: "has_starred", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateHasStarred, Transform: transform.FromValue(), Description: "If true, you have starred this repository."},
		{Name: "possible_commit_emails", Type: proto.ColumnType_JSON, Hydrate: repoHydratePossibleCommitEmails, Transform: transform.FromValue().NullIfZero().NullIfEmptySlice(), Description: "A list of emails you can commit to this repository with."},
		{Name: "subscription", Type: proto.ColumnType_STRING, Hydrate: repoHydrateSubscription, Transform: transform.FromValue(), Description: "Identifies if the current user is watching, not watching, or ignoring the repository."},
		{Name: "visibility", Type: proto.ColumnType_STRING, Hydrate: repoHydrateVisibility, Transform: transform.FromValue(), Description: "Indicates the repository's visibility level."},
		{Name: "your_permission", Type: proto.ColumnType_STRING, Hydrate: repoHydrateYourPermission, Transform: transform.FromValue(), Description: "Your permission level on the repository. Will return null if authenticated as an GitHub App."},
		{Name: "web_commit_signoff_required", Type: proto.ColumnType_BOOL, Hydrate: repoHydrateWebCommitSignoffRequired, Transform: transform.FromValue(), Description: "If true, contributors are required to sign off on web-based commits in this repository."},
		{Name: "repository_topics_total_count", Type: proto.ColumnType_INT, Hydrate: repoHydrateRepositoryTopicsCount, Transform: transform.FromValue(), Description: "Count of topics associated with the repository."},
		{Name: "open_issues_total_count", Type: proto.ColumnType_INT, Hydrate: repoHydrateOpenIssuesCount, Transform: transform.FromValue(), Description: "Count of issues open on the repository."},
		{Name: "watchers_total_count", Type: proto.ColumnType_INT, Hydrate: repoHydrateWatchersCount, Transform: transform.FromValue(), Description: "Count of watchers on the repository."},
		// // Columns from v3 api - hydrates
		{Name: "hooks", Type: proto.ColumnType_JSON, Description: "The API Hooks URL.", Hydrate: hydrateRepositoryHooksFromV3, Transform: transform.FromValue()},
		{Name: "topics", Type: proto.ColumnType_JSON, Description: "The topics (similar to tags or labels) associated with the repository.", Hydrate: hydrateRepositoryDataFromV3},
		{Name: "subscribers_count", Type: proto.ColumnType_INT, Description: "The number of users who have subscribed to the repository.", Hydrate: hydrateRepositoryDataFromV3},
		{Name: "has_downloads", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Downloads feature is enabled on the repository.", Hydrate: hydrateRepositoryDataFromV3},
		{Name: "has_pages", Type: proto.ColumnType_BOOL, Description: "If true, the GitHub Pages feature is enabled on the repository.", Hydrate: hydrateRepositoryDataFromV3},
		{Name: "network_count", Type: proto.ColumnType_INT, Description: "The number of member repositories in the network.", Hydrate: hydrateRepositoryDataFromV3},
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
	appendRepoColumnIncludes(&variables, d.QueryContext.Columns)

	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_repository", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_repository", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, query.Repository)

	return nil, nil
}

func hydrateRepositoryDataFromV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repo, err := extractRepoFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	owner := repo.Owner.Login
	repoName := repo.Name

	client := connect(ctx, d)
	r, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil, nil
		}
		return nil, err
	}
	if r == nil {
		return nil, nil
	}

	return r, nil
}

func hydrateRepositoryHooksFromV3(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	repo, err := extractRepoFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	owner := repo.Owner.Login
	repoName := repo.Name

	client := connect(ctx, d)
	var repositoryHooks []*github.Hook
	opt := &github.ListOptions{PerPage: 100}

	for {
		hooks, resp, err := client.Repositories.ListHooks(ctx, owner, repoName, opt)
		if err != nil && strings.Contains(err.Error(), "Not Found") {
			return nil, nil
		} else if err != nil {
			return nil, err
		}
		repositoryHooks = append(repositoryHooks, hooks...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return repositoryHooks, nil
}
