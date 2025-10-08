---
title: "Steampipe Table: github_actions_repository_variable - Query GitHub Actions Repository Variables using SQL"
description: "Allows users to query GitHub Actions Repository Variables, specifically to retrieve information about variables stored in a GitHub repository for use in GitHub Actions workflows."
folder: "Actions"
---

# Table: github_actions_repository_variable - Query GitHub Actions Repository Variables using SQL

GitHub Actions is a service provided by GitHub that allows you to automate, customize, and execute your software development workflows right in your repository. Variables allow you to store non-sensitive information that can be referenced in GitHub Actions workflows.

## Table Usage Guide

The `github_actions_repository_variable` table provides insights into variables stored within a GitHub repository. As a DevOps engineer or developer, explore variable-specific details through this table, including the names, values, and timestamps of when they were created or updated. Utilize it to uncover information about variables used in your GitHub Actions workflows, helping you understand and manage your CI/CD pipeline configuration at the repository level.

For more information, see [Store information in variables](https://docs.github.com/en/actions/learn-github-actions/variables) in the GitHub Actions documentation.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Metadata (Read-only): Required to access general repository metadata.
    - Variables (Read-only): Required to access all columns.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List repository variables
Explore the variables configured for a specific repository.

```sql+postgres
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe';
```

### Get a specific repository variable
Retrieve details about a specific variable in a repository.

```sql+postgres
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe'
  and name = 'MY_VARIABLE';
```

```sql+sqlite
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe'
  and name = 'MY_VARIABLE';
```

### List recently updated variables
Find variables that have been modified recently to track configuration changes.

```sql+postgres
select
  repository_full_name,
  name,
  value,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe'
  and updated_at > now() - interval '7 days'
order by
  updated_at desc;
```

```sql+sqlite
select
  repository_full_name,
  name,
  value,
  updated_at
from
  github_actions_repository_variable
where
  repository_full_name = 'turbot/steampipe'
  and updated_at > datetime('now', '-7 days')
order by
  updated_at desc;
```

### Count variables per repository
Analyze the number of variables configured across multiple repositories.

```sql+postgres
select
  repository_full_name,
  count(*) as variable_count
from
  github_actions_repository_variable
where
  repository_full_name in ('turbot/steampipe', 'turbot/steampipe-plugin-sdk')
group by
  repository_full_name;
```

```sql+sqlite
select
  repository_full_name,
  count(*) as variable_count
from
  github_actions_repository_variable
where
  repository_full_name in ('turbot/steampipe', 'turbot/steampipe-plugin-sdk')
group by
  repository_full_name;
```

### List all repository variables across your repositories
Get a comprehensive view of all repository variables across your repositories.

```sql+postgres
select
  r.name_with_owner as repository_full_name,
  v.name,
  v.value,
  v.updated_at
from
  github_my_repository r
  left join github_actions_repository_variable v on r.name_with_owner = v.repository_full_name
order by
  r.name_with_owner, v.name;
```

```sql+sqlite
select
  r.name_with_owner as repository_full_name,
  v.name,
  v.value,
  v.updated_at
from
  github_my_repository r
  left join github_actions_repository_variable v on r.name_with_owner = v.repository_full_name
order by
  r.name_with_owner, v.name;
```

