package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractCommunityProfileFromHydrateItem(h *plugin.HydrateData) (models.CommunityProfile, error) {
	if cp, ok := h.Item.(models.CommunityProfile); ok {
		return cp, nil
	} else {
		return models.CommunityProfile{}, fmt.Errorf("unable to parse hydrate item %v as an CommunityProfile", h.Item)
	}
}

func appendCommunityProfileColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeCPLicense"] = githubv4.Boolean(slices.Contains(cols, "license_info"))
	(*m)["includeCPCodeOfConduct"] = githubv4.Boolean(slices.Contains(cols, "code_of_conduct"))
	(*m)["includeCPIssueTemplates"] = githubv4.Boolean(slices.Contains(cols, "issue_templates"))
	(*m)["includeCPPullRequestTemplates"] = githubv4.Boolean(slices.Contains(cols, "pull_request_templates"))
	(*m)["includeCPReadme"] = githubv4.Boolean(slices.Contains(cols, "readme"))
	(*m)["includeCPContributing"] = githubv4.Boolean(slices.Contains(cols, "contributing"))
	(*m)["includeCPSecurity"] = githubv4.Boolean(slices.Contains(cols, "security"))
}

func cpHydrateLicense(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return cp.LicenseInfo, nil
}

func cpHydrateCodeOfConduct(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return cp.CodeOfConduct, nil
}

func cpHydrateIssueTemplates(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return cp.IssueTemplates, nil
}

func cpHydratePullRequestTemplates(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return cp.PullRequestTemplates, nil
}

func cpHydrateReadme(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	if cp.ReadMeUpper.Blob.Text != "" {
		return cp.ReadMeUpper.Blob, nil
	} else {
		return cp.ReadMeLower.Blob, nil
	}
}

func cpHydrateContributing(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	if cp.ContributingUpper.Blob.Text != "" {
		return cp.ContributingUpper.Blob, nil
	} else if cp.ContributingLower.Blob.Text != "" {
		return cp.ContributingLower.Blob, nil
	} else {
		return cp.ContributingTitle.Blob, nil
	}
}

func cpHydrateSecurity(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cp, err := extractCommunityProfileFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	if cp.SecurityUpper.Blob.Text != "" {
		return cp.SecurityUpper.Blob, nil
	} else if cp.SecurityLower.Blob.Text != "" {
		return cp.SecurityLower.Blob, nil
	} else {
		return cp.SecurityTitle.Blob, nil
	}
}
