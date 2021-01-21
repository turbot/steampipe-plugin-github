# Table: github_user

The `github_user` table does not list all users via the API - there is not currently an efficient way to limit the results in a useable way. As a result, you **must specify a user `login` in a `where`** or you will get no results.

## Examples

### Get information for a user

```sql
select
  *
from
  github_user
where
  login = 'torvalds';
```


### List of users which do not have two factor authentication turned on

```sql
select
  login,
  id,
  name,
  two_factor_authentication
from
  github_user
where
  not two_factor_authentication;
```
