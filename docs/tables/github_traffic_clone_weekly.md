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

### List clone statistics
Explore the weekly clone statistics of the 'turbot/steampipe' repository to assess developer interest and cloning activity. This can help you understand the adoption and popularity of your project over time.

```sql
select
  timestamp,
  count,
  uniques
from
  github_traffic_clone_weekly
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
