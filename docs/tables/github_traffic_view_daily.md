---
title: "Steampipe Table: github_traffic_view_daily - Query GitHub Traffic Views using SQL"
description: "Allows users to query Traffic Views on GitHub, specifically the daily view count, providing insights into repository traffic patterns and potential anomalies."
folder: "Repository"
---

# Table: github_traffic_view_daily - Query GitHub Traffic Views using SQL

GitHub Traffic Views is a feature within GitHub that allows you to monitor and respond to traffic patterns across your repositories. It provides a centralized way to set up and monitor views for various GitHub repositories, including the number of visitors, unique visitors, and the number of views per day. GitHub Traffic Views helps you stay informed about the health and performance of your GitHub repositories and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `github_traffic_view_daily` table provides insights into daily repository views within GitHub. As a repository owner, explore view-specific details through this table, including the number of views, unique visitors, and the timestamp of the views. Utilize it to uncover information about views, such as those with high traffic, the peak times for views, and the verification of view patterns.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List view statistics
Explore the daily traffic statistics of the 'turbot/steampipe' repository to assess its popularity and unique visits. This can help you understand the reach and impact of your project over time.

```sql+postgres
select
  timestamp,
  count,
  uniques
from
  github_traffic_view_daily
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
```

```sql+sqlite
select
  timestamp,
  count,
  uniques
from
  github_traffic_view_daily
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
```