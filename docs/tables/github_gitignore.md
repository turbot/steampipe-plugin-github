# Table: github_gitignore

GitHub allows you to associate a gitignore with your repository. The `github_gitignore` table lists information about the available gitignore templates.

## Examples

### List basic gitignore info

```sql
select
  *
from
  github_gitignore
order by
  name;
```

### View the source of the Go template

```sql
select
  source
from
  github_gitignore
where
  name = 'Go';
```
