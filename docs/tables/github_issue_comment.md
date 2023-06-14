# Table: github_issue_comment

The `github_issue_comment` table can be used to query comments from a specific issue.

**You must specify `repository_full_name` (repository including org/user prefix) and `number` (of the issue) in the WHERE or JOIN clause.**

## Examples

### List comments for a specific issue

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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201;
```

### List comments for your issues

```sql
select
  c.repository_full_name,
  c.number as issue,
  i.author_login as issue_author,
  c.author_login as comment_author,
  c.body_text as content,
  c.created_at as created,
  c.url
from 
  github_issue_comment c
join 
    github_my_issue i
on 
  i.repository_full_name = c.repository_full_name
and 
  i.number = c.number;
```

### List comments for a specific issue which match a certain body content

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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201
and
  body_text ~~* '%branch%';
```