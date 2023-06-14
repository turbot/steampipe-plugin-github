# Table: github_pull_request_comment

The `github_pull_request_comment` table can be used to query comments from a specific pull request.

**You must specify `repository_full_name` (repository including org/user prefix) and `number` (of the issue) in the WHERE or JOIN clause.**

## Examples

### List all comments for a specific pull request

```sql
select
  id,
  author_login,
  author_association,
  body_text,
  created_at,
  updated_at,
  published_at,
  last_edited_at,
  editor_login,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 207;
```

### List comments for a specific pull request which match a certain body content

```sql
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body_text as content,
  created_at,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 207
and
  body_text ~~* '%DELAY%';
```

### List comments for all open pull requests from a specific repository
```sql
select
  c.*
from
  github_pull_request r
join
  github_pull_request_comment c
on
  r.repository_full_name = c.repository_full_name
and
  r.number = c.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
and
  r.state = 'OPEN';
```