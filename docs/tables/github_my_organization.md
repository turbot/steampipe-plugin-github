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
  total_private_repos,
  public_repos,
  plan_name,
  plan_seats,
  plan_filled_seats
from
  github_my_organization;
```


### Show members of an organization

```sql
select
  login as organization,
  name,
  m ->> 'login' as member_login,
  m ->> 'type' as member_type
from
  github_my_organization,
  jsonb_array_elements(members) as m
where
  login = 'turbot';
```


### Show Organization security settings

```sql
select
  login as organization,
  jsonb_array_length(members) as num_members,
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


### Show collaborators in your organization's repositories that are not members of the organization

```sql
select
  r.name,
  collaborator_login
from
  github_my_repository as r,
  jsonb_array_elements_text(r.collaborator_logins) as collaborator_login,
  github_my_organization as o
where
  r.owner_login = o.login
  and collaborator_login not in (
    select m from github_my_organization, jsonb_array_elements_text(member_logins) as m
  ) ;
```
