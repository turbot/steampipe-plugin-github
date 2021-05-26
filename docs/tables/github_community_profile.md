# Table: github_community_profile

Community profile measures and information (e.g. README location) for a repository.

The `github_community_profile` table can be used to query information about any tag, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### Get community profile information for the repository

```sql
select
  *
from
  github_community_profile
where
  repository_full_name = 'turbot/steampipe'
```
