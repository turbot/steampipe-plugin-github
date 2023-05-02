# Table: github_license

GitHub allows you to associate a license with your repository. The `github_license` table lists information about the available licenses.

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
  github_license;
```

### Count your repositories by license

```sql
select
  l.name,
  count(r.license_key) as num_repos
from
  github_license as l
  left join github_my_repository as r on l.key = r.license_key
group by
  l.name
order by
  num_repos desc;
```

### View conditions for a specific license

```sql
select
  name,
  key,
  c ->> 'Key' as condition,
  c ->> 'Description' as condition_desc
from
  github_license,
  jsonb_array_elements(conditions) as c
where
  key = 'gpl-3.0';
```

### View limitations for a specific license

```sql
select
  name,
  key,
  l ->> 'Key' as limitation,
  l ->> 'Description' as limitation_desc
from
  github_license,
  jsonb_array_elements(limitations) as l
where
  key = 'gpl-3.0';
```

### View permissions for a specific license

```sql
select
  name,
  key,
  p ->> 'Key' as permission,
  p ->> 'Description' as permission_desc
from
  github_license,
  jsonb_array_elements(permissions) as p
where
  key = 'gpl-3.0';
```