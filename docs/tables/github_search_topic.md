---
title: "Steampipe Table: github_search_topic - Query GitHub Topics using SQL"
description: "Allows users to query Topics on GitHub, specifically to find repositories associated with a specific topic, providing insights into the popularity and usage of certain topics across repositories."
---

# Table: github_search_topic - Query GitHub Topics using SQL

GitHub Topics is a feature that allows users to explore repositories by technology, industry, and more. It helps developers to discover projects, learn from their peers, and better understand the landscape of possibilities within the platform. Topics are labels that repository owners can use to categorize their projects, so they're more discoverable.

## Table Usage Guide

The `github_search_topic` table provides insights into topics within GitHub. As a developer or project manager, explore topic-specific details through this table, including associated repositories, the number of stars, and the language used. Utilize it to uncover information about popular topics, the repositories that use them, and the overall popularity of different technologies and industries.

**Important Notes**
- You must always include at least one search term when searching topics in the where or join clause using the `query` column. See [Searching topics](https://docs.github.com/search-github/searching-on-github/searching-topics) for details on the GitHub query syntax.

## Examples

### List topics that are not curated
Identify topics that lack curation within a specific GitHub repository. This is useful for maintaining content quality and relevance by pinpointing areas needing further attention.

```sql+postgres
select
  name,
  created_at,
  curated
from
  github_search_topic
where
  query = 'target-searching is:not-curated repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  name,
  created_at,
  curated
from
  github_search_topic
where
  query = 'target-searching is:not-curated repo:turbot/steampipe-plugin-github';
```

### List featured topics
Explore which topics have been highlighted on the GitHub repository for the Steampipe GitHub plugin. This can be useful for understanding what areas of the plugin are being emphasized or promoted.

```sql+postgres
select
  name,
  created_at,
  featured
from
  github_search_topic
where
  query = 'markdown is:featured repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  name,
  created_at,
  featured
from
  github_search_topic
where
  query = 'markdown is:featured repo:turbot/steampipe-plugin-github';
```

### List topics that have more than 5000 repositories
Determine the areas in which popular topics on GitHub exceed 5000 repositories, providing insights into the most widely-used topics and enabling a focus on areas of high activity for potential collaboration or learning opportunities.

```sql+postgres
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'repositories:>5000 react-redux';
```

```sql+sqlite
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'repositories:>5000 react-redux';
```

### List topics that were created after a specific timestamp
Discover the segments that have been established after a certain date, specifically within the context of react-redux. This is useful for tracking the evolution and growth of react-redux related topics over time.

```sql+postgres
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'created:>2021-01-01 react-redux';
```

```sql+sqlite
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'created:>2021-01-01 react-redux';
```