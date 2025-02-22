---
title: "Steampipe Table: github_package - Query GitHub Packages using SQL"
description: "Allows users to query GitHub Packages, including package metadata, versions, owner details, and associated repositories, providing insights into the package management within GitHub repositories."
---

# Table: github_package - Query GitHub Packages using SQL

GitHub Packages allow you to store and manage container images and other packages directly within your GitHub repositories. With this table, you can query details about GitHub packages within an organization, including information about the package itself, its owner, associated repositories, and visibility.

## Table Usage Guide

The `github_package` table provides detailed insights into packages hosted in GitHub's container registry or other package managers like npm. You can retrieve information about the package owner, repository, creation time, versioning details, and more. This is particularly useful for monitoring and managing package versions and repository associations in an organization.

**Important Notes**
- You must specify the `organization` column in the `where` clause to query the table.
- OAuth app tokens and personal access tokens (classic) need the `read:packages` scope to use this endpoint. If the package_type belongs to a GitHub Packages registry that only supports repository-scoped permissions, the `repo` scope is also required. Please refer [Permissions for repository-scoped packages](https://docs.github.com/en/packages/learn-github-packages/about-permissions-for-github-packages#permissions-for-repository-scoped-packages) for more information.

## Examples

### List all packages for a specific organization
Retrieve the details of all packages available in a specified GitHub organization, including information about the owner, repository, and package type.

```sql+postgres
select
  id,
  name,
  package_type,
  repository_full_name,
  visibility,
  created_at,
  updated_at,
  owner_login,
  url
from
  github_package
where
  organization = 'turbot';
```

```sql+sqlite
select
  id,
  name,
  package_type,
  repository_full_name,
  visibility,
  created_at,
  updated_at,
  owner_login,
  url
from
  github_package
where
  organization = 'turbot';
```

### List all public packages for an organization
Filter and list only the packages that are publicly visible in a specific GitHub organization. This is helpful for managing public packages in your organization.

```sql+postgres
select
  id,
  name,
  package_type,
  repository_full_name,
  visibility,
  html_url,
  owner_login
from
  github_package
where
  organization = 'turbot'
  and visibility = 'public';
```

```sql+sqlite
select
  id,
  name,
  package_type,
  repository_full_name,
  visibility,
  html_url,
  owner_login
from
  github_package
where
  organization = 'turbot'
  and visibility = 'public';
```

### Find packages associated with private repositories
Identify packages that are tied to private repositories within a GitHub organization, allowing you to manage internal packages.

```sql+postgres
select
  name,
  repository_full_name,
  repository_private,
  owner_login
from
  github_package
where
  organization = 'turbot'
  and repository_private = true;
```

```sql+sqlite
select
  name,
  repository_full_name,
  repository_private,
  owner_login
from
  github_package
where
  organization = 'turbot'
  and repository_private = 1;
```

### Get the details of a specific package by name
Retrieve comprehensive details about a specific package by filtering based on the package name. It is useful for analyzing a single package's metadata, owner, and repository association.

```sql+postgres
select
  id,
  name,
  package_type,
  repository_full_name,
  owner_login,
  created_at,
  updated_at,
  url
from
  github_package
where
  organization = 'turbot'
  and name = 'steampipe/plugin/turbot/aws';
```

```sql+sqlite
select
  id,
  name,
  package_type,
  repository_full_name,
  owner_login,
  created_at,
  updated_at,
  url
from
  github_package
where
  organization = 'turbot'
  and name = 'steampipe/plugin/turbot/aws';
```

### List all versions of a package for an organization
Explore the versioning details for a specific package, showing the available versions and metadata for each version.

```sql+postgres
select
  name,
  jsonb_array_elements_text(package_version->'versions') as version
from
  github_package
where
  organization = 'turbot'
  and name = 'steampipe/plugin/turbot/aws';
```

```sql+sqlite
select
  name,
  json_extract(package_version, '$.versions[0]') as version
from
  github_package
where
  organization = 'turbot'
  and name = 'steampipe/plugin/turbot/aws';
```