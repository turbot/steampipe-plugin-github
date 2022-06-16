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
  name,
  permissions,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  description
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team';
```

### List visible teams and repositories they have admin permissions to

```sql
select
  t.organization as organization,
  t.name as team_name,
  t.slug as team_slug,
  t.privacy as team_privacy,
  t.description as team_description,
  tr.name as repo_name,
  tr.permissions as team_repo_permissions,
  tr.fork as repo_is_fork,
  tr.private as repo_is_private,
  tr.archived as repo_is_archived,
  tr.language as repo_primary_language
from
  github_team as t,
  github_team_repository AS tr
where
  t.organization = tr.organization
  and t.slug = tr.slug
  and permissions ? 'admin';
```
