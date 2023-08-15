# Table: github_pull_request

GitHub pull requests let you tell others about changes you've pushed to a branch in a repository on GitHub. Once a pull request is opened, you can discuss and review the potential changes with collaborators and add follow-up commits before your changes are merged into the base branch.

The `github_pull_request` table can be used to query issues belonging to a repository. **You must specify which repository** in a `where` or `join` clause (`where repository_full_name='`, `join github_pull_request on repository_full_name=`).

## Examples

### List open pull requests in a repository

```sql
select
  repository_full_name,
  number,
  title,
  state,
  mergeable
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN';
```

### List the pull requests for a repository that have been closed in the last week

```sql
select
  repository_full_name,
  number,
  title,
  state,
  closed_at,
  merged_at,
  merged_by_login
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'CLOSED'
  and closed_at >= (current_date - interval '7' day)
order by
  closed_at desc;
```

### List the open PRs in a repository with a given label

```sql
select
  repository_full_name,
  number,
  title,
  state,
  labels,
  tags
from
  github_pull_request
where
  repository_full_name = 'turbot/steampipe'
  and state = 'OPEN'
  and tags ? 'bug';
```

### List the open PRs in a repository assigned to a specific user

```sql
select
  repository_full_name,
  number,
  title,
  state,
  assigned_to
from
  github_pull_request,
  jsonb_array_elements_text(assignee_logins) as assigned_to
where
  repository_full_name = 'turbot/steampipe'
  and assigned_to = 'binaek89'
  and state = 'OPEN';
```

### Join with github_my_repository to find open PRs in multiple repos

```sql
select
  i.repository_full_name,
  i.number,
  i.title
from
  github_my_repository as r,
  github_pull_request as i
where 
  r.full_name like 'turbot/steampip%'
  and i.state = 'OPEN'
  and i.repository_full_name = r.full_name;
```

### List open PRs in a repository with an array of associated labels

```sql
select
  r.repository_full_name
  r.number,
  r.title,
  jsonb_agg(l ->> 'name') as labels
from
  github_pull_request r,
  jsonb_array_elements(r.labels_src) as l
where
  repository_full_name = 'turbot/steampipe'
and
  state = 'OPEN'
group by
  r.repository_full_name, r.number, r.title;
```

OR

```sql
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_pull_request r,
  jsonb_object_keys(r.labels) as t
where
  repository_full_name = 'turbot/steampipe'
and 
  state = 'OPEN'
group by
  repository_full_name, number, title;
```

### List all open PRs in a repository with a specific label

```sql
select
  repository_full_name,
  number,
  title,
  json_agg(t) as labels
from
  github_pull_request r,
  jsonb_object_keys(labels) as t
where
  repository_full_name = 'turbot/steampipe'
and
  state = 'OPEN'
and
  labels ? 'bug'
group by
  repository_full_name, number, title;
```