# Table: github_repository_content

Gets the contents of a file or directory in a repository.

Specify the file path or directory in `repository_content_path`.
If you omit `repository_content_path`, you will receive the contents of the repository's root directory.
See the description below regarding what the response includes for directories.

The `github_repository_content` table can be used to query information about **ANY** repository, and **you must specify which repository** in the where or join clause (`where repository_full_name=`, `join github_repository_content on repository_full_name=`).

## Examples

### List a repository

```sql
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs';
```

### List a directory in a repository

```sql
select
  repository_full_name,
  path,
  content,
  type,
  size,
  sha
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = 'docs';
```

### Get a file in a repository

```sql
select
  repository_full_name,
  path,
  type,
  size,
  sha,
  content
from
  github_repository_content
where
  repository_full_name = 'github/docs'
  and repository_content_path = '.vscode/settings.json';
```
