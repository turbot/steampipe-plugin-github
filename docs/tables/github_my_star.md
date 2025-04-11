---
title: "Steampipe Table: github_my_star - Query GitHub My Stars using SQL"
description: "Allows users to query My Stars in GitHub, specifically the repositories starred by the authenticated user, providing insights into user preferences and potential areas of interest."
folder: "Repository"
---

# Table: github_my_star - Query GitHub My Stars using SQL

GitHub Stars is a feature within GitHub that allows users to bookmark repositories for later reference. Users can star repositories to keep track of projects they find interesting, even if they do not directly contribute to them. This feature serves as a way to show appreciation for the repository maintainers' work and also to keep track of repositories for later use.

## Table Usage Guide

The `github_my_star` table provides insights into the repositories starred by the authenticated GitHub user. As a developer or project manager, explore details through this table, including repository names, owners, and star creation dates. Utilize it to analyze user preferences, discover potential areas of interest, and manage your starred repositories effectively.

## Examples

### List of your starred repositories
Explore which repositories you've starred on Github, allowing you to quickly access and review your favorite projects. This is particularly useful for keeping track of repositories you find interesting or intend to contribute to in the future.

```sql+postgres
select
  starred_at,
  repository_full_name,
  url
from
  github_my_star;
```

```sql+sqlite
select
  starred_at,
  repository_full_name,
  url
from
  github_my_star;
```