---
title: "Steampipe Table: github_actions_repository_secret - Query GitHub Actions Repository Secrets using SQL"
description: "Allows users to query GitHub Actions Repository Secrets, specifically to retrieve information about secrets stored in a GitHub repository, providing insights into the security measures in place."
---

# Table: github_actions_repository_secret - Query GitHub Actions Repository Secrets using SQL

GitHub Actions is a service provided by GitHub that allows you to automate, customize, and execute your software development workflows right in your repository. GitHub Actions makes it easy to automate all your software workflows, now with world-class CI/CD. It enables you to build, test, and deploy your code right from GitHub.

## Table Usage Guide

The `github_actions_repository_secret` table provides insights into secrets stored within a GitHub repository. As a security engineer, explore secret-specific details through this table, including the names of secrets and the dates they were created or updated. Utilize it to uncover information about secrets, such as those that may be outdated or unused, providing a comprehensive view of the repository's security measures.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List secrets
Explore the hidden aspects of a specific repository within the GitHub Actions environment. This is useful in assessing the security and integrity of the repository.

```sql+postgres
select
  *
from
  github_actions_repository_secret
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  *
from
  github_actions_repository_secret
where
  repository_full_name = 'turbot/steampipe';
```