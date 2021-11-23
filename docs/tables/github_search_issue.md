# Table: github_search_issue

The `github_search_issue` table helps to find issues by state and keyword. You can search for issues and pull requests globally across all of GitHub, or search for issues and pull requests within a particular organization.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. You can narrow the results using these search qualifiers in any combination.

## Examples

### List issues

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'github_search_commits is:issue';
```

### List pull requests

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'github_search_commits is:pr';
```

### List pull requests by the title, body, or comments

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'github_search_commits in:title in:body in:comments is:pr';
```

### List issues in open state assigned to a specific user

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'is:open is:issue assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List pull requests with public visibility assigned to a specific user

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'is:public is:pr assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List pull requests linked to a issue

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'is:open linked:issue repo:turbot/steampipe-plugin-github';
```

### List issues not linked to a pull request

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
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
  html_url
from
  github_search_issue
where
  query = 'label:blocked is:issue repo:turbot/steampipe-plugin-github';
```

### List pull requests with over 50 comments

```sql
select
  title,
  id,
  comments,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'org:turbot comments:>50 is:pr';
```

### List open draft pull requests

```sql
select
  title,
  id,
  state,
  created_at,
  html_url
from
  github_search_issue
where
  query = 'org:turbot draft:true state:open is:pr';
```

### List pull requests that took more than 30 days to close

```sql
select
  title,
  id,
  state,
  created_at,
  closed_at,
  html_url
from
  github_search_issue
where
  query = 'org:turbot state:closed is:pr'
  and closed_at > (created_at + interval '30' day);
```
