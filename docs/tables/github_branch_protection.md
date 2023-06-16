# Table: github_branch_protection

Branch protection is a set of rules protecting the branch from inappropriate changes.

The `github_branch_protection` table can be used to query information about any branch, and **you must specify which repository (or node_id for a single item)** in the where or join clause using the `repository_full_name` column.

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

### Get a single branch protection rule by node id

```sql
select
  node_id,
  matching_branches,
  is_admin_enforced,
  allows_deletions,
  allows_force_pushes,
  blocks_creations,
  creator_login,
  dismisses_stale_reviews,
  lock_allows_fetch_and_merge,
  lock_branch,
  require_last_push_approval,
  requires_approving_reviews,
  requires_commit_signatures,
  restricts_pushes,
  push_allowance_apps,
  push_allowance_apps,
  push_allowance_users
from
  github_branch_protection
where
  node_id = 'BPR_xxXXXX0X0X0XXXX0';
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
