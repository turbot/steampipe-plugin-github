# Table: github_issue

GitHub Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.  

The `github_issue` table can be used to query issues belonging to a repository, and **you must specify which repository** with `where repository_full_name='owner/repository'`.  To list all the issues **assigned to you across all repositories** use the `github_my_issue` table instead.

Note that pull requests are technically also issues in GitHub, however we do not include them in the `github_issue` table; You should use the `github_pull_request` table to query PRs.  


## Examples

### List the issues in a repository
```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  author_login,
  assignee_logins
from
  github_issue
where
  repository_full_name = 'turbot/steampipe';
```


### List the unassigned open issues in a repository

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  author_login,
  assignee_logins
from
  github_issue
where
  repository_full_name = 'turbot/steampipe'
  and jsonb_array_length(assignee_logins) = 0
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
  github_issue
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
  assigned_to
from
  github_issue,
  jsonb_array_elements_text(assignee_logins) as assigned_to
where
  repository_full_name = 'turbot/steampipe'
  and assigned_to = 'binaek89'
  and state = 'open';
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
  i.issue_number,
  i.title
from
  github_my_repository as r,
  github_issue as i
where 
  r.full_name like 'turbot/steampip%'
  and i.state = 'open'
  and i.repository_full_name = r.full_name;
```

