---
title: "Steampipe Table: github_workflow - Query GitHub Workflows using SQL"
description: "Allows users to query GitHub Workflows, specifically the details of workflows in GitHub repositories, enabling the tracking of workflow configurations, their state, and the associated branches."
---

# Table: github_workflow - Query GitHub Workflows using SQL

GitHub Workflows is a feature within GitHub Actions that allows you to automate, customize, and execute your software development workflows right in your repository. It provides a flexible way to build an automated software development lifecycle workflow. With GitHub Workflows, you can build, test, and deploy your code right from GitHub.

## Table Usage Guide

The `github_workflow` table provides insights into Workflows within GitHub Actions. As a DevOps engineer, explore workflow-specific details through this table, including workflow configurations, status, and associated branches. Utilize it to monitor and manage workflows, such as those with specific event triggers, the branches associated with a workflow, and the verification of workflow configurations.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List workflows
Explore the characteristics and details of workflows within a specific GitHub repository. This can help in understanding the workflow structure and any specific patterns or anomalies, thereby aiding in effective repository management.

```sql+postgres
select
  repository_full_name,
  name,
  path,
  text,
  line_count,
  size,
  language,
  node_id,
  is_truncated,
  is_generated,
  is_binary,
  text_json,
  pipeline
from
  github_workflow
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  repository_full_name,
  name,
  path,
  text,
  line_count,
  size,
  language,
  node_id,
  is_truncated,
  is_generated,
  is_binary,
  text_json,
  pipeline
from
  github_workflow
where
  repository_full_name = 'turbot/steampipe';
```