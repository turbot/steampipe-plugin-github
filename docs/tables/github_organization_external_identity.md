---
title: "Steampipe Table: github_organization_external_identity - Query GitHub Organization External Identities using SQL"
description: "Allows users to query GitHub Organization External Identities, providing information about the external identity of users in a GitHub organization."
folder: "Organization"
---

# Table: github_organization_external_identity - Query GitHub Organization External Identities using SQL

GitHub Organization External Identities represent the mapping between a GitHub user and their linked identity at an Identity Provider. It provides information about the external identity of users in a GitHub organization. This data is particularly useful for organizations that use SAML single sign-on (SSO) in conjunction with GitHub.

## Table Usage Guide

The `github_organization_external_identity` table provides insights into the external identities of users within a GitHub organization. As a GitHub organization administrator, this table can be used to gain a comprehensive understanding of the linked identities of users at your Identity Provider. This can be particularly useful when managing users in organizations that use SAML single sign-on (SSO) with GitHub.

**Important Notes**
- You must specify the `organization` column in `where` or `join` clause to query the table.
- To query this table using Fine-grained access tokens, the following permissions are required(The Fine-Grained access token should be created in Organization level):
  - **"Members" organization permissions (read)** – Required to access general organization metadata.

## Examples

### List external identities of an organization
This query is useful for gaining insights into the external identities associated with a specific organization. It allows you to identify the roles and usernames of these external identities, which can help in assessing the organization's security and access management structure.

```sql+postgres
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

```sql+sqlite
select
  guid,
  user_login,
  json_extract(saml_identity, '$.username') as saml_user,
  json_extract(scim_identity, '$.username') as scim_user,
  json_extract(organization_invitation, '$.role') as invited_role
from
  github_organization_external_identity
where
  organization = 'turbot';
```

### List external identities for all your organizations
This query is useful for gaining insights into the external identities linked to your organizations. It provides a comprehensive view of each user's login details, email, and role, which can be particularly beneficial for managing access and understanding user behavior across different platforms.

```sql+postgres
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

```sql+sqlite
select
  o.login as org,
  json_extract(o.saml_identity_provider, '$.sso_url') as sso_url,
  e.user_login,
  json_extract(e.user_detail, '$.email') as user_email,
  json_extract(e.saml_identity, '$.username') as saml_user,
  json_extract(e.scim_identity, '$.username') as scim_user,
  json_extract(e.organization_invitation, '$.role') as invited_role
from
  github_my_organization o
join
  github_organization_external_identity e
on 
  o.login = e.organization;
```