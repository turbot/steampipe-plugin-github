# Table: github_search_user

The `github_search_user` table helps to find users and organizations via various criteria. You can filter your search to the personal user or organization account name with `user` or `org` qualifiers.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column.

## Examples

### List users

```sql
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbot in:name type:user';
```

### List organizations

```sql
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'turbotio in:login type:org';
```

### Get user with specific username

```sql
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'user:c0d3r-arnab';
```

### List organizations with over 10000 repositories

```sql
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'repos:>10000 type:org';
```

### List users and organizations created between specific timestamp

```sql
select
  id,
  login,
  type,
  url
from
  github_search_user
where
  query = 'created:2021-01-01..2021-01-31 turbot';
```
