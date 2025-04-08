---
title: "Steampipe Table: github_license - Query GitHub Licenses using SQL"
description: "Allows users to query GitHub Licenses, specifically providing detailed information about the various open-source licenses used across GitHub repositories."
folder: "Repository"
---

# Table: github_license - Query GitHub Licenses using SQL

GitHub Licenses are a set of permissions that developers grant to others to study, distribute, and modify their software. These licenses allow the software to be freely used, modified, and shared. They are crucial to the open-source community as they ensure the software remains open-source, even after modifications.

## Table Usage Guide

The `github_license` table provides insights into the different licenses used across GitHub repositories. As a software developer or open-source enthusiast, you can explore the specifics of these licenses through this table, including their permissions, conditions, and limitations. Use it to understand the terms under which you can use, modify, or distribute the software in question.

## Examples

### List basic license info
Explore the fundamental details of licensing in your GitHub repositories. This can help you ensure compliance with open source licensing requirements and understand the permissions associated with different licenses.

```sql+postgres
select
  spdx_id,
  name,
  description
from
  github_license;
```

```sql+sqlite
select
  spdx_id,
  name,
  description
from
  github_license;
```

### View license permissions
Explore the specific permissions associated with different licenses on GitHub. This helps in understanding the rights granted by each license, assisting users in making an informed choice when selecting a license for their project.

```sql+postgres
select
  name,
  jsonb_pretty(permissions)
from
  github_license;
```

```sql+sqlite
select
  name,
  permissions
from
  github_license;
```

### Count your repositories by license
Determine the number of your repositories grouped by their respective licenses. This is useful for understanding the distribution of license usage across your repositories.

```sql+postgres
with license_key as (
  select
    license_info ->> 'key' as key
  from
    github_my_repository
)
select
  l.name,
  count(k.key) as num_repos
from
  github_license as l
  left join license_key as k on l.key = k.key
group by
  l.name
order by
  num_repos desc;
```

```sql+sqlite
with license_key as (
  select
    license_info ->> 'key' as key
  from
    github_my_repository
)
select
  l.name,
  count(k.key) as num_repos
from
  github_license as l
  left join license_key as k on l.key = k.key
group by
  l.name
order by
  num_repos desc;
```

### View conditions for a specific license
Explore the specific conditions and their descriptions associated with a particular license on GitHub. This is particularly useful for understanding the terms of use and restrictions tied to a license before integrating it into your project.

```sql+postgres
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

```sql+sqlite
select
  name,
  l.key,
  json_extract(c.value, '$.Key') as condition,
  json_extract(c.value, '$.Description') as condition_desc
from
  github_license as l,
  json_each(conditions) as c
where
  l.key = 'gpl-3.0';
```

### View limitations for a specific license
Determine the restrictions associated with a specific software license, such as the 'gpl-3.0'. This is useful for understanding the terms and conditions that govern the use of the licensed software.

```sql+postgres
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

```sql+sqlite
select
  name,
  g.key,
  json_extract(l.value, '$.Key') as limitation,
  json_extract(l.value, '$.Description') as limitation_desc
from
  github_license as g,
  json_each(limitations) as l
where
  g.key = 'gpl-3.0';
```

### View permissions for a specific license
Explore the specific permissions associated with a particular software license. This can be useful when understanding the scope and limitations of a license before using or distributing software under its terms.

```sql+postgres
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

```sql+sqlite
select
  name,
  l.key,
  json_extract(p.value, '$.Key') as permission,
  json_extract(p.value, '$.Description') as permission_desc
from
  github_license as l,
  json_each(permissions) as p
where
  l.key = 'gpl-3.0';
```