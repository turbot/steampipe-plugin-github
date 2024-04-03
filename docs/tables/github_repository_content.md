---
title: "Steampipe Table: github_repository_content - Query File and Directory Contents in GitHub Repositories using SQL"
description: "Facilitates querying the contents of files and directories within GitHub repositories, offering insights into file paths, types, sizes, and more."
---

# Table: github_repository_content - Query File and Directory Contents in GitHub Repositories using SQL

The `github_repository_content` table is designed to fetch the contents of files or directories within a GitHub repository. It provides a detailed view of file paths, types, contents, sizes, and other related information.

## Table Usage Guide

To utilize this table effectively, specify the file path or directory within `repository_content_path`. If `repository_content_path` is not specified, the table will return the contents of the repository's root directory. This feature allows for comprehensive exploration of repository contents, from individual files to entire directories.

**Important Notes**
- It's mandatory to specify the `repository_full_name` (including the organization/user prefix) in the `where` or `join` clause when querying this table.

## Examples

### List the root directory contents of a repository
This query is useful for obtaining an overview of the root directory of a specific repository, helping users quickly identify the initial set of files and directories it contains.

```sql+postgres
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs';
```

```sql+sqlite
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs';
```

### List contents of a specific directory within a repository
This query facilitates a deeper inspection into a specific directory within a repository, enabling users to understand its structure and the types of files it contains.

```sql+postgres
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = 'docs';
```

```sql+sqlite
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = 'docs';
```

### Retrieve a specific file within a repository
Targeting a specific file within a repository, this query is particularly useful for extracting detailed information about a file, such as its content, type, and size, which is essential for analysis or integration purposes.

```sql+postgres
select
  repository_full_name,
  path,
  type,
  size,
  sha,
  content
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = '.vscode/settings.json';
```

```sql+sqlite
select
  repository_full_name,
  path,
  type,
  size,
  sha,
  content
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = '.vscode/settings.json';
```