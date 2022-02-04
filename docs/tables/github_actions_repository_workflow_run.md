# Table: github_actions_repository_workflow_run

WorkflowRun represents a repository action workflow run.

The `github_actions_repository_workflow_run` table can be used to query information about any organization secret, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List workflow runs

```sql
select
  *
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe';
```

### List failure workflow runs

```sql
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch
from
    github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and conclusion = 'failure';
```

### List manual workflow runs

```sql
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch
from
    github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and event = 'workflow_dispatch';
```