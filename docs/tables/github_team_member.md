# Table: github_team_member

The `github_team_member` table can be used to query information about members of a team. **You must specify the organization and team slug** in the where or join clause (`where organization= AND slug=`, `join github_team_member on organization= AND slug=`).

## Examples

### List team members for a specific team

```sql
select
  organization,
  slug as team_slug,
  login,
  role,
  status_message
from
  github_team_member
where
  organization = 'my_org'
  and slug = 'my-team';
```

### List active team members with maintainer role for a specific team

```sql
select
  organization,
  slug as team_slug,
  login,
  role,
  status_message
from
  github_team_member
where
  organization = 'my_org'
  and slug = 'my-team'
  and role = 'maintainer';
```

### List team members with maintainer role for visible teams

```sql
select
  t.organization as organization,
  t.name as team_name,
  t.slug as team_slug,
  t.privacy as team_privacy,
  t.description as team_description,
  tm.login as member_login,
  tm.role as member_role
from
  github_team as t,
  github_team_member as tm
where
  t.organization = tm.organization
  and t.slug = tm.slug
  and tm.role = 'maintainer';
```
