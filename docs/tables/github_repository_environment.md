---
title: "Steampipe Table: github_repository_environment - Query GitHub Repository Environments using SQL"
description: "Allows users to query GitHub Repository Environments, specifically details about each environment associated with a repository, providing insights into environment names, URLs, and protection rules."
folder: "Deployments"
---

# Table: github_repository_environment - Query GitHub Repository Environments using SQL

GitHub Repository Environments are a feature of GitHub that allows developers to manage and track their software deployment pipelines. Each environment represents a stage in the deployment workflow, such as production, staging, or testing. Environments can have specific deployment policies and protection rules, providing a controlled and secure workflow for deploying software.

## Table Usage Guide

The `github_repository_environment` table provides insights into GitHub repository environments. As a DevOps engineer or a repository manager, explore environment-specific details through this table, including environment names, URLs, and protection rules. Utilize it to manage and monitor the deployment workflow, ensuring controlled and secure software deployment.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Actions (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List environments for a repository
Explore which environments are associated with a particular repository. This can be useful for understanding the different settings where your code is being tested or deployed.

```sql+postgres
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name = 'turbot/steampipe';
```

### List environments all your repositories
Discover the segments that contain all your repository environments. This allows you to identify and manage the environments in which your code is running, helping you enhance your project's efficiency and security.

```sql+postgres
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name in (select name_with_owner from github_my_repository);
```

```sql+sqlite
select
  id,
  node_id,
  name
from
  github_repository_environment
where
  repository_full_name in (select name_with_owner from github_my_repository);
```