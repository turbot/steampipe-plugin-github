# Table: github_pull_request_review

The `github_pull_request_review` table can be used to query reviews from a specific pull request.

**You must specify `repository_full_name` (repository including org/user prefix) and `number` (of the pull request) in the WHERE or JOIN clause.**

## Examples

### List all reviews for a specific pull request

```sql
select
  id,
  author_login,
  author_association,
  state,
  body,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207;
```

### List reviews for a specific pull request which match a certain body content

```sql
select
  id,
  number as issue,
  author_login as comment_author,
  author_association,
  body as content,
  submitted_at,
  url
from
  github_pull_request_review
where
  repository_full_name = 'turbot/steampipe-plugin-github'
  and number = 207
  and body ~~ * '%minor%';
```

### List reviews for all open pull requests from a specific repository

```sql
select
  rv.*
from
  github_pull_request r
  join
    github_pull_request_review rv
    on r.repository_full_name = rv.repository_full_name
    and r.number = rv.number
where
  r.repository_full_name = 'turbot/steampipe-plugin-github'
  and r.state = 'OPEN';
```
