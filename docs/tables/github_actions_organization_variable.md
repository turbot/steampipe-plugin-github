---
title: "Steampipe Table: github_actions_organization_variable - Query GitHub Actions Organization Variables using SQL"
description: "Allows users to query GitHub Actions Organization Variables, specifically to retrieve information about variables stored in a GitHub organization for use in GitHub Actions workflows."
folder: "Actions"
---

# Table: github_actions_organization_variable - Query GitHub Actions Organization Variables using SQL

GitHub Actions is a service provided by GitHub that allows you to automate, customize, and execute your software development workflows right in your repository. Organization-level variables can be shared across multiple repositories within the organization.

## Table Usage Guide

The `github_actions_organization_variable` table provides insights into variables stored within a GitHub organization. As a DevOps engineer or organization administrator, explore variable-specific details through this table, including the names, values, visibility settings, and timestamps of when they were created or updated. Utilize it to uncover information about variables used across your organization's GitHub Actions workflows, helping you understand and manage your CI/CD pipeline configuration at the organization level.

For more information, see [Store information in variables](https://docs.github.com/en/actions/learn-github-actions/variables) in the GitHub Actions documentation.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Organization permissions:
    - Variables (Read-only): Required to access all columns.

**Important Notes**
- You must specify the `organization` column in `where` or `join` clause to query the table.
- The `visibility` column indicates whether the variable is visible to `all`, `private`, or `selected` repositories.

## Examples

### List organization variables
Explore the variables configured for a specific organization.

```sql+postgres
select
  name,
  value,
  visibility,
  created_at,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org';
```

```sql+sqlite
select
  name,
  value,
  visibility,
  created_at,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org';
```

### Get a specific organization variable
Retrieve details about a specific variable in an organization.

```sql+postgres
select
  name,
  value,
  visibility,
  created_at,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and name = 'MY_VARIABLE';
```

```sql+sqlite
select
  name,
  value,
  visibility,
  created_at,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and name = 'MY_VARIABLE';
```

### List recently updated organization variables
Find organization variables that have been modified recently to track configuration changes.

```sql+postgres
select
  organization,
  name,
  value,
  visibility,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and updated_at > now() - interval '7 days'
order by
  updated_at desc;
```

```sql+sqlite
select
  organization,
  name,
  value,
  visibility,
  updated_at
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and updated_at > datetime('now', '-7 days')
order by
  updated_at desc;
```

### List variables by visibility
Analyze which variables are visible to all repositories vs. selected repositories in your organization.

```sql+postgres
select
  name,
  value,
  visibility,
  selected_repositories_url
from
  github_actions_organization_variable
where
  organization = 'my-org'
order by
  visibility, name;
```

```sql+sqlite
select
  name,
  value,
  visibility,
  selected_repositories_url
from
  github_actions_organization_variable
where
  organization = 'my-org'
order by
  visibility, name;
```

### Count variables per organization
Analyze the number of variables configured across multiple organizations.

```sql+postgres
select
  organization,
  count(*) as variable_count
from
  github_actions_organization_variable
where
  organization in ('org-1', 'org-2', 'org-3')
group by
  organization;
```

```sql+sqlite
select
  organization,
  count(*) as variable_count
from
  github_actions_organization_variable
where
  organization in ('org-1', 'org-2', 'org-3')
group by
  organization;
```

### List all organization variables across your organizations
Get a comprehensive view of all variables across organizations you have access to.

```sql+postgres
select
  o.login as organization,
  v.name,
  v.value,
  v.visibility
from
  github_my_organization o
  left join github_actions_organization_variable v on o.login = v.organization
order by
  o.login, v.name;
```

```sql+sqlite
select
  o.login as organization,
  v.name,
  v.value,
  v.visibility
from
  github_my_organization o
  left join github_actions_organization_variable v on o.login = v.organization
order by
  o.login, v.name;
```

### List variables with 'selected' visibility
Find organization variables that are only available to selected repositories.

```sql+postgres
select
  name,
  value,
  visibility,
  selected_repositories_url
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and visibility = 'selected';
```

```sql+sqlite
select
  name,
  value,
  visibility,
  selected_repositories_url
from
  github_actions_organization_variable
where
  organization = 'my-org'
  and visibility = 'selected';
```

