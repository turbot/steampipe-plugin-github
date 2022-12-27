# Table: github_organization_dependabot_alert

The `github_organization_dependabot_alert` table can be used to query information about dependabot alerts from an organization. You must be an owner or security manager for the organization to successfully query dependabot alerts.

**You must specify the organization** in the where or join clause (`where organization=`, `join github_organization_depedanbot_alert on organization=`).

## Examples

### List dependabot alerts

```sql
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org';
```

### List open dependabot alerts

```sql
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open';
```

### List open critical dependabot alerts

```sql
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open'
  and security_advisory_severity = 'critical';
```
