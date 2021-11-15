# Table: github_my_star

Stars are used to have of list of repositories you have interest on.

The `github_my_star` table can be used to query the repositories you have starred.

## Examples

### List of your starred repositories

```sql
select
  starred_at,
  repository_full_name
from
  github_my_star;
```

### Show all details from your starred repositories

```sql
select
  gr.homepage
from
  github_my_star gs
inner join
  github_repository gr
on
  gs.repository_full_name = gr.full_name;
```
