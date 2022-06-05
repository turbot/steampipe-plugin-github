# Table: github_team_member

The `github_team_member` table can be used to query information about members of a team. **You must specify the organization and team slug** in the where or join clause (`where organization= AND slug=`, `join github_repository on organization= AND slug=`).

## Examples

### Get information about a specific organization's team members

```sql
select
  login,
  role
from
  github_team_member
where
  organization = 'my_org'
  AND slug = 'my-team';
```

### To get members for all teams

```sql
select
  t.name        as team_name,
  t.privacy     as team_privacy,
  t.description as team_description,
  tm.login      as user_login,
  tm.role       as user_role
from
    github.github_team as t
inner join
    github.github_team_member as tm
    on t.organization = tm.organization
    and t.slug = tm.slug
where t.organization = 'my_org'
```
