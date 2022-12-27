# Table: github_repository_dependabot_alert

The `github_repository_dependabot_alert` table can be used to query information about dependabot alerts from a repository.

**You must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List dependabot alerts

```sql
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe';
```

### List open dependabot alerts

```sql
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open';
```

### List open critical dependabot alerts

```sql
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open'
  and security_advisory_severity = 'critical';
```
