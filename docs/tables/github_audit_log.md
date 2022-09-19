# Table: github_audit_log

The audit log lists events triggered by activities that affect your organization. Only owners can access an organization's audit log.

The `github_audit_log` table helps to find all audit events for an organization, and **you must always specify the organization** in the where or join clause (`where organization=`, `join github_audit_log on organization=`).

**Note**: This table only works for organizations on an [GitHub Enterprise plan](https://docs.github.com/en/enterprise-cloud@latest/admin/overview/about-enterprise-accounts).

This table supports optional quals. Queries with optional quals are optimised to use GitHub query filters. Optional quals are supported for the following columns:
  - `action`
  - `actor`
  - `created_at`
  - `include`
  - `organization`
  - `phrase`

## Examples

### List recent audit events for an organization

```sql
select
  id,
  created_at,
  actor,
  action,
  data
from
  github_audit_log
where
  organization = 'my_org'
order by
  created_at
limit 10;
```

### List audit events in a specific date range

```sql
select
  id,
  created_at,
  actor,
  action,
  data
from
  github_audit_log
where
  organization = 'my_org'
  and created_at between '2022-06-27' and '2022-06-29'
order by
  created_at
```

### List repository creation and deletion audit events on a specific date

```sql
select
  id,
  created_at,
  actor,
  action,
  data
from
  github_audit_log
where
  organization = 'my_org'
  and action IN ('repo.create', 'repo.destroy')
  and created_at = '2022-01-01'
order by
  created_at;
```

### List audit events by a specific actor (user) in the last 30 days

```sql
select
  id,
  created_at,
  actor,
  action,
  data
from
  github_audit_log
where
  organization = 'my_org'
  and actor = 'some_user'
  and created_at > now() - interval '30 day'
order by
  created_at;
```

### List branch protection override audit events on a specific date using a search phrase

```sql
select
  id,
  created_at,
  actor,
  action,
  data
from
  github_audit_log
where
  organization = 'my_org'
  phrase = 'action:protected_branch.policy_override created:2022-06-28'
order by
  created_at;
```
