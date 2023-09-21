package github

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"slices"
)

func extractRepoFromHydrateItem(h *plugin.HydrateData) (models.Repository, error) {
	if repo, ok := h.Item.(models.Repository); ok {
		return repo, nil
	} else if searchResult, ok := h.Item.(models.SearchRepositoryResult); ok {
		return searchResult.Node.Repository, nil
	} else if teamResult, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return teamResult.Node, nil
	} else {
		return models.Repository{}, fmt.Errorf("unable to parse hydrate item %v as a Repository", h.Item)
	}
}

func appendRepoColumnIncludes(m *map[string]interface{}, cols []string) {
	optionals := map[string]string{
		"allow_update_branch":              "includeAllowUpdateBranch",
		"archived_at":                      "includeArchivedAt",
		"auto_merge_allowed":               "includeAutoMergeAllowed",
		"can_administer":                   "includeCanAdminister",
		"can_create_projects":              "includeCanCreateProjects",
		"can_subscribe":                    "includeCanSubscribe",
		"can_update_topics":                "includeCanUpdateTopics",
		"code_of_conduct":                  "includeCodeOfConduct",
		"contact_links":                    "includeContactLinks",
		"created_at":                       "includeCreatedAt",
		"default_branch_ref":               "includeDefaultBranchRef",
		"delete_branch_on_merge":           "includeDeleteBranchOnMerge",
		"description":                      "includeDescription",
		"disk_usage":                       "includeDiskUsage",
		"fork_count":                       "includeForkCount",
		"forking_allowed":                  "includeForkingAllowed",
		"funding_links":                    "includeFundingLinks",
		"has_discussions_enabled":          "includeHasDiscussionsEnabled",
		"has_issues_enabled":               "includeHasIssuesEnabled",
		"has_projects_enabled":             "includeHasProjectsEnabled",
		"has_starred":                      "includeHasStarred",
		"has_vulnerability_alerts_enabled": "includeHasVulnerabilityAlertsEnabled",
		"has_wiki_enabled":                 "includeHasWikiEnabled",
		"homepage_url":                     "includeHomepageUrl",
		"interaction_ability":              "includeInteractionAbility",
		"is_archived":                      "includeIsArchived",
		"is_blank_issues_enabled":          "includeIsBlankIssuesEnabled",
		"is_disabled":                      "includeIsDisabled",
		"is_empty":                         "includeIsEmpty",
		"is_fork":                          "includeIsFork",
		"is_in_organization":               "includeIsInOrganization",
		"is_locked":                        "includeIsLocked",
		"is_mirror":                        "includeIsMirror",
		"is_private":                       "includeIsPrivate",
		"is_security_policy_enabled":       "includeIsSecurityPolicyEnabled",
		"is_template":                      "includeIsTemplate",
		"is_user_configuration_repository": "includeIsUserConfigurationRepository",
		"issue_templates":                  "includeIssueTemplates",
		"license_info":                     "includeLicenseInfo",
		"lock_reason":                      "includeLockReason",
		"merge_commit_allowed":             "includeMergeCommitAllowed",
		"merge_commit_message":             "includeMergeCommitMessage",
		"merge_commit_title":               "includeMergeCommitTitle",
		"mirror_url":                       "includeMirrorUrl",
		"open_graph_image_url":             "includeOpenGraphImageUrl",
		"open_issues_total_count":          "includeOpenIssues",
		"possible_commit_emails":           "includePossibleCommitEmails",
		"primary_language":                 "includePrimaryLanguage",
		"projects_url":                     "includeProjectsUrl",
		"pull_request_templates":           "includePullRequestTemplates",
		"pushed_at":                        "includePushedAt",
		"rebase_merge_allowed":             "includeRebaseMergeAllowed",
		"repository_topics_total_count":    "includeRepositoryTopics",
		"security_policy_url":              "includeSecurityPolicyUrl",
		"squash_merge_allowed":             "includeSquashMergeAllowed",
		"squash_merge_commit_message":      "includeSquashMergeCommitMessage",
		"squash_merge_commit_title":        "includeSquashMergeCommitTitle",
		"ssh_url":                          "includeSshUrl",
		"stargazer_count":                  "includeStargazerCount",
		"subscription":                     "includeSubscription",
		"updated_at":                       "includeUpdatedAt",
		"url":                              "includeUrl",
		"uses_custom_open_graph_image":     "includeUsesCustomOpenGraphImage",
		"visibility":                       "includeVisibility",
		"watchers_total_count":             "includeWatchers",
		"web_commit_signoff_required":      "includeWebCommitSignoffRequired",
		"your_permission":                  "includeYourPermission",
	}

	for key, value := range optionals {
		(*m)[value] = githubv4.Boolean(slices.Contains(cols, key))
	}
}

func repoHydrateAllowUpdateBranch(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.AllowUpdateBranch, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.AllowUpdateBranch, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.AllowUpdateBranch, nil
	}
	return nil, nil
}

func repoHydrateArchivedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.ArchivedAt, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.ArchivedAt, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.ArchivedAt, nil
	}
	return nil, nil
}

