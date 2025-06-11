---
title: "Steampipe Table: github_my_repository - Query GitHub Repositories using SQL"
description: "Allows users to query their GitHub Repositories, specifically focusing on repository details such as name, description, owner, and visibility status."
folder: "Repository"
---

# Table: github_my_repository - Query GitHub Repositories using SQL

GitHub is a version control system that allows developers to collaborate on projects. Repositories in GitHub contain all of the project files and each file's revision history. They also contain a README file which provides information about the project. Repositories can be public or private, and they can be accessed and managed on GitHub's website.

## Table Usage Guide

The `github_my_repository` table provides insights into personal GitHub repositories. As a developer or project manager, explore repository-specific details through this table, including repository name, description, owner, and visibility status. Utilize it to manage and monitor your repositories, such as checking the visibility status, reviewing the repository description, and identifying the repository owner.

**Important Notes**
- You can own repositories individually, or you can share ownership of repositories with other people in an organization. The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.
- To query **ANY** repository, including public repos, use the `github_repository` table.
- To query this table using Fine-grained access tokens, the following permissions are required:
   - **"Metadata" repository permission (read)** – Required to access general repository metadata.
  - **"Contents" repository permission (read)** – Required for the `issue_templates` and `pull_request_templates` columns.
  - **"Webhooks" repository permission (read)** – Required for the `hooks` column.
  - **"Issues" repository permission (read)** – Required for the `open_issues_total_count` column.

## Examples

### List of repositories that you or your organizations own or contribute to
Explore which repositories you or your organizations have ownership of or contribute to, in order to better understand your coding involvement and collaborations. This could be particularly useful for tracking project participation or evaluating the spread of your contributions.

```sql+postgres
select
  name,
  owner_login,
  name_with_owner
from
  github_my_repository
order by
  name_with_owner;
```

```sql+sqlite
select
  name,
  owner_login,
  name_with_owner
from
  github_my_repository
order by
  name_with_owner;
```

### Show repository stats
Explore your GitHub repositories' statistics to gain insights into popularity and user engagement. This query is useful in understanding the overall performance and reach of your repositories, including the primary language used, number of forks, stars, subscribers, watchers, and the last updated date.

```sql+postgres
select
  name,
  owner_login,
  primary_language ->> 'name' as language,
  fork_count,
  stargazer_count,
  subscribers_count,
  watchers_total_count,
  updated_at as last_updated,
  description
from
  github_my_repository;
```

```sql+sqlite
select
  name,
  owner_login,
  json_extract(primary_language, '$.name') as language,
  fork_count,
  stargazer_count,
  subscribers_count,
  watchers_total_count,
  updated_at as last_updated,
  description
from
  github_my_repository;
```

### List your public repositories
Explore which of your repositories on GitHub are publicly accessible. This is useful for ensuring your private work remains confidential while sharing the projects you wish to showcase.

```sql+postgres
select
  name,
  is_private,
  visibility,
  owner_login
from
  github_my_repository
where
  not is_private;
```

```sql+sqlite
select
  name,
  is_private,
  visibility,
  owner_login
from
  github_my_repository
where
  not is_private;
```

OR

```sql+postgres
select
  name,
  is_private,
  visibility
from
  github_my_repository
where
  visibility = 'PUBLIC';
```

```sql+sqlite
select
  name,
  is_private,
  visibility
from
  github_my_repository
where
  visibility = 'PUBLIC';
```

### List all your repositories and their collaborators
Gain insights into the collaboration dynamics of your GitHub repositories by identifying the collaborators and their respective permissions. This can be useful in assessing access levels and ensuring proper project management.

```sql+postgres
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
 ,github_repository_collaborator c
where
  r.name_with_owner = c.repository_full_name;
```

```sql+sqlite
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
 ,github_repository_collaborator c
where
  r.name_with_owner = c.repository_full_name;
```

### List all your repository collaborators with admin or maintainer permissions
Identify individuals with elevated access rights within your repositories. This can help enhance security by ensuring only necessary permissions are granted.

```sql+postgres
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
 ,github_repository_collaborator c
where
  r.name_with_owner = c.repository_full_name
and
  permission in ('ADMIN', 'MAINTAIN');
```

```sql+sqlite
select
  r.name_with_owner as repository_full_name,
  c.user_login,
  c.permission
from
  github_my_repository r
join
  github_repository_collaborator c
on
  r.name_with_owner = c.repository_full_name
where
  c.permission in ('ADMIN', 'MAINTAIN');
```

### List repository hooks that are insecure
Discover the segments that have insecure settings in your GitHub repository hooks. This query is useful to identify potential security vulnerabilities in your repositories' hooks configuration.

```sql+postgres
select
  name as repository,
  hook
from
  github_my_repository,
  jsonb_array_elements(hooks) as hook
where
  hook -> 'config' ->> 'insecure_ssl' = '1'
  or hook -> 'config' ->> 'secret' is null
  or hook -> 'config' ->> 'url' not like '%https:%';
```

```sql+sqlite
select
  name as repository,
  hook.value as hook
from
  github_my_repository,
  json_each(hooks) as hook
where
  json_extract(hook.value, '$.config.insecure_ssl') = '1'
  or json_extract(hook.value, '$.config.secret') is null
  or json_extract(hook.value, '$.config.url') not like '%https:%';
```