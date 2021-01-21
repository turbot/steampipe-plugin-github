# Table: github_organization

Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once. Owners and administrators can manage member access to the organization's data and projects with sophisticated security and administrative features.

The `github_organization` table will list the organization **that you are a member of**.  You can query **ANY** organization that you have access to by specifying its `login` explicitly in the where clause with  `where login=`  .

## Examples

### Basis info for the Github Organizations to which you belong

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
  github_organization;
```


### Show members of an organization

```sql
select
  login as organization,
  name,
  m ->> 'login' as member_login,
  m ->> 'type' as member_type
from
  github_organization,
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
  members_can_create_repositories,
  default_repo_permission,
  two_factor_requirement_enabled
from
  github_organization;
```
