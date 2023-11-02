package github

import (
	"context"
	"fmt"
	"slices"

	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func extractLicenseFromHydrateItem(h *plugin.HydrateData) (models.License, error) {
	if license, ok := h.Item.(models.License); ok {
		return license, nil
	} else {
		return models.License{}, fmt.Errorf("unable to parse hydrate item %v as an License", h.Item)
	}
}

func appendLicenseColumnIncludes(m *map[string]interface{}, cols []string) {
	(*m)["includeLicenseName"] = githubv4.Boolean(slices.Contains(cols, "name"))
	(*m)["includeLicenseSpdxId"] = githubv4.Boolean(slices.Contains(cols, "spdx_id"))
	(*m)["includeLicenseUrl"] = githubv4.Boolean(slices.Contains(cols, "url"))
	(*m)["includeLicenseConditions"] = githubv4.Boolean(slices.Contains(cols, "conditions"))
	(*m)["includeLicenseDescription"] = githubv4.Boolean(slices.Contains(cols, "description"))
	(*m)["includeLicenseFeatured"] = githubv4.Boolean(slices.Contains(cols, "featured"))
	(*m)["includeLicenseHidden"] = githubv4.Boolean(slices.Contains(cols, "hidden"))
	(*m)["includeLicenseImplementation"] = githubv4.Boolean(slices.Contains(cols, "implementation"))
	(*m)["includeLicenseLimitations"] = githubv4.Boolean(slices.Contains(cols, "limitations"))
	(*m)["includeLicensePermissions"] = githubv4.Boolean(slices.Contains(cols, "permissions"))
	(*m)["includeLicenseNickname"] = githubv4.Boolean(slices.Contains(cols, "nickname"))
	(*m)["includeLicensePseudoLicense"] = githubv4.Boolean(slices.Contains(cols, "pseudo_license"))
}

func licenseHydrateSpdxId(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.SpdxId, nil
}

func licenseHydrateName(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Name, nil
}

func licenseHydrateUrl(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Url, nil
}

func licenseHydrateConditions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Conditions, nil
}

func licenseHydrateDescription(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Description, nil
}

func licenseHydrateFeatured(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Featured, nil
}

func licenseHydrateHidden(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Hidden, nil
}

func licenseHydrateImplementation(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Implementation, nil
}

func licenseHydrateLimitations(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Limitations, nil
}

func licenseHydratePermissions(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Permissions, nil
}

func licenseHydrateNickname(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.Nickname, nil
}

func licenseHydratePseudoLicense(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	license, err := extractLicenseFromHydrateItem(h)
	if err != nil {
		return nil, err
	}
	return license.PseudoLicense, nil
}
