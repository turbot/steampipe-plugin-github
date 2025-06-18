---
title: "Steampipe Table: github_organization_collaborator - Query GitHub Organization Collaborators using SQL"
description: "Allows users to query GitHub Organization Collaborators, providing detailed information about collaboration permissions and the status of invitations."
folder: "Collaborator"
---

# Table: github_organization_collaborator - Query GitHub Organization Collaborators using SQL

GitHub Organization Collaborators are users who have been granted access to contribute to a repository within an organization. Collaborators can be given different levels of access, from read-only to full admin rights. The status of their invitation, whether it's pending acceptance or has been accepted, is also tracked.

## Table Usage Guide

The `github_organization_collaborator` table provides insights into the collaborators associated with GitHub repository within an organization. As a repository manager, you can use this table to explore details about collaborators, including their permissions and the status of their invitations. This can be particularly useful for managing access control and ensuring the appropriate level of access is granted to each collaborator.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required (the token must be created under the resource owner organization):
  - Repository permissions:
    - Contents (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `organization` column in the `where` or `join` clause to query the table.

## Examples

### List all contributors of an organization
Determine the areas in which various users have permissions within a specific project. This is useful for project managers who need to understand the roles and access levels of different contributors to better manage project resources and responsibilities.

```sql+postgres
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot';
```

### List all outside collaborators on a repository
Determine the areas in which outside collaborators have access within a specific repository. This is useful for ensuring appropriate access levels and identifying potential security risks.

```sql+postgres
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot'
  and affiliation = 'OUTSIDE';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot'
  and affiliation = 'OUTSIDE';
```

### List all collaborators with admins role
Identify instances where users have administrative access in a specific GitHub repository. This could be useful in managing access control and ensuring the right people have the appropriate permissions.

```sql+postgres
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot'
  and permission = 'ADMIN';
```

```sql+sqlite
select
  user_login,
  permission
from
  github_organization_collaborator
where
  organization = 'turbot'
  and permission = 'ADMIN';
```