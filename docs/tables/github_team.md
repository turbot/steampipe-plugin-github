---
title: "Steampipe Table: github_team - Query GitHub Teams using SQL"
description: "Allows users to query GitHub Teams, specifically providing details about each team within a GitHub organization. This information includes team ID, name, description, privacy level, and more."
folder: "Team"
---

# Table: github_team - Query GitHub Teams using SQL

GitHub Teams is a feature within GitHub that allows organizations to create teams, manage permissions, and simplify @mentions. Teams are groups of organization members that reflect the company or project's structure. They can be used to create nested teams, mentionable as a single unit, and provide a social graph of an organization's repo permissions.

## Table Usage Guide

The `github_team` table provides insights into the teams within GitHub organizations. As a project manager or team lead, you can explore team-specific details through this table, including team ID, name, description, and privacy level. Utilize it to manage permissions, simplify @mentions, and understand the social graph of your organization's repo permissions.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Organization permissions:
    - Members (Read-only): Required to access all columns.
  - Repository permissions:
    - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `organization` column in the `where` or `join` clause to query the table.

## Examples

### List all visible teams
Explore which teams are visible on your GitHub account, including details like their privacy settings and descriptions, to better manage your collaborations and understand team dynamics.

```sql+postgres
select
  name,
  slug,
  privacy,
  description
from
  github_team
where
  organization = 'turbot';
```

```sql+sqlite
select
  name,
  slug,
  privacy,
  description
from
  github_team
where
  organization = 'turbot';
```

### List all visible teams in an organization
Explore which teams are publicly visible within a specific organization on GitHub. This is useful for understanding the structure and privacy settings of your organization's teams.

```sql+postgres
select
  name,
  slug,
  privacy,
  description
from
  github_team
where
  organization = 'turbot';
```

```sql+sqlite
select
  name,
  slug,
  privacy,
  description
from
  github_team
where
  organization = 'turbot';
```

### Get the number of members for a single team
Explore the size of a specific team within your organization on Github. This can be useful for resource allocation and understanding team dynamics.

```sql+postgres
select
  name,
  slug,
  members_total_count
from
  github_team
where
  organization = 'my_org'
  and slug = 'my_team';
```

```sql+sqlite
select
  name,
  slug,
  members_total_count
from
  github_team
where
  organization = 'my_org'
  and slug = 'my_team';
```

### Get the number of repositories for a single team
Determine the total number of repositories associated with a specific team within your organization. This can be useful for understanding the team's workload or for assessing the distribution of resources within the organization.

```sql+postgres
select
  name,
  slug,
  repositories_total_count
from
  github_team
where
  organization = 'my_org'
  and slug = 'my_team';
```

```sql+sqlite
select
  name,
  slug,
  repositories_total_count
from
  github_team
where
  organization = 'my_org'
  and slug = 'my_team';
```

### Get parent team details for child teams
Determine the hierarchical relationships within your organization's teams on Github. This query is useful for understanding team structures and identifying which teams are sub-teams of larger, parent teams.

```sql+postgres
select
  slug,
  organization,
  parent_team ->> 'id' as parent_team_id,
  parent_team ->> 'node_id' as parent_team_node_id,
  parent_team ->> 'slug' as parent_team_slug
from
  github_team
where
  organization = 'turbot'
  and parent_team is not null;
```

```sql+sqlite
select
  slug,
  organization,
  parent_team ->> 'id' as parent_team_id,
  parent_team ->> 'node_id' as parent_team_node_id,
  parent_team ->> 'slug' as parent_team_slug
from
  github_team
where
  organization = 'turbot'
  and parent_team is not null;
```

### List teams with pending user invitations
Identify teams that have outstanding invitations to users. This can help manage and expedite the onboarding process by pinpointing where follow-ups may be needed.

```sql+postgres
select
  name,
  slug,
  invitations_total_count
from
  github_team
where
  organization = 'turbot'
  and invitations_total_count > 0;
```

```sql+sqlite
select
  name,
  slug,
  invitations_total_count
from
  github_team
where
  organization = 'turbot'
  and invitations_total_count > 0;
```