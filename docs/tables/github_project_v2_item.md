---
title: "Steampipe Table: github_project_v2_item - Query GitHub Project (V2) items using SQL"
description: "Allows users to query the items in GitHub Projects (V2), including issues, pull requests and draft issues, with their project field values such as Status."
folder: "Project"
---

# Table: github_project_v2_item - Query GitHub Project (V2) items using SQL

GitHub Projects (V2) track work as items. Each item is an issue, a pull request or a draft issue (a note that only exists in the project), along with the values of the project's fields, such as Status or Iteration.

## Table Usage Guide

The `github_project_v2_item` table provides insights into the items in a GitHub Project (V2) owned by an organization or user. As a project manager or developer, explore what's on a project board through this table, including each item's type, status and other custom field values, and the details of the underlying issue, pull request or draft issue. Utilize it to report on project progress, such as which items are in progress, which issues are still open, or what's planned for the current iteration.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Organization permissions:
    - Projects (Read-only): Required to access all columns.

**Important Notes**
- You must specify the `project_number` column, and either the `organization` or the `user_login` column, in a `where` or `join` clause to query the table.
- Each row has a `type` of `ISSUE`, `PULL_REQUEST`, `DRAFT_ISSUE` or `REDACTED`. Columns that don't apply to an item's type are null, e.g. draft issues have no `repository_full_name` or `number`. `REDACTED` items are items the token can't access, and only have the item columns (`id`, `type`, `created_at`, etc.) set.
- The `field_values` column contains the values of the project's custom fields (Text, Number, Date, Single select, Multi select and Iteration), keyed by field name. Built-in fields, such as Title, Assignees, Labels and Milestone, are available as their own columns.
- Archived items are included. Use `is_archived = false` to only list the items shown on the project board.
- This table supports optional quals. Queries with optional quals are optimised to use GitHub's project item filters. Optional quals are supported for the following columns:
  - `type`
  - `is_archived`
  - `repository_full_name`

## Examples

### List the items in a project
Explore the issues, pull requests and draft issues in a project, along with their status, to get an overview of the work being tracked.

```sql+postgres
select
  type,
  number,
  title,
  state,
  field_values ->> 'Status' as status,
  url
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1;
```

```sql+sqlite
select
  type,
  number,
  title,
  state,
  json_extract(field_values, '$.Status') as status,
  url
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1;
```

### List the items with a specific status
Identify the items that are currently in progress, to see what the team is working on.

```sql+postgres
select
  type,
  number,
  title,
  repository_full_name,
  assignees
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and field_values ->> 'Status' = 'In Progress';
```

```sql+sqlite
select
  type,
  number,
  title,
  repository_full_name,
  assignees
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and json_extract(field_values, '$.Status') = 'In Progress';
```

### Count the items in each status
Summarize the number of items in each status, to understand the overall progress of a project.

```sql+postgres
select
  field_values ->> 'Status' as status,
  count(*) as items
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and not is_archived
group by
  status
order by
  items desc;
```

```sql+sqlite
select
  json_extract(field_values, '$.Status') as status,
  count(*) as items
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and is_archived = 0
group by
  status
order by
  items desc;
```

### List the open issues in a project
Explore the open issues in a project, along with their labels and assignees, to help with triage.

```sql+postgres
select
  number,
  title,
  repository_full_name,
  labels,
  assignees_total_count,
  field_values ->> 'Status' as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and type = 'ISSUE'
  and state = 'OPEN';
```

```sql+sqlite
select
  number,
  title,
  repository_full_name,
  labels,
  assignees_total_count,
  json_extract(field_values, '$.Status') as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and type = 'ISSUE'
  and state = 'OPEN';
```

### List the draft issues in a project
Discover the draft issues in a project, which haven't been converted to issues yet, along with who created them.

```sql+postgres
select
  title,
  author ->> 'login' as author,
  content_created_at,
  field_values ->> 'Status' as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and type = 'DRAFT_ISSUE';
```

```sql+sqlite
select
  title,
  json_extract(author, '$.login') as author,
  content_created_at,
  json_extract(field_values, '$.Status') as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and type = 'DRAFT_ISSUE';
```

### List the items from a repository
Explore the items in a project that belong to a specific repository.

```sql+postgres
select
  type,
  number,
  title,
  state,
  field_values ->> 'Status' as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  type,
  number,
  title,
  state,
  json_extract(field_values, '$.Status') as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and repository_full_name = 'turbot/steampipe';
```

### List the items in the current iteration
Identify the items planned for the current iteration, based on an iteration field named `Iteration`.

