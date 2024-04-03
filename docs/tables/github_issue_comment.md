---
title: "Steampipe Table: github_issue_comment - Query GitHub Issue Comments using SQL"
description: "Allows users to query GitHub Issue Comments, specifically the details of comments made on issues across all repositories, providing insights into user interactions and feedback."
---

# Table: github_issue_comment - Query GitHub Issue Comments using SQL

GitHub Issue Comments are integral parts of the GitHub platform that enable users to discuss and provide feedback on particular issues within a repository. They serve as a platform for collaboration and can contain important information, updates, or solutions related to the issue at hand. GitHub Issue Comments are particularly useful for tracking the progress of issue resolution and fostering collaborative problem-solving.

## Table Usage Guide

The `github_issue_comment` table provides in-depth insights into Issue Comments within GitHub. As a project manager or developer, explore comment-specific details through this table, including the author, creation time, body of the comment, and associated metadata. Utilize it to track user interactions, gather feedback, monitor issue resolution progress, and encourage collaborative problem-solving.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List comments for a specific issue
This query is useful for gaining insights into the discussion surrounding a specific issue within a specified GitHub repository. It helps identify the contributors, their comments, and the timeline of their interactions, which can be beneficial for understanding the development and resolution process of the issue.

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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201;
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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201;
```

### List comments for your issues
Explore the dialogue associated with your issues on GitHub. This query helps you identify instances where a comment has been made on your issue, providing insights into who made the comment, the content of the comment, and when it was created.

```sql+postgres
select
  c.repository_full_name,
  c.number as issue,
  i.author_login as issue_author,
  c.author_login as comment_author,
  c.body_text as content,
  c.created_at as created,
  c.url
from 
  github_issue_comment c
join 
    github_my_issue i
on 
  i.repository_full_name = c.repository_full_name
and 
  i.number = c.number;
```

```sql+sqlite
select
  c.repository_full_name,
  c.number as issue,
  i.author_login as issue_author,
  c.author_login as comment_author,
  c.body_text as content,
  c.created_at as created,
  c.url
from 
  github_issue_comment c
join 
    github_my_issue i
on 
  i.repository_full_name = c.repository_full_name
and 
  i.number = c.number;
```

### List comments for a specific issue which match a certain body content
This query is useful when you want to pinpoint specific comments within a particular issue in a GitHub repository that contain a certain keyword in their content. It can help manage and analyze the discussion within the issue, especially when looking for mentions of a specific topic or term.

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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201
and
  body_text ~~* '%branch%';
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
  github_issue_comment
where
  repository_full_name = 'turbot/steampipe-plugin-github'
and
  number = 201
and
  body_text like '%branch%';
```