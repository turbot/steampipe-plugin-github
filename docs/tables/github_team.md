# Table: github_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions. The `github_team` table lists all teams you have visibility to across your organizations.

To list the teams that you're a member of across your organizations, use the `github_my_team` table.

## Examples

## List all visible teams

```sql
select
  id,
  name,
  privacy,
  description
from
  github_team;
```

## List all visible teams in an organization

```sql
select
  id,
  name,
  privacy,
  description
from
  github_team
where
  organization = 'my_org';
```
