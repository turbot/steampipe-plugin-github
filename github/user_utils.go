package github

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractUserFromHydrateItem(h *plugin.HydrateData) (models.UserWithCounts, error) {
	if user, ok := h.Item.(models.UserWithCounts); ok {
		return user, nil
	} else if teamMember, ok := h.Item.(models.TeamMemberWithRole); ok {
		return models.UserWithCounts{User: teamMember.Node}, nil
	} else if orgMember, ok := h.Item.(memberWithRole); ok {
		return models.UserWithCounts{User: orgMember.Node}, nil
	} else {
		return models.UserWithCounts{}, fmt.Errorf("unable to parse hydrate item %v as an User", h.Item)
	}
}

// With Fine-grained access token we are getting field error even though we have proper access. For the column "assignees" column table "github_issue". so we need to avoid this error
// https://spec.graphql.org/October2021/#sec-Errors.Field-errors
// https://spec.graphql.org/October2021/#sec-Handling-Field-Errors
func appendUserInteractionAbilityForIssue(m *map[string]interface{}, cols []string, d *plugin.QueryData) {
	githubConfig := GetConfig(d.Connection)
	token := os.Getenv("GITHUB_TOKEN")
	if slices.Contains(cols, "assignees") && (strings.HasPrefix(token, "github_pat") || strings.HasPrefix(*githubConfig.Token, "github_pat")) {
		(*m)["includeUserInteractionAbility"] = githubv4.Boolean(false)
	} else {
		(*m)["includeUserInteractionAbility"] = githubv4.Boolean(true)
	}
}

func appendUserColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeUserAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeUserAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeUserBio"] = githubv4.Boolean(slices.Contains(cols, "bio"))
	(*m)["includeUserCompany"] = githubv4.Boolean(slices.Contains(cols, "company"))
	(*m)["includeUserEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeUserHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeUserInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeUserIsBountyHunter"] = githubv4.Boolean(slices.Contains(cols, "is_bounty_hunter"))
	(*m)["includeUserIsCampusExpert"] = githubv4.Boolean(slices.Contains(cols, "is_campus_expert"))
	(*m)["includeUserIsDeveloperProgramMember"] = githubv4.Boolean(slices.Contains(cols, "is_developer_program_member"))
	(*m)["includeUserIsEmployee"] = githubv4.Boolean(slices.Contains(cols, "is_employee"))
	(*m)["includeUserIsFollowingYou"] = githubv4.Boolean(slices.Contains(cols, "is_following_you"))
	(*m)["includeUserIsGitHubStar"] = githubv4.Boolean(slices.Contains(cols, "is_github_star"))
	(*m)["includeUserIsHireable"] = githubv4.Boolean(slices.Contains(cols, "is_hireable"))
	(*m)["includeUserIsSiteAdmin"] = githubv4.Boolean(slices.Contains(cols, "is_site_admin"))
	(*m)["includeUserIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeUserIsYou"] = githubv4.Boolean(slices.Contains(cols, "is_you"))
	(*m)["includeUserLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeUserMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeUserPinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeUserProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeUserPronouns"] = githubv4.Boolean(slices.Contains(cols, "pronouns"))
	(*m)["includeUserSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeUserStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeUserTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeUserCanChangedPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "can_changed_pinned_items"))
	(*m)["includeUserCanCreateProjects"] = githubv4.Boolean(slices.Contains(cols, "can_create_projects"))
	(*m)["includeUserCanFollow"] = githubv4.Boolean(slices.Contains(cols, "can_follow"))
	(*m)["includeUserCanSponsor"] = githubv4.Boolean(slices.Contains(cols, "can_sponsor"))
	(*m)["includeUserIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeUserIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeUserWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))
}

func appendUserWithCountColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeUserAnyPinnableItems"] = githubv4.Boolean(slices.Contains(cols, "any_pinnable_items"))
	(*m)["includeUserAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeUserBio"] = githubv4.Boolean(slices.Contains(cols, "bio"))
	(*m)["includeUserCompany"] = githubv4.Boolean(slices.Contains(cols, "company"))
	(*m)["includeUserEstimatedNextSponsorsPayoutInCents"] = githubv4.Boolean(slices.Contains(cols, "estimated_next_sponsors_payout_in_cents"))
	(*m)["includeUserHasSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "has_sponsors_listing"))
	(*m)["includeUserInteractionAbility"] = githubv4.Boolean(slices.Contains(cols, "interaction_ability"))
	(*m)["includeUserIsBountyHunter"] = githubv4.Boolean(slices.Contains(cols, "is_bounty_hunter"))
	(*m)["includeUserIsCampusExpert"] = githubv4.Boolean(slices.Contains(cols, "is_campus_expert"))
	(*m)["includeUserIsDeveloperProgramMember"] = githubv4.Boolean(slices.Contains(cols, "is_developer_program_member"))
	(*m)["includeUserIsEmployee"] = githubv4.Boolean(slices.Contains(cols, "is_employee"))
	(*m)["includeUserIsFollowingYou"] = githubv4.Boolean(slices.Contains(cols, "is_following_you"))
	(*m)["includeUserIsGitHubStar"] = githubv4.Boolean(slices.Contains(cols, "is_github_star"))
	(*m)["includeUserIsHireable"] = githubv4.Boolean(slices.Contains(cols, "is_hireable"))
	(*m)["includeUserIsSiteAdmin"] = githubv4.Boolean(slices.Contains(cols, "is_site_admin"))
	(*m)["includeUserIsSponsoringYou"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring_you"))
	(*m)["includeUserIsYou"] = githubv4.Boolean(slices.Contains(cols, "is_you"))
	(*m)["includeUserLocation"] = githubv4.Boolean(slices.Contains(cols, "location"))
	(*m)["includeUserMonthlyEstimatedSponsorsIncomeInCents"] = githubv4.Boolean(slices.Contains(cols, "monthly_estimated_sponsors_income_in_cents"))
	(*m)["includeUserPinnedItemsRemaining"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_remaining"))
	(*m)["includeUserProjectsUrl"] = githubv4.Boolean(slices.Contains(cols, "projects_url"))
	(*m)["includeUserPronouns"] = githubv4.Boolean(slices.Contains(cols, "pronouns"))
	(*m)["includeUserSponsorsListing"] = githubv4.Boolean(slices.Contains(cols, "sponsors_listing"))
	(*m)["includeUserStatus"] = githubv4.Boolean(slices.Contains(cols, "status"))
	(*m)["includeUserTwitterUsername"] = githubv4.Boolean(slices.Contains(cols, "twitter_username"))
	(*m)["includeUserCanChangedPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "can_changed_pinned_items"))
	(*m)["includeUserCanCreateProjects"] = githubv4.Boolean(slices.Contains(cols, "can_create_projects"))
	(*m)["includeUserCanFollow"] = githubv4.Boolean(slices.Contains(cols, "can_follow"))
	(*m)["includeUserCanSponsor"] = githubv4.Boolean(slices.Contains(cols, "can_sponsor"))
	(*m)["includeUserIsFollowing"] = githubv4.Boolean(slices.Contains(cols, "is_following"))
	(*m)["includeUserIsSponsoring"] = githubv4.Boolean(slices.Contains(cols, "is_sponsoring"))
	(*m)["includeUserWebsiteUrl"] = githubv4.Boolean(slices.Contains(cols, "website_url"))

	(*m)["includeUserRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_disk_usage"))
	(*m)["includeUserFollowers"] = githubv4.Boolean(slices.Contains(cols, "followers_total_count"))
	(*m)["includeUserFollowing"] = githubv4.Boolean(slices.Contains(cols, "following_total_count"))
	(*m)["includeUserPublicRepositories"] = githubv4.Boolean(slices.Contains(cols, "public_repositories_total_count"))
	(*m)["includeUserPrivateRepositories"] = githubv4.Boolean(slices.Contains(cols, "private_repositories_total_count"))
	(*m)["includeUserPublicGists"] = githubv4.Boolean(slices.Contains(cols, "public_gists_total_count"))
	(*m)["includeUserIssues"] = githubv4.Boolean(slices.Contains(cols, "issues_total_count"))
	(*m)["includeUserOrganizations"] = githubv4.Boolean(slices.Contains(cols, "organizations_total_count"))
	(*m)["includeUserPublicKeys"] = githubv4.Boolean(slices.Contains(cols, "public_keys_total_count"))
	(*m)["includeUserOpenPullRequests"] = githubv4.Boolean(slices.Contains(cols, "open_pull_requests_total_count"))
	(*m)["includeUserMergedPullRequests"] = githubv4.Boolean(slices.Contains(cols, "merged_pull_requests_total_count"))
	(*m)["includeUserClosedPullRequests"] = githubv4.Boolean(slices.Contains(cols, "closed_pull_requests_total_count"))
	(*m)["includeUserPackages"] = githubv4.Boolean(slices.Contains(cols, "packages_total_count"))
	(*m)["includeUserPinnedItems"] = githubv4.Boolean(slices.Contains(cols, "pinned_items_total_count"))
	(*m)["includeUserSponsoring"] = githubv4.Boolean(slices.Contains(cols, "sponsoring_total_count"))
	(*m)["includeUserSponsors"] = githubv4.Boolean(slices.Contains(cols, "sponsors_total_count"))
	(*m)["includeUserStarredRepositories"] = githubv4.Boolean(slices.Contains(cols, "starred_repositories_total_count"))
	(*m)["includeUserWatching"] = githubv4.Boolean(slices.Contains(cols, "watching_total_count"))
}

func userHydrateAnyPinnableItems(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.AnyPinnableItems, nil
}

func userHydrateAvatarUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.AvatarUrl, nil
}

func userHydrateBio(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Bio, nil
}

func userHydrateCompany(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Company, nil
}

func userHydrateEstimatedNextSponsorsPayoutInCents(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.EstimatedNextSponsorsPayoutInCents, nil
}

func userHydrateHasSponsorsListing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.HasSponsorsListing, nil
}

func userHydrateInteractionAbility(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.InteractionAbility, nil
}

func userHydrateIsBountyHunter(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsBountyHunter, nil
}

func userHydrateIsCampusExpert(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsCampusExpert, nil
}

func userHydrateIsDeveloperProgramMember(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsDeveloperProgramMember, nil
}

func userHydrateIsEmployee(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsEmployee, nil
}

func userHydrateIsFollowingYou(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsFollowingYou, nil
}

func userHydrateIsGitHubStar(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsGitHubStar, nil
}

func userHydrateIsHireable(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsHireable, nil
}

func userHydrateIsSiteAdmin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsSiteAdmin, nil
}

