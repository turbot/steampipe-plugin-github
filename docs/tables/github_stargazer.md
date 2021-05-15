# Table: github_stargazer

Stargazers are users who have starred the repository.

The `github_stargazer` table can be used to query stargazers belonging to a repository, and **you must specify which repository** with `where repository_full_name='owner/repository'`.


## Examples

### List the stargazers of a repository

```sql
select
  user_login,
  starred_at
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
order by
  starred_at desc;
```

### New stargazers by month

```sql
select
  to_char(starred_at, 'YYYY-MM') as month,
  count(*)
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe'
group by
  month
order by
  month;
```

### List stargazers with their contact information

```sql
select
  u.login,
  s.starred_at,
  u.name,
  u.company,
  u.email,
  u.html_url,
  u.twitter_username,
  u.blog,
  u.location,
  u.bio
from
  github_stargazer as s,
  github_user as u
where
  s.repository_full_name = 'turbot/steampipe'
  and s.user_login = u.login
order by
  s.starred_at desc;
```
