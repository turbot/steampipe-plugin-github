# Table: github_audit_log

The `github_audit_log` table helps to find all audit events for an organization.

**You must always specify the organization** in the where or join clause using the `organization` column. Additionally, you can filter the logs by using a search phrase (`phrase`), event types (`include`), and before/after a timestamp (`created_at`).

## Examples

### Get all audit events for an organization

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
  organization = 'my_org';
```

### Get specific audit events

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
  and phrase = "action:repo.create action:repo.destroy"
  and created_at >= '2022-01-01';
```
