# Table: github_team_repository

A repository contains all of your project's files and each file's revision history.

The `github_team_repository` table can be used to query information about repositories that a team have access to. **You must specify the organization and team slug** in the where or join clause (`where organization= AND slug=`, `join github_team_repository on organization= AND slug=`).

To list all **your** repositories use the `github_my_repository` table instead. To get information about **any** repository, use the `github_repository` table instead.

## Examples

### Get information about a specific organization's team

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count,
  description
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team';
```

### To get repositories for all teams

```sql
select
  t.id          as team_id,
  t.name        as team_name,
  t.privacy     as team_privacy,
  t.description as team_description,
  tr.name       as repo_name,
  tr.fork       as repo_is_fork,
  tr.private    as repo_is_private,
  tr.archived   as repo_is_archived,
  tr.language   as repo_primary_language
from
  github.github_team as t,
  github.github_team_repository AS tr
where
  t.organization = 'my_org'
  and t.organization = tr.organization
  and t.slug = tr.slug
```
