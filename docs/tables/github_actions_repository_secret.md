# Table: github_actions_repository_secret

Secrets are encrypted environment variables that you create in an organization.

The `github_actions_repository_secret` table can be used to query information about any organization secret, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List secrets

```sql
select
  *
from
  github_actions_repository_secret
where
  repository_full_name = 'turbot/steampipe';
```