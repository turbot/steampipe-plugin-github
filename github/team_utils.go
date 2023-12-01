package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractTeamFromHydrateItem(h *plugin.HydrateData) (models.TeamWithCounts, error) {
	if team, ok := h.Item.(models.TeamWithCounts); ok {
		return team, nil
	} else {
		return models.TeamWithCounts{}, fmt.Errorf("unable to parse hydrate item %v as an Team", h.Item)
	}
}

func appendTeamColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeTeamAvatarUrl"] = githubv4.Boolean(slices.Contains(cols, "avatar_url"))
	(*m)["includeTeamCombinedSlug"] = githubv4.Boolean(slices.Contains(cols, "combined_slug"))
	(*m)["includeTeamCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeTeamDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeTeamDiscussionsUrl"] = githubv4.Boolean(slices.Contains(cols, "discussions_url"))
	(*m)["includeTeamEditTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "edit_team_url"))
	(*m)["includeTeamMembersUrl"] = githubv4.Boolean(slices.Contains(cols, "members_url"))
	(*m)["includeTeamNewTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "new_team_url"))
	(*m)["includeTeamParentTeam"] = githubv4.Boolean(slices.Contains(cols, "parent_team"))
	(*m)["includeTeamPrivacy"] = githubv4.Boolean(slices.Contains(cols, "privacy"))
	(*m)["includeTeamRepositoriesUrl"] = githubv4.Boolean(slices.Contains(cols, "repositories_url"))
	(*m)["includeTeamTeamsUrl"] = githubv4.Boolean(slices.Contains(cols, "teams_url"))
	(*m)["includeTeamUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeTeamUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeTeamCanAdminister"] = githubv4.Boolean(slices.Contains(cols, "can_administer"))
	(*m)["includeTeamCanSubscribe"] = githubv4.Boolean(slices.Contains(cols, "can_subscribe"))
	(*m)["includeTeamSubscription"] = githubv4.Boolean(slices.Contains(cols, "subscription"))
	(*m)["includeTeamAncestors"] = githubv4.Boolean(slices.Contains(cols, "ancestors_total_count"))
	(*m)["includeTeamChildTeams"] = githubv4.Boolean(slices.Contains(cols, "child_teams_total_count"))
	(*m)["includeTeamDiscussions"] = githubv4.Boolean(slices.Contains(cols, "discussions_total_count"))
	(*m)["includeTeamInvitations"] = githubv4.Boolean(slices.Contains(cols, "invitations_total_count"))
	(*m)["includeTeamMembers"] = githubv4.Boolean(slices.Contains(cols, "members_total_count"))
	(*m)["includeTeamProjectsV2"] = githubv4.Boolean(slices.Contains(cols, "projects_v2_total_count"))
	(*m)["includeTeamRepositories"] = githubv4.Boolean(slices.Contains(cols, "repositories_total_count"))
}

func teamHydrateAvatarUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.AvatarUrl, nil
}

func teamHydrateCombinedSlug(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.CombinedSlug, nil
}

func teamHydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.CreatedAt, nil
}

func teamHydrateDescription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.Description, nil
}

func teamHydrateDiscussionsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.DiscussionsUrl, nil
}

func teamHydrateEditTeamUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.EditTeamUrl, nil
}

func teamHydrateMembersUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.MembersUrl, nil
}

func teamHydrateNewTeamUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.NewTeamUrl, nil
}

func teamHydrateParentTeam(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.ParentTeam, nil
}

func teamHydratePrivacy(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.Privacy, nil
}

func teamHydrateRepositoriesUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.RepositoriesUrl, nil
}

func teamHydrateTeamsUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.TeamsUrl, nil
}

func teamHydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.UpdatedAt, nil
}

func teamHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.Url, nil
}

func teamHydrateCanAdminister(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.CanAdminister, nil
}

func teamHydrateCanSubscribe(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.CanSubscribe, nil
}

func teamHydrateSubscription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	team, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return team.Subscription, nil
}

func teamHydrateAncestorsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.Ancestors.TotalCount, nil
}

func teamHydrateChildTeamsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.ChildTeams.TotalCount, nil
}

func teamHydrateDiscussionsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.Discussions.TotalCount, nil
}

func teamHydrateInvitationsTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.Invitations.TotalCount, nil
}

func teamHydrateMembersTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.Members.TotalCount, nil
}

func teamHydrateProjectsV2TotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.ProjectsV2.TotalCount, nil
}

func teamHydrateRepositoriesTotalCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	teamWithCounts, err := extractTeamFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return teamWithCounts.Repositories.TotalCount, nil
}
