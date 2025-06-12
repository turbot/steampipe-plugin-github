---
title: "Steampipe Table: github_branch_protection - Query GitHub Branch Protections using SQL"
description: "Allows users to query GitHub Branch Protections, specifically the protection rules applied to branches in a GitHub repository."
folder: "Branch"
---

# Table: github_branch_protection - Query GitHub Branch Protections using SQL

GitHub Branch Protection is a feature within GitHub that allows you to define certain rules for branches, particularly those that are part of the project's deployment process. These rules can include required status checks, required pull request reviews, and restrictions on who can push to the branch. Branch Protection helps maintain code integrity by enforcing workflow policies and preventing force pushes and accidental deletions.

## Table Usage Guide

The `github_branch_protection` table provides insights into branch protection rules within GitHub. As a DevOps engineer or a repository manager, explore branch-specific details through this table, including the enforcement of status checks, pull request reviews, and push restrictions. Utilize it to uncover information about branch protections, such as those with strict requirements, the enforcement of signed commits, and the restrictions on who can push to the branch.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
  - Repository permissions:
    - Administration (Read-only): Required to access general branch protection details.
    - Contents (Read-only): Required to access the `matching_branches` column.
    - Metadata (Read-only): Required to access general repository metadata.

## Examples

### List all branch protection rules for a repository
Explore the safety measures applied to a specific repository to understand its level of protection against unauthorized changes.

```sql+postgres
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe';
```

### Get a single branch protection rule by node id
Explore the specific rules and restrictions applied to a particular branch in a GitHub repository. This can be useful for understanding how different branches are managed and protected, and to ensure compliance with best practices for code review and collaboration.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which branch protection rules are not being used within the 'turbot/steampipe' repository. This can help in identifying unused rules and optimizing the repository's security configuration.

```sql+postgres
select
  *
from
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe'
and 
  matching_branches = 0;
```

```sql+sqlite
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
Explore which branches in the 'turbot/steampipe' repository have a policy requiring signed commits for merging. This can help maintain code integrity by ensuring only verified changes are merged.

```sql+postgres
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

```sql+sqlite
select 
  repository_full_name,
  pattern,
  matching_branches
from 
  github_branch_protection
where
  repository_full_name = 'turbot/steampipe'
and
  requires_commit_signatures = 1;
```