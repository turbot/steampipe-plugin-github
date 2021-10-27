# Table: github_organization

Organizations are shared accounts where businesses and open-source projects can collaborate across many projects at once. Owners and administrators can manage member access to the organization's data and projects with sophisticated security and administrative features.

You can query details for **ANY** organization with the `github_organization` table, but you must specify the `login` explicitly in the where or join clause (`where login=`, `join github_organization on login=`).

To list organizations **that you are a member of**, use the `github_my_organization` table.

## Examples

### Basic info for a GitHub Organization

```sql
select
  login as organization,
  name,
  twitter_username,
  public_repos,
  public_gists,
  member_logins
from
  github_organization
where
  login = 'postgres';
```

### List members of an organization

```sql
select
  login as organization,
  name,
  member_login
from
  github_organization,
  jsonb_array_elements_text(member_logins) as member_login
where
  login = 'google';
```
