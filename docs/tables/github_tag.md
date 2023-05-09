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
  string_to_array(regexp_replace(name, '[^0-9\.]', '', 'g'), '.'),
  name;
```

### Get commit details for each tag

```sql
select
  name,
  commit_sha,
  commit_short_sha,
  commit_date,
  commit_author,
  commit_message
from
  github_tag
where
  t.repository_full_name = 'turbot/steampipe';
```
