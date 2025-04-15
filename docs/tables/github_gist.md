---
title: "Steampipe Table: github_gist - Query GitHub Gists using SQL"
description: "Allows users to query GitHub Gists, specifically retrieving details such as gist id, description, comments, files, forks, history, owner, and public status."
folder: "Gist"
---

# Table: github_gist - Query GitHub Gists using SQL

GitHub Gists are a simple way to share snippets and pastes with others. Gists are version controlled, forkable, and embeddable, making them an ideal tool for sharing and discussing code. They can be public or secret, and can be associated with a specific GitHub user or anonymous.

## Table Usage Guide

The `github_gist` table provides insights into Gists within GitHub. As a developer or team lead, explore Gist-specific details through this table, including its description, comments, files, forks, history, owner, and public status. Utilize it to uncover information about Gists, such as their version history, fork details, and associated comments.

**Important Notes**
- You must specify the `id` column in `where` or `join` clause to query the table.

## Examples

### Get details about ANY public gist (by id)
Explore the specifics of a publicly shared code snippet on Github by providing its unique ID. This is useful to understand the content and context of the code snippet without having to navigate through the Github platform.

```sql+postgres
select
  *
from
  github_gist
where
  id='633175';
```

```sql+sqlite
select
  *
from
  github_gist
where
  id='633175';
```

### Get file details about ANY public gist (by id)
Explore the contents of any public 'gist' on GitHub by specifying its unique ID. This can be useful for understanding the content and structure of shared code snippets without having to navigate to the GitHub site.

```sql+postgres
select
  id,
  jsonb_pretty(files)
from
  github_gist
where
  id = 'e85a3d8e7a23c247f672aaf95b6c3da9';
```

```sql+sqlite
select
  id,
  files
from
  github_gist
where
  id = 'e85a3d8e7a23c247f672aaf95b6c3da9';
```