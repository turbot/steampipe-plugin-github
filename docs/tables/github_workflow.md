---
title: "Steampipe Table: github_workflow - Query GitHub Workflows using SQL"
description: "Allows users to query GitHub Workflows, specifically the details of workflows in GitHub repositories, enabling the tracking of workflow configurations, their state, and the associated branches."
folder: "Actions"
---

# Table: github_workflow - Query GitHub Workflows using SQL

GitHub Workflows is a feature within GitHub Actions that allows you to automate, customize, and execute your software development workflows right in your repository. It provides a flexible way to build an automated software development lifecycle workflow. With GitHub Workflows, you can build, test, and deploy your code right from GitHub.

## Table Usage Guide

The `github_workflow` table provides insights into Workflows within GitHub Actions. As a DevOps engineer, explore workflow-specific details through this table, including workflow configurations, status, and associated branches. Utilize it to monitor and manage workflows, such as those with specific event triggers, the branches associated with a workflow, and the verification of workflow configurations.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.
- To query this table using Fine-grained access tokens, the following permissions are required:
  - **"Actions" repository permission (read)** – Required to access the all columns.
  - **"Contents" repository permission (read)** – Required to access the columns `workflow_file_content`, `workflow_file_content_json`, and `pipeline`.

## Examples

### List workflows
Explore the characteristics and details of workflows within a specific GitHub repository. This can help in understanding the workflow structure and any specific patterns or anomalies, thereby aiding in effective repository management.

```sql+postgres
select
  repository_full_name,
  name,
  path,
  node_id,
  state,
  url
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
  node_id,
  state,
  url
from
  github_workflow
where
  repository_full_name = 'turbot/steampipe';
```