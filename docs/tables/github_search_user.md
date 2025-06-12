---
title: "Steampipe Table: github_search_user - Query GitHub Users using SQL"
description: "Allows users to query GitHub Users, specifically the user's login name, type, score, and other details, providing insights into user activity and profile information."
folder: "Search"
---

# Table: github_search_user - Query GitHub Users using SQL

GitHub Users is a resource within GitHub that represents an individual or organization account. It holds information about the user's profile, including login name, type, and other details. GitHub Users provide a way to interact with the user's repositories, gists, followers, and more.

## Table Usage Guide

The `github_search_user` table provides insights into user profiles within GitHub. As a developer or a security analyst, explore user-specific details through this table, including login name, type, and score. Utilize it to uncover information about users, such as their activity, repositories they have access to, and their general profile information.

**Important Notes**
- You must always include at least one search term when searching users in the where or join clause using the `query` column. See [Searching users](https://docs.github.com/search-github/searching-on-github/searching-users) for details on the GitHub query syntax.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### List users
Identify instances where users have 'turbot' in their name within GitHub. This can help you find specific user profiles related to 'turbot' for further analysis or contact.

```sql+postgres
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbot in:name type:user';
```

```sql+sqlite
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbot in:name type:user';
```

### List organizations
Determine the areas in which specific organizations are operating by exploring their login details and types. This is useful in understanding the organizations' online presence and their categorization.

```sql+postgres
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbotio in:login type:org';
```

```sql+sqlite
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbotio in:login type:org';
```

### Get user with specific username
Explore which GitHub user corresponds to a specific username. This is useful for finding detailed information about a particular user, such as their unique ID, login type, and URL.

```sql+postgres
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'user:c0d3r-arnab';
```

```sql+sqlite
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'user:c0d3r-arnab';
```

### List organizations with over 10000 repositories
This query is useful for identifying large-scale organizations on GitHub, specifically those that have more than 10,000 repositories. It can be used to understand the scale of open-source contributions or to target potential collaborations with prolific organizations.

```sql+postgres
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'repos:>10000 type:org';
```

```sql+sqlite
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'repos:>10000 type:org';
```

### List users and organizations created between specific timestamp
Discover the segments that were added to your GitHub user base within a specific timeframe. This can help you assess the growth and monitor the activity on your platform during that period.

```sql+postgres
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'created:2021-01-01..2021-01-31 turbot';
```

```sql+sqlite
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'created:2021-01-01..2021-01-31 turbot';
```