# Table: github_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions.


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
  github_team;
```


### Permission for the teams

```sql
select
  name,
  organization,
  permission
from
  github_team;
```