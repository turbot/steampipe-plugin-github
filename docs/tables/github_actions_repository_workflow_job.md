---
title: "Steampipe Table: github_actions_repository_workflow_job - Query GitHub Actions Repository Workflow Jobs using SQL"
description: "Allows users to query GitHub Actions Repository Workflow Jobs, specifically the details of each workflow job in a repository, providing insights into the status, conclusion, and other metadata of the jobs."
folder: "Actions"
---

# Table: github_actions_repository_workflow_job - Query GitHub Actions Repository Workflow Jobs using SQL

GitHub Actions is a CI/CD solution that allows you to automate how you build, test, and deploy your projects on any platform, including Linux, macOS, and Windows. It lets you run a series of commands in response to events on GitHub. With GitHub Actions, you can build end-to-end continuous integration (CI) and continuous deployment (CD) capabilities directly in your repository.

## Table Usage Guide

The `github_actions_repository_workflow_job` table provides insights into GitHub Actions Repository Workflow Jobs. As a software developer or DevOps engineer, explore details of each workflow job in a repository through this table, including its status, conclusion, and other metadata. Utilize it to monitor and analyze the performance and results of your CI/CD workflows, ensuring your software development process is efficient and effective.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Actions (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name`  column, along with either the `run_id` or the `id` column, in the `where` or `join` clause to query the table.

## Examples

### List workflow jobs
Analyze the settings to understand the operations and status of the workflow jobs in a specific GitHub repository. This can be beneficial in assessing the efficiency and effectiveness of workflows, identifying potential issues, and making informed decisions on workflow optimization.

```sql+postgres
select
  *
from
  github_actions_repository_workflow_job
where
  repository_full_name = 'turbot/steampipe'
  and run_id = 26404053809;
```

```sql+sqlite
select
  *
from
  github_actions_repository_workflow_job
where
  repository_full_name = 'turbot/steampipe'
  and run_id = 26404053809;
```

### List failed workflow jobs
Identify instances where workflow jobs have failed within the 'turbot/steampipe' repository. This can be useful for debugging and identifying problematic jobs.

```sql+postgres
select
  id,
  steps,
  runner_id,
  conclusion,
  status,
  run_attempt,
  run_url,
  head_sha,
  head_branch
from
  github_actions_repository_workflow_job
where
  repository_full_name = 'turbot/steampipe'
  and run_id = 26404053809
  and conclusion = 'failure';
```

```sql+sqlite
select
  id,
  steps,
  runner_id,
  conclusion,
  status,
  run_attempt,
  run_url,
  head_sha,
  head_branch
from
  github_actions_repository_workflow_job
where
  repository_full_name = 'turbot/steampipe'
  and run_id = 26404053809
  and conclusion = 'failure';
```