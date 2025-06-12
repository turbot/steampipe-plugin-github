---
title: "Steampipe Table: github_commit - Query GitHub Commits using SQL"
description: "Allows users to query GitHub Commits, specifically the details of individual commits made in a repository, providing insights into code changes and development patterns."
folder: "Repository"
---

# Table: github_commit - Query GitHub Commits using SQL

GitHub Commits are an integral part of the GitHub service, which provides a collaborative platform for software development. Each commit represents a single point in the git history, storing the changes made to the codebase along with metadata such as the author, timestamp, and associated comments. Commits serve as a record of what changes were made, when, and by whom, providing a comprehensive history of a project's development.

## Table Usage Guide

The `github_commit` table provides insights into individual commits within a GitHub repository. As a developer or project manager, explore commit-specific details through this table, including author information, commit messages, and timestamps. Utilize it to understand the development history, track code changes, and analyze development patterns within your projects.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Contents (Read-only): Required to access all columns.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### Recent commits
Explore the latest updates made to a specific GitHub repository. This allows you to track changes and progress over time, providing valuable insights into the development process.

```sql+postgres
select
  sha,
  author_login,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  authored_date desc;
```

```sql+sqlite
select
  sha,
  author_login,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  authored_date desc;
```

### Commits by a given author
Explore which contributions have been made by a specific user to a particular GitHub repository, sorted by most recent. This can be useful for tracking individual productivity or reviewing the history of changes made by a particular developer.

```sql+postgres
select
  sha,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
  and author_login = 'e-gineer'
order by
  authored_date desc;
```

```sql+sqlite
select
  sha,
  authored_date,
  message
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
  and author_login = 'e-gineer'
order by
  authored_date desc;
```

### Contributions by author
Discover the segments that feature the most activity by tracking the frequency of contributions from individual authors in a specific GitHub repository. This is useful for understanding who are the most active contributors in a project.

```sql+postgres
select
  author_login,
  count(*)
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  count desc;
```

```sql+sqlite
select
  author_login,
  count(*)
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  count(*) desc;
```

### Commits that were not verified
Discover the segments that include unverified commits in the Steampipe repository. This can be particularly useful for tracking and improving security by identifying potentially malicious or unauthorised changes to the codebase.

```sql+postgres
select
  sha,
  author_login,
  authored_date
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
and
  signature is null
order by
  authored_date desc;
```

```sql+sqlite
select
  sha,
  author_login,
  authored_date
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
and
  signature is null
order by
  authored_date desc;
```

### Commits with most file changes
Explore which commits have resulted in the most file changes within a specific GitHub repository. This can help identify instances where significant modifications have been made, providing insights into the evolution and maintenance of the project.

```sql+postgres
select
  sha,
  message,
  author_login,
  changed_files,
  additions,
  deletions
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  changed_files desc;
```

```sql+sqlite
select
  sha,
  message,
  author_login,
  changed_files,
  additions,
  deletions
from
  github_commit
where
  repository_full_name = 'turbot/steampipe'
order by
  changed_files desc;
```