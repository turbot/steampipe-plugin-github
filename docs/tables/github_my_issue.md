---
title: "Steampipe Table: github_my_issue - Query GitHub Issues using SQL"
description: "Allows users to query personal issues in GitHub, specifically focusing on the details of issues assigned to or created by the authenticated user."
folder: "Issue"
---

# Table: github_my_issue - Query GitHub Issues using SQL

GitHub Issues is a feature in GitHub that provides a platform to track bugs, enhancements, or other requests. It allows users to collaborate on tasks, discuss project details, and manage project timelines. Issues are a great way to keep track of tasks, improvements, and bugs for your projects.

## Table Usage Guide

The `github_my_issue` table provides insights into personal issues within GitHub. As a project manager or developer, explore issue-specific details through this table, including the issue title, state, assignee, and associated metadata. Utilize it to manage and track tasks, improvements, and bugs for your projects.

**Important Notes**
- To view **all the issues belonging to a repository**, use the `github_issue` table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Issues (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### List all of the open issues assigned to you
Explore which open issues are currently assigned to you on GitHub. This is useful for managing your workload and prioritizing tasks.

```sql+postgres
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  author_login
from
  github_my_issue
where
  state = 'OPEN';
```

```sql+sqlite
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  author_login
from
  github_my_issue
where
  state = 'OPEN';
```

### List your 10 oldest open issues
Explore which of your open issues have been unresolved the longest to help prioritize your workflow and manage your project effectively.

```sql+postgres
select
  repository_full_name,
  number,
  created_at,
  age(created_at),
  title,
  state
from
  github_my_issue
where
  state = 'OPEN'
order by
  created_at
limit 10;
```

```sql+sqlite
select
  repository_full_name,
  number,
  created_at,
  julianday('now') - julianday(created_at) as age,
  title,
  state
from
  github_my_issue
where
  state = 'OPEN'
order by
  created_at
limit 10;
```