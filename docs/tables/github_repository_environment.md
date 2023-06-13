# Table: github_repository_environment

The `github_repository_environment` table can be used to query environments belonging to a repository.

**You must specify `repository_full_name` in the WHERE or JOIN clause.**

## Examples

### List environments for a repository

```sql
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name = 'turbot/steampipe';
```

### List environments all your repositories

```sql
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name IN (select name_with_owner from github_my_repository);
```