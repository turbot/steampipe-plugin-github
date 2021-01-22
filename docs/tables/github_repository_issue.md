# Table: github_repository_issue

Github Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.  Note that Pull Requests are ALSO issues in Github.

The `github_repository_issue` table can be used to query issues belonging to a repository, and **you must specify which repository** with `where repository_full_name='owner/repository'`.   


## Examples

### List the issues in a repository
```sql
select
  repository_full_name,
  issue_number,
  title,
  state
from
  github_repository_issue
where
  repository_full_name = 'turbot/steampipe';
```

### List pull requests in a repository

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  pull_request_links ->> 'html_url' as pr_link
from
  github_repository_issue
where
  repository_full_name = 'turbot/steampipe'
  and pull_request_links is not null;
```

### List the unassigned open issues in a repository

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  assignees
from
  github_repository_issue
where
  repository_full_name = 'turbot/steampipe'
  and jsonb_array_length(assignees) = 0
  and state = 'open';
```

### List the open issues in a repository with a given label

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  tags
from
  github_repository_issue
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open'
  and tags ? 'bug';
```


### List the open issues in a repository assigned to a specific user

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  a ->> 'login' as assigned_to
from
  github_repository_issue,
  jsonb_array_elements(assignees) as a
where
  repository_full_name = 'turbot/steampipe'
  and a ->> 'login' = 'binaek89'
  and state = 'open';
```


### Report of the number issues in a repository by author

```sql
select
  author ->> 'login' as author,
  count(*) as num_issues
from
  github_repository_issue
where
  repository_full_name = 'turbot/steampipe'
group by
  author ->> 'login'
order by
  num_issues desc;
```






