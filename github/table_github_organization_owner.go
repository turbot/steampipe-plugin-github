package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"github.com/turbot/steampipe-plugin-github/github/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"strings"
)

func gitHubOrganizationOwnerColumns() []*plugin.Column {
	ownerCols := []*plugin.Column{
		{Name: "ip_allow_list_enabled_setting", Type: proto.ColumnType_STRING, Transform: transform.FromField("IpAllowListEnabledSetting", "Node.IpAllowListEnabledSetting"), Description: "The setting value for whether the organization has an IP allow list enabled."},
		{Name: "ip_allow_list_for_installed_apps_enabled_setting", Type: proto.ColumnType_STRING, Transform: transform.FromField("IpAllowListForInstalledAppsEnabledSetting", "Node.IpAllowListForInstalledAppsEnabledSetting"), Description: "The setting value for whether the organization has IP allow list configuration for installed GitHub Apps enabled."},
		{Name: "members_can_fork_private_repositories", Type: proto.ColumnType_BOOL, Transform: transform.FromField("MembersCanForkPrivateRepositories", "Node.MembersCanForkPrivateRepositories"), Description: "If true, members can fork private repositories in this organization."},
		{Name: "organization_billing_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("OrganizationBillingEmail", "Node.OrganizationBillingEmail"), Description: "The billing email for the organization."},
		{Name: "notification_delivery_restriction_enabled_setting", Type: proto.ColumnType_STRING, Transform: transform.FromField("NotificationDeliveryRestrictionEnabledSetting", "Node.NotificationDeliveryRestrictionEnabledSetting"), Description: "Indicates if email notification delivery for this organization is restricted to verified or approved domains."},
		{Name: "requires_two_factor_authentication", Type: proto.ColumnType_BOOL, Transform: transform.FromField("RequiresTwoFactorAuthentication", "Node.RequiresTwoFactorAuthentication"), Description: "If true, the organization requires all members, billing managers, and outside collaborators to enable two-factor authentication."},
		{Name: "web_commit_signoff_required", Type: proto.ColumnType_BOOL, Transform: transform.FromField("WebCommitSignoffRequired", "Node.WebCommitSignoffRequired"), Description: "Whether contributors are required to sign off on web-based commits for repositories in this organization."},
	}

	return append(gitHubOrganizationColumns(), ownerCols...)
}

func tableGitHubOrganizationOwner() *plugin.Table {
	return &plugin.Table{
		Name:        "github_organization_owner",
		Description: "An extended version of the github_organization table that returns extra data for organization owners.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("login"),
			Hydrate:    tableGitHubOrganizationOwnerList,
		},
		Columns: gitHubOrganizationOwnerColumns(),
	}
}

func tableGitHubOrganizationOwnerList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	login := d.EqualsQuals["login"].GetStringValue()

	var query struct {
		RateLimit    models.RateLimit
		Organization models.OrganizationWithOwnerPropertiesAndCounts `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}

	dummyRow := models.BasicOrganization{
		Login:       login,
		Description: "INSUFFICIENT_PERMISSIONS",
	}

	client := connectV4(ctx, d)
	err := client.Query(ctx, &query, variables)
	plugin.Logger(ctx).Debug(rateLimitLogString("github_organization", &query.RateLimit))
	if err != nil {
		plugin.Logger(ctx).Error("github_organization", "api_error", err)
		if strings.Contains(err.Error(), "does not have the right permission to retrieve") {
			d.StreamListItem(ctx, dummyRow)
			return nil, nil
		}
		if strings.Contains(err.Error(), "Could not resolve to an Organization with the login of") {
			return nil, nil
		}
		return nil, err
	}

	d.StreamListItem(ctx, query.Organization)

	return nil, nil
}
