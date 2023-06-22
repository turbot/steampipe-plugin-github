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
  user_login,
  starred_at,
  user_detail ->> 'name' as name,
  user_detail ->> 'company' as company,
  user_detail ->> 'email' as email,
  user_detail ->> 'url' as url,
  user_detail ->> 'twitter_username' as twitter_username,
  user_detail ->> 'website_url' as website,
  user_detail ->> 'location' as location,
  user_detail ->> 'bio' as bio
from
  github_stargazer
where
  repository_full_name = 'turbot/steampipe-plugin-github'
order by
  starred_at desc;
```
