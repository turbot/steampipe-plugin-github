# Table: github_search_repository

The `github_search_repository` table helps to find repositories via various criteria. You can search for repositories on GitHub and narrow the results using these repository search qualifiers in any combination.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. You can search for repositories globally across all of GitHub.com, or search for repositories within a particular organization. See [Searching for repositories](https://docs.github.com/search-github/searching-on-github/searching-for-repositories) for details on the GitHub query syntax.


## Examples

### Get a specific repository

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'repo:turbot/steampipe-plugin-github';
```

### List repositories based on contents of a repository

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'stargazers in:readme repo:turbot/steampipe-plugin-github';
```

### List repositories with more than 100000 followers

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'followers:>=100000';
```

### List forked repositories created within specific timestamp

```sql
select
  name,
  owner_login,
  language,
  forks_count,
  stargazers_count,
  subscribers_count,
  watchers_count
from
  github_search_repository
where
  query = 'tinyspotifyr in:name created:2021-01-01..2021-01-05 fork:only';
```
