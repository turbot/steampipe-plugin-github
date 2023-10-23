# Table: github_repository_sbom

The `github_repository_sbom` table can be used to query information about packages listed in the SBOM of a repository.

**You must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List SBOM packages

```sql
select
  github_repository_sbom,
  p->>'name' as package_name,
  p->>'versionInfo' as package_version,
  p->>'licenseConcluded' as package_license
from
  github_repository_sbom,
  jsonb_array_elements(packages) p
where
  repository_full_name = 'turbot/steampipe';
```
