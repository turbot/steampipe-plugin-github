# Table: github_rate_limit

With the Rate Limit API, you can check the current rate limit status of various REST APIs.

## Examples

### List rate limit of rest apis

```sql
select
  core_limit,
  core_remaining,
  search_limit,
  search_remaining
from
  github_rate_limit;
```