package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractStargazerFromHydrateItem(h *plugin.HydrateData) (Stargazer, error) {
	if str, ok := h.Item.(Stargazer); ok {
		return str, nil
	} else {
		return Stargazer{}, fmt.Errorf("unable to parse hydrate item %v as a Stargazer", h.Item)
	}
}

func appendStargazerColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeStargazerStarredAt"] = githubv4.Boolean(slices.Contains(cols, "starred_at"))
	(*m)["includeStargazerNode"] = githubv4.Boolean(slices.Contains(cols, "user_login") || slices.Contains(cols, "user_detail"))
}

func strHydrateStarredAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	str, err := extractStargazerFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return str.StarredAt, nil
}

func strHydrateUserLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	str, err := extractStargazerFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return str.Node.Login, nil
}

func strHydrateUser(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	str, err := extractStargazerFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return str.Node, nil
}