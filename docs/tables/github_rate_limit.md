---
title: "Steampipe Table: github_rate_limit - Query GitHub Rate Limits using SQL"
description: "Allows users to query GitHub Rate Limits, providing insights into the maximum number of requests that the user can make to the GitHub API within an hour."
folder: "API"
---

# Table: github_rate_limit - Query GitHub Rate Limits using SQL

GitHub Rate Limits are a part of the GitHub API that sets the maximum number of requests that a user can make within an hour. These limits are designed to prevent abuse and ensure fair usage of the GitHub API. They apply to authenticated and unauthenticated users, and vary depending on the type of authentication used.

## Table Usage Guide

The `github_rate_limit` table provides insights into the rate limits set by GitHub for API usage. As a developer or system administrator, you can use this table to monitor your application's API usage, ensuring it stays within the prescribed limits to avoid service disruption. This table is also useful for planning and optimizing the distribution of API requests to maximize efficiency and avoid hitting the rate limit.

**Important Notes**
If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### List rate limit of rest apis
Explore the usage of your REST APIs by identifying the remaining and maximum allowed requests. This is beneficial to avoid hitting rate limits and ensuring smooth operation of your services.

```sql+postgres
select
  core_limit,
  core_remaining,
  search_limit,
  search_remaining
from
  github_rate_limit;
```

```sql+sqlite
select
  core_limit,
  core_remaining,
  search_limit,
  search_remaining
from
  github_rate_limit;
```