# Table: github_organization_owner

You can use this table to retrieve all information about an Organization you have OWNER permissions on, but you must specify the `login` explicitly in the where or join clause.

This table is an extended variant of the `github_organization` table with the following additional fields:
- ip_allow_list_enabled_setting
- ip_allow_list_for_installed_apps_enabled_setting
- members_can_fork_private_repositories
- organization_billing_email
- notification_delivery_restriction_enabled_setting
- requires_two_factor_authentication
- web_commit_signoff_required

> Note: If you do not have OWNER permissions on the organization you will receive a dummy row of data.

## Examples

### Owner info for a GitHub Organization

```sql
select
  login as organization,
  name,
  ip_allow_list_enabled_setting,
  ip_allow_list_for_installed_apps_enabled_setting,
  members_can_fork_private_repositories,
  organization_billing_email,
  notification_delivery_restriction_enabled_setting,
  requires_two_factor_authentication,
  web_commit_signoff_required
from
  github_organization_owner
where
  login = 'turbot';
```