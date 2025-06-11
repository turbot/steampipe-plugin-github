---
title: "Steampipe Table: github_actions_repository_workflow_run - Query GitHub Actions Repository Workflow Runs using SQL"
description: "Allows users to query GitHub Actions Repository Workflow Runs, specifically the details of each workflow run in a repository, providing insights into the status, conclusion, and other metadata of the runs."
folder: "Actions"
---

# Table: github_actions_repository_workflow_run - Query GitHub Actions Repository Workflow Runs using SQL

GitHub Actions is a CI/CD solution that allows you to automate how you build, test, and deploy your projects on any platform, including Linux, macOS, and Windows. It lets you run a series of commands in response to events on GitHub. With GitHub Actions, you can build end-to-end continuous integration (CI) and continuous deployment (CD) capabilities directly in your repository.

## Table Usage Guide

The `github_actions_repository_workflow_run` table provides insights into GitHub Actions Repository Workflow Runs. As a software developer or DevOps engineer, explore details of each workflow run in a repository through this table, including its status, conclusion, and other metadata. Utilize it to monitor and analyze the performance and results of your CI/CD workflows, ensuring your software development process is efficient and effective.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table. 
- To query this table using Fine-grained access tokens, the following permissions are required:
  - **"Actions" repository permission (read)** – Required to access the all columns.

## Examples

### List workflow runs
Analyze the settings to understand the operations and status of the workflow runs in a specific GitHub repository. This can be beneficial in assessing the efficiency and effectiveness of workflows, identifying potential issues, and making informed decisions on workflow optimization.

```sql+postgres
select
  *
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  *
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe';
```

### List failure workflow runs
Identify instances where workflow runs have failed within the 'turbot/steampipe' repository. This can be useful for debugging and identifying problematic workflows.

```sql+postgres
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch
from
    github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and conclusion = 'failure';
```

```sql+sqlite
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and conclusion = 'failure';
```

### List manual workflow runs
This query helps you gain insights into the manually triggered workflow runs in the 'turbot/steampipe' repository. It's particularly useful for tracking and analyzing the performance and results of these workflow runs, allowing you to identify any potential issues or areas for improvement.

```sql+postgres
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch,
  actor_login,
  triggering_actor_login
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and event = 'workflow_dispatch';
```

```sql+sqlite
select
  id,
  event,
  workflow_id,
  conclusion,
  status,
  run_number,
  workflow_url,
  head_commit,
  head_branch,
  actor_login,
  triggering_actor_login
from
  github_actions_repository_workflow_run
where
  repository_full_name = 'turbot/steampipe' and event = 'workflow_dispatch';
```