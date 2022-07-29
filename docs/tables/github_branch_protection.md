# Table: github_branch_protection

Branch protection is a set of rules protecting the branch from inappropriate changes.

The `github_branch_protection` table can be used to query information about any branch, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

GitHub users with admin permissions to a repository can manage branch protection rules.

## Examples

### List branches and their protection for a repository

```sql
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe';
```

### Get branch protection for a specific repo

```sql
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe'
  and name = 'main';
```

## Get repositories where coversation resolution is required for merging

```sql
select 
  repository_full_name,
  b.name as branch_name,
  signatures_protected_branch
from 
  github_branch_protection b 
  join github_my_repository r on r.full_name = b.repository_full_name 
  where required_conversation_resolution = true;
```

## Get repositories that require signed commits for merging

```sql
select 
  repository_full_name,
  b.name as branch_name,
  required_conversation_resolution
from 
  github_branch_protection b 
  join github_my_repository r on r.full_name = b.repository_full_name 
  where signatures_protected_branch = true;
```