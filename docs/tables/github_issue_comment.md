# Table: github_issue_comment

The `github_issue_comment` table can be used to query comments from a specific issue.

**You must specify `repository_full_name` (repository including org/user prefix) and `number` (of the issue) in the WHERE or JOIN clause.**

## Examples

### List all comments for a specific issue

```sql
select
  id,
  node_id,
  author,
  author_login,
  author_association,
  body,
  body_text,
  created_at,
  updated_at,
  published_at,
  last_edited_at,
  created_via_email,
  editor,
  editor_login,
  includes_created_edit,
  is_minimized,
  minimized_reason,
  url,
  can_delete,
  can_minimize,
  can_react,
  can_update,
  cannot_update_reasons,
  did_author
from
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201;
```