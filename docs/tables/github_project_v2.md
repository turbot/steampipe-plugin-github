---
title: "Steampipe Table: github_project_v2 - Query GitHub Projects (V2) using SQL"
description: "Allows users to query GitHub Projects (V2), providing insights into the organization and user projects used to plan and track work."
folder: "Project"
---

# Table: github_project_v2 - Query GitHub Projects (V2) using SQL

GitHub Projects (V2) is GitHub's flexible, table and board based tool for planning and tracking work across issues and pull requests. It allows teams to organize work items, track status, and view progress across one or more repositories using customizable views, fields, and workflows.

## Table Usage Guide

The `github_project_v2` table provides insights into ProjectsV2 owned by a GitHub organization or user. As a project manager or developer, explore project-specific details through this table, including title, description, visibility, status, linked repositories, and linked teams. Utilize it to uncover information about projects, such as which ones are public, which repositories and teams are associated with them, and when they were last updated.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Organization permissions:
    - Projects (Read-only): Required to access all columns.
    - Members (Read-only): Required for the `teams` column to return team slugs; without it, GitHub redacts team nodes the token can't see and returns empty strings (e.g. `["", ""]`) instead of the real slugs.

**Important Notes**
- You must specify either the `organization` or the `user_login` column in a `where` or `join` clause to query the table (projects owned by a repository are not currently supported).

## Examples

### List the projects in an organization
Explore the title, state, and visibility of the ProjectsV2 owned by a specific GitHub organization to get an overview of ongoing work.

```sql+postgres
select
  organization,
  number,
  title,
  state,
  is_public,
  created_at
from
  github_project_v2
where
  organization = 'turbot';
```

```sql+sqlite
select
  organization,
  number,
  title,
  state,
  is_public,
  created_at
from
  github_project_v2
where
  organization = 'turbot';
```

### List open projects in an organization
Identify the projects that are still open in a specific organization, to help focus attention on active planning boards.

```sql+postgres
select
  organization,
  number,
  title,
  created_at,
  updated_at
from
  github_project_v2
where
  organization = 'turbot'
  and state = 'open';
```

```sql+sqlite
select
  organization,
  number,
  title,
  created_at,
  updated_at
from
  github_project_v2
where
  organization = 'turbot'
  and state = 'open';
```

### Get a specific project by number
Retrieve the details of a single project using its project number, useful when you already know which project you want to inspect.

```sql+postgres
select
  number,
  title,
  description,
  owner,
  creator
from
  github_project_v2
where
  organization = 'turbot'
  and number = 1;
```

```sql+sqlite
select
  number,
  title,
  description,
  owner,
  creator
from
  github_project_v2
where
  organization = 'turbot'
  and number = 1;
```

### List projects updated in the last 30 days
Discover the projects that have had recent activity, useful for tracking which planning boards are actively being maintained.

```sql+postgres
select
  number,
  title,
  updated_at
from
  github_project_v2
where
  organization = 'turbot'
  and updated_at >= now() - interval '30 days'
order by
  updated_at desc;
```

```sql+sqlite
select
  number,
  title,
  updated_at
from
  github_project_v2
where
  organization = 'turbot'
  and updated_at >= datetime('now', '-30 days')
order by
  updated_at desc;
```

### List repositories and teams linked to each project
Explore which repositories and teams are linked to each project, to understand the scope of collaboration around a project.

```sql+postgres
select
  number,
  title,
  repositories,
  repositories_total_count,
  teams,
  teams_total_count
from
  github_project_v2
where
  organization = 'turbot';
```

```sql+sqlite
select
  number,
  title,
  repositories,
  repositories_total_count,
  teams,
  teams_total_count
from
  github_project_v2
where
  organization = 'turbot';
```

### List the latest status update for each project
Explore the most recent status update posted on each project, useful for quickly checking the reported health and progress of a project.

```sql+postgres
select
  number,
  title,
  latest_status_update ->> 'status' as status,
  latest_status_update ->> 'body' as body,
  latest_status_update ->> 'created_at' as reported_at
from
  github_project_v2
where
  organization = 'turbot';
```

```sql+sqlite
select
  number,
  title,
  json_extract(latest_status_update, '$.status') as status,
  json_extract(latest_status_update, '$.body') as body,
  json_extract(latest_status_update, '$.created_at') as reported_at
from
  github_project_v2
where
  organization = 'turbot';
```

### List the projects owned by a user
Explore the title, state, and visibility of the ProjectsV2 owned by a specific GitHub user, e.g. your own personal projects.

```sql+postgres
select
  user_login,
  number,
  title,
  state,
  is_public,
  created_at
from
  github_project_v2
where
  user_login = 'octocat';
```

```sql+sqlite
select
  user_login,
  number,
  title,
  state,
  is_public,
  created_at
from
  github_project_v2
where
  user_login = 'octocat';
```
