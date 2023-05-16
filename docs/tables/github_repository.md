# Table: github_repository

A repository contains all of your project's files and each file's revision history.

The `github_repository` table can be used to query information about **ANY** repository, and **you must specify which repository** in the where or join clause (`where full_name=`, `join github_repository on full_name=`).

To list all of **your** repositories use the `github_my_repository` table instead. The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.

## Examples

### Get information about a specific repository

```sql
select
  name,
  node_id,
  id,
  created_at,
  updated_at,
  disk_usage,
  owner_login,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  url,
  license_info ->> 'spdx_id' as license,
  description
from
  github_repository
where
  full_name = 'postgres/postgres';
```

### Get your permissions for a specific repository

```sql
select
  name,
  your_permission,
  can_administer,
  can_create_projects,
  can_subscribe,
  can_update_topics,
  possible_commit_emails
from
  github_repository
where
  full_name = 'turbot/steampipe';
```