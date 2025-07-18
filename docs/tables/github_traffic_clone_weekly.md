---
title: "Steampipe Table: github_traffic_clone_weekly - Query GitHub Traffic Clone Weekly using SQL"
description: "Allows users to query weekly repository clones on GitHub, specifically the weekly clone count, providing insights into repository cloning patterns and developer interest."
folder: "Repository"
---

# Table: github_traffic_clone_weekly - Query GitHub Traffic Clone Weekly using SQL

GitHub Weekly Traffic Clones is a feature within GitHub that allows repository owners to monitor and understand how frequently their repositories are cloned over time. It provides a weekly summary of the number of clones and unique clones of a repository. This feature aids in understanding the adoption, popularity, and developer interest in a repository over time.

## Table Usage Guide

The `github_traffic_clone_weekly` table provides insights into weekly repository clone activity within GitHub. As a repository owner, explore clone-specific details through this table, including the number of clones, unique clones, and the timestamp of the clone activity. Utilize it to uncover information about repository adoption trends, spikes in developer interest, and general repository engagement.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Administration (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List clone statistics for all your repositories
Explore the weekly clone statistics for all your repositories to assess developer interest and cloning activity. This can help you understand the adoption and popularity of your projects over time.

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
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
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
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
order by
  r.name_with_owner,
  c.timestamp;
```

### Find your repositories with the highest weekly clones
Summarize weekly clone counts to highlight your most visited repositories and understand which projects are gaining the most traction.

```sql+postgres
select
  r.name_with_owner,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.stargazer_count,
  r.fork_count,
  c.timestamp as week_ending,
  c.count as weekly_clones,
  c.uniques as unique_weekly_cloners,
  r.visibility,
  round(c.count::numeric / nullif(c.uniques, 0), 2) as clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
order by
  c.count desc,
  c.timestamp desc;
```

```sql+sqlite
select
  r.name_with_owner,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.stargazer_count,
  r.fork_count,
  c.timestamp as week_ending,
  c.count as weekly_clones,
  c.uniques as unique_weekly_cloners,
  r.visibility,
  case
    when c.uniques = 0 then null
    else cast(c.count as real) / c.uniques
  end as clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
order by
  c.count desc,
  c.timestamp desc;
```

### Compare weekly clone counts for your repositories over time
Identify how your repositories' clone activity is trending week-over-week by analyzing historical patterns.

```sql+postgres
with weekly_trends as (
  select
    r.name_with_owner,
    r.name,
    r.primary_language->>'name' as language,
    c.timestamp,
    c.count as weekly_clones,
    c.uniques as unique_cloners,
    lag(c.count) over (partition by r.name_with_owner order by c.timestamp) as previous_week_clones,
    lag(c.uniques) over (partition by r.name_with_owner order by c.timestamp) as previous_week_uniques
  from
    github_my_repository r
  join
    github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
)
select
  name_with_owner,
  name,
  language,
  timestamp as week_ending,
  weekly_clones,
  unique_cloners,
  previous_week_clones,
  previous_week_uniques,
  case
    when previous_week_clones is null then null
    else weekly_clones - previous_week_clones
  end as clone_change,
  case
    when previous_week_clones is null or previous_week_clones = 0 then null
    else round(((weekly_clones - previous_week_clones)::numeric / previous_week_clones) * 100, 2)
  end as percent_change
from
  weekly_trends
where
  timestamp >= current_date - interval '8 weeks'
order by
  name_with_owner,
  timestamp;
```

```sql+sqlite
with weekly_trends as (
  select
    r.name_with_owner,
    r.name,
    r.primary_language->>'name' as language,
    c.timestamp,
    c.count as weekly_clones,
    c.uniques as unique_cloners,
    lag(c.count) over (partition by r.name_with_owner order by c.timestamp) as previous_week_clones,
    lag(c.uniques) over (partition by r.name_with_owner order by c.timestamp) as previous_week_uniques
  from
    github_my_repository r
  join
    github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
)
select
  name_with_owner,
  name,
  language,
  timestamp as week_ending,
  weekly_clones,
  unique_cloners,
  previous_week_clones,
  previous_week_uniques,
  case
    when previous_week_clones is null then null
    else weekly_clones - previous_week_clones
  end as clone_change,
  case
    when previous_week_clones is null or previous_week_clones = 0 then null
    else ((weekly_clones - previous_week_clones) * 100.0) / previous_week_clones
  end as percent_change
from
  weekly_trends
where
  timestamp >= '2025-05-01' 
order by
  name_with_owner,
  timestamp;
```

### Get the latest weekly clone metrics for your repositories
Quickly view the most recent weekly clone statistics for all your repositories, along with repository characteristics and performance indicators.

```sql+postgres
select
  r.name_with_owner,
  r.name,
  r.description,
  r.primary_language->>'name' as language,
  r.stargazer_count,
  r.fork_count,
  r.visibility,
  r.created_at as repo_created,
  latest_clone.timestamp as latest_week_ending,
  latest_clone.count as latest_weekly_clones,
  latest_clone.uniques as latest_unique_cloners,
  round(latest_clone.count::numeric / nullif(latest_clone.uniques, 0), 2) as clones_per_unique_user
from
  github_my_repository r
join lateral (
  select
    timestamp,
    count,
    uniques
  from
    github_traffic_clone_weekly c
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
  r.visibility,
  r.created_at as repo_created,
  c.timestamp as latest_week_ending,
  c.count as latest_weekly_clones,
  c.uniques as latest_unique_cloners,
  case
    when c.uniques = 0 then null
    else cast(c.count as real) / c.uniques
  end as clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
where
  c.timestamp = (
    select max(timestamp)
    from github_traffic_clone_weekly c2
    where c2.repository_full_name = r.name_with_owner
  )
order by
  c.count desc;
```

### Analyze weekly performance patterns for your repositories
Examine weekly clone performance over time to understand engagement patterns and identify trends across all your repositories.

```sql+postgres
select
  r.primary_language->>'name' as language,
  r.visibility,
  count(distinct r.name_with_owner) as repository_count,
  sum(c.count) as total_weekly_clones,
  sum(c.uniques) as total_unique_cloners,
  avg(c.count) as avg_weekly_clones_per_repo,
  avg(c.uniques) as avg_unique_cloners_per_repo,
  max(c.count) as max_weekly_clones,
  round(avg(c.count::numeric / nullif(c.uniques, 0)), 2) as avg_clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= current_date - interval '4 weeks'
  and r.primary_language is not null
group by
  r.primary_language->>'name', r.visibility
having
  count(distinct r.name_with_owner) >= 1
order by
  total_weekly_clones desc;
```

```sql+sqlite
select
  r.primary_language->>'name' as language,
  r.visibility,
  count(distinct r.name_with_owner) as repository_count,
  sum(c.count) as total_weekly_clones,
  sum(c.uniques) as total_unique_cloners,
  avg(c.count) as avg_weekly_clones_per_repo,
  avg(c.uniques) as avg_unique_cloners_per_repo,
  max(c.count) as max_weekly_clones,
  avg(case
    when c.uniques = 0 then null
    else cast(c.count as real) / c.uniques
  end) as avg_clones_per_unique_user
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
where
  c.timestamp >= '2025-06-01' 
  and r.primary_language is not null
group by
  r.primary_language->>'name', r.visibility
having
  count(distinct r.name_with_owner) >= 1
order by
  total_weekly_clones desc;
```

### Compare weekly clone and view metrics for your repositories
Analyze the relationship between clones and views on a weekly basis to understand user engagement depth across all your repositories.

```sql+postgres
select
  c.timestamp as week_ending,
  r.name_with_owner,
  r.name,
  r.primary_language->>'name' as language,
  c.count as weekly_clones,
  c.uniques as unique_cloners,
  v.count as weekly_views,
  v.uniques as unique_viewers,
  round((c.count::numeric / nullif(v.count, 0)) * 100, 2) as clone_to_view_ratio,
  round((c.uniques::numeric / nullif(v.uniques, 0)) * 100, 2) as unique_clone_to_view_ratio
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
left join
  github_traffic_view_weekly v on c.repository_full_name = v.repository_full_name
  and c.timestamp = v.timestamp
where
  c.timestamp >= current_date - interval '12 weeks'
order by
  r.name_with_owner,
  c.timestamp;
```

```sql+sqlite
select
  c.timestamp as week_ending,
  r.name_with_owner,
  r.name,
  r.primary_language->>'name' as language,
  c.count as weekly_clones,
  c.uniques as unique_cloners,
  v.count as weekly_views,
  v.uniques as unique_viewers,
  case
    when v.count = 0 then null
    else (cast(c.count as real) / v.count) * 100
  end as clone_to_view_ratio,
  case
    when v.uniques = 0 then null
    else (cast(c.uniques as real) / v.uniques) * 100
  end as unique_clone_to_view_ratio
from
  github_my_repository r
join
  github_traffic_clone_weekly c on r.name_with_owner = c.repository_full_name
left join
  github_traffic_view_weekly v on c.repository_full_name = v.repository_full_name
  and c.timestamp = v.timestamp
where
  c.timestamp >= '2025-04-01' 
order by
  r.name_with_owner,
  c.timestamp;
```
