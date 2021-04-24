# Table: github_release

A release is a package of software, along with release notes and links to binary files, for other people to use.

The `github_release` table can be used to query information about any release, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List releases

```sql
select
  name,
  published_at
from
  github_release
where
  repository_full_name = 'turbot/steampipe'
order by
  published_at desc
```

### Download statistics per release

```sql
select
  r.name as release_name,
  r.published_at,
  a ->> 'name' as asset_name,
  a ->> 'download_count' as download_count
from
  github_release as r,
  jsonb_array_elements(assets) as a
where
  r.repository_full_name = 'turbot/steampipe'
  and a ->> 'content_type' in ('application/zip', 'application/gzip')
order by
  r.published_at desc,
  asset_name
```
