package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractStarFromHydrateItem(h *plugin.HydrateData) (myStar, error) {
	if star, ok := h.Item.(myStar); ok {
		return star, nil
	} else {
		return myStar{}, fmt.Errorf("unable to parse hydrate item %v as an Star", h.Item)
	}
}

func appendStarColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeStarNode"] = githubv4.Boolean(slices.Contains(cols, "repository_full_name") || slices.Contains(cols, "url"))
	(*m)["includeStarEdges"] = githubv4.Boolean(slices.Contains(cols, "starred_at"))
}

func starHydrateNameWithOwner(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	star, err := extractStarFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return star.NameWithOwner, nil
}

func starHydrateStarredAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	star, err := extractStarFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return star.StarredAt, nil
}

func starHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	star, err := extractStarFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return star.Url, nil
}
