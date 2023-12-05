package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractRateLimitFromHydrateItem(h *plugin.HydrateData) (models.BaseRateLimit, error) {
	if rl, ok := h.Item.(models.BaseRateLimit); ok {
		return rl, nil
	} else {
		return models.BaseRateLimit{}, fmt.Errorf("unable to parse hydrate item %v as an RateLimit", h.Item)
	}
}

func appendRateLimitColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeRLRemaining"] = githubv4.Boolean(slices.Contains(cols, "remaining"))
	(*m)["includeRLUsed"] = githubv4.Boolean(slices.Contains(cols, "used"))
	(*m)["includeRLCost"] = githubv4.Boolean(slices.Contains(cols, "cost"))
	(*m)["includeRLLimit"] = githubv4.Boolean(slices.Contains(cols, "limit"))
	(*m)["includeRLResetAt"] = githubv4.Boolean(slices.Contains(cols, "reset_at"))
	(*m)["includeRLNodeCount"] = githubv4.Boolean(slices.Contains(cols, "node_count"))
}

func rateLimitHydrateRemaining(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.Remaining, nil
}

func rateLimitHydrateUsed(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.Used, nil
}

func rateLimitHydrateCost(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.Cost, nil
}

func rateLimitHydrateLimit(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.Limit, nil
}

func rateLimitHydrateResetAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.ResetAt, nil
}

func rateLimitHydrateNodeCount(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	rl, err := extractRateLimitFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return rl.NodeCount, nil
}
