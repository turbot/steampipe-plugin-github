# Table: github_community_profile

Community profile measures and information (e.g. README location) for a public repository.

Note:

- A `repository_full_name` must be provided in all queries to this table.
- Community profile data is only available for public GitHub repositories.

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
