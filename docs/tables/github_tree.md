---
title: "Steampipe Table: github_tree - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories, specifically the tree structures, providing insights into the file and directory hierarchy of each repository."
folder: "Repository"
---

# Table: github_tree - Query GitHub Repositories using SQL

GitHub Repositories are a fundamental resource in GitHub. They allow users to host and manage their codebase, track changes, and collaborate with other users. Each repository contains a tree structure that represents the file and directory hierarchy.

## Table Usage Guide

The `github_tree` table provides insights into the tree structures within GitHub Repositories. As a developer or project manager, explore each repository's file and directory hierarchy through this table, including file names, types, and associated metadata. Utilize it to uncover information about the organization and structure of repositories, such as the distribution of file types and the depth of directory nesting.

**Important Notes**
- You must specify the `repository_full_name` and `tree_sha` columns in `where` or `join` clause to query the table.

## Examples

### List tree entries non-recursively
Explore the specific elements within the 'turbot/steampipe' repository by pinpointing specific locations using a unique identifier. This allows for a non-recursive view of the repository's structure, enabling easier navigation and understanding of the repository's layout and content.

```sql+postgres
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a';
```

```sql+sqlite
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a';
```

### List tree entries for a subtree recursively
Determine the areas in which you can explore the components of a specific subtree within the 'turbot/steampipe' repository. This is useful for gaining insights into the structure and elements of the subtree in a recursive manner.

```sql+postgres
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '5622172b528cd38438c52ecfa3c20ac3f71dd2df'
  and recursive = true;
```

```sql+sqlite
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '5622172b528cd38438c52ecfa3c20ac3f71dd2df'
  and recursive = 1;
```

### List executable files
This query allows you to identify all the executable files within a specified repository. It's particularly useful for understanding the structure and content of a repository, and for identifying potential security risks associated with executable files.

```sql+postgres
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = true
  and mode = '100755';
```

```sql+sqlite
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = 1
  and mode = '100755';
```

### List JSON files
This query is useful for identifying all JSON files within a specific GitHub repository. It can help developers or project managers to quickly locate and manage all JSON files in the repository, aiding in tasks such as code review, debugging, or configuration management.

```sql+postgres
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = true
  and path like '%.json';
```

```sql+sqlite
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = 1
  and path like '%.json';
```