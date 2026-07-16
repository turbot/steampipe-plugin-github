package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractProjectV2FromHydrateItem(h *plugin.HydrateData) (models.ProjectV2, error) {
	if project, ok := h.Item.(models.ProjectV2); ok {
		return project, nil
	} else {
		return models.ProjectV2{}, fmt.Errorf("unable to parse hydrate item %v as a ProjectV2", h.Item)
	}
}
func appendProjectV2ColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeId"] = githubv4.Boolean(slices.Contains(cols, "id"))
	(*m)["includeNodeId"] = githubv4.Boolean(slices.Contains(cols, "node_id"))
	(*m)["includeOwner"] = githubv4.Boolean(slices.Contains(cols, "owner"))
	(*m)["includeCreator"] = githubv4.Boolean(slices.Contains(cols, "creator"))
	(*m)["includeTitle"] = githubv4.Boolean(slices.Contains(cols, "title"))
	(*m)["includeDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeIsPublic"] = githubv4.Boolean(slices.Contains(cols, "is_public"))
	(*m)["includeClosedAt"] = githubv4.Boolean(slices.Contains(cols, "closed_at"))
	(*m)["includeCreatedAt"] = githubv4.Boolean(slices.Contains(cols, "created_at"))
	(*m)["includeUpdatedAt"] = githubv4.Boolean(slices.Contains(cols, "updated_at"))
	(*m)["includeState"] = githubv4.Boolean(slices.Contains(cols, "state"))
	(*m)["includeLatestStatusUpdate"] = githubv4.Boolean(slices.Contains(cols, "latest_status_update"))
	(*m)["includeIsTemplate"] = githubv4.Boolean(slices.Contains(cols, "is_template"))
}

func projectV2HydrateId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.Id, nil
}

func projectV2HydrateNodeId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.NodeId, nil
}

func projectV2HydrateOwner(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.Owner, nil
}

func projectV2HydrateCreator(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.Creator, nil
}

func projectV2HydrateTitle(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.Title, nil
}

func projectV2HydrateDescription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.Description, nil
}

func projectV2HydrateIsPublic(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.IsPublic, nil
}

func projectV2HydrateClosedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.ClosedAt, nil
}

func projectV2HydrateCreatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.CreatedAt, nil
}

func projectV2HydrateUpdatedAt(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.UpdatedAt, nil
}

// projectV2HydrateState derives the REST-compatible "open"/"closed" state string from the GraphQL boolean.
func projectV2HydrateState(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	if project.Closed {
		return "closed", nil
	}
	return "open", nil
}

// projectV2HydrateLatestStatusUpdate returns the most recent status update from the statusUpdates connection.
func projectV2HydrateLatestStatusUpdate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	if len(project.LatestStatusUpdate.Nodes) > 0 {
		return project.LatestStatusUpdate.Nodes[0], nil
	}
	return nil, nil
}

func projectV2HydrateIsTemplate(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project, err := extractProjectV2FromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return project.IsTemplate, nil
}
