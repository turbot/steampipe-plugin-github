---
title: "Steampipe Table: github_community_profile - Query Github Community Profiles using SQL"
description: "Allows users to query Github Community Profiles, providing insights into the health and community activity of repositories."
folder: "Community"
---

# Table: github_community_profile - Query Github Community Profiles using SQL

A Github Community Profile is a measure of the health and community activity of a repository. It includes factors such as the presence of a code of conduct, contributing guidelines, and issue and pull request templates. A healthy community profile indicates a well-maintained, welcoming, and sustainable open-source project.

## Table Usage Guide

The `github_community_profile` table provides insights into the health and community activity of Github repositories. As a repository owner or contributor, explore details through this table, including the presence of contributing guidelines, code of conduct, and issue templates. Utilize it to evaluate the health of repositories, understand their community engagement, and identify areas for improvement.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Contents (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### Get community profile information for the repository
Explore the community profile details for a specific repository to gain insights into the repository's health and activity. This is useful for assessing the repository's contribution guidelines, code of conduct, and other community-related aspects.

```sql+postgres
select
  *
from
  github_community_profile
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  *
from
  github_community_profile
where
  repository_full_name = 'turbot/steampipe';
```

### List repositories having the security file
This example helps identify the GitHub repositories that have a security file in place. This is useful for understanding the security measures implemented in each repository, which is crucial for maintaining project integrity and preventing unauthorized access.

```sql+postgres
select
  repository_full_name,
  security ->> 'text' as security_file_content
from
  github_community_profile c
  join github_my_repository r on r.name_with_owner = c.repository_full_name
  where security is not null;
```

```sql+sqlite
select
  repository_full_name,
  json_extract(security, '$.text') as security_file_content
from
  github_community_profile c
  join github_my_repository r on r.name_with_owner = c.repository_full_name
  where security is not null;
```