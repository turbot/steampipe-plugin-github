---
title: "Steampipe Table: github_organization_member - Query GitHub Organization Members using SQL"
description: "Allows users to query GitHub Organization Members, specifically the details of members within an organization, providing insights into member profiles, roles, and status."
folder: "Organization"
---

# Table: github_organization_member - Query GitHub Organization Members using SQL

GitHub Organization Members is a feature within GitHub that allows you to manage and coordinate teams and repositories within your organization. It provides a centralized way to manage permissions for various repositories, assign roles to members, and monitor the status of each member. GitHub Organization Members helps you maintain control over the resources within your organization and manage member access effectively.

## Table Usage Guide

The `github_organization_member` table provides insights into members within a GitHub organization. As a project manager or team leader, explore member-specific details through this table, including roles, permissions, and status. Utilize it to uncover information about members, such as their roles within the organization, their access permissions, and their activity status.

**Important Notes**

- You must specify the `organization` column in `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required (the token must be created under the resource owner organization):

- Organization permissions:
  - Members (Read-only): Required to access general organization metadata.
- Repository permissions:
  - Metadata (Read-only): Required to access general repository metadata.

## Examples

### List organization members

Explore which members belong to your organization and their respective roles, while also identifying if they have two-factor authentication enabled. This can enhance your organization's security by ensuring all members have this additional layer of protection.

```sql+postgres
select
  organization,
  login,
  role,
  has_two_factor_enabled
from
  github_organization_member
where
  organization = 'my_org';
```

```sql+sqlite
select
  organization,
  login,
  role,
  has_two_factor_enabled
from
  github_organization_member
where
  organization = 'my_org';
```

### List admin members with two factor authentication disabled

Identify instances where administrative members in your organization have not enabled two-factor authentication, allowing you to enhance your organization's security by addressing these vulnerabilities.

```sql+postgres
select
  organization,
  login,
  role,
  has_two_factor_enabled
from
  github_organization_member
where
  organization = 'my_org'
  and role = 'ADMIN'
  and not has_two_factor_enabled;
```

```sql+sqlite
select
  organization,
  login,
  role,
  has_two_factor_enabled
from
  github_organization_member
where
  organization = 'my_org'
  and role = 'ADMIN'
  and not has_two_factor_enabled;
```
