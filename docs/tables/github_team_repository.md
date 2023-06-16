# Table: github_team_repository

A repository contains all of your project's files and each file's revision history.

The `github_team_repository` table can be used to query information about repositories that a team has access to. **You must specify the organization and team slug** in the where or join clause (`where organization= AND slug=`, `join github_team_repository on organization= AND slug=`).

To list all **your** repositories use the `github_my_repository` table instead. To get information about **any** repository, use the `github_repository` table instead.

## Examples

### List a specific team's repositories

```sql
select
  organization,
  slug as team_slug,
  name as team_name,
  permission,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  license_info ->> 'spdx_id' as license,
  description,
  url
from
  github_team_repository
where
  organization = 'my_org'
and 
  slug = 'my-team';
```

### List visible teams and repositories they have admin permissions to

```sql
select
  organization,
  slug as team_slug,
  name as name,
  description,
  permission,
  is_fork,
  is_private,
  is_archived,
  primary_language ->> 'name' as language
from
  github_team_repository
where
  organization = 'my_org'
and 
  slug = 'my-team'
and 
  permission = 'ADMIN';
```