```sql+postgres
select
  type,
  number,
  title,
  field_values -> 'Iteration' ->> 'title' as iteration,
  field_values ->> 'Status' as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and (field_values -> 'Iteration' ->> 'start_date')::date <= current_date
  and (field_values -> 'Iteration' ->> 'start_date')::date + (field_values -> 'Iteration' ->> 'duration')::int > current_date;
```

```sql+sqlite
select
  type,
  number,
  title,
  json_extract(field_values, '$.Iteration.title') as iteration,
  json_extract(field_values, '$.Status') as status
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and date(json_extract(field_values, '$.Iteration.start_date')) <= date('now')
  and date(json_extract(field_values, '$.Iteration.start_date'), '+' || json_extract(field_values, '$.Iteration.duration') || ' days') > date('now');
```

### List the archived items in a project
Discover the items that have been archived from a project.

```sql+postgres
select
  type,
  number,
  title,
  updated_at
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and is_archived;
```

```sql+sqlite
select
  type,
  number,
  title,
  updated_at
from
  github_project_v2_item
where
  organization = 'turbot'
  and project_number = 1
  and is_archived = 1;
```

### List the in progress items across all open projects in an organization
Explore what's in progress across every open project in an organization. The `materialized` keyword makes Postgres list the projects first, then the items in each project. Specify the organization on both tables, since it isn't passed from the `projects` query to the items.

```sql+postgres
with projects as materialized (
  select
    number,
    title
  from
    github_project_v2
  where
    organization = 'turbot'
    and state = 'open'
)
select
  p.number as project_number,
  p.title as project_title,
  i.type,
  i.title,
  i.url
from
  projects p
  join github_project_v2_item i on i.project_number = p.number
where
  i.organization = 'turbot'
  and i.field_values ->> 'Status' = 'In Progress';
```

```sql+sqlite
select
  p.number as project_number,
  p.title as project_title,
  i.type,
  i.title,
  i.url
from
  github_project_v2 p
  join github_project_v2_item i on i.project_number = p.number
where
  p.organization = 'turbot'
  and p.state = 'open'
  and i.organization = 'turbot'
  and json_extract(i.field_values, '$.Status') = 'In Progress';
```

### Get the state reason of the closed issues in a project
Explore why each closed issue in a project was closed, by joining with the `github_issue` table for details that aren't in this table. The `materialized` keyword makes Postgres list the project items first, then get each issue.

```sql+postgres
with items as materialized (
  select
    repository_full_name,
    number,
    title,
    field_values
  from
    github_project_v2_item
  where
    organization = 'turbot'
    and project_number = 1
    and type = 'ISSUE'
    and state = 'CLOSED'
)
select
  i.number,
  i.title,
  i.field_values ->> 'Status' as status,
  gi.state_reason
from
  items i
  join github_issue gi on gi.repository_full_name = i.repository_full_name and gi.number = i.number;
```

```sql+sqlite
select
  i.number,
  i.title,
  json_extract(i.field_values, '$.Status') as status,
  gi.state_reason
from
  github_project_v2_item i
  join github_issue gi on gi.repository_full_name = i.repository_full_name and gi.number = i.number
where
  i.organization = 'turbot'
  and i.project_number = 1
  and i.type = 'ISSUE'
  and i.state = 'CLOSED';
```

### Get the review status of the open pull requests in a project
Explore the review decision of each open pull request in a project, by joining with the `github_pull_request` table.

```sql+postgres
with items as materialized (
  select
    repository_full_name,
    number,
    title
  from
    github_project_v2_item
  where
    organization = 'turbot'
    and project_number = 1
    and type = 'PULL_REQUEST'
    and state = 'OPEN'
)
select
  i.number,
  i.title,
  pr.review_decision,
  pr.is_draft
from
  items i
  join github_pull_request pr on pr.repository_full_name = i.repository_full_name and pr.number = i.number;
```

```sql+sqlite
select
  i.number,
  i.title,
  pr.review_decision,
  pr.is_draft
from
  github_project_v2_item i
  join github_pull_request pr on pr.repository_full_name = i.repository_full_name and pr.number = i.number
where
  i.organization = 'turbot'
  and i.project_number = 1
  and i.type = 'PULL_REQUEST'
  and i.state = 'OPEN';
```

### List the items in a project owned by a user
Explore the items in a project owned by a specific GitHub user, e.g. one of your own personal projects.

```sql+postgres
select
  type,
  number,
  title,
  field_values ->> 'Status' as status
from
  github_project_v2_item
where
  user_login = 'octocat'
  and project_number = 1;
```

```sql+sqlite
select
  type,
  number,
  title,
  json_extract(field_values, '$.Status') as status
from
  github_project_v2_item
where
  user_login = 'octocat'
  and project_number = 1;
```
