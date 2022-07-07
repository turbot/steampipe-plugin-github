# Table: github_audit_log

The `github_audit_log` table helps to find all audit events for an organization. Note: this only works for organizations on an **GitHub Enterprise plan**.

**You must always specify the organization** in the where or join clause using the `organization` column. Additionally, you can filter the logs by using a search phrase (`phrase`), event types (`include`), and before/after a timestamp (`created_at`).

## Examples

### Get recent audit events for an organization

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
order by created_at
limit 10;
```

### Get specific audit events

For example, find out which repos have been created or deleted on the first of January.

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
  and created_at = '2022-01-01';
```

### Get specific events by a specific actor (user) in the last 30 days

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
  and created_at > (created_at - interval '30' day)
order by created_at;
```
