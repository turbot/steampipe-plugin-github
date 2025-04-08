---
title: "Steampipe Table: github_pull_request_review - Query GitHub Pull Request Reviews using SQL"
description: "Allows users to query Pull Request Reviews in GitHub, specifically the review comments, review status, and the reviewer details, providing insights into the review process of pull requests."
folder: "Pull Request"
---

# Table: github_pull_request_review - Query GitHub Pull Request Reviews using SQL

A GitHub Pull Request Review is a feature within GitHub that allows users to provide feedback on pull requests. It provides a collaborative platform for code review where users can comment, approve or request changes on the proposed code changes. GitHub Pull Request Reviews help to ensure code quality and maintainability by facilitating peer review before code merging.

## Table Usage Guide

The `github_pull_request_review` table provides insights into the review process of pull requests within GitHub. As a developer or a team lead, explore review-specific details through this table, including the review comments, review status, and the reviewer details. Utilize it to understand the feedback on pull requests, the approval process, and to gain insights into the code review practices in your projects.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) and `number` (of the PR) columns in the `where` or `join` clause to query the table.

## Examples

### List all reviews for a specific pull request
Explore all feedback for a specific project update. This is particularly useful for developers and project managers who want to understand the team's thoughts, concerns, and suggestions regarding a particular code change or feature addition.

```sql+postgres
select
  id,
  author_login,
  author_association,
  state,
  body,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207;
```

```sql+sqlite
select
  id,
  author_login,
  author_association,
  state,
  body,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207;
```

### List reviews for a specific pull request which match a certain body content
This query is useful for identifying specific feedback within the reviews of a particular pull request. It can help you to pinpoint comments that match a certain keyword, enabling you to quickly find and address relevant concerns or suggestions.

```sql+postgres
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body as content,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207
  and body like '%minor%';
```

```sql+sqlite
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body as content,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207
  and body like '%minor%';
```

### List reviews for all open pull requests from a specific repository
Determine the areas in which feedback has been provided for all active changes proposed in a specific project. This can be useful to understand the type of improvements or modifications suggested by contributors during the development process.

```sql+postgres
select
  rv.*
from
  github_pull_request r
  join github_pull_request_review rv on r.repository_full_name = rv.repository_full_name and r.number = rv.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
  and r.state = 'OPEN';
```

```sql+sqlite
select
  rv.*
from
  github_pull_request r
  join github_pull_request_review rv on r.repository_full_name = rv.repository_full_name and r.number = rv.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
  and r.state = 'OPEN';
```