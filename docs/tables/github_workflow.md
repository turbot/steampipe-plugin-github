# Table: github_workflow

A workflow is a configurable automated process made up of one or more jobs.

The `github_workflow` table can be used to query information about any workflow, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List workflows

```sql
select
  repository_full_name,
  name,
  path,
  text,
  line_count,
  size,
  language,
  node_id,
  is_truncated,
  is_generated,
  is_binary,
  text_json,
  pipeline
from
  github_workflow
where
  repository_full_name = 'turbot/steampipe';
```
