---
title: "Steampipe Table: github_package_version - Query GitHub Package Versions using SQL"
description: "Allows users to query GitHub Package Versions, providing insights into the metadata, release information, visibility, and creation details of each version."
folder: "Package"
---

# Table: github_package_version - Query GitHub Package Versions using SQL

GitHub Package Versions represent different versions of packages stored in GitHub, such as container images, npm packages, and more. These versions hold critical details about the state of a package at a specific point in time, including release information, tags, metadata, and visibility.

## Table Usage Guide

The `github_package_version` table allows you to query detailed information about different versions of packages in GitHub's package registry. This includes data such as the package author, digest, release information, and the visibility of the version (whether public or private).

**Important Notes**
- You must specify the `organization` column in the `where` clause to query the table.
- OAuth app tokens and personal access tokens (classic) need the `read:packages` scope to use this endpoint. If the package_type belongs to a GitHub Packages registry that only supports repository-scoped permissions, the `repo` scope is also required. Please refer [Permissions for repository-scoped packages](https://docs.github.com/en/packages/learn-github-packages/about-permissions-for-github-packages#permissions-for-repository-scoped-packages) for more information.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### List all versions of a specific package
This query retrieves all the versions for a specific package in an organization. It includes details such as the version's digest, creation date, and visibility status.

```sql+postgres
select
  id,
  package_name,
  name as digest,
  prerelease,
  created_at,
  visibility
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws';
```

```sql+sqlite
select
  id,
  package_name,
  name as digest,
  prerelease,
  created_at,
  visibility
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws';
```

### List all public package versions for an organization
This query lists all publicly visible versions of packages in a GitHub organization, including their package type and associated metadata.

```sql+postgres
select
  package_name,
  package_type,
  name as digest,
  visibility,
  created_at,
  updated_at
from
  github_package_version
where
  organization = 'turbot'
  and visibility = 'public';
```

```sql+sqlite
select
  package_name,
  package_type,
  name as digest,
  visibility,
  created_at,
  updated_at
from
  github_package_version
where
  organization = 'turbot'
  and visibility = 'public';
```

### List pre-release package versions for a specific package
This query retrieves all the pre-release versions of a specific package in an organization. Pre-release versions are used for testing before a package is officially released.

```sql+postgres
select
  id,
  package_name,
  prerelease,
  created_at,
  html_url
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and prerelease = true;
```

```sql+sqlite
select
  id,
  package_name,
  prerelease,
  created_at,
  html_url
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and prerelease = 1;
```

### Get metadata of a specific package version
This query retrieves metadata details for a specific package version. Metadata can include additional version-specific information, such as description, version number, and other properties.

```sql+postgres
select
  id,
  package_name,
  jsonb_pretty(metadata) as metadata
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and id = 12345;
```

```sql+sqlite
select
  id,
  package_name,
  json_extract(metadata, '$') as metadata
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and id = 12345;
```

### List the tags associated with a package version
This query retrieves the tags associated with a specific package version. Tags help identify different aspects of a version and can assist with version management.

```sql+postgres
select
  id,
  package_name,
  jsonb_array_elements_text(tags) as tag
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and id = 12345;
```

```sql+sqlite
select
  id,
  package_name,
  json_extract(tags, '$[0]') as tag
from
  github_package_version
where
  organization = 'turbot'
  and package_name = 'steampipe/plugin/turbot/aws'
  and id = 12345;
```
