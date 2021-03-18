# Table: github_my_repository

A repository contains all of your project's files and each file's revision history.

You can own repositories individually, or you can share ownership of repositories with other people in an organization.  The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.

To query **ANY** repository, including public repos, use the `github_repository` table.

## Examples

### List of repositories that you or your organizations own or contribute to

```sql
select
  name,
  owner_login,
  full_name
from
  github_my_repository
order by
  full_name;
```

### Show repository stats

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
  github_my_repository;
```

### List your public repositories

```sql
select
  name,
  private,
  visibility,
  owner_login
from
  github_my_repository
where
  not private;
```
OR

```sql
select
  name,
  private,
  visibility
from
  github_my_repository
where
  visibility = 'public';
```

### List all your repositories and their collaborators 

```sql
select
  name,
  collaborator_logins
from
  github_my_repository;
```

### List collaborators and their permissions in your repositories

```sql
select
  name,
  c ->> 'login' as login,
  c -> 'permissions' -> 'pull' as can_pull,
  c -> 'permissions' -> 'push' as can_push,
  c -> 'permissions' -> 'admin' as is_admin
from
  github_my_repository,
  jsonb_array_elements(collaborators) as c
order by
  name,
  c ->> 'login';
```


### List collaborators who have "push" or "admin" to a specific repository

In this case, collaborators who have "push" or "admin" to the `turbot/steampipe-plugin-aws repository`:
```sql
select
  name,
  c ->> 'login' as login,
  c -> 'permissions' -> 'pull' as can_pull,
  c -> 'permissions' -> 'push' as can_push,
  c -> 'permissions' -> 'admin' as is_admin
from
  github_my_repository,
  jsonb_array_elements(collaborators) as c
where
  name = 'steampipe-plugin-aws'
  and owner_login = 'turbot'
  and (
    (c -> 'permissions' -> 'admin') :: bool
    or (c -> 'permissions' -> 'push') :: bool
  );
  ```


### List collaborators for organization repositories that are not organization members

In this case, for the `turbot` org:
```sql
select
  name,
  owner_login as owner,
  c ->> 'login' as login,
  c -> 'permissions' -> 'pull' as can_pull,
  c -> 'permissions' -> 'push' as can_push,
  c -> 'permissions' -> 'admin' as is_admin
from
  github_my_repository as r,
  jsonb_array_elements(collaborators) as c
where
  owner_login = 'turbotio'
  and c ->> 'login' not in (
    select
      m ->> 'login' as member_login
    from
      github_organization,
      jsonb_array_elements(members) as m
    where
      login = 'turbotio'
  );
```
