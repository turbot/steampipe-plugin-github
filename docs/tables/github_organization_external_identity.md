# Table: github_organization_external_identity

The `github_organization_external_identity` used to query information about external identities of an organization.

**You must specify the organization** in the where or join clause (`where organization=`, `join github_organization_external_identity on organization=`).

## Examples

### List external identities of an organization

```sql
select
  guid,
  user_login,
  saml_identity ->> 'username' as saml_user,
  scim_identity ->> 'username' as scim_user,
  organization_invitation ->> 'role' as invited_role
from
  github_organization_external_identity
where
  organization = 'turbot';
```

### List external identities for all your organizations

```sql
select
  o.login as org,
  o.saml_identity_provider ->> 'sso_url' as sso_url,
  e.user_login,
  e.user_detail ->> 'email' as user_email,
  e.saml_identity ->> 'username' as saml_user,
  e.scim_identity ->> 'username' as scim_user,
  e.organization_invitation ->> 'role' as invited_role
from
  github_my_organization o
join
  github_organization_external_identity e
on 
  o.login = e.organization;
```