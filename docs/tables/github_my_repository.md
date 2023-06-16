# Table: github_my_repository

A repository contains all of your project's files and each file's revision history.

You can own repositories individually, or you can share ownership of repositories with other people in an organization. The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.

To query **ANY** repository, including public repos, use the `github_repository` table.

## Examples

### List of repositories that you or your organizations own or contribute to

```sql
select
  name,
  owner_login,
  name_with_owner
from
  github_my_repository
order by
  name_with_owner;
```

### Show repository stats

```sql
select
  name,
  owner_login,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  subscribers_count,
  watchers_total_count,
  updated_at as last_updated,
  description
from
  github_my_repository;
```

### List your public repositories

```sql
select
  name,
  is_private,
  visibility,
  owner_login
from
  github_my_repository
where
  not is_private;
```

OR

```sql
select
  name,
  is_private,
  visibility
from
  github_my_repository
where
  visibility = 'PUBLIC';
```

### List all your repositories and their collaborators

```sql
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
 ,github_repository_collaborator c
where
  r.name_with_owner = c.repository_full_name;
```

### List all your repository collaborators with admin or maintainer permissions

```sql
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
 ,github_repository_collaborator c
where
  r.name_with_owner = c.repository_full_name
and
  permission in ('ADMIN', 'MAINTAIN');
```

### List repository hooks that are insecure

```sql
select
  name as repository,
  hook
from
  github_my_repository,
  jsonb_array_elements(hooks) as hook
where
  hook -> 'config' ->> 'insecure_ssl' = '1'
    or hook -> 'config' ->> 'secret' is null
    or hook -> 'config' ->> 'url' not like '%https:%';
```