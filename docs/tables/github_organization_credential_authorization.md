---
title: "Steampipe Table: github_organization_credential_authorization - Query GitHub Organization Credential Authorizations using SQL"
description: "Allows users to query classic personal access tokens (and SSH keys) that members have authorized to access an organization's resources through SAML single sign-on (SSO)."
folder: "Organization"
---

# Table: github_organization_credential_authorization - Query GitHub Organization Credential Authorizations using SQL

GitHub organizations that enforce SAML single sign-on (SSO) require members to authorize their classic personal access tokens (PATs) and SSH keys before those credentials can access organization resources. A credential authorization records the token (or key) that a member has approved for use against the organization, including who it belongs to, its scopes, when it was authorized, when it was last used, and when it expires.

## Table Usage Guide

The `github_organization_credential_authorization` table provides insights into the classic personal access tokens authorized against a GitHub organization. As a security or compliance engineer, use this table to inventory the classic PATs that can reach your organization's resources, review the scopes they have been granted, and identify tokens that are stale, never used, or approaching expiration. The same endpoint also returns authorized SSH keys, which can be distinguished using the `credential_type` column.

To query this table, the credential must be created under an organization owner. Classic personal access tokens require the `admin:org` scope (specifically `read:org`).

**Important Notes**
- You must specify the `organization` column in the `where` or `join` clause to query the table.
- This table is only available for organizations on **GitHub Enterprise Cloud** that have **SAML single sign-on (SSO)** enabled. Organizations without SAML SSO return no rows.
- The underlying [SAML SSO authorizations API](https://docs.github.com/en/enterprise-cloud@latest/rest/orgs/orgs#list-saml-sso-authorizations-for-an-organization) returns both personal access tokens and SSH keys. Filter on `credential_type = 'personal access token'` to return only classic PATs.

## Examples

### List all credential authorizations for an organization
Inventory every classic personal access token and SSH key that has been authorized to access your organization's resources.

```sql+postgres
select
  login,
  credential_id,
  credential_type,
  scopes,
  credential_authorized_at,
  credential_accessed_at,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org';
```

```sql+sqlite
select
  login,
  credential_id,
  credential_type,
  scopes,
  credential_authorized_at,
  credential_accessed_at,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org';
```

### List only classic personal access tokens
Focus on classic PATs (excluding SSH keys) to review the tokens that can access your organization.

```sql+postgres
select
  login,
  credential_id,
  token_last_eight,
  scopes,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and credential_type = 'personal access token';
```

```sql+sqlite
select
  login,
  credential_id,
  token_last_eight,
  scopes,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and credential_type = 'personal access token';
```

### Find tokens that have expired or expire within 30 days
Identify personal access tokens that need to be rotated or removed because they are expired or about to expire.

```sql+postgres
select
  login,
  credential_id,
  scopes,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and authorized_credential_expires_at is not null
  and authorized_credential_expires_at < now() + interval '30 days'
order by
  authorized_credential_expires_at;
```

```sql+sqlite
select
  login,
  credential_id,
  scopes,
  authorized_credential_expires_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and authorized_credential_expires_at is not null
  and authorized_credential_expires_at < datetime('now', '+30 days')
order by
  authorized_credential_expires_at;
```

### Find tokens that have never been used
Spot credentials that were authorized but never accessed, which are good candidates for revocation.

```sql+postgres
select
  login,
  credential_id,
  credential_type,
  credential_authorized_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and credential_accessed_at is null;
```

```sql+sqlite
select
  login,
  credential_id,
  credential_type,
  credential_authorized_at
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and credential_accessed_at is null;
```

### Find tokens with highly privileged scopes
Surface classic PATs that have been granted broad access such as full `repo` or organization administration scopes.

```sql+postgres
select
  login,
  credential_id,
  scopes
from
  github_organization_credential_authorization
where
  organization = 'my_org'
  and scopes ?| array['repo', 'admin:org', 'delete_repo'];
```

```sql+sqlite
select
  c.login,
  c.credential_id,
  c.scopes
from
  github_organization_credential_authorization as c,
  json_each(c.scopes) as s
where
  c.organization = 'my_org'
  and s.value in ('repo', 'admin:org', 'delete_repo');
```

### Count authorized tokens per member
Understand which members have the most credentials authorized against the organization.

```sql+postgres
select
  login,
  count(*) as authorized_credentials
from
  github_organization_credential_authorization
where
  organization = 'my_org'
group by
  login
order by
  authorized_credentials desc;
```

```sql+sqlite
select
  login,
  count(*) as authorized_credentials
from
  github_organization_credential_authorization
where
  organization = 'my_org'
group by
  login
order by
  authorized_credentials desc;
```
