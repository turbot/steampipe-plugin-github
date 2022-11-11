# Table: github_organization_member

The `github_organization_member` table can be used to query information about members of an organization. You must be an owner of the organization in order to successfully query member role and two factor authentication information. If you are not an owner of the organization, these columns will be returned as `null`.

**You must specify the organization** in the where or join clause (`where organization=`, `join github_organization_member on organization=`).

## Examples

### List organization members

```sql
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

```sql
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
