---
title: "Steampipe Table: github_repository - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories, providing detailed information about each repository including its owner, name, description, visibility, and more."
folder: "Repository"
---

# Table: github_repository - Query GitHub Repositories using SQL

GitHub is a web-based hosting service for version control using Git. It offers all of the distributed version control and source code management (SCM) functionality of Git, as well as adding its own features. GitHub provides access control and several collaboration features such as bug tracking, feature requests, task management, and wikis for every project.

## Table Usage Guide

The `github_repository` table provides insights into repositories within GitHub. As a developer or project manager, explore repository-specific details through this table, including owner details, repository names, descriptions, and visibility status. Utilize it to uncover information about repositories, such as those with public or private visibility, the owner of each repository, and the description of what each repository contains.

**Important Notes**
- You must specify the `full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.
- To list all of your repositories use the `github_my_repository` table instead. The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.

## Examples

### Get information about a specific repository
Discover the details of a specific GitHub repository, such as its creation date, update history, disk usage, owner, primary language, number of forks, star count, URL, license information, and description. This is beneficial for gaining insights into the repository's overall status, usage, and popularity.

```sql+postgres
select
  name,
  node_id,
  id,
  created_at,
  updated_at,
  disk_usage,
  owner_login,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  url,
  license_info ->> 'spdx_id' as license,
  description
from
  github_repository
where
  full_name = 'postgres/postgres';
```

```sql+sqlite
select
  name,
  node_id,
  id,
  created_at,
  updated_at,
  disk_usage,
  owner_login,
  json_extract(primary_language, '$.name') as language,
  fork_count,
  stargazer_count,
  url,
  json_extract(license_info, '$.spdx_id') as license,
  description
from
  github_repository
where
  full_name = 'postgres/postgres';
```

### Get your permissions for a specific repository
This query allows you to understand the level of access you have to a particular repository, such as whether you can administer, create projects, subscribe, or update topics. It's useful to ensure you have the correct permissions for your intended actions, helping to avoid unexpected access issues.

```sql+postgres
select
  name,
  your_permission,
  can_administer,
  can_create_projects,
  can_subscribe,
  can_update_topics,
  possible_commit_emails
from
  github_repository
where
  full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  your_permission,
  can_administer,
  can_create_projects,
  can_subscribe,
  can_update_topics,
  possible_commit_emails
from
  github_repository
where
  full_name = 'turbot/steampipe';
```