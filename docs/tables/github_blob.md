---
title: "Steampipe Table: github_blob - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories, specifically the blob content."
---

# Table: github_blob - Query GitHub Repositories using SQL

GitHub Repositories are a fundamental resource in GitHub. They allow users to host and manage their codebase, track changes, and collaborate with other users. Each repository contains a tree structure that represents the file and directory hierarchy.

## Table Usage Guide

The `github_blob` table provides the contents of files within GitHub Repositories. As a developer or project manager, explore each repository's contents. Utilize it to uncover information about specific files in the repository, such as configuration files.

**Important Notes**
- You must specify the `repository_full_name` and `blob_sha` columns in `where` or `join` clause to query the table.

## Examples

### List blob contents for a specific file
Explore the specific elements within the 'turbot/steampipe' repository by pinpointing specific files using a unique identifier. This helps understand the contents of the repository and makes it easy to access and analyze the file's content.

```sql+postgres
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha,
  decode(content, encoding) as content
from
  github_tree t
left outer join
  github_blob b on b.repository_full_name = t.repository_full_name and b.blob_sha = t.sha
where
  t.repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and path = 'Makefile';
```

```sql+sqlite
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha,
  decode(content, encoding) as content
from
  github_tree t
left outer join
  github_blob b on b.repository_full_name = t.repository_full_name and b.blob_sha = t.sha
where
  t.repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and path = 'Makefile';
```