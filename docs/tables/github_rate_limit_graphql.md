# Table: github_rate_limit_graphql

This table allows you to check the current rate limit status for the GitHub GraphQL API endpoint.

## Examples

### List rate limit info for GraphQL 

```sql
select
  limit,
  used,
  remaining,
  reset_at
from
  github_rate_limit_graphql;
```