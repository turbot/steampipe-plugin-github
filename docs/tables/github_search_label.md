# Table: github_search_label

The `github_search_label` table helps to find labels in a repository with names or descriptions that match search keywords. 

 **You must always include at least one search term and repository ID when searching source code** in the where or join clause using the `query` and `repository_id` columns respectively.

## Examples

### List labels for bug, enhancement and blocked

```sql
select
  id,
  repository_id,
  name,
  repository_full_name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'bug enhancement blocked';
```

### List labels where specific text matches in name or description

```sql
select
  id,
  repository_id,
  name,
  description
from
  github_search_label
where
  repository_id = 331646306 and query = 'work';
```