func repoHydrateAutoMergeAllowed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.AutoMergeAllowed, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.AutoMergeAllowed, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.AutoMergeAllowed, nil
	}
	return nil, nil
}

func repoHydrateCodeOfConduct(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CodeOfConduct, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CodeOfConduct, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CodeOfConduct, nil
	}
	return nil, nil
}

func repoHydrateContactLinks(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.ContactLinks, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.ContactLinks, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.ContactLinks, nil
	}
	return nil, nil
}

func repoHydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CreatedAt, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CreatedAt, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CreatedAt, nil
	}
	return nil, nil
}

func repoHydrateDefaultBranchRef(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.DefaultBranchRef, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.DefaultBranchRef, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.DefaultBranchRef, nil
	}
	return nil, nil
}

func repoHydrateDeleteBranchOnMerge(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.DeleteBranchOnMerge, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.DeleteBranchOnMerge, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.DeleteBranchOnMerge, nil
	}
	return nil, nil
}

func repoHydrateDescription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.Description, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.Description, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.Description, nil
	}
	return nil, nil
}

func repoHydrateDiskUsage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.DiskUsage, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.DiskUsage, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.DiskUsage, nil
	}
	return nil, nil
}

func repoHydrateForkCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.ForkCount, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.ForkCount, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.ForkCount, nil
	}
	return nil, nil
}

func repoHydrateForkingAllowed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.ForkingAllowed, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.ForkingAllowed, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.ForkingAllowed, nil
	}
	return nil, nil
}

func repoHydrateFundingLinks(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.FundingLinks, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.FundingLinks, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.FundingLinks, nil
	}
	return nil, nil
}

func repoHydrateHasDiscussionsEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasDiscussionsEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasDiscussionsEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasDiscussionsEnabled, nil
	}
	return nil, nil
}

func repoHydrateHasIssuesEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasIssuesEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasIssuesEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasIssuesEnabled, nil
	}
	return nil, nil
}

func repoHydrateHasProjectsEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasProjectsEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasProjectsEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasProjectsEnabled, nil
	}
	return nil, nil
}

func repoHydrateHasVulnerabilityAlertsEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasVulnerabilityAlertsEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasVulnerabilityAlertsEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasVulnerabilityAlertsEnabled, nil
	}
	return nil, nil
}

func repoHydrateHasWikiEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasWikiEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasWikiEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasWikiEnabled, nil
	}
	return nil, nil
}

func repoHydrateHomepageUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HomepageUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HomepageUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HomepageUrl, nil
	}
	return nil, nil
}

func repoHydrateInteractionAbility(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.InteractionAbility, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.InteractionAbility, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.InteractionAbility, nil
	}
	return nil, nil
}

func repoHydrateIsArchived(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsArchived, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsArchived, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsArchived, nil
	}
	return nil, nil
}

func repoHydrateIsBlankIssuesEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsBlankIssuesEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsBlankIssuesEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsBlankIssuesEnabled, nil
	}
	return nil, nil
}

func repoHydrateIsDisabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsDisabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsDisabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsDisabled, nil
	}
	return nil, nil
}

func repoHydrateIsEmpty(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsEmpty, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsEmpty, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsEmpty, nil
	}
	return nil, nil
}

func repoHydrateIsFork(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsFork, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsFork, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsFork, nil
	}
	return nil, nil
}

func repoHydrateIsInOrganization(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsInOrganization, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsInOrganization, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsInOrganization, nil
	}
	return nil, nil
}

func repoHydrateIsLocked(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsLocked, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsLocked, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsLocked, nil
	}
	return nil, nil
}

func repoHydrateIsMirror(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsMirror, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsMirror, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsMirror, nil
	}
	return nil, nil
}

func repoHydrateIsPrivate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsPrivate, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsPrivate, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsPrivate, nil
	}
	return nil, nil
}

func repoHydrateIsSecurityPolicyEnabled(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsSecurityPolicyEnabled, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsSecurityPolicyEnabled, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsSecurityPolicyEnabled, nil
	}
	return nil, nil
}

func repoHydrateIsTemplate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsTemplate, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsTemplate, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsTemplate, nil
	}
	return nil, nil
}

func repoHydrateIsUserConfigurationRepository(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IsUserConfigurationRepository, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IsUserConfigurationRepository, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IsUserConfigurationRepository, nil
	}
	return nil, nil
}

func repoHydrateIssueTemplates(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.IssueTemplates, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.IssueTemplates, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.IssueTemplates, nil
	}
	return nil, nil
}

func repoHydrateLicenseInfo(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.LicenseInfo, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.LicenseInfo, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.LicenseInfo, nil
	}
	return nil, nil
}

func repoHydrateLockReason(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.LockReason, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.LockReason, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.LockReason, nil
	}
	return nil, nil
}

func repoHydrateMergeCommitAllowed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.MergeCommitAllowed, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.MergeCommitAllowed, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.MergeCommitAllowed, nil
	}
	return nil, nil
}

func repoHydrateMergeCommitMessage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.MergeCommitMessage, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.MergeCommitMessage, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.MergeCommitMessage, nil
	}
	return nil, nil
}

