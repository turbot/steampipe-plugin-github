# Table: github_branch_protection

Branch protection is a set of rules protecting the branch from inappropriate changes.

The `github_branch_protection` table can be used to query information about any branch, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

GitHub users with admin permissions to a repository can manage branch protection rules.

## Examples

### List all branch protection rules for a repository

```sql
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe';
```

### List branch protection rules which are not currently utilised

```sql
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe'
and 
  matching_branches = 0;
```

### Get repositories that require signed commits for merging

```sql
select 
  repository_full_name,
  pattern,
  matching_branches
from 
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe'
and
  requires_commit_signatures = true;
```
