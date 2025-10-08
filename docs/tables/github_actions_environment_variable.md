---
title: "Steampipe Table: github_actions_environment_variable - Query GitHub Actions Environment Variables using SQL"
description: "Allows users to query GitHub Actions Environment Variables, specifically to retrieve information about variables stored in a repository environment for use in GitHub Actions workflows."
folder: "Actions"
---

# Table: github_actions_environment_variable - Query GitHub Actions Environment Variables using SQL

GitHub Actions is a service provided by GitHub that allows you to automate, customize, and execute your software development workflows right in your repository. Environment-level variables are scoped to a specific environment within a repository, allowing you to have different values for different deployment environments (e.g., production, staging, development).

## Table Usage Guide

The `github_actions_environment_variable` table provides insights into variables stored within a repository environment. As a DevOps engineer or developer, explore variable-specific details through this table, including the names, values, and timestamps of when they were created or updated. Utilize it to uncover information about variables used in specific environments within your GitHub Actions workflows, helping you understand and manage your CI/CD pipeline configuration at the environment level.

For more information, see [Store information in variables](https://docs.github.com/en/actions/learn-github-actions/variables) and [Using environments for deployment](https://docs.github.com/en/actions/deployment/targeting-different-environments/using-environments-for-deployment) in the GitHub Actions documentation.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Metadata (Read-only): Required to access general repository metadata.
    - Environments (Read-only): Required to access environment variables.

**Important Notes**
- You must specify both the `repository_full_name` and `environment_name` columns in `where` or `join` clause to query the table.
- Environment names in GitHub can contain special characters and should be URL-encoded if needed.

## Examples

### List environment variables
Explore the variables configured for a specific environment in a repository.

```sql+postgres
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production';
```

```sql+sqlite
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production';
```

### Get a specific environment variable
Retrieve details about a specific variable in a repository environment.

```sql+postgres
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production'
  and name = 'API_URL';
```

```sql+sqlite
select
  name,
  value,
  created_at,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production'
  and name = 'API_URL';
```

### List recently updated environment variables
Find environment variables that have been modified recently to track configuration changes.

```sql+postgres
select
  repository_full_name,
  environment_name,
  name,
  value,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production'
  and updated_at > now() - interval '7 days'
order by
  updated_at desc;
```

```sql+sqlite
select
  repository_full_name,
  environment_name,
  name,
  value,
  updated_at
from
  github_actions_environment_variable
where
  repository_full_name = 'turbot/steampipe'
  and environment_name = 'production'
  and updated_at > datetime('now', '-7 days')
order by
  updated_at desc;
```

### List all variables for all environments in a repository
Get a comprehensive view of all environment variables in a repository.

```sql+postgres
select
  e.name as environment_name,
  v.name as variable_name,
  v.value,
  v.updated_at
from
  github_repository_environment e
  left join github_actions_environment_variable v
    on e.repository_full_name = v.repository_full_name
    and e.name = v.environment_name
where
  e.repository_full_name = 'turbot/steampipe'
order by
  e.name, v.name;
```

```sql+sqlite
select
  e.name as environment_name,
  v.name as variable_name,
  v.value,
  v.updated_at
from
  github_repository_environment e
  left join github_actions_environment_variable v
    on e.repository_full_name = v.repository_full_name
    and e.name = v.environment_name
where
  e.repository_full_name = 'turbot/steampipe'
order by
  e.name, v.name;
```
