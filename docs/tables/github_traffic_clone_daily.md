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

### List clone statistics
Explore the daily clone statistics of the 'turbot/steampipe' repository to assess developer interest and cloning activity. This can help you understand the adoption and popularity of your project over time.

```sql
select
  timestamp,
  count,
  uniques
from
  github_traffic_clone_daily
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
