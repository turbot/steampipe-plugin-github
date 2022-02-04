# Table: github_actions_organization_secret

Secrets are encrypted environment variables that you create in an organization.

The `github_actions_organization_secret` table can be used to query information about any organization secret, and **you must specify which organization** in the where or join clause using the `organization_name` column.

## Examples

### List secrets

```sql
select
  *
from
  github_actions_organization_secret
where
  organization_name = 'turbot';
```