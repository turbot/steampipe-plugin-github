---
title: "Steampipe Table: github_team_member - Query GitHub Team Members using SQL"
description: "Allows users to query GitHub Team Members, providing insights into team member's details and their roles in specific GitHub teams."
---

# Table: github_team_member - Query GitHub Team Members using SQL

GitHub Team Members are part of GitHub's Team feature which allows users to organize and manage user groups within a repository. These Team Members can be assigned different roles and permissions within the team. This feature is part of GitHub's collaboration tools designed to facilitate project management and team coordination.

## Table Usage Guide

The `github_team_member` table provides insights into team members within GitHub. As a project manager or team lead, explore team member-specific details through this table, including their roles and permissions within teams. Utilize it to manage and coordinate your team effectively, ensuring that each team member has the appropriate permissions and roles for their tasks.

**Important Notes**
- You must specify the `organization` and `slug` columns in the `where` or `join` clause to query the table.

## Examples

### List team members for a specific team
Explore which team members belong to a specific team in your organization. This is useful for understanding team composition and roles within your organization.

```sql+postgres
select
  organization,
  slug as team_slug,
  login,
  role,
  status
from
  github_team_member
where
  organization = 'my_org'
and
  slug = 'my-team';
```

```sql+sqlite
select
  organization,
  slug as team_slug,
  login,
  role,
  status
from
  github_team_member
where
  organization = 'my_org'
and
  slug = 'my-team';
```

### List active team members with maintainer role for a specific team
This query helps to identify active team members who hold the 'Maintainer' role within a specific team in your organization. It is useful for managing team roles and ensuring that every team has an active maintainer.

```sql+postgres
select
  organization,
  slug as team_slug,
  login,
  role,
  status
from
  github_team_member
where
  organization = 'my_org'
and 
  slug = 'my-team'
and 
  role = 'MAINTAINER';
```

```sql+sqlite
select
  organization,
  slug as team_slug,
  login,
  role,
  status
from
  github_team_member
where
  organization = 'my_org'
and 
  slug = 'my-team'
and 
  role = 'MAINTAINER';
```

### List team members with maintainer role for visible teams
Discover the segments that consist of team members with a maintainer role in visible teams. This can be useful for understanding the distribution of roles within your organization's teams, particularly in identifying those who have the authority to manage team settings.

```sql+postgres
select
  t.organization as organization,
  t.name as team_name,
  t.slug as team_slug,
  t.privacy as team_privacy,
  t.description as team_description,
  tm.login as member_login,
  tm.role as member_role,
  tm.status as member_status
from
  github_team as t,
  github_team_member as tm
where
  t.organization = tm.organization
  and t.slug = tm.slug
  and tm.role = 'MAINTAINER';
```

```sql+sqlite
select
  t.organization as organization,
  t.name as team_name,
  t.slug as team_slug,
  t.privacy as team_privacy,
  t.description as team_description,
  tm.login as member_login,
  tm.role as member_role,
  tm.status as member_status
from
  github_team as t
join 
  github_team_member as tm
on
  t.organization = tm.organization
  and t.slug = tm.slug
where
  tm.role = 'MAINTAINER';
```