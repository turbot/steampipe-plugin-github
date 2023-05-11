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
  commit_authored_date,
  commit_author_login,
  commit_committed_date,
  commit_committer_login,
  commit_message,
  commit_url,
  commit_additions,
  commit_deletions,
  commit_changed_files,
  commit_committed_via_web,
  commit_signature_is_valid,
  commit_signature_email,
  commit_signature_login,
  commit_tarball_url,
  commit_zipball_url,
  commit_tree_url,
  commit_status
from
  github_tag
where
  t.repository_full_name = 'turbot/steampipe';
```
