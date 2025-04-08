---
title: "Steampipe Table: github_tag - Query GitHub Tags using SQL"
description: "Allows users to query Tags in GitHub, specifically the tag name, ID, commit details, and associated repository information, providing insights into version control and release management."
folder: "Repository"
---

# Table: github_tag - Query GitHub Tags using SQL

A GitHub Tag is a pointer to a specific commit in a repository, often used to capture a point in history that is used for a marked version release (i.e. v1.0.0). Tags are ref's that point to specific points in Git history. Tagging is generally used to capture a point in history that is used for a marked version release (i.e. v1.0.0).

## Table Usage Guide

The `github_tag` table provides insights into tags within GitHub repositories. As a developer or release manager, explore tag-specific details through this table, including commit details, tag name, and associated repository information. Utilize it to uncover information about tags, such as those associated with specific releases, the commit history of tags, and the management of version control.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List tags
Explore which versions of the 'turbot/steampipe' repository have been tagged on GitHub. This can be useful to understand the evolution of the project and identify specific versions for troubleshooting or reference.

```sql+postgres
select
  name,
  commit ->> 'sha' as commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  json_extract("commit", '$.sha') as commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```

### Order tags by semantic version
Discover the segments that are most relevant in the 'turbot/steampipe' repository by arranging tags in order of their semantic version. This can be helpful in understanding the progression and development of the project over time.

```sql+postgres
select
  name,
  commit ->> 'sha' as commit_sha
from
  github_tag
where
  repository_full_name = 'turbot/steampipe'
order by
  string_to_array(regexp_replace(name, '[^0-9\.]', '', 'g'), '.'),
  name;
```

```sql+sqlite
Error: SQLite does not support string_to_array and regexp_replace functions.
```

### Get commit details for each tag
Explore the specifics of each commit for different tags in a GitHub repository. This can help in tracking changes, understanding authorship and verifying the validity of commits, which can be crucial for code review and version control.

```sql+postgres
select
  name,
  commit ->> 'sha' as commit_sha,
  commit ->> 'message' as commit_message,
  commit ->> 'url' as commit_url,
  commit -> 'author' -> 'user' ->> 'login' as author,
  commit ->> 'authored_date' as authored_date,
  commit -> 'committer' -> 'user' ->> 'login' as committer,
  commit ->> 'committed_date' as committed_date,
  commit ->> 'additions' as additions,
  commit ->> 'deletions' as deletions,
  commit ->> 'changed_files' as changed_files,
  commit -> 'signature' ->> 'is_valid' as commit_signed,
  commit -> 'signature' ->> 'email' as commit_signature_email,
  commit -> 'signature' -> 'signer' ->> 'login' as commit_signature_login,
  commit ->> 'tarball_url' as tarball_url,
  commit ->> 'zipball_url' as zipball_url
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  json_extract(commit, '$.sha') as commit_sha,
  json_extract(commit, '$.message') as commit_message,
  json_extract(commit, '$.url') as commit_url,
  json_extract(commit, '$.author.user.login') as author,
  json_extract(commit, '$.authored_date') as authored_date,
  json_extract(commit, '$.committer.user.login') as committer,
  json_extract(commit, '$.committed_date') as committed_date,
  json_extract(commit, '$.additions') as additions,
  json_extract(commit, '$.deletions') as deletions,
  json_extract(commit, '$.changed_files') as changed_files,
  json_extract(commit, '$.signature.is_valid') as commit_signed,
  json_extract(commit, '$.signature.email') as commit_signature_email,
  json_extract(commit, '$.signature.signer.login') as commit_signature_login,
  json_extract(commit, '$.tarball_url') as tarball_url,
  json_extract(commit, '$.zipball_url') as zipball_url
from
  github_tag
where
  repository_full_name = 'turbot/steampipe';
```