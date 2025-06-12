---
title: "Steampipe Table: github_repository_deployment - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories, specifically the deployment information, providing insights into the deployment status, environment, and associated metadata."
folder: "Deployments"
---

# Table: github_repository_deployment - Query GitHub Repositories using SQL

GitHub Repositories are a fundamental resource in GitHub that allow users to manage and store revisions of projects. Repositories contain all of the project files and revision history. Users can work together within a repository to edit and manage the project.

## Table Usage Guide

The `github_repository_deployment` table offers insights into GitHub repositories' deployment details. As a developer or project manager, you can use this table to retrieve deployment status, environment, and related metadata for each repository. This can be particularly useful for monitoring deployment progress, identifying deployment patterns, and troubleshooting deployment issues.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Deployments (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### List deployments for a repository
Explore the deployment history of a specific repository to understand its version updates and changes made by different contributors. This is useful in tracking the evolution of a project and assessing its progress over time.

```sql+postgres
select
  id,
  node_id,
  commit_sha,
  created_at,
  creator ->> 'login' as creator_login,
  description,
  environment,
  latest_status,
  payload,
  ref ->> 'prefix' as ref_prefix,
  ref ->> 'name' as ref_name,
  state,
  task,
  updated_at
from
  github_repository_deployment
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  id,
  node_id,
  commit_sha,
  created_at,
  json_extract(creator, '$.login') as creator_login,
  description,
  environment,
  latest_status,
  payload,
  json_extract(ref, '$.prefix') as ref_prefix,
  json_extract(ref, '$.name') as ref_name,
  state,
  task,
  updated_at
from
  github_repository_deployment
where
  repository_full_name = 'turbot/steampipe';
```

### List deployments for all your repositories
Explore the deployment history across all your GitHub repositories. This query helps you assess the details of each deployment such as its creator, status, environment, and more, providing a comprehensive view of your repositories' deployment activities.

```sql+postgres
select
  id,
  node_id,
  created_at,
  creator ->> 'login' as creator_login,
  commit_sha,
  description,
  environment,
  latest_status,
  payload,
  ref ->> 'prefix' as ref_prefix,
  ref ->> 'name' as ref_name,
  state,
  task,
  updated_at
from
  github_repository_deployment
where
  repository_full_name in (select name_with_owner from github_my_repository);
```

```sql+sqlite
select
  id,
  node_id,
  created_at,
  json_extract(creator, '$.login') as creator_login,
  commit_sha,
  description,
  environment,
  latest_status,
  payload,
  json_extract(ref, '$.prefix') as ref_prefix,
  json_extract(ref, '$.name') as ref_name,
  state,
  task,
  updated_at
from
  github_repository_deployment
where
  repository_full_name in (select name_with_owner from github_my_repository);
```