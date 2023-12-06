---
title: "Steampipe Table: github_audit_log - Query GitHub Audit Logs using SQL"
description: "Allows users to query GitHub Audit Logs, specifically the actions, users, and repositories involved, providing insights into user activity and potential security concerns."
---

# Table: github_audit_log - Query GitHub Audit Logs using SQL

GitHub Audit Logs is a feature within GitHub that allows you to keep track of what's happening in your organization, repositories, and teams. It provides a record of actions taken by users, whether they're adding new members, changing repository settings, or deleting branches. GitHub Audit Logs helps you stay informed about the activities within your GitHub resources and take appropriate actions when needed.

## Table Usage Guide

The `github_audit_log` table provides insights into user activity within GitHub. As a Security Analyst, explore user-specific actions through this table, including performed actions, involved repositories, and action timestamps. Utilize it to uncover information about user actions, such as repository changes, team membership alterations, and other potential security risks.

**Important Notes**
- You must specify the `organization` column in `where` or `join` clause to query the table.
- This table only works for organizations on an [GitHub Enterprise plan](https://docs.github.com/en/enterprise-cloud@latest/admin/overview/about-enterprise-accounts).
- This table supports optional quals. Queries with optional quals are optimised to use GitHub query filters. Optional quals are supported for the following columns:
  - `action`
  - `actor`
  - `created_at`
  - `include`
  - `organization`
  - `phrase`

## Examples

### List recent audit events for an organization
Explore the recent audit activities within your organization to gain insights into actions taken and by whom, which can aid in understanding behavioral patterns and identifying potential security issues.

```sql+postgres
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

```sql+sqlite
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
Explore which audit events occurred within your organization over a specific date range. This can help you understand the activity and changes made during that period, allowing for better tracking and management.

```sql+postgres
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

```sql+sqlite
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
Explore which repository creation and deletion events occurred on a specific date within your organization. This is useful for tracking changes and maintaining a record of repository actions for potential audit or review purposes.

```sql+postgres
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

```sql+sqlite
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
This query is useful for tracking the activities of a particular user within your organization on Github over the past month. It helps in monitoring user behavior, identifying any unusual actions, and maintaining a safe and secure environment.

```sql+postgres
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

```sql+sqlite
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
  and created_at > datetime('now', '-30 day')
order by
  created_at;
```

### List branch protection override audit events on a specific date using a search phrase
Gain insights into the audit events that occurred on a specific date, particularly those related to branch protection overrides. This is useful for organizations that want to monitor and assess potential security risks or policy violations within their GitHub repositories.

```sql+postgres
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

```sql+sqlite
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
  and phrase = 'action:protected_branch.policy_override created:2022-06-28'
order by
  created_at;
```