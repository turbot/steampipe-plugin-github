# Table: github_repository_sbom

The `github_repository_sbom` table can be used to query information about packages listed in the SBOM of a repository.

**You must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List SBOM packages with a specific package version

```sql
select
  spdx_id,
  spdx_version,
  p ->> 'name' as package_name,
  p ->> 'versionInfo' as package_version,
  p ->> 'licenseConcluded' as package_license
from
  github_repository_sbom,
  jsonb_array_elements(packages) p
where
  p ->> 'versionInfo' = '2.6.0'
  and repository_full_name = 'turbot/steampipe';
```

### Find SBOMs conforming to a specific SPDX version

```sql
select
  name,
  spdx_version
from
  github_repository_sbom
where
  spdx_version = '2.2'
  and repository_full_name = 'turbot/steampipe';
```

### Retrieve SBOMs under a specific data license

```sql
select
  name,
  data_license
from
  github_repository_sbom
where
  data_license = 'CC0-1.0'
  and repository_full_name = 'turbot/steampipe';
```

### Find SBOMs created by a specific user or at a specific time

```sql
select
  repository_full_name,
  creation_info
from
  github_repository_sbom
where
  (creation_info ->> 'created_by' = 'madhushreeray30' or creation_info ->> 'created_at' = '2023-11-16')
  and repository_full_name = 'turbot/steampipe';
```