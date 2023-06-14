# Table: github_repository_deployment

The `github_repository_deployment` table can be used to query deployments for a specific repository.

**You must specify `repository_full_name` in the WHERE or JOIN clause.**

## Examples

### List deployments for a repository

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

### List deployments for all your repositories

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
  repository_full_name IN (select name_with_owner from github_my_repository);
```