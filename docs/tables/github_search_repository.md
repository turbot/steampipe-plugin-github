---
title: "Steampipe Table: github_search_repository - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories based on specific search criteria, providing insights into repository details such as name, owner, description, and more."
---

# Table: github_search_repository - Query GitHub Repositories using SQL

GitHub Repositories represent a location where all the files for a project are stored. Each project corresponds to one repository. Repositories can contain folders and files, images, videos, spreadsheets, and data sets - anything your project needs. 

## Table Usage Guide

The `github_search_repository` table provides insights into repositories within GitHub. As a developer or project manager, explore repository-specific details through this table, including owner, name, description, and more. Utilize it to uncover information about repositories, such as those with specific topics, languages, or other search criteria, and to manage and organize your project's files and resources.

**Important Notes**
- You must always include at least one search term when searching repositories in the where or join clause using the `query` column. You can search for repositories globally across all of GitHub.com, or search for repositories within a particular organization. See [Searching for repositories](https://docs.github.com/search-github/searching-on-github/searching-for-repositories) for details on the GitHub query syntax.

## Examples

### Get a specific repository
Identify specific details about a given repository, such as the owner, language used, and various user interaction metrics. This can be useful for understanding the popularity and reach of a repository within the GitHub community.

```sql+postgres
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'repo:turbot/steampipe-plugin-github';
```

### List repositories based on contents of a repository
Discover the segments that are popular among users by identifying repositories based on the content of a specific repository. This can help in understanding user preferences and trends in the open-source community.

```sql+postgres
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'stargazers in:readme repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'stargazers in:readme repo:turbot/steampipe-plugin-github';
```

### List repositories with more than 100000 followers
Determine the areas in which popular repositories on GitHub have amassed a large following. This query is useful for identifying trends and patterns among the most followed repositories, offering insights into what makes a repository popular.

```sql+postgres
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'followers:>=100000';
```

```sql+sqlite
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'followers:>=100000';
```

### List forked repositories created within specific timestamp
Discover the segments that have forked a specific repository within a particular time frame. This can be particularly useful for understanding the popularity and reach of a project during that period.

```sql+postgres
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'tinyspotifyr in:name created:2021-01-01..2021-01-05 fork:only';
```

```sql+sqlite
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'tinyspotifyr in:name created:2021-01-01..2021-01-05 fork:only';
```