---
title: "Steampipe Table: github_team_repository - Query GitHub Teams Repositories using SQL"
description: "Allows users to query GitHub Teams Repositories, specifically the association between teams and repositories within a GitHub organization, providing insights into team access to specific repositories."
---

# Table: github_team_repository - Query GitHub Teams Repositories using SQL

GitHub Teams Repositories represent the association between teams and repositories within a GitHub organization. Teams are groups of organization members that reflect the company or project structure in your organization. Repositories are where all your project files (including documentation) reside and teams can be granted access to these repositories.

## Table Usage Guide

The `github_team_repository` table provides insights into the association between teams and repositories within a GitHub organization. As a project manager or team lead, explore team-specific access details through this table, including repository permissions and associated metadata. Utilize it to uncover information about team access to repositories, such as those with admin permissions, the relationship between teams and repositories, and the verification of access policies.

**Important Notes**
- You must specify the `organization` and `slug` columns in the `where` or `join` clause to query the table.
- To list all your repositories use the `github_my_repository` table instead. To get information about any repository, use the `github_repository` table instead.

## Examples

### List a specific team's repositories
This query is designed to help you gain insights into the details of a specific team's repositories within your organization on GitHub. It is useful for understanding the team's repository permissions, primary language, fork count, stargazer count, license information, and other relevant details, which can help in managing team resources and identifying areas for improvement.

```sql+postgres
select
  organization,
  slug as team_slug,
  name as team_name,
  permission,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  license_info ->> 'spdx_id' as license,
  description,
  url
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team';
```

```sql+sqlite
select
  organization,
  slug as team_slug,
  name as team_name,
  permission,
  (primary_language ->> 'name') as language,
  fork_count,
  stargazer_count,
  (license_info ->> 'spdx_id') as license,
  description,
  url
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team';
```

### List visible teams and repositories they have admin permissions to
Explore the teams and associated repositories within your organization that have administrative permissions. This is useful to ensure appropriate access rights and maintain security within your GitHub organization.

```sql+postgres
select
  organization,
  slug as team_slug,
  name as name,
  description,
  permission,
  is_fork,
  is_private,
  is_archived,
  primary_language ->> 'name' as language
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team'
  and permission = 'ADMIN';
```

```sql+sqlite
select
  organization,
  slug as team_slug,
  name as name,
  description,
  permission,
  is_fork,
  is_private,
  is_archived,
  json_extract(primary_language, '$.name') as language
from
  github_team_repository
where
  organization = 'my_org'
  and slug = 'my-team'
  and permission = 'ADMIN';
```