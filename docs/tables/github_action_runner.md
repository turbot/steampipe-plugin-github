# Table: github_action_runner

A runner is a server that runs your workflows when they're triggered. Each runner can run a single job at a time. Self-hosted runners offer more control of hardware, operating system, and software tools than GitHub-hosted runners provide.

The `github_action_runner` table can be used to query information about any self-hosted runner, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List runners

```sql
select
  *
from
  github_action_runner
where
  repository_full_name = 'turbot/steampipe';
```

### List runners with mac operating system

```sql
select
  repository_full_name,
  id,
  name,
  os
from
  github_action_runner
where
  repository_full_name = 'turbot/steampipe' and os = 'macos';
```

### List runners which are in use currently

```sql
select
  repository_full_name,
  id,
  name,
  os,
  busy
from
  github_action_runner
where
  repository_full_name = 'turbot/steampipe' and busy;
```