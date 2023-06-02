# Table: github_my_organization

Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once. Owners and administrators can manage member access to the organization's data and projects with sophisticated security and administrative features.

The `github_my_organization` table will list the organization **that you are a member of**. To view details of **ANY** organization, use the `github_organization` table.

## Examples

### Basic info for the GitHub Organizations to which you belong

```sql
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

```sql
select
  o.login as organization,
  m.login as member_login
from
  github_my_organization o
join github_organization_member m
on o.login = m.organization;
```

### Show your permissions on the Organization

```sql
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

```sql
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

```sql
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