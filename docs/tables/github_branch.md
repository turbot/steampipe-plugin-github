# Table: github_branch

A branch is essentially is a unique set of code changes with a unique name.

The `github_branch` table can be used to query information about any branch, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List branches

```sql
select
  name,
  commit_sha,
  protected
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

### List commit details for each branch

```sql
select
  name,
  commit_sha,
  commit_short_sha,
  commit_author_login,
  commit_authored_date,
  commit_committer_login,
  commit_committed_date,
  commit_message,
  commit_url,
  commit_additions,
  commit_deletions,
  commit_changed_files,
  commit_signature_is_valid
from
  github_branch
where
  repository_full_name = 'turbot/steampipe'
order by
  commit_authored_date desc;
```

### List branch protection information for each protected branch

```sql
select
  name,
  protected,
  protection_rule_node_id,
  protection_rule_pattern,
  protection_rule_is_admin_enforced,
  protection_rule_allows_deletions,
  protection_rule_allows_force_pushes,
  protection_rule_blocks_creations,
  protection_rule_creator_login,
  protection_rule_dismisses_stale_reviews,
  protection_rule_lock_allows_fetch_and_merge,
  protection_rule_lock_branch,
  protection_rule_require_last_push_approval,
  protection_rule_requires_approving_reviews,
  protection_rule_requires_commit_signatures,
  protection_rule_restricts_pushes
from
  github_branch
where
  repository_full_name = 'turbot/steampipe'
and
  protected = true;
```
