# Table: github_search_issue

The `github_search_issue` table helps to find issues by state and keyword. You can search for issues globally across all of GitHub, or search for issues within a particular organization.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. You can narrow the results using these search qualifiers in any combination. See [Searching issues and pull requests](https://docs.github.com/search-github/searching-on-github/searching-issues-and-pull-requests) for details on the GitHub query syntax.

## Examples

### List issues by the title, body, or comments

```sql
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'github_search_commit in:title in:body';
```

### List issues in open state assigned to a specific user

```sql
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List issues with public visibility assigned to a specific user

```sql
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:public assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List issues not linked to a pull request

```sql
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'is:open -linked:pr repo:turbot/steampipe-plugin-github';
```

### List blocked issues

```sql
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'label:blocked repo:turbot/steampipe-plugin-github';
```

### List issues with over 10 comments

```sql
select
  title,
  id,
  comments_total_count,
  state,
  created_at,
  url
from
  github_search_issue
where
  query = 'org:turbot comments:>10';
```

### List issues that took more than 30 days to close

```sql
select
  title,
  id,
  state,
  created_at,
  closed_at,
  url
from
  github_search_issue
where
  query = 'org:turbot state:closed'
  and closed_at > (created_at + interval '30' day);
```
