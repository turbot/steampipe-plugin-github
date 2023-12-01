package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractTagFromHydrateItem(h *plugin.HydrateData) (tagRow, error) {
	if tag, ok := h.Item.(tagRow); ok {
		return tag, nil
	} else {
		return tagRow{}, fmt.Errorf("unable to parse hydrate item %v as an Tag", h.Item)
	}
}

func appendTagColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeTagTarget"] = githubv4.Boolean(slices.Contains(cols, "tagger_date") || slices.Contains(cols, "tagger_name") || slices.Contains(cols, "tagger_login") || slices.Contains(cols, "message") || slices.Contains(cols, "commit"))
	(*m)["includeTagName"] = githubv4.Boolean(slices.Contains(cols, "name"))
}

func tagHydrateName(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.Name, nil
}

func tagHydrateTaggerName(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.TaggerName, nil
}

func tagHydrateTaggerDate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.TaggerDate, nil
}

func tagHydrateTaggerLogin(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.TaggerLogin, nil
}

func tagHydrateMessage(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.Message, nil
}

func tagHydrateCommit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	tag, err := extractTagFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return tag.Commit, nil
}