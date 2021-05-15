# Table: github_branch

Tags mark specific commits in a repository history.

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
  repository_full_name = 'turbot/steampipe'
```

### Get commit details for each branch (Not working yet)

Note: This example is intended to join branches with commit information to return the
full details of each item. Currently joins with multiple columns are not
working pending a solution to [#47](https://github.com/turbot/steampipe-postgres-fdw/issues/47).

```sql
select
  t.name,
  t.commit_sha,
  c.author_date,
  c.message
from
  github_branch as t,
  github_commit as c
where
  t.repository_full_name = 'turbot/steampipe'
  and t.repository_full_name = c.repository_full_name
  and t.commit_sha = c.sha
order by
  c.author_date desc;
```
