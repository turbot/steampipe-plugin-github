# Table: github_commit

Commits for a repository.

The `github_commit` table can be used to query information about any commit, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### Recent commits

```sql
select
  sha,
  author_login,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  authored_date desc;
```

### Commits by a given author

```sql
select
  sha,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
  and author_login = 'e-gineer'
order by
  authored_date desc;
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
  authored_date
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
and
  signature is null
order by
  authored_date desc;
```

### Commits with most file changes

```sql
select
  sha,
  message,
  author_login,
  changed_files,
  additions,
  deletions
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  changed_files desc;
```