func userHydrateIsSponsoringYou(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsSponsoringYou, nil
}

func userHydrateIsYou(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsYou, nil
}

func userHydrateLocation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Location, nil
}

func userHydrateMonthlyEstimatedSponsorsIncomeInCents(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.MonthlyEstimatedSponsorsIncomeInCents, nil
}

func userHydratePinnedItemsRemaining(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PinnedItemsRemaining, nil
}

func userHydrateProjectsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.ProjectsUrl, nil
}

func userHydratePronouns(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Pronouns, nil
}

func userHydrateSponsorsListing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.SponsorsListing, nil
}

func userHydrateStatus(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Status, nil
}

func userHydrateTwitterUsername(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.TwitterUsername, nil
}

func userHydrateCanChangedPinnedItems(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.CanChangedPinnedItems, nil
}

func userHydrateCanCreateProjects(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.CanCreateProjects, nil
}

func userHydrateCanFollow(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.CanFollow, nil
}

func userHydrateCanSponsor(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.CanSponsor, nil
}

func userHydrateIsFollowing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsFollowing, nil
}

func userHydrateIsSponsoring(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.IsSponsoring, nil
}

func userHydrateWebsiteUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.WebsiteUrl, nil
}

func userHydrateRepositoriesTotalDiskUsage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Repositories.TotalDiskUsage, nil
}

func userHydrateFollowersTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Followers.TotalCount, nil
}

func userHydrateFollowingTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Following.TotalCount, nil
}

func userHydratePublicRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PublicRepositories.TotalCount, nil
}

func userHydratePrivateRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PrivateRepositories.TotalCount, nil
}

func userHydratePublicGistsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PublicGists.TotalCount, nil
}

func userHydrateIssuesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Issues.TotalCount, nil
}

func userHydrateOrganizationsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Organizations.TotalCount, nil
}

func userHydratePublicKeysTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PublicKeys.TotalCount, nil
}

func userHydrateOpenPullRequestsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.OpenPullRequests.TotalCount, nil
}

func userHydrateMergedPullRequestsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.MergedPullRequests.TotalCount, nil
}

func userHydrateClosedPullRequestsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.ClosedPullRequests.TotalCount, nil
}

func userHydratePackagesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Packages.TotalCount, nil
}

func userHydratePinnedItemsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.PinnedItems.TotalCount, nil
}

func userHydrateSponsoringTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Sponsoring.TotalCount, nil
}

func userHydrateSponsorsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Sponsors.TotalCount, nil
}

func userHydrateStarredRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.StarredRepositories.TotalCount, nil
}

func userHydrateWatchingTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user, err := extractUserFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return user.Watching.TotalCount, nil
}
