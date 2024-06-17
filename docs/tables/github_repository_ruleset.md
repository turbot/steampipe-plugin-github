---
title: "Steampipe Table: github_repository_ruleset - Query GitHub Repository Rulesets using SQL"
description: "Allows users to query GitHub Repository Rulesets, providing details about each ruleset within a repository. This information includes ruleset ID, name, enforcement level, bypass actors, and more."
---

# Table: github_repository_ruleset - Query GitHub Repository Rulesets using SQL

GitHub Repository Rulesets is a feature within GitHub that allows organizations to enforce rules and conditions on repositories. These rulesets help manage repository settings, permissions, and enforce best practices.

## Table Usage Guide

The `github_repository_ruleset` table provides insights into the rulesets within GitHub repositories. As a project manager or team lead, you can explore ruleset-specific details through this table, including ruleset ID, name, enforcement level, bypass actors, and conditions. Utilize it to enforce repository policies, manage permissions, and ensure compliance with organizational standards.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List all rulesets in a repository
Explore all rulesets within a specific repository, including their enforcement levels and creation dates, to understand and manage repository policies.

```sql+postgres
select
  name,
  enforcement,
  created_at
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule';
```

```sql+sqlite
select
  name,
  enforcement,
  created_at
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule';
```

### Get rules from a specific ruleset
Retrieve the detailed rules of a specific ruleset within your repository. This can be useful for reviewing the rules enforced and ensuring they align with your project requirements.

```sql+postgres
select
  name,
  r->>'id' as rule_id,
  r->>'type' as rule_type,
  r->>'parameters' as rule_parameters
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and name = 'test34';
```

```sql+sqlite
select
  name,
  json_extract(r.value, '$.id') as rule_id,
  json_extract(r.value, '$.type') as rule_type,
  json_extract(r.value, '$.parameters') as rule_parameters
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and name = 'test34';
```

### Get bypass actors for a specific ruleset
Identify the actors who can bypass the ruleset within your repository. This information is crucial for managing exceptions and understanding who has elevated permissions.

```sql+postgres
select
  name,
  b ->>'id' as bypass_actor_id,
  b ->>'deploy_key' as deploy_key,
  b ->>'bypass_mode' as bypass_mode,
  b ->>'repository_role_name' as repository_role_name,
  b ->>'repository_role_database_id' as repository_role_database_id
from
  github_repository_ruleset,
  jsonb_array_elements(bypass_actors) as b
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and name = 'test34';
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
  github_repository_ruleset,
  json_each(bypass_actors) as b
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and name = 'test34';
```

### List rulesets with specific enforcement levels
Identify rulesets within a repository that have specific enforcement levels, helping to understand the compliance and security posture of the repository.

```sql+postgres
select
  name,
  enforcement
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and enforcement = 'strict';
```

```sql+sqlite
select
  name,
  enforcement
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and enforcement = 'strict';
```

### List all rulesets created after a specific date
Retrieve all rulesets that were created after a specified date, useful for auditing and tracking recent changes in repository policies.

```sql+postgres
select
  name,
  created_at
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and created_at > '2023-01-01T00:00:00Z';
```

```sql+sqlite
select
  name,
  created_at
from
  github_repository_ruleset
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and created_at > '2023-01-01T00:00:00Z';
```

### Example Query Based on Parameters in Rules

Explore specific rules within the ruleset based on their parameter types. Below are queries for different parameter types:

#### 1. **UpdateParameters**

List rules with update parameters, focusing on the `update_allows_fetch_and_merge` setting.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'UpdateParameters' ->> 'update_allows_fetch_and_merge' as update_allows_fetch_and_merge
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'UpdateParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.UpdateParameters.update_allows_fetch_and_merge') as update_allows_fetch_and_merge
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'UpdateParameters';
```

#### 2. **WorkflowsParameters**

List rules with workflow parameters, focusing on the workflow configurations.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'WorkflowsParameters' ->> 'workflows' as workflows
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'WorkflowsParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.WorkflowsParameters.workflows') as workflows
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'WorkflowsParameters';
```

#### 3. **PullRequestParameters**

List rules with pull request parameters, including various settings such as code owner review requirements.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'PullRequestParameters' ->> 'require_code_owner_review' as require_code_owner_review,
  r -> 'parameters' -> 'PullRequestParameters' ->> 'required_approving_review_count' as required_approving_review_count
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'PullRequestParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.PullRequestParameters.require_code_owner_review') as require_code_owner_review,
  json_extract(r.value, '$.parameters.PullRequestParameters.required_approving_review_count') as required_approving_review_count
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'PullRequestParameters';
```

#### 4. **CodeScanningParameters**

List rules with code scanning parameters, showing tool configurations and thresholds.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'CodeScanningParameters' ->> 'code_scanning_tools' as code_scanning_tools
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'CodeScanningParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.CodeScanningParameters.code_scanning_tools') as code_scanning_tools
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'CodeScanningParameters';
```

#### 5. **TagNamePatternParameters**

List rules with tag name pattern parameters, including pattern configurations.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'TagNamePatternParameters' ->> 'pattern' as pattern
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'TagNamePatternParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.TagNamePatternParameters.pattern') as pattern
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'TagNamePatternParameters';
```

#### 6. **MaxFilePathLengthParameters**

List rules with max file path length parameters, specifying the maximum length allowed.

```sql+postgres
select
  id,
  name,
  r -> 'parameters' ->> 'Type' as type,
  r -> 'parameters' -> 'MaxFilePathLengthParameters' ->> 'max_file_path_length' as max_file_path_length
from
  github_repository_ruleset,
  jsonb_array_elements(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
and
  (r -> 'parameters' ->> 'Type') = 'MaxFilePathLengthParameters';
```

```sql+sqlite
select
  id,
  name,
  json_extract(r.value, '$.parameters.Type') as type,
  json_extract(r.value, '$.parameters.MaxFilePathLengthParameters.max_file_path_length') as max_file_path_length
from
  github_repository_ruleset,
  json_each(rules) as r
where
  repository_full_name = 'pro-cloud-49/test-rule'
  and json_extract(r.value, '$.parameters.Type') = 'MaxFilePathLengthParameters';
```