# Table: github_repository_collaborator

A collaborator is a user who has permissions to contribute to a repository.

**You must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List all contributors of a repository

```sql
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe';
```

### List all outside collaborators on a repository

```sql
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
and
  filter = 'OUTSIDE';
```

### List all repository admins

```sql
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
and
  permission = 'ADMIN';
```

### Obtain a JSON array of admins for all your repositories

```sql
with repos as (
  select 
    name_with_owner 
  from 
    github_my_repository
)
select
  r.name_with_owner as repo,
  json_agg(user_login) as admins
from 
  repos as r
inner join 
  github_repository_collaborator as c
on 
  r.name_with_owner = c.repository_full_name
and
  c.permission = 'ADMIN'
group by 
  r.name_with_owner;
```
