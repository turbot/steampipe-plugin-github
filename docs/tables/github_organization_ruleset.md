---
title: "Steampipe Table: github_organization_ruleset - Query GitHub Organization Rulesets using SQL"
description: "Allows users to query GitHub Organization Rulesets, providing details about each ruleset within an organization. This information includes ruleset ID, name, enforcement level, bypass actors, and more."
folder: "Organization"
---

# Table: github_organization_ruleset - Query GitHub Organization Rulesets using SQL

GitHub Organization Rulesets is a feature within GitHub that allows organizations to enforce rules and conditions across repositories. These rulesets help manage repository settings, permissions, and enforce best practices at the organization level.

## Table Usage Guide

The `github_organization_ruleset` table provides insights into the rulesets within a GitHub organization. As a security engineer or team lead, you can explore ruleset-specific details through this table, including ruleset ID, name, enforcement level, bypass actors, and conditions. Utilize it to enforce organization-wide policies, manage permissions, and ensure compliance with organizational standards.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required [according to the doc](https://docs.github.com/en/rest/orgs/rules?apiVersion=2026-03-10):
  - Organization permissions:
    - Administration (Write): Required to access organization rulesets.

**Important Notes**
- You must specify the `organization` column in the `where` or `join` clause to query the table.

## Examples

### List all rulesets in an organization
Explore all rulesets within a specific organization, including their enforcement levels and creation dates, to understand and manage organization-wide policies.

```sql+postgres
select
  name,
  enforcement,
  created_at
from
  github_organization_ruleset
where
  organization = 'my-org';
```

```sql+sqlite
select
  name,
  enforcement,
  created_at
from
  github_organization_ruleset
where
  organization = 'my-org';
```

### Get rules from a specific ruleset
Retrieve the detailed rules of a specific ruleset within your organization.

```sql+postgres
select
  name,
  r->>'id' as rule_id,
  r->>'type' as rule_type,
  r->>'parameters' as rule_parameters
from
  github_organization_ruleset,
  jsonb_array_elements(rules) as r
where
  organization = 'my-org'
  and name = 'my-ruleset';
```

```sql+sqlite
select
  name,
  json_extract(r.value, '$.id') as rule_id,
  json_extract(r.value, '$.type') as rule_type,
  json_extract(r.value, '$.parameters') as rule_parameters
from
  github_organization_ruleset,
  json_each(rules) as r
where
  organization = 'my-org'
  and name = 'my-ruleset';
```

### Get bypass actors for a specific ruleset
Identify the actors who can bypass the ruleset within your organization.

```sql+postgres
select
  name,
  b->>'id' as bypass_actor_id,
  b->>'deploy_key' as deploy_key,
  b->>'bypass_mode' as bypass_mode,
  b->>'repository_role_name' as repository_role_name,
  b->>'repository_role_database_id' as repository_role_database_id
from
  github_organization_ruleset,
  jsonb_array_elements(bypass_actors) as b
where
  organization = 'my-org'
  and name = 'my-ruleset';
```

```sql+sqlite
select
  name,
  json_extract(b.value, '$.id') as bypass_actor_id,
  json_extract(b.value, '$.deploy_key') as deploy_key,
  json_extract(b.value, '$.bypass_mode') as bypass_mode,
  json_extract(b.value, '$.repository_role_name') as repository_role_name,
  json_extract(b.value, '$.repository_role_database_id') as repository_role_database_id
from
  github_organization_ruleset,
  json_each(bypass_actors) as b
where
  organization = 'my-org'
  and name = 'my-ruleset';
```

### List rulesets with specific enforcement levels
Identify rulesets within an organization that have specific enforcement levels (`ACTIVE`, `EVALUATE`, `DISABLED`).

```sql+postgres
select
  name,
  enforcement
from
  github_organization_ruleset
where
  organization = 'my-org'
  and enforcement = 'ACTIVE';
```

```sql+sqlite
select
  name,
  enforcement
from
  github_organization_ruleset
where
  organization = 'my-org'
  and enforcement = 'EVALUATE';
```

### List all rulesets created after a specific date
Retrieve all rulesets created after a specified date, useful for auditing recent policy changes.

```sql+postgres
select
  name,
  created_at
from
  github_organization_ruleset
where
  organization = 'my-org'
  and created_at > '2023-01-01T00:00:00Z';
```

```sql+sqlite
select
  name,
  created_at
from
  github_organization_ruleset
where
  organization = 'my-org'
  and created_at > '2023-01-01T00:00:00Z';
```

### List rulesets with pull request parameters
List rules with pull request parameters, including code owner review requirements.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'PullRequestParameters' ->> 'require_code_owner_review' as require_code_owner_review,
  r -> 'parameters' -> 'PullRequestParameters' ->> 'required_approving_review_count' as required_approving_review_count
from
  github_organization_ruleset,
  jsonb_array_elements(rules) as r
where
  organization = 'my-org'
  and (r -> 'parameters' ->> 'Type') = 'PullRequestParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.PullRequestParameters.require_code_owner_review') as require_code_owner_review,
  json_extract(r.value, '$.parameters.PullRequestParameters.required_approving_review_count') as required_approving_review_count
from
  github_organization_ruleset,
  json_each(rules) as r
where
  organization = 'my-org'
  and json_extract(r.value, '$.parameters.Type') = 'PullRequestParameters';
```

### List rulesets with required status check parameters
List rules with required status check parameters.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'RequiredStatusChecksParameters' ->> 'required_status_checks' as required_status_checks
from
  github_organization_ruleset,
  jsonb_array_elements(rules) as r
where
  organization = 'my-org'
  and (r -> 'parameters' ->> 'Type') = 'RequiredStatusChecksParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.RequiredStatusChecksParameters.required_status_checks') as required_status_checks
from
  github_organization_ruleset,
  json_each(rules) as r
where
  organization = 'my-org'
  and json_extract(r.value, '$.parameters.Type') = 'RequiredStatusChecksParameters';
```
