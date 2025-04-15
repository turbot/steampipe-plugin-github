---
title: "Steampipe Table: github_stargazer - Query GitHub Stargazers using SQL"
description: "Allows users to query GitHub Stargazers, specifically the users who have starred a particular repository, providing insights into user engagement and repository popularity."
folder: "Repository"
---

# Table: github_stargazer - Query GitHub Stargazers using SQL

GitHub Stargazers is a feature within GitHub that allows users to show appreciation to the repositories they find interesting. Users can star repositories to keep track of projects they find fascinating or useful. This feature provides a simple way to express interest in a project or to bookmark it for later reference.

## Table Usage Guide

The `github_stargazer` table provides insights into GitHub Stargazers within a specific repository. As a repository owner or collaborator, explore stargazer-specific details through this table, including user profiles, star creation timestamps, and associated metadata. Utilize it to uncover information about who is interested in your repository, when they starred it, and how your repository's popularity is growing over time.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List the stargazers of a repository
Discover who has shown interest in a specific Github repository and when, by identifying the users who have starred it and the time they did so. This can be useful for understanding the popularity and reach of the repository over time.

```sql+postgres
select
  user_login,
  starred_at
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
order by
  starred_at desc;
```

```sql+sqlite
select
  user_login,
  starred_at
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
order by
  starred_at desc;
```

### New stargazers by month
Discover the popularity trend of the 'Steampipe' project on Github by counting new stargazers each month. This helps in understanding the project's growth and tracking its community engagement over time.

```sql+postgres
select
  to_char(starred_at, 'YYYY-MM') as month,
  count(*)
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
group by
  month
order by
  month;
```

```sql+sqlite
select
  strftime('%Y-%m', starred_at) as month,
  count(*)
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
group by
  month
order by
  month;
```

### List stargazers with their contact information
Discover the segments that are actively showing interest in your GitHub repository by identifying users who have starred it and gathering their contact information. This can help in understanding your user base, facilitating community engagement or reaching out for feedback.

```sql+postgres
select
  user_login,
  starred_at,
  user_detail ->> 'name' as name,
  user_detail ->> 'company' as company,
  user_detail ->> 'email' as email,
  user_detail ->> 'url' as url,
  user_detail ->> 'twitter_username' as twitter_username,
  user_detail ->> 'website_url' as website,
  user_detail ->> 'location' as location,
  user_detail ->> 'bio' as bio
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe-plugin-github'
order by
  starred_at desc;
```

```sql+sqlite
select
  user_login,
  starred_at,
  json_extract(user_detail, '$.name') as name,
  json_extract(user_detail, '$.company') as company,
  json_extract(user_detail, '$.email') as email,
  json_extract(user_detail, '$.url') as url,
  json_extract(user_detail, '$.twitter_username') as twitter_username,
  json_extract(user_detail, '$.website_url') as website,
  json_extract(user_detail, '$.location') as location,
  json_extract(user_detail, '$.bio') as bio
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe-plugin-github'
order by
  starred_at desc;
```