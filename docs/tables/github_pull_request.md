---
title: "Steampipe Table: github_pull_request - Query GitHub Pull Requests using SQL"
description: "Allows users to query GitHub Pull Requests, providing detailed insights into pull requests across repositories."
---

# Table: github_pull_request - Query GitHub Pull Requests using SQL

GitHub Pull Requests are a feature of GitHub, a web-based hosting service for version control using Git. Pull requests let you tell others about changes you've pushed to a branch in a repository on GitHub. Once a pull request is opened, you can discuss and review the potential changes with collaborators and add follow-up commits before your changes are merged into the base branch.

## Table Usage Guide

The `github_pull_request` table provides insights into pull requests within GitHub. As a developer or project manager, explore pull request-specific details through this table, including the status, assignees, reviewers, and associated metadata. Utilize it to track the progress of pull requests, identify bottlenecks in the review process, and ensure timely merging of approved changes.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) `where` or `join` clause to query the table.


## Examples

### List open pull requests in a repository
Determine the areas in which there are active discussions about code changes in a specific project. This is useful for project managers and contributors to track ongoing development efforts and understand the status of proposed modifications.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  mergeable
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  state,
  mergeable
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN';
```

### List the pull requests for a repository that have been closed in the last week
This example provides a way to monitor recent activity in a specific GitHub repository. It's particularly useful for project managers who want to track the progress of their projects by identifying which pull requests have been closed in the last week.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  closed_at,
  merged_at,
  merged_by
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'CLOSED'
  and closed_at >= (current_date - interval '7' day)
order by
  closed_at desc;
```

```sql+sqlite
## Runtime error: parsing time "2023-12-01" as "2006-01-02 15:04:05.999": cannot parse "" as "15"
select
  repository_full_name,
  number,
  title,
  state,
  closed_at,
  merged_at,
  merged_by
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'CLOSED'
  and closed_at >= date('now','-7 day')
order by
  closed_at desc;
```

### List the open PRs in a repository with a given label
Explore which open pull requests in a specific repository are tagged as 'bug'. This can help prioritize bug-fixing efforts and manage the project more efficiently.

```sql+postgres
select
  repository_full_name,
  number,
  state,
  labels
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe-plugin-aws'
  and labels -> 'bug' = 'true';
```

```sql+sqlite
select
  repository_full_name,
  number,
  state,
  labels
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe-plugin-aws'
  and json_extract(labels, '$.bug') = 1;
```

### List the open PRs in a repository assigned to a specific user
This query can be used to identify all the open pull requests in a specific repository that have been assigned to a particular user. This is useful in tracking and managing the workload of individual contributors within a project.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  assignee_data ->> 'login' as assignee_login
from
  github_pull_request,
  jsonb_array_elements(assignees) as assignee_data
where
  repository_full_name = 'turbot/steampipe-plugin-aws'
  and assignee_data ->> 'login' = 'madhushreeray30'
  and state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  state,
  json_extract(assignee_data.value, '$.login') as assignee_login
from
  github_pull_request,
  json_each(assignees) as assignee_data
where
  repository_full_name = 'turbot/steampipe-plugin-aws'
  and json_extract(assignee_data.value, '$.login') = 'madhushreeray30'
  and state = 'OPEN';

```

### Join with github_my_repository to find open PRs in multiple repos
This query allows you to identify open pull requests across multiple repositories within the 'turbot/steampipe' project. It's particularly useful for project managers who need to track ongoing contributions and updates across various parts of the project.

```sql+postgres
select
  i.repository_full_name,
  i.number,
  i.title
from
  github_my_repository as r,
  github_pull_request as i
where
  r.name_with_owner like 'turbot/steampip%'
  and i.state = 'OPEN'
  and i.repository_full_name = r.name_with_owner;
```

```sql+sqlite
select
  i.repository_full_name,
  i.number,
  i.title
from
  github_my_repository as r,
  github_pull_request as i
where
  r.name_with_owner like 'turbot/steampip%'
  and i.state = 'OPEN'
  and i.repository_full_name = r.name_with_owner;
```

### List open PRs in a repository with an array of associated labels
This query is useful for exploring open pull requests in a specific repository, along with their associated labels. This can help in managing and prioritizing work by understanding the context and importance of each pull request.

```sql+postgres
select
  r.repository_full_name,
  r.number,
  r.title,
  jsonb_agg(l ->> 'name') as labels
from
  github_pull_request r,
  jsonb_array_elements(r.labels_src) as l
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
group by
  r.repository_full_name, r.number, r.title;
```

```sql+sqlite
select
  r.repository_full_name,
  r.number,
  r.title,
  json_group_array(json_extract(l.value, '$.name')) as labels
from
  github_pull_request r,
  json_each(r.labels_src) as l
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
group by
  r.repository_full_name, r.number, r.title;
```

OR

```sql+postgres
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_pull_request r,
  jsonb_object_keys(r.labels) as t
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
group by
  repository_full_name, number, title;
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  json_group_array(t.value) as labels
from
  github_pull_request r,
  json_each(r.labels) as t
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
group by
  repository_full_name, number, title;
```