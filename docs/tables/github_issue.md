---
title: "Steampipe Table: github_issue - Query GitHub Issues using SQL"
description: "Allows users to query GitHub Issues, providing insights into open, closed, and in-progress issues across repositories."
folder: "Issue"
---

# Table: github_issue - Query GitHub Issues using SQL

GitHub Issues is a feature within GitHub that allows users to track bugs, enhancements, or other requests. It provides a platform for users to collaborate on problems, discuss complex details, and manage updates on ongoing issues. GitHub Issues helps in tracking individual tasks within a project, linking tasks to each other, and organizing tasks into manageable units.

## Table Usage Guide

The `github_issue` table provides insights into issues within GitHub repositories. As a project manager or developer, explore issue-specific details through this table, including status, assignees, labels, and associated metadata. Utilize it to uncover information about issues, such as those that are overdue, the collaboration between team members on certain issues, and the overall progress of issues within a project.

**Important Notes**
- You must specify the `repository_full_name` (owner/repository) column in `where` or `join` clause to query the table.
- The pull requests are technically also issues in GitHub, however we do not include them in the `github_issue` table; You should use the `github_pull_request` table to query PRs.

## Examples

### List the issues in a repository
Explore the status and assignment of issues within a specific GitHub repository to better manage project tasks and responsibilities. This can help in tracking task progress and identifying bottlenecks in the project workflow.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe';
```

### List the unassigned open issues in a repository
Identify instances where there are open issues in a specific repository that have not been assigned to anyone. This is useful to ensure all issues are being addressed and no task is left unattended.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
and
  assignees_total_count = 0
and
  state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
and
  assignees_total_count = 0
and
  state = 'OPEN';
```

### Report of the number issues in a repository by author
Determine the areas in which specific authors are contributing to the number of issues in a particular project, in this case, 'turbot/steampipe'. This can be useful for understanding individual contributions and identifying key contributors or problematic areas based on the number of issues raised by each author.

```sql+postgres
select
  author_login,
  count(*) as num_issues
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  num_issues desc;
```

```sql+sqlite
select
  author_login,
  count(*) as num_issues
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  num_issues desc;
```

### Join with github_my_repository to find open issues in multiple repos that you own or contribute to
Discover the open issues across multiple repositories that you own or contribute to, particularly those related to 'turbot/steampipe'. This can help manage and prioritize your workflow by providing a clear overview of outstanding tasks.

```sql+postgres
select
  i.repository_full_name,
  i.number,
  i.title
from
  github_my_repository as r,
  github_issue as i
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
  github_issue as i
where
  r.name_with_owner like 'turbot/steampip%'
  and i.state = 'OPEN'
  and i.repository_full_name = r.name_with_owner;
```

### List all issues with labels as a string array (instead of JSON objects)
Explore which issues on the 'turbot/steampipe' repository have been tagged with specific labels. This can help in identifying trends or patterns in issue categorization, aiding in more efficient issue management and resolution.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_issue i,
  jsonb_object_keys(i.labels) as t
where
  repository_full_name = 'turbot/steampipe'
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
  github_issue i,
  json_each(i.labels) as t
where
  repository_full_name = 'turbot/steampipe'
group by
  repository_full_name, number, title;
```

OR

```sql+postgres
select
  i.repository_full_name,
  i.number,
  i.title,
  json_agg(l ->> 'name') as labels
from
  github_issue i,
  jsonb_array_elements(i.labels_src) as l
where
  repository_full_name = 'turbot/steampipe'
group by
  i.repository_full_name, i.number, i.title;
```

```sql+sqlite
select
  i.repository_full_name,
  i.number,
  i.title,
  json_group_array(json_extract(l.value, '$.name')) as labels
from
  github_issue i,
  json_each(i.labels_src) as l
where
  repository_full_name = 'turbot/steampipe'
group by
  i.repository_full_name, i.number, i.title;
```

### List all open issues in a repository with a specific label
This query is useful for identifying open issues tagged with a specific label within a designated repository. This can help in prioritizing bug fixes and managing project workflows effectively.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_issue i,
  jsonb_object_keys(labels) as t
where
  repository_full_name = 'turbot/steampipe'
and
  state = 'OPEN'
and
  labels ? 'bug'
group by
  repository_full_name, number, title;
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  (select json_group_array(labels.value) FROM json_each(i.labels) as labels) as labels
from
  github_issue i
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
  and json_extract(i.labels, '$.bug') is not null
group by
  repository_full_name, number, title;
```