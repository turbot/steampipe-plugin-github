# Table: github_user

The `github_user` table does not list all users via the API - there is not currently an efficient way to limit the results in a useable way. As a result, you **must specify a user `login` in a `where`** or you will get no results.

## Examples

### Get information for a user

```sql
select
  *
from
  github_user
where
  login = 'torvalds';
```

### List of users in your organizations

```sql
select
  u.login,
  o.login as organization,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_organization as o,
  jsonb_array_elements_text(o.member_logins) as member_login
where
  u.login = member_login;
```

### List of users that collaborate on a repository that you own

```sql
select
  r.full_name as repository,
  u.login,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_repository as r,
  jsonb_array_elements_text(r.collaborator_logins) as collaborator_login
where
  u.login = collaborator_login
  and r.full_name = 'turbot/steampipe';
```
