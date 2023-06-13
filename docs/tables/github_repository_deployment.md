# Table: github_repository_deployment

The `github_repository_deployment` table can be used to query deployments for a specific repository.

**You must specify `repository_full_name` in the WHERE or JOIN clause.**

## Examples

### List all deployments for a repository

```sql
select
  id,
  node_id,
  sha,
  created_at,
  creator ->> 'login' as creator_login,
  commit_sha,
  description,
  environment,
  latest_status,
  payload,
  ref ->> 'prefix' as ref_prefix,
  ref ->> 'name' as ref_name,
  state,
  task,
  updated_at
from
  github_repository_deployment
where
  repository_full_name = 'turbot/steampipe';
```