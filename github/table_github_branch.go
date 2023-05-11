package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableGitHubBranch() *plugin.Table {
	return &plugin.Table{
		Name:        "github_branch",
		Description: "Branches in the given repository.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "repository_full_name", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"404"}),
			Hydrate:           tableGitHubBranchList,
		},
		Columns: []*plugin.Column{
			{Name: "repository_full_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("repository_full_name"), Description: "Full name of the repository that contains the branch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the branch.", Transform: transform.FromField("Node.Name")},
			{Name: "commit_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Sha"), Description: "Commit SHA the branch refers to."},
			{Name: "commit_short_sha", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.ShortSha"), Description: "Commit short SHA the branch refers to."},
			{Name: "commit_authored_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Node.Target.Commit.AuthoredDate"), Description: "Date commit was authored."},
			{Name: "commit_author_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Author.Name"), Description: "Commit authors display name."},
			{Name: "commit_author_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Author.User.Login"), Description: "Commit authors login."},
			{Name: "commit_committed_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Node.Target.Commit.CommittedDate"), Description: "Date commit was committed."},
			{Name: "commit_committer_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Committer.Name"), Description: "Commit committers display name."},
			{Name: "commit_committer_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Committer.User.Login"), Description: "Commit committers login."},
			{Name: "commit_message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Message"), Description: "Commit message."},
			{Name: "commit_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Url"), Description: "Commit URL the branch refers to."},
			{Name: "commit_additions", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.Target.Commit.Additions"), Description: "Number of additions in the commit."},
			{Name: "commit_deletions", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.Target.Commit.Deletions"), Description: "Number of deletions in the commit."},
			{Name: "commit_changed_files", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.Target.Commit.ChangedFiles").NullIfZero(), Description: "Number of files changed in the commit if available (null if not available)."},
			{Name: "commit_authored_by_committer", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.Target.Commit.AuthoredByCommitter"), Description: "If true, the commits committer and author are the same."},
			{Name: "commit_committed_via_web", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.Target.Commit.CommittedViaWeb"), Description: "If true, the commit was from the GitHub web app."},
			{Name: "commit_signature_is_valid", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.Target.Commit.Signature.IsValid"), Description: "If true, commit was signed by a valid signature."},
			{Name: "commit_signature_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Signature.Email"), Description: "Email associated with the commit signature."},
			{Name: "commit_signature_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Signature.Signer.Login"), Description: "Login associated with the commit signature."},
			{Name: "commit_tarball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.TarballUrl"), Description: "URL to download a tar file of the code for this commit."},
			{Name: "commit_zipball_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.ZipballUrl"), Description: "URL to download a zip file of the code for this commit."},
			{Name: "commit_tree_url", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.TreeUrl"), Description: "URL for the tree of this commit."},
			{Name: "commit_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.Target.Commit.Status.State"), Description: "Status of the commit (ERROR, EXPECTED, FAILURE, PENDING, SUCCESS)."},
			{Name: "protected", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.NodeId").Transform(HasValue), Description: "If true, the branch is protected."},
			{Name: "protection_rule_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.BranchProtectionRule.Id").NullIfZero(), Description: "Branch protection rule id, null if not protected."},
			{Name: "protection_rule_node_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.BranchProtectionRule.NodeId"), Description: "The Node ID of the branch protection rule."},
			{Name: "protection_rule_matching_branches", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.BranchProtectionRule.MatchingBranches.TotalCount"), Description: "Count of branches which match this rule."},
			{Name: "protection_rule_is_admin_enforced", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.IsAdminEnforced"), Description: "If true, enforce all configured restrictions for administrators."},
			{Name: "protection_rule_allows_deletions", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.AllowsDeletions"), Description: "If true, allow users with push access to delete matching branches."},
			{Name: "protection_rule_allows_force_pushes", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.AllowsForcePushes"), Description: "If true, permit force pushes for all users with push access."},
			{Name: "protection_rule_blocks_creations", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.BlocksCreations"), Description: "If true, indicates that branch creation is a protected operation."},
			{Name: "protection_rule_creator_login", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.BranchProtectionRule.Creator.Login"), Description: "The login of the user whom created the branch protection rule."},
			{Name: "protection_rule_dismisses_stale_reviews", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.DismissesStaleReviews"), Description: "If true, new commits pushed to matching branches dismiss pull request review approvals."},
			{Name: "protection_rule_lock_allows_fetch_and_merge", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.LockAllowsFetchAndMerge"), Description: "If true, users can pull changes from upstream when the branch is locked."},
			{Name: "protection_rule_lock_branch", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.LockBranch"), Description: "If true, matching branches are read-only and cannot be pushed to."},
			{Name: "protection_rule_pattern", Type: proto.ColumnType_STRING, Transform: transform.FromField("Node.BranchProtectionRule.Pattern"), Description: "The protection rule pattern."},
			{Name: "protection_rule_require_last_push_approval", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequireLastPushApproval"), Description: "If true, he most recent push must be approved by someone other than the person who pushed it."},
			{Name: "protection_rule_requires_approving_reviews", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresApprovingReviews"), Description: "If true, approving reviews required to update matching branches."},
			{Name: "protection_rule_required_approving_review_count", Type: proto.ColumnType_INT, Transform: transform.FromField("Node.BranchProtectionRule.RequiredApprovingReviewCount"), Description: "Number of approving reviews required to update matching branches."},
			{Name: "protection_rule_requires_conversation_resolution", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresConversationResolution"), Description: "If true, requires all comments on the pull request to be resolved before it can be merged to a protected branch."},
			{Name: "protection_rule_requires_code_owner_reviews", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresCodeOwnerReviews"), Description: "If true, reviews from code owners are required to update matching branches."},
			{Name: "protection_rule_requires_commit_signatures", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresCommitSignatures"), Description: "If true, commits are required to be signed by verified signatures."},
			{Name: "protection_rule_requires_deployments", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresDeployments"), Description: "If true, matching branches require deployment to specific environments before merging."},
			{Name: "protection_rule_required_deployment_environments", Type: proto.ColumnType_JSON, Transform: transform.FromField("Node.BranchProtectionRule.RequiredDeploymentEnvironments"), Description: "List of required deployment environments that must be deployed successfully to update matching branches."},
			{Name: "protection_rule_requires_linear_history", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresLinearHistory"), Description: "If true, prevent merge commits from being pushed to matching branches."},
			{Name: "protection_rule_requires_status_checks", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresStatusChecks"), Description: "If true, status checks are required to update matching branches."},
			{Name: "protection_rule_required_status_checks", Type: proto.ColumnType_JSON, Transform: transform.FromField("Node.BranchProtectionRule.RequiredStatusChecks"), Description: "Status checks that must pass before a branch can be merged into branches matching this rule."},
			{Name: "protection_rule_requires_strict_status_checks", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RequiresStrictStatusChecks"), Description: "If true, branches required to be up to date before merging."},
			{Name: "protection_rule_restricts_pushes", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Node.BranchProtectionRule.RestrictsPushes"), Description: "If true, pushing to matching branches is restricted."},
		},
	}
}

func tableGitHubBranchList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client := connectV4(ctx, d)

	fullName := d.EqualsQuals["repository_full_name"].GetStringValue()
	owner, repo := parseRepoFullName(fullName)

	pageSize := adjustPageSize(100, d.QueryContext.Limit)

	var query struct {
		RateLimit  models.RateLimit
		Repository struct {
			Refs struct {
				TotalCount int
				PageInfo   models.PageInfo
				Edges      []struct {
					Node models.Branch
				}
			} `graphql:"refs(refPrefix: \"refs/heads/\", first: $pageSize, after: $cursor)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":    githubv4.String(owner),
		"repo":     githubv4.String(repo),
		"pageSize": githubv4.Int(pageSize),
		"cursor":   (*githubv4.String)(nil),
	}

	for {
		err := client.Query(ctx, &query, variables)
		plugin.Logger(ctx).Debug(rateLimitLogString("github_branch", &query.RateLimit))
		if err != nil {
			plugin.Logger(ctx).Error("github_branch", "api_error", err)
			return nil, err
		}

		for _, branch := range query.Repository.Refs.Edges {
			d.StreamListItem(ctx, branch)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if !query.Repository.Refs.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.Repository.Refs.PageInfo.EndCursor)
	}

	return nil, nil
}

// Note: if useful to other tables, move to utils.go
func HasValue(_ context.Context, input *transform.TransformData) (interface{}, error) {
	if input.Value == nil || input.Value.(string) == "" {
		return false, nil
	}

	return true, nil
}
