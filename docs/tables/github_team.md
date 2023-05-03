# Table: github_team

Teams are groups of organization members that reflect your company or group's structure with cascading access permissions and mentions. The `github_team` table lists all teams you have visibility to across your organizations.

To list the teams that you're a member of across your organizations, use the `github_my_team` table.

## Examples

## List all visible teams

```sql
select
  name,
  slug,
  privacy,
  description
from
  github_team;
```

## List all visible teams in an organization

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

## Get the number of members for a single team

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

### Get parent team details for child teams

```sql
select
  slug,
  organization,
  parent ->> 'DatabaseId' as parent_team_id,
  parent ->> 'Slug' as parent_team_slug
from
  github_team
where
  parent is not null;
```

### Get the number of repositories for a single team

```sql
select
  name,
  slug,
  repos_count
from
  github_team
where
  organization = 'my_org'
and
  slug = 'my_team';
```