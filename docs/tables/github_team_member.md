# Table: github_team_member

The `github_team_member` table can be used to query information about members of a team. **You must specify the organization and team slug** in the where or join clause (`where organization= AND slug=`, `join github_team_member on organization= AND slug=`).

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
  and slug = 'my-team';
```

### To get members for all teams and all organizations

```sql
select
  t.organization,
  t.name,
  t.privacy,
  t.description,
  tm.login,
  tm.role
from
  github_team as t,
  github_team_member as tm
where
  t.organization = tm.organization
  and t.slug = tm.slug
```
