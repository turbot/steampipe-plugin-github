# Table: github_issue

GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.

The `github_issue` table can be used to query issues belonging to a repository, and **you must specify which repository** with `where repository_full_name='owner/repository'`. To list all the issues **assigned to you across all repositories** use the `github_my_issue` table instead.

Note that pull requests are technically also issues in GitHub, however we do not include them in the `github_issue` table; You should use the `github_pull_request` table to query PRs.

## Examples

### List the issues in a repository

```sql
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe';
```

### List the unassigned open issues in a repository

```sql
select
  repository_full_name,
  number,
  title,
  state,
  author_login,
  assignees_total_count
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
and 
  assignees_total_count = 0
and 
  state = 'OPEN';

```

### Report of the number issues in a repository by author

```sql
select
  author_login,
  count(*) as num_issues
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
group by
  author_login
order by
  num_issues desc;
```

### Join with github_my_repository to find open issues in multiple repos that you own or contribute to

```sql
select
  i.repository_full_name,
  i.number,
  i.title
from
  github_my_repository as r,
  github_issue as i
where 
  r.full_name like 'turbot/steampip%'
  and i.state = 'OPEN'
  and i.repository_full_name = r.full_name;
```

### List all issues with labels as a string array (instead of JSON objects)

```sql
select
  i.repository_full_name
  i.number,
  i.title,
  json_agg(l ->> 'name') as labels
from
  github_issue i,
  jsonb_array_elements(i.labels) as l
where
  repository_full_name = 'turbot/steampipe'
group by
  i.repository_full_name, i.number, i.title;
```

OR

```sql
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_issue i,
  jsonb_object_keys(i.tags) as t
where
  repository_full_name = 'turbot/steampipe'
and
  state = 'OPEN'
group by
  repository_full_name, number, title;
```

### List all issues in a repository with a specific label

```sql
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_issue i,
  jsonb_object_keys(i.tags) as t
where
  repository_full_name = 'turbot/steampipe'
and
  state = 'OPEN'
and
  tags ? 'bug'
group by
  repository_full_name, number, title;
```