---
title: "Steampipe Table: github_search_issue - Query GitHub Issues using SQL"
description: "Allows users to query GitHub Issues, specifically to retrieve and analyze issue data related to any GitHub repository, providing insights into repository management and development activities."
folder: "Issue"
---

# Table: github_search_issue - Query GitHub Issues using SQL

GitHub Issues is a feature of GitHub, a web-based hosting service for version control, that allows users to track and manage tasks, enhancements, and bugs for projects. It provides a platform for collaboration, enabling developers to work together on projects from anywhere. GitHub Issues helps users stay informed about the progress and performance of their projects, and take appropriate actions when required.

## Table Usage Guide

The `github_search_issue` table provides insights into issues within GitHub repositories. As a project manager or developer, explore issue-specific details through this table, including status, assignees, labels, and associated metadata. Utilize it to uncover information about issues, such as those that are open, the assignees working on them, and the labels attached to them.

**Important Notes**
- You must always include at least one search term when searching source code in the where or join clause using the `query` column. You can narrow the results using these search qualifiers in any combination. See [Searching issues and pull requests](https://docs.github.com/search-github/searching-on-github/searching-issues-and-pull-requests) for details on the GitHub query syntax.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### List issues by the title, body, or comments
Discover the segments that contain issues based on their title, body or comments. This can be beneficial for understanding and managing the issues more effectively.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'github_search_commit in:title in:body';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'github_search_commit in:title in:body';
```

### List issues in open state assigned to a specific user
This query allows you to identify open issues assigned to a specific user within a particular GitHub repository. This is useful for monitoring a user's workload or tracking the progress of issue resolution.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List issues with public visibility assigned to a specific user
Discover the segments that include public issues assigned to a specific user on GitHub. This can be useful to monitor the work of a specific developer or track public issues in a particular repository.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:public assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:public assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List issues not linked to a pull request
Determine the areas in which active issues exist that are not linked to any pull request in the GitHub repository for the Steampipe plugin. This is useful to identify potential tasks that may need attention or further investigation.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open -linked:pr repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open -linked:pr repo:turbot/steampipe-plugin-github';
```

### List blocked issues
Identify instances where issues have been blocked in the 'turbot/steampipe-plugin-github' repository. This can help in understanding project bottlenecks and prioritizing tasks.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'label:blocked repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'label:blocked repo:turbot/steampipe-plugin-github';
```

### List issues with over 10 comments
Identify instances where GitHub issues in the Turbot organization have garnered more than 10 comments. This is useful for tracking popular discussions and understanding the issues that are generating significant community engagement.

```sql+postgres
select
  title,
  id,
  comments_total_count,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'org:turbot comments:>10';
```

```sql+sqlite
select
  title,
  id,
  comments_total_count,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'org:turbot comments:>10';
```

### List issues that took more than 30 days to close
Discover the segments that took over a month to resolve within a specific organization. This allows for an analysis of efficiency in issue resolution and can highlight areas for improvement in the workflow process.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  closed_at,
  url
from
  github_search_issue
where
  query = 'org:turbot state:closed'
  and closed_at > (created_at + interval '30' day);
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  closed_at,
  url
from
  github_search_issue
where
  query = 'org:turbot state:closed'
  and closed_at > datetime(created_at, '+30 day');
```