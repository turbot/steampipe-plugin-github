---
title: "Steampipe Table: github_traffic_clone_daily - Query GitHub Traffic Clones using SQL"
description: "Allows users to query repository clones on GitHub, specifically the daily clone count, providing insights into repository cloning patterns and developer interest."
folder: "Repository"
---

# Table: github_traffic_clone_daily - Query GitHub Traffic Clones using SQL

GitHub Traffic Clones is a feature within GitHub that allows you to monitor and analyze how often your repository is cloned, and how many unique developers are cloning it. It provides a centralized way to track developer interest in your repositories, including the number of clones and unique clones per day. This feature aids in understanding the adoption, popularity, and developer interest in a repository over time.

## Table Usage Guide

The `github_traffic_clone_daily` table provides insights into daily repository clones within GitHub. As a repository owner, explore clone-specific details through this table, including the number of clones, unique clones, and the timestamp of the clone activity. Utilize it to uncover information about repository adoption trends, spikes in developer interest, and general repository engagement.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Administration (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List clone statistics for all your repositories
Explore the daily clone statistics for all your repositories to assess developer interest and cloning activity. This can help you understand the adoption and popularity of your projects over time.

```sql+postgres
select
  r.name_with_owner,
  r.name,
  c.timestamp,
  c.count,
  c.uniques
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  r.name_with_owner like '%turbot/%'
order by
  r.name_with_owner,
  c.timestamp;
```

```sql+sqlite
select
  r.name_with_owner,
  r.name,
  c.timestamp,
  c.count,
  c.uniques
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
order by
  r.name_with_owner,
  c.timestamp;
```

### Show daily clone trends for your repositories
Track the number of clones and unique users per day to analyze traffic over time for all your repositories.

```sql+postgres
select
  r.name_with_owner,
  c.timestamp::date as date,
  c.count as total_clones,
  c.uniques as unique_cloners,
  r.name as repository_name,
  r.primary_language->>'name' as primary_language,
  r.visibility,
  round(c.count::numeric / nullif(c.uniques, 0), 2) as clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= current_date - interval '30 days'
order by
  r.name_with_owner,
  c.timestamp;
```

```sql+sqlite
select
  r.name_with_owner,
  date(c.timestamp) as date,
  c.count as total_clones,
  c.uniques as unique_cloners,
  r.name as repository_name,
  r.primary_language->>'name' as primary_language,
  r.visibility,
  round(cast(c.count as real) / nullif(c.uniques, 0), 2) as clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= date('now', '-30 days')
order by
  r.name_with_owner,
  c.timestamp;
```

### Compare clone and view activity for your repositories
Analyze the correlation between repository clones and views to understand user engagement patterns across all your repositories.

```sql+postgres
select
  c.timestamp::date as date,
  r.name_with_owner as full_name,
  r.name,
  c.count as daily_clones,
  c.uniques as unique_cloners,
  v.count as daily_views,
  v.uniques as unique_viewers,
  round((c.count::numeric / nullif(v.count, 0)) * 100, 2) as clone_to_view_ratio
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
left join
  github_traffic_view_daily v on c.repository_full_name = v.repository_full_name
  and c.timestamp::date = v.timestamp::date
where
  c.timestamp >= current_date - interval '14 days'
order by
  r.name_with_owner,
  c.timestamp;
```

```sql+sqlite
select
  date(c.timestamp) as date,
  r.name_with_owner as full_name,
  r.name,
  c.count as daily_clones,
  c.uniques as unique_cloners,
  v.count as daily_views,
  v.uniques as unique_viewers,
  round((cast(c.count as real) / nullif(v.count, 0)) * 100, 2) as clone_to_view_ratio
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
left join
  github_traffic_view_daily v on c.repository_full_name = v.repository_full_name
  and date(c.timestamp) = date(v.timestamp)
where
  c.timestamp >= date('now', '-14 days')
order by
  r.name_with_owner,
  c.timestamp;
```

### Get the latest clone metrics for your repositories
Quickly view the most recent daily clone statistics for all your repositories, along with repository metadata.

```sql+postgres
select
  r.name_with_owner as full_name,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.stargazer_count,
  r.fork_count,
  latest_clone.timestamp as latest_clone_date,
  latest_clone.count as latest_daily_clones,
  latest_clone.uniques as latest_unique_cloners
from
  github_my_repository r
join lateral (
  select
    timestamp,
    count,
    uniques
  from
    github_traffic_clone_daily c
  where
    c.repository_full_name = r.name_with_owner
  order by
    c.timestamp desc
  limit 1
) latest_clone on true
order by
  latest_clone.count desc;
```

```sql+sqlite
select
  r.name_with_owner as full_name,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.stargazer_count,
  r.fork_count,
  c.timestamp as latest_clone_date,
  c.count as latest_daily_clones,
  c.uniques as latest_unique_cloners
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  c.timestamp = (
    select max(timestamp)
    from github_traffic_clone_daily c2
    where c2.repository_full_name = r.name_with_owner
  )
order by
  c.count desc;
```

### Identify your repositories with highest clone activity in past week
Find your repositories with the most clone activity over the last 7 days to understand which projects are trending.

```sql+postgres
select
  r.name_with_owner as full_name,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.visibility,
  sum(c.count) as total_clones_week,
  sum(c.uniques) as total_unique_cloners_week,
  avg(c.count) as avg_daily_clones,
  max(c.count) as peak_daily_clones,
  count(*) as days_with_data
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= current_date - interval '7 days'
group by
  r.name_with_owner, r.name, r.description, r.primary_language->>'name', r.visibility
having
  sum(c.count) > 0
order by
  total_clones_week desc;
```

```sql+sqlite
select
  r.name_with_owner as full_name,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.visibility,
  sum(c.count) as total_clones_week,
  sum(c.uniques) as total_unique_cloners_week,
  avg(c.count) as avg_daily_clones,
  max(c.count) as peak_daily_clones,
  count(*) as days_with_data
from
  github_my_repository r
join
  github_traffic_clone_daily c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= date('now', '-7 days')
group by
  r.name_with_owner, r.name, r.description, r.primary_language->>'name', r.visibility
having
  sum(c.count) > 0
order by
  total_clones_week desc;
```
