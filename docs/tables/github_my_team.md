---
title: "Steampipe Table: github_my_team - Query GitHub Teams using SQL"
description: "Allows users to query Teams in GitHub, specifically the details of teams that the authenticated user is a part of, providing insights into team structure, access rights and associated repositories."
folder: "Repository"
---

# Table: github_my_team - Query GitHub Teams using SQL

GitHub Teams is a feature within GitHub that allows for easy collaboration and access management within repositories. Teams can have different access rights to repositories and can consist of any number of users. Teams are a convenient way to manage large groups of users, both for assigning access rights and for mentioning multiple users at once.

## Table Usage Guide

The `github_my_team` table provides insights into Teams within GitHub. As a developer or project manager, explore team-specific details through this table, including access rights, team structure and associated repositories. Utilize it to uncover information about teams, such as those with admin access to repositories, the distribution of access rights within a team, and the verification of team members.

**Important Notes**
- To view **all teams you have visibility to across your organizations,** use the `github_team` table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Organization permissions:
    - Members (Read-only): Required to access all columns.
  - Repository permissions:
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### Basic info
Discover the segments that make up your GitHub team, including the number of members and repositories, to gain insights into the team's structure and activities. This allows for better management and allocation of resources within the team.

```sql+postgres
select
  name,
  slug,
  description,
  organization,
  members_total_count,
  repositories_total_count
from
  github_my_team;
```

```sql+sqlite
select
  name,
  slug,
  description,
  organization,
  members_total_count,
  repositories_total_count
from
  github_my_team;
```

### Get organization permission for each team
Explore which permissions each team in your organization has on GitHub. This can help you manage access controls and ensure the right teams have the right permissions.

```sql+postgres
select
  name,
  organization,
  privacy
from
  github_my_team;
```

```sql+sqlite
select
  name,
  organization,
  privacy
from
  github_my_team;
```

### Get parent team details for child teams
Determine the hierarchical structure within your GitHub organization by identifying which sub-teams have a parent team. This can help in understanding the team dynamics and collaboration structure within your organization.

```sql+postgres
select
  slug,
  organization,
  parent_team ->> 'id' as parent_team_id,
  parent_team ->> 'slug' as parent_team_slug
from
  github_my_team
where
  parent_team is not null;
```

```sql+sqlite
select
  slug,
  organization,
  json_extract(parent_team, '$.id') as parent_team_id,
  json_extract(parent_team, '$.slug') as parent_team_slug
from
  github_my_team
where
  parent_team is not null;
```