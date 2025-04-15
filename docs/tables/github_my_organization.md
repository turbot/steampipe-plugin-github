---
title: "Steampipe Table: github_my_organization - Query GitHub Organizations using SQL"
description: "Allows users to query GitHub Organizations, specifically details about a user's primary organization on GitHub. This table provides insights into the organization's details, members, teams, repositories, and more."
folder: "Organization"
---

# Table: github_my_organization - Query GitHub Organizations using SQL

GitHub Organizations are a shared workspace where businesses and open-source projects can collaborate across many projects at once. Owners and administrators can manage member access to the organization's data and projects with sophisticated security and administrative features. GitHub Organizations help streamline project management and boost productivity in your organization.

## Table Usage Guide

The `github_my_organization` table provides insights into a user's primary organization on GitHub. As a project manager or team lead, explore organization-specific details through this table, including member access, repository details, and administrative features. Utilize it to manage and streamline your organization's projects and enhance team productivity.

**Important Notes**
- The `github_my_organization` table will list the organization **that you are a member of**. To view details of **ANY** organization, use the `github_organization` table.

## Examples

### Basic info for the GitHub Organizations to which you belong
Determine the areas in which you hold membership across various GitHub organizations. This query is useful in understanding your involvement and role within these organizations, including details such as the number of private and public repositories, team counts, and member counts.

```sql+postgres
select
  login as organization,
  name,
  twitter_username,
  private_repositories_total_count as private_repos,
  public_repositories_total_count as public_repos,
  created_at,
  updated_at,
  is_verified,
  teams_total_count as teams_count,
  members_with_role_total_count as member_count,
  url
from
  github_my_organization;
```

```sql+sqlite
select
  login as organization,
  name,
  twitter_username,
  private_repositories_total_count as private_repos,
  public_repositories_total_count as public_repos,
  created_at,
  updated_at,
  is_verified,
  teams_total_count as teams_count,
  members_with_role_total_count as member_count,
  url
from
  github_my_organization;
```

### Show all members for the GitHub Organizations to which you belong
Determine the areas in which you are a member of a GitHub organization, offering insights into your collaborative coding environments and affiliations. This can be useful in managing and understanding your participation in various coding projects and teams.

```sql+postgres
select
  o.login as organization,
  m.login as member_login
from
  github_my_organization o
  join github_organization_member m on o.login = m.organization;
```

```sql+sqlite
select
  o.login as organization,
  m.login as member_login
from
  github_my_organization o
  join github_organization_member m on o.login = m.organization;
```

### Show your permissions on the Organization
Explore your access level and permissions within your GitHub organization. This can help in understanding what actions you are authorized to perform, such as administering the organization, changing pinned items, creating projects, repositories, or teams, and whether you are currently a member.

```sql+postgres
select
  login as organization,
  members_with_role_total_count as members_count,
  can_administer,
  can_changed_pinned_items,
  can_create_projects,
  can_create_repositories,
  can_create_teams,
  is_a_member as current_member
from
  github_my_organization;
```

```sql+sqlite
select
  login as organization,
  members_with_role_total_count as members_count,
  can_administer,
  can_changed_pinned_items,
  can_create_projects,
  can_create_repositories,
  can_create_teams,
  is_a_member as current_member
from
  github_my_organization;
```

### Show Organization security settings
Gain insights into your organization's security settings, such as member permissions and two-factor authentication requirements. This can help ensure your organization's GitHub repositories and pages are appropriately protected.

```sql+postgres
select
  login as organization,
  members_with_role_total_count as members_count,
  members_allowed_repository_creation_type,
  members_can_create_internal_repos,
  members_can_create_pages,
  members_can_create_private_repos,
  members_can_create_public_repos,
  members_can_create_repos,
  default_repo_permission,
  two_factor_requirement_enabled
from
  github_my_organization;
```

```sql+sqlite
select
  login as organization,
  members_with_role_total_count as members_count,
  members_allowed_repository_creation_type,
  members_can_create_internal_repos,
  members_can_create_pages,
  members_can_create_private_repos,
  members_can_create_public_repos,
  members_can_create_repos,
  default_repo_permission,
  two_factor_requirement_enabled
from
  github_my_organization;
```

### List organization hooks that are insecure
Explore which organization hooks are potentially insecure due to specific settings, such as lack of SSL security, absence of a secret, or non-HTTPS URLs. This is particularly useful in identifying and mitigating potential security vulnerabilities within your organization's GitHub configuration.

```sql+postgres
select
  login as organization,
  hook
from
  github_my_organization,
  jsonb_array_elements(hooks) as hook
where
  hook -> 'config' ->> 'insecure_ssl' = '1'
  or hook -> 'config' ->> 'secret' is null
  or hook -> 'config' ->> 'url' not like '%https:%';
```

```sql+sqlite
select
  login as organization,
  hook.value as hook
from
  github_my_organization,
  json_each(hooks) as hook
where
  json_extract(hook.value, '$.config.insecure_ssl') = '1'
  or json_extract(hook.value, '$.config.secret') is null
  or json_extract(hook.value, '$.config.url') not like '%https:%';
```