# Table: github_repository

A repository contains all of your project's files and each file's revision history.

The `github_repository` table can be used to query information about ANY repository, and **you must specify which repository** in the where or join clause  (`where full_name=`, `join github_repository on full_name=`).

To list all the **your** repositories use the `github_my_repository` table instead.  The `github_my_repository` table will list tables you own, you collaborate on, or that belong to your organizations.


## Examples

### Get information about a specific repository

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count,
  description
from
  github_repository
where 
  full_name = 'postgres/postgres'
```