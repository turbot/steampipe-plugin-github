# Table: github_branch

A branch is essentially is a unique set of code changes with a unique name.

The `github_branch` table can be used to query information about any branch, and **you must specify which repository** in the where or join clause using the `repository_full_name` column.

## Examples

### List branches

```sql
select
  name,
  commit ->> 'sha' as commit_sha,
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
  commit ->> 'sha' as commit_sha,
  commit ->> 'message' as commit_message,
  commit ->> 'url' as commit_url,
  commit -> 'author' -> 'user' ->> 'login' as author,
  commit ->> 'authored_date' as authored_date,
  commit -> 'committer' -> 'user' ->> 'login' as committer,
  commit ->> 'committed_date' as committed_date,
  commit ->> 'additions' as additions,
  commit ->> 'deletions' as deletions,
  commit ->> 'changed_files' as changed_files
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

### List branch protection information for each protected branch

```sql
select
  name,
  protected,
  branch_protection_rule ->> 'id' as rule_id,
  branch_protection_rule ->> 'node_id' as rule_node_id,
  branch_protection_rule ->> 'allows_deletions' as allows_deletions,
  branch_protection_rule ->> 'allows_force_pushes' as allows_force_pushes,
  branch_protection_rule ->> 'creator_login' as rule_creator,
  branch_protection_rule ->> 'requires_commit_signatures' as requires_signatures,
  branch_protection_rule ->> 'restricts_pushes' as restricts_pushes
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
and
  protected = true;
```
