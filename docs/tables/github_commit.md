# Table: github_commit

Commits for a repository.

The `github_commit` table can be used to query information about any commit, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### Recent commits

```sql
select
  sha,
  author_login,
  author_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  author_date desc;
```

### Commits by a given author

```sql
select
  sha,
  author_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
  and author_login = 'e-gineer'
order by
  author_date desc;
```

### Contributions by author

```sql
select
  author_login,
  count(*)
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  count desc;
```

### Commits that were not verified

```sql
select
  sha,
  author_login,
  author_date
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
  and not verified
order by
  author_date desc;
```
