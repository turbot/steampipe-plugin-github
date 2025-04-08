---
title: "Steampipe Table: github_my_gist - Query GitHub Gists using SQL"
description: "Allows users to query GitHub Gists, specifically those owned by the authenticated user, providing insights into the user's Gist details."
folder: "Gist"
---

# Table: github_my_gist - Query GitHub Gists using SQL

GitHub Gists are a simple way to share snippets and pastes with others. Gists are Git repositories, which means that you can fork and clone them. Gists can be public or secret, and you can share them with specific people using the GitHub interface.

## Table Usage Guide

The `github_my_gist` table provides insights into Gists within GitHub. As a developer or GitHub user, explore gist-specific details through this table, including file content, comments, and associated metadata. Utilize it to manage and analyze your Gists, such as those with certain content, the number of comments, and the details of the files in the Gists.

**Important Notes**
- To query **ANY** gist that you have access to (including any public gists), use the `github_gist` table.

## Examples

### List your gists
Identify and review all the gists you have created on GitHub. This is useful for managing and keeping track of your code snippets shared through GitHub.

```sql+postgres
select
  *
from
  github_my_gist;
```

```sql+sqlite
select
  *
from
  github_my_gist;
```

### List your public gists
Explore which of your code snippets on GitHub are publicly accessible. This is useful for reviewing your open-source contributions and ensuring you're not unintentionally sharing private code.

```sql+postgres
select
  *
from
  github_my_gist
where
  public;
```

```sql+sqlite
select
  *
from
  github_my_gist
where
  public = 1;
```

### Summarize your gists by language.
Determine the distribution of your GitHub gists by programming language. This will help you understand the languages you use most frequently when creating gists.

```sql+postgres
select
  file ->> 'language' as language,
  count(*)
from
  github_my_gist g
cross join
  jsonb_array_elements(g.files) file
group by
  language
order by
  count desc;
```

```sql+sqlite
select
  json_extract(file.value, '$.language') as language,
  count(*)
from
  github_my_gist g,
  json_each(g.files) as file
group by
  language
order by
  count(*) desc;
```