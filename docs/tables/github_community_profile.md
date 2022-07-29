# Table: github_community_profile

Community profile measures and information (e.g. README location) for a public repository.

Note:

- A `repository_full_name` must be provided in all queries to this table.
- Community profile data is only available for public GitHub repositories.

## Examples

### Get community profile information for the repository

```sql
select
  *
from
  github_community_profile
where
  repository_full_name = 'turbot/steampipe';
```

### List repositories having the security file

```sql
select 
  repository_full_name, 
  security ->> 'html_url' as security_file_url,
  security ->> 'name' as security_file_name
from
  github_community_profile c 
  join github_my_repository r on r.full_name = c.repository_full_name
  where security is not null;
```