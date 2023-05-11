# Table: github_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions. The `github_team` table lists all teams you have visibility to across your organizations.

To list the teams that you're a member of across your organizations, use the `github_my_team` table.

## Examples

### List all visible teams

```sql
select
  name,
  slug,
  privacy,
  description
from
  github_team;
```

### List all visible teams in an organization

```sql
select
  name,
  slug,
  privacy,
  description
from
  github_team
where
  organization = 'my_org';
```

### Get the number of members for a single team

```sql
select
  name,
  slug,
  members_count
from
  github_team
where
  organization = 'my_org'
and
  slug = 'my_team';
```

### Get the number of repositories for a single team

```sql
select
  name,
  slug,
  repositories_count
from
  github_team
where
  organization = 'my_org'
and
  slug = 'my_team';
```

### Get parent team details for child teams

```sql
select
  slug,
  organization,
  parent ->> 'id' as parent_team_id,
  parent ->> 'node_id' as parent_team_node_id,
  parent ->> 'slug' as parent_team_slug
from
  github_team
where
  parent is not null;
```

### List teams with pending user invitations

```sql
select
  name,
  slug,
  invitations_count
from
  github_team
where
  invitations_count > 0;
```