---
title: "Steampipe Table: github_search_label - Query GitHub Labels using SQL"
description: "Allows users to query Labels in GitHub, specifically the metadata and details of labels that exist within a repository."
---

# Table: github_search_label - Query GitHub Labels using SQL

GitHub Labels are a feature within GitHub that allows users to categorize and filter issues and pull requests. They are customizable and can be used in a variety of ways, including to indicate priority, type of work, or status. Labels are a powerful tool for managing work and communicating about issues and pull requests across the team.

## Table Usage Guide

The `github_search_label` table provides insights into Labels within GitHub. As a project manager or developer, explore label-specific details through this table, including color, default status, and associated metadata. Utilize it to uncover information about labels, such as their usage across issues and pull requests, and to facilitate efficient project management and issue tracking.

**Important Notes**
- You must always include at least one search term and repository ID when searching source code in the where or join clause using the `query` and `repository_id` columns respectively.

## Examples

### List labels for bug, enhancement and blocked
Determine the areas in which specific labels such as 'bug', 'enhancement', and 'blocked' are used within a particular GitHub repository. This allows for a better understanding of issue categorization and priority setting within the project.

```sql+postgres
select
  id,
  repository_id,
  name,
  repository_full_name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'bug enhancement blocked';
```

```sql+sqlite
select
  id,
  repository_id,
  name,
  repository_full_name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'bug enhancement blocked';
```

### List labels where specific text matches in name or description
Determine the areas in which specific text matches in labels' name or description within a particular GitHub repository. This can be used to quickly locate and organize labels related to a specific topic or task.

```sql+postgres
select
  id,
  repository_id,
  name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'work';
```

```sql+sqlite
select
  id,
  repository_id,
  name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'work';
```