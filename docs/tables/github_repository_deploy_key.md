---
title: "Steampipe Table: github_repository_deploy_key - Query GitHub Deploy Keys using SQL"
description: "Allows users to query deploy keys in GitHub repositories, providing insights into which SSH keys can read or push to a repository."
folder: "Repository"
---

# Table: github_repository_deploy_key - Query GitHub Deploy Keys using SQL

A GitHub deploy key is an SSH key that grants access to a single repository. It is read-only unless write access is granted when the key is added, in which case it can also push to the repository.

## Table Usage Guide

The `github_repository_deploy_key` table provides insights into the deploy keys of GitHub repositories. As a security engineer, explore key-specific details through this table, including whether each key can push, who added it and when it was last used.

The public key material is not exposed as a column.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Administration (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List deploy keys
Explore the deploy keys of a repository and whether each can push to it.

```sql+postgres
select
  id,
  title,
  read_only,
  added_by,
  created_at,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo';
```

```sql+sqlite
select
  id,
  title,
  read_only,
  added_by,
  created_at,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo';
```

### List deploy keys with write access
Identify deploy keys that can push to a repository, which should be read-only unless write access is required.

```sql+postgres
select
  id,
  title,
  added_by,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo'
  and not read_only;
```

```sql+sqlite
select
  id,
  title,
  added_by,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo'
  and read_only = 0;
```

### List deploy keys not used in the last 90 days
Find deploy keys that may no longer be needed.

```sql+postgres
select
  id,
  title,
  created_at,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo'
  and (last_used is null or last_used < now() - interval '90 days');
```

```sql+sqlite
select
  id,
  title,
  created_at,
  last_used
from
  github_repository_deploy_key
where
  repository_full_name = 'my-org/my-repo'
  and (last_used is null or last_used < datetime('now', '-90 days'));
```
