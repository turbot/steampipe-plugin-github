# Table: github_workflow

A workflow is a configurable automated process made up of one or more jobs.

The `github_workflow` table can be used to query information about any workflow, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List workflows

```sql
select
  *
from
  github_workflow
where
  repository_full_name = 'turbot/steampipe';
```

### List build jobs in workflows for a specific repository

```sql
with pipelines as (
  select
    name,
    repository_full_name,
    pipeline
  from
    github_workflow
  where
    repository_full_name = 'turbot/steampipe'
)
select distinct
  p.repository_full_name,
  p.name as workflow_name,
  j ->> 'name' as job_name
from
  pipelines as p,
  jsonb_array_elements(pipeline -> 'jobs') as j
where
  (j -> 'metadata' -> 'build')::bool
```
