---
title: "Steampipe Table: github_search_commit - Query GitHub Commits using SQL"
description: "Allows users to query GitHub Commits, specifically providing the ability to search for commits across repositories based on specified criteria."
folder: "Repository"
---

# Table: github_search_commit - Query GitHub Commits using SQL

GitHub Commits is a feature within the GitHub service that allows users to make changes to a repository. Each commit corresponds to a specific change made to the repository, and includes information about the author, the date of the commit, and an associated message describing the change. GitHub Commits form the core of GitHub's version control functionality.

## Table Usage Guide

The `github_search_commit` table enables insights into commit activities within GitHub repositories. As a developer or project manager, leverage this table to track changes, monitor project progress, and maintain version control. Utilize it to search for specific commits based on various criteria, such as author, date, and associated message.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

**Important Notes**
- You must always include at least one search term when searching source code in the where or join clause using the `query` column. You can narrow the results using these commit search qualifiers in any combination. See [Searching commits](https://docs.github.com/search-github/searching-on-github/searching-commits) for details on the GitHub query syntax.

## Examples

### List searched commits within organization's repositories
This query is useful for tracking specific changes or additions made within an organization's repositories, such as identifying when a certain table was added. It provides a way to monitor and review the evolution of the codebase, enhancing oversight and control over the development process.

```sql+postgres
select
  sha,
  query,
  html_url,
  repository_full_name,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star org:turbot';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  repository_full_name,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star org:turbot';
```

### List searched commits within a user's repositories
Explore the specific commits in a user's GitHub repositories that match a particular search term. This is useful for tracking changes related to specific features or issues.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star user:turbot';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star user:turbot';
```

### List searched commits within a user's specific repository
Determine the areas in which specific changes have been made within a user's specific repository. This can be beneficial for tracking the development and modifications of a project over time.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star repo:turbot/steampipe-plugin-github';
```

### List searched commits by author
Explore which commits have been made by a specific author in a particular GitHub repository. This can be useful to understand an author's contribution history and impact on a project.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'author:VincentHardouin repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'author:VincentHardouin repo:turbot/steampipe-plugin-github';
```

### List searched commits committed within the specified date range
Explore specific commits made within a certain date range on the 'turbot/steampipe-plugin-github' repository. This can be useful for tracking project progress or reviewing changes made during a specific period.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'committer-date:2021-11-16..2021-11-23 repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'committer-date:2021-11-16..2021-11-23 repo:turbot/steampipe-plugin-github';
```

### List searched commits by hash
Explore the specific commits in a GitHub repository by searching for a particular hash. This is useful for tracking changes and understanding the history of a project.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```

### List searched commits by parent
Explore which commits on GitHub have a specific parent commit. This can be useful to track changes made in a project over time, understand the evolution of code, and identify the impact of a particular commit.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'parent:b0566ea';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'parent:b0566ea';
```

### List searched commits by merge commits
Explore which commits have been merged into the 'turbot/steampipe-plugin-azure' repository. This can be useful in identifying the changes made to the repository and tracking the progress of the project.

```sql+postgres
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'merge:true repo:turbot/steampipe-plugin-azure';
```

```sql+sqlite
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'merge:true repo:turbot/steampipe-plugin-azure';
```

### List repository details
Explore the specifics of a GitHub repository, including its name, ID, and owner, by using this query. This is particularly useful for gaining insights into the repository's details without having to manually sift through GitHub's interface.

```sql+postgres
select
  sha,
  repository -> 'name' as repo_name,
  repository ->> 'id' as repo_id,
  repository ->> 'html_url' as repo_html_url,
  jsonb_pretty(repository -> 'owner') as repo_owner
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  sha,
  json_extract(repository, '$.name') as repo_name,
  json_extract(repository, '$.id') as repo_id,
  json_extract(repository, '$.html_url') as repo_html_url,
  repository -> 'owner' as repo_owner
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```