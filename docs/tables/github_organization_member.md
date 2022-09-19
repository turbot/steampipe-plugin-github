# Table: github_organization_member

The `github_organization_member` table can be used to query information about members of a organization. **You must specify the organization** in the where or join clause (`where organization=`, `join github_organization_member on organization=`).

## Examples

### List organization members

```sql
select
  organization,
  login,
  role,
  state
from
  github_organization_member
where
  organization = 'my_org';
```

### List active organization members who are admin

```sql
select
  organization,
  login,
  role,
  state
from
  github_organization_member
where
  organization = 'my_org'
  and role = 'admin'
  and state = 'active';
```

### List members with pending invitations

```sql
select
  organization,
  login,
  role,
  state
from
  github_organization_member
where
  organization = 'my_org'
  and state = 'pending';
```