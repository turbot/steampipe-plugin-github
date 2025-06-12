---
title: "Steampipe Table: github_package - Query GitHub Packages using SQL"
description: "Allows users to query GitHub Packages, including package metadata, versions, owner details, and associated repositories, providing insights into the package management within GitHub repositories."
folder: "Package"
---

# Table: github_package - Query GitHub Packages using SQL

GitHub Packages allow you to store and manage container images and other packages directly within your GitHub repositories. With this table, you can query details about GitHub packages within an organization, including information about the package itself, its owner, associated repositories, and visibility.

## Table Usage Guide

The `github_package` table provides detailed insights into packages hosted in GitHub's container registry or other package managers like npm. You can retrieve information about the package owner, repository, creation time, versioning details, and more. This is particularly useful for monitoring and managing package versions and repository associations in an organization.

**Important Notes**
- You must specify the `organization` column in the `where` clause to query the table.
- OAuth app tokens and personal access tokens (classic) need the `read:packages` scope to use this endpoint. If the package_type belongs to a GitHub Packages registry that only supports repository-scoped permissions, the `repo` scope is also required. Please refer [Permissions for repository-scoped packages](https://docs.github.com/en/packages/learn-github-packages/about-permissions-for-github-packages#permissions-for-repository-scoped-packages) for more information.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

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
  html_url
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
  html_url
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
  (repository ->> 'private') as repository_private
from
  github_package
where
  organization = 'turbot'
  and (repository ->> 'private')::bool = true;
```

```sql+sqlite
select
  name,
  repository_full_name,
  json_extract(repository, '$.private') as repository_private
from
  github_package
where
  organization = 'turbot'
  and json_extract(repository, '$.private') = true;
```

### Get the details of a specific package by name
Retrieve comprehensive details about a specific package by filtering based on the package name. It is useful for analyzing a single package's metadata, owner, and repository association.

```sql+postgres
select
  id,
  name,
  package_type,
  repository_full_name,
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

### Get owner information for the packages
Retrieve package names along with owner information (login, ID, URL, and HTML URL) for all GitHub packages under the turbot organization.

```sql+postgres
select
  name,
  (owner ->> 'login') as owner_login,
  (owner ->> 'id') as owner_id,
  (owner ->> 'url') as owner_url,
  (owner ->> 'html_url') as owner_html_url
from
  github_package
where
  organization = 'turbot';
```

```sql+sqlite
select
  name,
  json_extract(owner, '$.login') as owner_login,
  json_extract(owner, '$.id') as owner_id,
  json_extract(owner, '$.url') as owner_url,
  json_extract(owner, '$.html_url') as owner_html_url
from
  github_package
where
  organization = 'turbot';
```

### Get repository info for the packages
Extract key fields (like repository name, full name, visibility, URL, and owner login) from the nested repository object in each GitHub package. This provides a flat, queryable view of essential repository information per package.

```sql+postgres
select
  name,
  repository ->> 'name' as repository_name,
  repository ->> 'id' as repository_id,
  repository ->> 'private' as repository_private,
  repository ->> 'html_url' as repository_html_url,
  repository ->> 'description' as repository_description,
  repository ->> 'fork' as repository_fork,
  repository -> 'owner' ->> 'login' as repository_owner_login,
  repository ->> 'stargazers_url' as repository_stargazers_url,
  repository ->> 'contents_url' as repository_contents_url
from
  github_package
where
  organization = 'turbot';
```

```sql+sqlite
select
  name,
  json_extract(repository, '$.name') as repository_name,
  json_extract(repository, '$.id') as repository_id,
  json_extract(repository, '$.private') as repository_private,
  json_extract(repository, '$.html_url') as repository_html_url,
  json_extract(repository, '$.description') as repository_description,
  json_extract(repository, '$.fork') as repository_fork,
  json_extract(repository, '$.owner.login') as repository_owner_login,
  json_extract(repository, '$.stargazers_url') as repository_stargazers_url,
  json_extract(repository, '$.contents_url') as repository_contents_url
from
  github_package
where
  organization = 'turbot';
```
