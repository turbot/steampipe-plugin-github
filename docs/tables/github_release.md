---
title: "Steampipe Table: github_release - Query GitHub Releases using SQL"
description: "Allows users to query GitHub Releases, specifically providing detailed information about each release of a repository, including release id, tag name, draft status, prerelease status, and more."
folder: "Release"
---

# Table: github_release - Query GitHub Releases using SQL

GitHub Releases is a feature of GitHub that allows you to present significant points in your repository's history, such as milestone versions, by associating them with tags. You can use GitHub Releases to manage and upload binary files, as well as providing release notes and links to binary files, directly from a repository. GitHub Releases is a great way to package software, release notes, and links to binary files for other people to use.

## Table Usage Guide

The `github_release` table provides insights into GitHub Releases within a repository. As a software developer or project manager, explore release-specific details through this table, including release id, tag name, draft status, prerelease status, and more. Utilize it to track the progress and status of different versions of your software, identify any prerelease versions, and manage your software releases more effectively.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Contents (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List releases
Explore the timeline of updates for the Steampipe project on Github. This allows you to track the progression of the project over time, helping you stay updated on new releases and changes.

```sql+postgres
select
  name,
  published_at
from
  github_release
where
  repository_full_name = 'turbot/steampipe'
order by
  published_at desc;
```

```sql+sqlite
select
  name,
  published_at
from
  github_release
where
  repository_full_name = 'turbot/steampipe'
order by
  published_at desc;
```

### Download statistics per release
Explore the popularity of different Steampipe releases by tracking the number of downloads. This can help in understanding user preferences and identifying the most successful releases.

```sql+postgres
select
  r.name as release_name,
  r.published_at,
  a ->> 'name' as asset_name,
  a ->> 'download_count' as download_count
from
  github_release as r,
  jsonb_array_elements(assets) as a
where
  r.repository_full_name = 'turbot/steampipe'
  and a ->> 'content_type' in ('application/zip', 'application/gzip')
order by
  r.published_at desc,
  asset_name;
```

```sql+sqlite
select
  r.name as release_name,
  r.published_at,
  json_extract(a.value, '$.name') as asset_name,
  json_extract(a.value, '$.download_count') as download_count
from
  github_release as r,
  json_each(assets) as a
where
  r.repository_full_name = 'turbot/steampipe'
  and json_extract(a.value, '$.content_type') in ('application/zip', 'application/gzip')
order by
  r.published_at desc,
  asset_name;
```