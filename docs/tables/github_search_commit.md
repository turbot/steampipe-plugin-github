# Table: github_search_commit

The `github_search_commit` table helps to find commits via various criteria on the default branch (usually master). You can search for commits globally across all of GitHub, or search for commits within a particular repository or organization.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. You can narrow the results using these commit search qualifiers in any combination.

## Examples

### List searched commits within organization's repositories

```sql
select
  sha,
  query,
  html_url,
  repository_full_name,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star org:turbot';
```

### List searched commits within a user's repositories

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star user:turbot';
```

### List searched commits within a user's specific repository

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'Add table github_my_star repo:turbot/steampipe-plugin-github';
```

### List searched commits by author

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'author:VincentHardouin repo:turbot/steampipe-plugin-github';
```

### List searched commits committed within the specified date range

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'committer-date:2021-11-16..2021-11-23 repo:turbot/steampipe-plugin-github';
```

### List searched commits by hash

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```

### List searched commits by parent

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'parent:b0566ea';
```

### List searched commits by merge commits

```sql
select
  sha,
  query,
  html_url,
  score
from
  github_search_commit
where
  query = 'merge:true repo:turbot/steampipe-plugin-azure';
```

### List repository details

```sql
select
  sha,
  repository -> 'name' as repo_name,
  repository ->> 'id' as repo_id,
  repository ->> 'html_url' as repo_html_url,
  jsonb_pretty(repository -> 'owner') as repo_owner
from
  github_search_commit
where
  query = 'hash:b0566eafb30e0595651c14a4c499b16e1c443767 repo:turbot/steampipe-plugin-github';
```
