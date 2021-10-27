# Table: github_tag

Tags mark specific commits in a repository history.

The `github_tag` table can be used to query information about any tag, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List tags

```sql
select
  name,
  commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```

### Order tags by semantic version

```sql
select
  name,
  commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe'
order by
  string_to_array(regexp_replace(name, '[^0-9\.]', '', 'g'), '.')::int[],
  name;
```

### Get commit details for each tag (Not working yet)

Note: This example is intended to join tags with commit information to return the
full details of each item. Currently joins with multiple columns are not
working pending a solution to [#47](https://github.com/turbot/steampipe-postgres-fdw/issues/47).

```sql
select
  t.name,
  t.commit_sha,
  c.author_date,
  c.message
from
  github_tag as t,
  github_commit as c
where
  t.repository_full_name = 'turbot/steampipe'
  and t.repository_full_name = c.repository_full_name
  and t.commit_sha = c.sha
order by
  c.author_date desc;
```
