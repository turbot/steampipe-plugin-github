---
title: "Steampipe Table: github_organization - Query GitHub Organizations using SQL"
description: "Allows users to query GitHub Organizations, specifically details about each organization's profile, including name, email, blog, location, and public repository count."
folder: "Organization"
---

# Table: github_organization - Query GitHub Organizations using SQL

GitHub Organizations is a feature within GitHub that allows users to collaborate across many projects at once. Organizations include features such as unified billing, access control, and multiple repositories. It is a way for businesses and open-source projects to manage their projects and teams.

## Table Usage Guide

The `github_organization` table provides insights into Organizations within GitHub. As a developer or project manager, explore organization-specific details through this table, including profile information, public repository count, and associated metadata. Utilize it to uncover information about organizations, such as their location, public repository count, and other profile details.

**Important Notes**

- You must specify the `login` column in `where` or `join` clause to query the table.
- To list organizations that you are a member of, use the `github_my_organization` table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required (the token must be created under the resource owner organization):

- Organization permissions:
  - Administration (Read-only): Required for the `interaction_ability` column.
  - Members (Read-only): Required for the `members_with_role_total_count` and `teams_total_count` columns.
  - Webhooks (Read-only): Required for the `hooks` column.

## Examples

### Basic info for a GitHub Organization

Explore essential details about a specific GitHub organization to understand its structure and activity. This is useful for gaining insights into the organization's verification status, team and member counts, and repository count.

```sql+postgres
select
  login as organization,
  name,
  twitter_username,
  created_at,
  updated_at,
  is_verified,
  teams_total_count as teams_count,
  members_with_role_total_count as member_count,
  repositories_total_count as repo_count
from
  github_organization
where
  login = 'postgres';
```

```sql+sqlite
select
  login as organization,
  name,
  twitter_username,
  created_at,
  updated_at,
  is_verified,
  teams_total_count as teams_count,
  members_with_role_total_count as member_count,
  repositories_total_count as repo_count
from
  github_organization
where
  login = 'postgres';
```

### List members of an organization

This query is used to identify members of a specific organization and check if they have two-factor authentication enabled. This can be useful for organizations looking to enforce security measures and ensure all members have additional protection for their accounts.

```sql+postgres
select
  o.login as organization,
  m.login as user_login,
  m.has_two_factor_enabled as mfa_enabled
from
  github_organization o,
  github_organization_member m
where
  o.login = 'turbot'
and
  o.login = m.organization;
```

```sql+sqlite
select
  o.login as organization,
  m.login as user_login,
  m.has_two_factor_enabled as mfa_enabled
from
  github_organization o
join
  github_organization_member m on o.login = m.organization
where
  o.login = 'turbot';
```

OR

```sql+postgres
select
  organization,
  login as user_login,
  has_two_factor_enabled as mfa_enabled
from
  github_organization_member
where
  organization = 'turbot';
```

```sql+sqlite
select
  organization,
  login as user_login,
  has_two_factor_enabled as mfa_enabled
from
  github_organization_member
where
  organization = 'turbot';
```
