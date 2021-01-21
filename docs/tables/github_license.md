# Table: github_license

Github allows you to associate a license with your repository.  The `github_license` table lists information about the available licenses.

## Examples

### List basic license info

```sql
select
  spdx_id,
  name,
  description
from
  github_license;
```


### View license permissions

```sql
select
  name,
  jsonb_pretty(permissions)
from
  github_license
```


### Count repositories by license

```sql
select
  l.name,
  count(r.license_key) as num_repos
from
  github_license as l
  left join github_repository as r on l.key = r.license_key
group by
  l.name
order by
  num_repos desc;
```
