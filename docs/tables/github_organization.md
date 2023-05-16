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
