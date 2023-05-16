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
