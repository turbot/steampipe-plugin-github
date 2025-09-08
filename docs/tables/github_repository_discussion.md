---
title: "Steampipe Table: github_repository_discussion - Query GitHub Repository Discussions using SQL"
description: "Allows users to query GitHub Repository Discussions, providing insights into discussions, comments, and replies within repositories."
folder: "Repository"
---

# Table: github_repository_discussion - Query GitHub Repository Discussions using SQL

GitHub Discussions are conversations that can be started by anyone and are organized into categories. They provide a place for having conversations that are not issues or pull requests. Discussions can be used to ask questions, share ideas, or have general conversations about a project.

## Table Usage Guide

The `github_repository_discussion` table provides insights into discussions within GitHub repositories. As a project manager or developer, explore discussion-specific details through this table, including comments, replies, categories, and associated metadata. Utilize it to uncover information about discussions, such as those with the most engagement, discussions by category, and the overall activity within a repository's discussion forum.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:

- Repository permissions:
  - Discussions (Read-only): Required to access all columns.

**Important Notes**
- You must specify the `repository_full_name` (owner/repository) column in the `where` or `join` clause to query the table.

## Examples

### List discussions in a repository
Explore the discussions within a specific GitHub repository to understand community engagement and topics being discussed. This can help in tracking community activity and identifying popular discussion topics.

```sql+postgres
select
  number,
  title,
  author_login,
  category_name,
  created_at,
  comments_total_count
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  number,
  title,
  author_login,
  category_name,
  created_at,
  comments_total_count
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe';
```

### Get a specific discussion by number
Retrieve detailed information about a specific discussion, including its full content, author details, and engagement metrics. This is useful for analyzing individual discussion threads and understanding their impact.

```sql+postgres
select
  number,
  title,
  author_login,
  category_name,
  created_at,
  comments_total_count,
  answer
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

```sql+sqlite
select
  number,
  title,
  author_login,
  category_name,
  created_at,
  comments_total_count,
  answer
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

### Find discussions with answers
Identify discussions that have been answered, which can help in understanding which topics have been resolved and which discussions have received helpful responses from the community.

```sql+postgres
select
  number,
  title,
  author_login,
  category_name,
  answer
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and answer is not null;
```

```sql+sqlite
select
  number,
  title,
  author_login,
  category_name,
  answer
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and answer is not null;
```

### Find discussions by category
Filter discussions by their category to focus on specific types of conversations, such as FAQs, feature requests, or general discussions. This helps in organizing and prioritizing community feedback.

```sql+postgres
select
  number,
  title,
  author_login,
  category_name,
  created_at
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and category_name = 'FAQ';
```

```sql+sqlite
select
  number,
  title,
  author_login,
  category_name,
  created_at
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and category_name = 'FAQ';
```

### Get all comments for a discussion
Retrieve all comments for a specific discussion to understand the full conversation thread and community engagement. This includes both the comment content and metadata about each comment.

```sql+postgres
select
  number,
  title,
  comments
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

```sql+sqlite
select
  number,
  title,
  comments
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

### Get all replies for a discussion
Retrieve all replies to comments in a specific discussion to understand the nested conversation structure and identify threads with high engagement.

```sql+postgres
select
  number,
  title,
  replies
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

```sql+sqlite
select
  number,
  title,
  replies
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and number = 1007;
```

### Find discussions with the most comments
Identify the most active discussions in a repository by sorting by comment count. This helps in understanding which topics generate the most community engagement and discussion.

```sql+postgres
select
  number,
  title,
  author_login,
  comments_total_count
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
order by
  comments_total_count desc
limit 10;
```

```sql+sqlite
select
  number,
  title,
  author_login,
  comments_total_count
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
order by
  comments_total_count desc
limit 10;
```

### Find recent discussions
Identify recently created discussions to stay updated on the latest community activity and new topics being discussed in the repository.

```sql+postgres
select
  number,
  title,
  author_login,
  created_at
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and created_at > now() - interval '30 days'
order by
  created_at desc;
```

```sql+sqlite
select
  number,
  title,
  author_login,
  created_at
from
  github_repository_discussion
where
  repository_full_name = 'turbot/steampipe'
  and created_at > datetime('now', '-30 days')
order by
  created_at desc;
```
