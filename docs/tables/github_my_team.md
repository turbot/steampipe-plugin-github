# Table: github_my_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions. The `github_my_team` table lists all teams you're a member of across your organizations.

To view **all teams you have visibility to across your organizations,** use the `github_team` table.

## Examples

### Basic info

```sql
select
  name,
  slug,
  description,
  organization,
  members_count,
  repos_count
from
  github_my_team;
```

### Get organization permission for each team

```sql
select
  name,
  organization,
  permission
from
  github_my_team;
```
