# Table: github_repository_sbom

The `github_repository_sbom` table can be used to query information about packages listed in the SBOM of a repository.

**You must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List SBOM packages

```sql
select
  name,
  version,
  license
from
  github_repository_sbom
where
  repository_full_name = 'turbot/steampipe';
```
