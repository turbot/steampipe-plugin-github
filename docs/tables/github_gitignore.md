---
title: "Steampipe Table: github_gitignore - Query GitHub Gitignore using SQL"
description: "Allows users to query GitHub Gitignore templates, providing a list of all templates available for use on GitHub."
folder: "Repository"
---

# Table: github_gitignore - Query GitHub Gitignore using SQL

GitHub Gitignore is a feature that specifies intentionally untracked files to ignore when using Git. These are typically files that are generated during the build process or at runtime, and are not part of the source code repository. GitHub provides a list of templates for different programming languages and environments, which can be used to generate a `.gitignore` file for a new project.

## Table Usage Guide

The `github_gitignore` table provides insights into Gitignore templates available on GitHub. As a developer or a DevOps engineer, explore template-specific details through this table, including the name of the template and its source. Utilize it to discover available templates, understand their structure, and decide which one to use for your project.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### List basic gitignore info
Explore the various types of gitignore files in your GitHub repositories, in an organized manner. This will help you understand the different ignore rules applied across your projects, promoting better code management.

```sql+postgres
select
  *
from
  github_gitignore
order by
  name;
```

```sql+sqlite
select
  *
from
  github_gitignore
order by
  name;
```

### View the source of the Go template
Uncover the details of the original Go template from GitHub. This query can help developers understand the foundation of Go programming language templates, providing a starting point for potential customization or improvement.

```sql+postgres
select
  source
from
  github_gitignore
where
  name = 'Go';
```

```sql+sqlite
select
  source
from
  github_gitignore
where
  name = 'Go';
```