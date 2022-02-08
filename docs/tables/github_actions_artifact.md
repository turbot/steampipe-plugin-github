# Table: github_actions_artifact

Artifacts allow you to persist data after a job has completed, and share that data with another job in the same workflow. An artifact is a file or collection of files produced during a workflow run.

The `github_actions_artifact` table can be used to query information about any artifacts, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List artifacts

```sql
select
  *
from
  github_actions_artifact
where
  repository_full_name = 'turbot/steampipe';
```

### List active artifacts

```sql
select
  id,
  node_id,
  name,
  archive_download_url,
  expired
from
  github_actions_artifact
where
  repository_full_name = 'turbot/steampipe' and not expired;
```