func repoHydrateMergeCommitTitle(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.MergeCommitTitle, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.MergeCommitTitle, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.MergeCommitTitle, nil
	}
	return nil, nil
}

func repoHydrateMirrorUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.MirrorUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.MirrorUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.MirrorUrl, nil
	}
	return nil, nil
}

func repoHydrateOpenGraphImageUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.OpenGraphImageUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.OpenGraphImageUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.OpenGraphImageUrl, nil
	}
	return nil, nil
}

func repoHydratePrimaryLanguage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.PrimaryLanguage, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.PrimaryLanguage, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.PrimaryLanguage, nil
	}
	return nil, nil
}

func repoHydrateProjectsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.ProjectsUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.ProjectsUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.ProjectsUrl, nil
	}
	return nil, nil
}

func repoHydratePullRequestTemplates(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.PullRequestTemplates, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.PullRequestTemplates, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.PullRequestTemplates, nil
	}
	return nil, nil
}

func repoHydratePushedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.PushedAt, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.PushedAt, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.PushedAt, nil
	}
	return nil, nil
}

func repoHydrateRebaseMergeAllowed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.RebaseMergeAllowed, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.RebaseMergeAllowed, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.RebaseMergeAllowed, nil
	}
	return nil, nil
}

func repoHydrateSecurityPolicyUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.SecurityPolicyUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.SecurityPolicyUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.SecurityPolicyUrl, nil
	}
	return nil, nil
}

func repoHydrateSquashMergeAllowed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.SquashMergeAllowed, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.SquashMergeAllowed, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.SquashMergeAllowed, nil
	}
	return nil, nil
}

func repoHydrateSquashMergeCommitMessage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.SquashMergeCommitMessage, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.SquashMergeCommitMessage, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.SquashMergeCommitMessage, nil
	}
	return nil, nil
}

func repoHydrateSquashMergeCommitTitle(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.SquashMergeCommitTitle, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.SquashMergeCommitTitle, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.SquashMergeCommitTitle, nil
	}
	return nil, nil
}

func repoHydrateSshUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.SshUrl, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.SshUrl, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.SshUrl, nil
	}
	return nil, nil
}

func repoHydrateStargazerCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.StargazerCount, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.StargazerCount, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.StargazerCount, nil
	}
	return nil, nil
}

func repoHydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.UpdatedAt, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.UpdatedAt, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.UpdatedAt, nil
	}
	return nil, nil
}

func repoHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.Url, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.Url, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.Url, nil
	}
	return nil, nil
}

func repoHydrateUsesCustomOpenGraphImage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.UsesCustomOpenGraphImage, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.UsesCustomOpenGraphImage, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.UsesCustomOpenGraphImage, nil
	}
	return nil, nil
}

func repoHydrateCanAdminister(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CanAdminister, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CanAdminister, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CanAdminister, nil
	}
	return nil, nil
}

func repoHydrateCanCreateProjects(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CanCreateProjects, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CanCreateProjects, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CanCreateProjects, nil
	}
	return nil, nil
}

func repoHydrateCanSubscribe(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CanSubscribe, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CanSubscribe, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CanSubscribe, nil
	}
	return nil, nil
}

func repoHydrateCanUpdateTopics(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.CanUpdateTopics, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.CanUpdateTopics, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.CanUpdateTopics, nil
	}
	return nil, nil
}

func repoHydrateHasStarred(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.HasStarred, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.HasStarred, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.HasStarred, nil
	}
	return nil, nil
}

func repoHydratePossibleCommitEmails(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.PossibleCommitEmails, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.PossibleCommitEmails, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.PossibleCommitEmails, nil
	}
	return nil, nil
}

func repoHydrateSubscription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.Subscription, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.Subscription, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.Subscription, nil
	}
	return nil, nil
}

func repoHydrateVisibility(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.Visibility, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.Visibility, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.Visibility, nil
	}
	return nil, nil
}

func repoHydrateYourPermission(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.YourPermission, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.YourPermission, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.YourPermission, nil
	}
	return nil, nil
}

func repoHydrateWebCommitSignoffRequired(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.WebCommitSignoffRequired, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.WebCommitSignoffRequired, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.WebCommitSignoffRequired, nil
	}
	return nil, nil
}

func repoHydrateRepositoryTopicsCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.RepositoryTopics.TotalCount, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.RepositoryTopics.TotalCount, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.RepositoryTopics.TotalCount, nil
	}
	return nil, nil
}

func repoHydrateOpenIssuesCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.OpenIssues.TotalCount, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.OpenIssues.TotalCount, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.OpenIssues.TotalCount, nil
	}
	return nil, nil
}

func repoHydrateWatchersCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	if r, ok := h.Item.(models.Repository); ok {
		return r.Watchers.TotalCount, nil
	} else if r, ok := h.Item.(models.SearchRepositoryResult); ok {
		return r.Node.Watchers.TotalCount, nil
	} else if r, ok := h.Item.(models.TeamRepositoryWithPermission); ok {
		return r.Node.Watchers.TotalCount, nil
	}
	return nil, nil
}
