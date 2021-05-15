# Table: github_traffic_view_daily

Daily views to the repository over the last 14 days.

The `github_traffic_view_daily` table can be used to query information about any tag, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List view statistics

```sql
select
  timestamp,
  count,
  uniques
from
  github_traffic_view_daily
where
  repository_full_name = 'turbot/steampipe'
order by
  timestamp;
```
