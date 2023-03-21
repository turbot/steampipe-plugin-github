package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v48/github"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGitHubBranchProtection(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch_protection",
		Description: "Branch protection defines rules for pushing to and managing a branch.",
		List: &plugin.ListConfig{
			KeyColumns:    plugin.SingleColumn("repository_full_name"),
			Hydrate:       tableGitHubRepositoryBranchProtectionGet,
			ParentHydrate: tableGitHubBranchList,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"repository_full_name", "name"}),
			Hydrate:    tableGitHubRepositoryBranchProtectionGet,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "The full name of the repository (login/repo-name)."},
			{Name: "name", Type: proto.ColumnType_STRING, Hydrate: branchNameQual, Transform: transform.FromValue(), Description: "The branch name."},
			{Name: "restrictions_apps", Type: proto.ColumnType_JSON, Transform: transform.FromField("Restrictions.Apps"), Description: "Applications can push to the branch only if in this list."},
			{Name: "restrictions_teams", Type: proto.ColumnType_JSON, Transform: transform.FromField("Restrictions.Teams"), Description: "Teams can push to the branch only if in this list."},
			{Name: "restrictions_users", Type: proto.ColumnType_JSON, Transform: transform.FromField("Restrictions.Users"), Description: "Users can push to the branch only if in this list."},
			{Name: "enforce_admins_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("EnforceAdmins.Enabled"), Description: "If true, enforce all configured restrictions for administrators."},
			{Name: "allow_deletions_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("AllowDeletions.Enabled"), Description: "If true, allow users with push access to delete matching branches."},
			{Name: "allow_force_pushes_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("AllowForcePushes.Enabled"), Description: "If true, permit force pushes for all users with push access."},
			{Name: "required_conversation_resolution_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("RequiredConversationResolution.Enabled"), Description: "If enabled, requires all comments on the pull request to be resolved before it can be merged to a protected branch."},
			{Name: "required_linear_history_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("RequireLinearHistory.Enabled"), Description: "If true, prevent merge commits from being pushed to matching branches."},
			{Name: "required_status_checks", Type: proto.ColumnType_JSON, Description: "Status checks that must pass before a branch can be merged into branches matching this rule."},
			{Name: "required_pull_request_reviews", Type: proto.ColumnType_JSON, Description: "Pull request reviews required before merging."},
			{Name: "signatures_protected_branch_enabled", Type: proto.ColumnType_BOOL, Description: "Commits pushed to matching branches must have verified signatures.", Hydrate: repositorySignaturesProtectedBranchGet, Transform: transform.FromValue()},
		},
	}
}

//// LIST FUNCTION

func tableGitHubRepositoryBranchProtectionGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	quals := d.EqualsQuals

	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	branchName := ""

	if h.Item != nil {
		b := h.Item.(*github.Branch)
		branchName = *b.Name
	} else {
		branchName = quals["name"].GetStringValue()
	}

	logger.Trace("tableGitHubRepositoryBranchProtectionGet", "owner", owner, "repo", repo, "branchName", branchName)

	client := connect(ctx, d)

	type GetResponse struct {
		protection *github.Protection
	}

	get := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		protection, _, err := client.Repositories.GetBranchProtection(ctx, owner, repo, branchName)
		if err != nil {
			// For private and archived repositories, users who do not have owner/admin access will get the below error
			// 403 Upgrade to GitHub Pro or make this repository public to enable this feature.
			// For repository owners the API will return nil if the repository is private and archived
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Upgrade to GitHub Pro") || strings.Contains(err.Error(), "branch is not protected") {
				return nil, nil
			}
			return nil, err
		}

		return GetResponse{
			protection: protection,
		}, err
	}

	getDetails, err := retryHydrate(ctx, d, h, get)
	if err != nil {
		return nil, err
	}

	if getDetails == nil {
		return nil, nil
	}
	getResp := getDetails.(GetResponse)
	protection := getResp.protection

	if protection != nil {
		d.StreamLeafListItem(ctx, protection)
	}
	return nil, nil
}

func repositorySignaturesProtectedBranchGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	quals := d.EqualsQuals
	fullName := quals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)
	branchName := ""

	if h.ParentItem != nil {
		b := h.ParentItem.(*github.Branch)
		branchName = *b.Name
	} else {
		branchName = quals["name"].GetStringValue()
	}

	logger.Trace("tableGitHubRepositoryBranchProtectionGet", "owner", owner, "repo", repo, "branchName", branchName)
	client := connect(ctx, d)

	type GetResponse struct {
		protectedBranch *github.SignaturesProtectedBranch
	}

	getDetails := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		protectedBranch, _, err := client.Repositories.GetSignaturesProtectedBranch(ctx, owner, repo, branchName)
		if err != nil {
			// For private and archived repositories, users who do not have owner/admin access will get the below error
			// 403 Upgrade to GitHub Pro or make this repository public to enable this feature.
			// For repository owners the API will return nil if the repository is private and archived
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "Upgrade to GitHub Pro") || strings.Contains(err.Error(), "branch is not protected") {
				return nil, nil
			}
			return nil, err
		}

		return GetResponse{
			protectedBranch: protectedBranch,
		}, err
	}

	getResponse, err := retryHydrate(ctx, d, h, getDetails)
	if err != nil {
		return nil, err
	}
	getResp := getResponse.(GetResponse)
	protectedBranch := getResp.protectedBranch

	if protectedBranch != nil {
		return protectedBranch.Enabled, nil
	}
	return nil, nil
}

func branchNameQual(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	if quals["name"] != nil {
		return quals["name"].GetStringValue(), nil
	}
	item := h.ParentItem.(*github.Branch)
	return *item.Name, nil
}
