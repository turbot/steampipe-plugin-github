# Table: github_my_star

Stars are used to have of list of repositories you have interest on.

The `github_my_star` table can be used to query the repositories you have starred.

## Examples

### List of your starred repositories

```sql
select
  starred_at,
  repository_full_name,
  url
from
  github_my_star;
```
