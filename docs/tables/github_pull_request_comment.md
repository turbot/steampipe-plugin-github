---
title: "Steampipe Table: github_pull_request_comment - Query GitHub Pull Request Comments using SQL"
description: "Allows users to query GitHub Pull Request Comments, specifically the details of each comment made on pull requests, providing insights into discussions and feedbacks on code changes."
folder: "Pull Request"
---

# Table: github_pull_request_comment - Query GitHub Pull Request Comments using SQL

GitHub Pull Request Comments are individual responses or feedbacks given on a pull request in a GitHub repository. These comments facilitate discussions on proposed changes in the codebase, allowing for collaborative decision-making and code review. They represent an integral part of the code review process in GitHub, fostering effective communication and quality control among contributors.

## Table Usage Guide

The `github_pull_request_comment` table provides insights into the comments made on pull requests within a GitHub repository. As a developer or project manager, explore comment-specific details through this table, including the comment content, author, creation date, and associated metadata. Utilize it to understand the discussions and feedback on pull requests, facilitating effective code reviews and collaborative decision-making.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) and `number` (of the issue) columns in the `where` or `join` clause to query the table.

## Examples

### List all comments for a specific pull request
Determine the areas in which user comments on a particular pull request can provide valuable insights. This query is useful for understanding user engagement and feedback on specific code changes in a GitHub repository.

```sql+postgres
select
  id,
  author_login,
  author_association,
  body_text,
  created_at,
  updated_at,
  published_at,
  last_edited_at,
  editor_login,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207;
```

```sql+sqlite
select
  id,
  author_login,
  author_association,
  body_text,
  created_at,
  updated_at,
  published_at,
  last_edited_at,
  editor_login,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207;
```

### List comments for a specific pull request which match a certain body content
Determine the comments for a specific project update that contain a particular keyword. This is useful for filtering and understanding discussions related to specific topics or issues in your project.

```sql+postgres
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body_text as content,
  created_at,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207
  and body_text ~~* '%DELAY%';
```

```sql+sqlite
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body_text as content,
  created_at,
  url
from
  github_pull_request_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207
  and body_text like '%DELAY%';
```

### List comments for all open pull requests from a specific repository
Explore the discussion around ongoing modifications in a specific project by viewing the comments on all open pull requests. This can aid in understanding the current issues, proposed solutions, and overall progress of the project.
```sql+postgres
select
  c.*
from
  github_pull_request r
  join github_pull_request_comment c on r.repository_full_name = c.repository_full_name and r.number = c.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
  and r.state = 'OPEN';
```

```sql+sqlite
select
  c.*
from
  github_pull_request r
  join github_pull_request_comment c on r.repository_full_name = c.repository_full_name and r.number = c.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
  and r.state = 'OPEN';
```