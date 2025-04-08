---
title: "Steampipe Table: github_rate_limit_graphql - Query GitHub GraphQL API Rate Limits using SQL"
description: "Allows users to query GitHub GraphQL API Rate Limits, providing insights into the current rate limit status for the GraphQL API."
folder: "Meta"
---

# Table: github_rate_limit_graphql - Query GitHub GraphQL API Rate Limits using SQL

GitHub's GraphQL API Rate Limiting is a feature that prevents excessive use of the API by limiting the number of requests that can be made within a certain time frame. This feature helps ensure fair usage and prevents any single user from monopolizing the API resources. It provides a mechanism to monitor and manage the rate at which applications can access the API.

## Table Usage Guide

The `github_rate_limit_graphql` table provides insights into the rate limiting status for GitHub's GraphQL API. As a developer or DevOps engineer, you can use this table to monitor the current rate limit status of your application's API usage. This can be particularly useful in managing and optimizing your application's API requests to ensure they stay within the allowed limits.

## Examples

### List rate limit info for GraphQL
Assess the usage and availability of your GraphQL rate limit on GitHub to manage your API requests effectively and avoid exceeding the limit. This helps in planning your application's interactions with GitHub's API and ensures uninterrupted service.

```sql+postgres
select
  used,
  remaining,
  reset_at
from
  github_rate_limit_graphql;
```

```sql+sqlite
select
  used,
  remaining,
  reset_at
from
  github_rate_limit_graphql;
```