# Table: github_my_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.  The `github_my_team` table lists all teams in your organizations,


## Examples

### Basic Team Info

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


### Permission for the teams

```sql
select
  name,
  organization,
  permission
from
  github_my_team;
```