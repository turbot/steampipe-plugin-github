package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func appendContributionColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeContributionCalendar"] = githubv4.Boolean(slices.Contains(cols, "contribution_calendar"))
	(*m)["includeCommitContributionsByRepository"] = githubv4.Boolean(slices.Contains(cols, "commit_contributions_by_repository"))
}

func extractContributionsCollectionFromHydrateItem(h *plugin.HydrateData) (models.ContributionsCollection, error) {
	if collection, ok := h.Item.(models.ContributionsCollection); ok {
		return collection, nil
	}
	return models.ContributionsCollection{}, fmt.Errorf("unable to parse hydrate item %v as ContributionsCollection", h.Item)
}

func contributionHydrateCalendar(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	collection, err := extractContributionsCollectionFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return collection.ContributionCalendar, nil
}

func contributionHydrateCommitContributionsByRepository(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	collection, err := extractContributionsCollectionFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return collection.CommitContributionsByRepository, nil
}
