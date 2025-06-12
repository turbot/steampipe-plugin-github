---
title: "Steampipe Table: github_repository_collaborator - Query GitHub Repository Collaborators using SQL"
description: "Allows users to query GitHub Repository Collaborators, providing detailed information about collaboration permissions and the status of invitations."
folder: "Collaborator"
---

# Table: github_repository_collaborator - Query GitHub Repository Collaborators using SQL

GitHub Repository Collaborators are users who have been granted access to contribute to a repository. Collaborators can be given different levels of access, from read-only to full admin rights. The status of their invitation, whether it's pending acceptance or has been accepted, is also tracked.

## Table Usage Guide

The `github_repository_collaborator` table provides insights into the collaborators associated with GitHub repositories. As a repository manager, you can use this table to explore details about collaborators, including their permissions and the status of their invitations. This can be particularly useful for managing access control and ensuring the appropriate level of access is granted to each collaborator.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Contents (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### List all contributors of a repository
Determine the areas in which various users have permissions within a specific project. This is useful for project managers who need to understand the roles and access levels of different contributors to better manage project resources and responsibilities.

```sql+postgres
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe';
```

### List all outside collaborators on a repository
Determine the areas in which outside collaborators have access within a specific repository. This is useful for ensuring appropriate access levels and identifying potential security risks.

```sql+postgres
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
  and affiliation = 'OUTSIDE';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
  and affiliation = 'OUTSIDE';
```

### List all repository admins
Identify instances where users have administrative access in a specific GitHub repository. This could be useful in managing access control and ensuring the right people have the appropriate permissions.

```sql+postgres
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
  and permission = 'ADMIN';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_repository_collaborator
where
  repository_full_name = 'turbot/steampipe'
  and permission = 'ADMIN';
```

### Obtain a JSON array of admins for all your repositories
Discover the segments that allow you to identify all the administrators for your GitHub repositories. This is useful for managing access and permissions across your repositories.

```sql+postgres
with repos as (
  select
    name_with_owner
  from
    github_my_repository
)
select
  r.name_with_owner as repo,
  json_agg(user_login) as admins
from
  repos as r
  inner join github_repository_collaborator as c on r.name_with_owner = c.repository_full_name and c.permission = 'ADMIN'
group by
  r.name_with_owner;
```

```sql+sqlite
with repos as (
  select
    name_with_owner
  from
    github_my_repository
)
select
  r.name_with_owner as repo,
  group_concat(user_login) as admins
from
  repos as r
  inner join github_repository_collaborator as c on r.name_with_owner = c.repository_full_name and c.permission = 'ADMIN'
group by
  r.name_with_owner;
```