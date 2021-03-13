# Table: github_my_issue

Github Issues are used to track ideas, enhancements, tasks, or bugs for work on GitHub.  The `github_my_issue` table lists issues that are assigned to you, across all repositories.

To view **all the issues belonging to a repository**, use the `github_issue` table.


## Examples

### List all of the open issues assigned to you
```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  author_login,
  assignee_logins
from
  github_my_issue
where
  state = 'open';
```


### List your open issues with a given label

```sql
select
  repository_full_name,
  issue_number,
  title,
  state,
  tags
from
  github_my_issue
where
  state = 'open'
  and tags ? 'bug';
```


### List your 10 oldest open issues

```sql
select
  repository_full_name,
  issue_number,
  created_at,
  age (created_at),
  title,
  state
from
  github_my_issue
where
  state = 'open'
order by
  created_at
limit 10;
```


