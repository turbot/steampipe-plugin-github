# Table: github_tag

Tags mark specific commits in a repository history.

The `github_tag` table can be used to query information about any tag, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List tags

```sql
select
  name,
  commit ->> 'sha' as commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```

### Order tags by semantic version

```sql
select
  name,
  commit ->> 'sha' as commit_sha
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
  commit ->> 'sha' as commit_sha,
  commit ->> 'message' as commit_message,
  commit ->> 'url' as commit_url,
  commit -> 'author' -> 'user' ->> 'login' as author,
  commit ->> 'authored_date' as authored_date,
  commit -> 'committer' -> 'user' ->> 'login' as committer,
  commit ->> 'committed_date' as committed_date,
  commit ->> 'additions' as additions,
  commit ->> 'deletions' as deletions,
  commit ->> 'changed_files' as changed_files,
  commit -> 'signature' ->> 'is_valid' as commit_signed,
  commit -> 'signature' ->> 'email' as commit_signature_email,
  commit -> 'signature' -> 'signer' ->> 'login' as commit_signature_login,
  commit ->> 'tarball_url' as tarball_url,
  commit ->> 'zipball_url' as zipball_url 
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```
