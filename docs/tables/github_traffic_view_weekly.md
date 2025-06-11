---
title: "Steampipe Table: github_traffic_view_weekly - Query GitHub Traffic View Weekly using SQL"
description: "Allows users to query weekly traffic views on GitHub, providing insights into repository visit patterns and potential trends."
folder: "Repository"
---

# Table: github_traffic_view_weekly - Query GitHub Traffic View Weekly using SQL

GitHub Weekly Traffic Views is a feature within GitHub that allows repository owners to monitor and understand the frequency and source of visits to their repositories. It provides a weekly summary of the number of views and unique visitors to a repository. This feature aids in understanding the popularity and reach of a repository over time.

## Table Usage Guide

The `github_traffic_view_weekly` table provides insights into weekly traffic views within GitHub. As a repository owner or contributor, explore weekly view details through this table, including the number of views, unique visitors, and timestamp of the views. Utilize it to uncover information about repository popularity, such as peak visit times, trending repositories, and the reach of your repository.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.
- To query this table using Fine-grained access tokens, the following permissions are required:
  - **"Administration" repository permissions (read)** – Required to access the all columns.

## Examples

### List view statistics
Explore the popularity and unique visitor count of the 'Steampipe' repository on Github over time. This can be beneficial in understanding the reach and impact of the project, helping to inform future development and marketing strategies.

```sql+postgres
select
  timestamp,
  count,
  uniques
from
  github_traffic_view_weekly
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
  github_traffic_view_weekly
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
```