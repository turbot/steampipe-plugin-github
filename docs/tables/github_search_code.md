---
title: "Steampipe Table: github_search_code - Query GitHub Code using SQL"
description: "Allows users to query Code in GitHub, specifically the details of the code files and their properties, providing insights into the codebase and its structure."
folder: "Search"
---

# Table: github_search_code - Query GitHub Code using SQL

GitHub is a web-based hosting service for version control using Git. It is primarily used for computer code and offers all of the distributed version control and source code management (SCM) functionality of Git as well as adding its own features. It provides access control and several collaboration features such as bug tracking, feature requests, task management, and wikis for every project.

## Table Usage Guide

The `github_search_code` table provides insights into the code files within GitHub repositories. As a developer or DevOps engineer, explore file-specific details through this table, including file names, paths, and associated metadata. Utilize it to uncover information about the codebase, such as the distribution of file types, the structure of the repositories, and the details of the code files.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

**Important Notes**
- You must always include at least one search term when searching source code in the where or join clause using the `query` column. The `query` contains one or more search keywords and qualifiers. Qualifiers allow you to limit your search to specific areas of GitHub. See [Searching code](https://docs.github.com/search-github/searching-on-github/searching-code) for details on the GitHub query syntax.

## Examples

### List searched codes by file name
Explore which codes have been searched by file name. This is useful for tracking the usage and relevance of certain codes within your organization's repository.

```sql+postgres
select
  name,
  query,
  html_url,
  repository_full_name,
  sha
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  repository_full_name,
  sha
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

### List searched codes by file location
Explore which codes have been searched by their file location on GitHub. This can be useful to understand the areas of interest or frequent search patterns within specific directories.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables';
```

### List searched codes by extension
Explore which codes have been searched by their extension to gain insights into commonly used or referenced codes in your GitHub projects. This can help in understanding code usage patterns and optimizing your project structure.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables extension:md';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables extension:md';
```

### List searched codes within organization's repositories
Explore which codes within your organization's repositories have been searched. This can be useful in understanding the areas of interest and focus within your team's coding projects.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_stargazer org:turbot extension:go';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_stargazer org:turbot extension:go';
```

### List searched codes within a user's repositories
Explore which codes within a user's repositories have been searched for. This is handy for understanding which parts of your codebase are garnering attention and could potentially inform future development decisions.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers user:turbot extension:go';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers user:turbot extension:go';
```

### List searched codes within a user's specific repository
Explore which codes have been searched within a specific user's repository. This is useful for identifying popular topics or areas of interest within a particular project.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers repo:turbot/steampipe-plugin-github extension:go';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers repo:turbot/steampipe-plugin-github extension:go';
```

### List searched codes by language
Explore which codes have been searched by language on GitHub. This can be particularly useful for understanding the popularity and usage of different programming languages within a specific organization.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_tag org:turbot language:markdown';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_tag org:turbot language:markdown';
```

### List searched codes while file size is greater than 40 KB
Explore which files in a particular organization on GitHub are larger than 40KB and written in Markdown. This can be useful for identifying potentially bloated files that may need to be optimized or split up.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'org:turbot size:>40000 language:markdown';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'org:turbot size:>40000 language:markdown';
```

### List searched codes by the file contents or file path
Discover the segments that have been searched in GitHub by either file contents or file path. This can be beneficial in pinpointing the specific locations where certain codes are used, providing insights into code usage and organization.

```sql+postgres
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers org:turbot in:file,path extension:go';
```

```sql+sqlite
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers org:turbot in:file,path extension:go';
```

### List text match details
This query is used to explore the details of text matches in a GitHub code search. It's useful for identifying instances where specific text, such as a filename or property, appears in your organization's codebase.

```sql+postgres
select
  name,
  jsonb_pretty(match -> 'matches') as matches,
  match ->> 'fragment' as fragment,
  match ->> 'property' as property,
  match ->> 'object_url' as object_url,
  match ->> 'object_type' as object_type
from
  github_search_code,
  jsonb_array_elements(text_matches) as match
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

```sql+sqlite
select
  name,
  match.value as matches,
  json_extract(match.value, '$.fragment') as fragment,
  json_extract(match.value, '$.property') as property,
  json_extract(match.value, '$.object_url') as object_url,
  json_extract(match.value, '$.object_type') as object_type
from
  github_search_code,
  json_each(text_matches) as match
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

### List repository details
This query helps you to explore and understand the details of specific repositories within your organization on GitHub. It is particularly useful in identifying and analyzing the repositories that contain a specific file, thereby aiding in efficient file management and organization.

```sql+postgres
select
  name,
  repository -> 'id' as repo_id,
  repository ->> 'name' as repo_name,
  repository ->> 'url' as repo_url,
  jsonb_pretty(repository -> 'owner') as repo_owner
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

```sql+sqlite
select
  name,
  json_extract(repository, '$.id') as repo_id,
  json_extract(repository, '$.name') as repo_name,
  json_extract(repository, '$.url') as repo_url,
  repository -> 'owner' as repo_owner
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```