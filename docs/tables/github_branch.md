# Table: github_branch

A branch is essentially is a unique set of code changes with a unique name.

The `github_branch` table can be used to query information about any branch, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List branches

```sql
select
  name,
  commit_sha,
  protected
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

### Get commit details for each branch

```sql
select
  name,
  commit_sha,
  commit_date,
  commit_author,
  commit_message
from
  github_branch
where
  repository_full_name = 'turbot/steampipe'
order by
  commit_date desc;
```
