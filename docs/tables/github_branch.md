---
title: "Steampipe Table: github_branch - Query GitHub Branches using SQL"
description: "Allows users to query GitHub Branches, specifically the details of branches for a repository, providing insights into branch information such as name, commit, protected status, and more."
---

# Table: github_branch - Query GitHub Branches using SQL

GitHub Branches are references within a repository that point to the state of a project at a specific point in time. They represent an isolated line of development and can be created, deleted, and manipulated without affecting the rest of the project. Branches are a crucial part of the GitHub workflow, enabling collaboration and iteration on different features or fixes.

## Table Usage Guide

The `github_branch` table provides insights into branches within GitHub repositories. As a developer or project manager, explore branch-specific details through this table, including name, commit, protected status, and associated metadata. Utilize it to uncover information about branches, such as those with protected status, the commit associated with each branch, and the verification of branch protections.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List branches
Discover the segments that represent specific branches of a GitHub repository, including their protection status and associated commit details. This is particularly useful for understanding the structure and security measures of a project.

```sql+postgres
select
  name,
  commit ->> 'sha' as commit_sha,
  protected
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  json_extract(commit, '$.sha') as commit_sha,
  protected
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

### List commit details for each branch
Discover the specifics of each commit made in different branches of a particular GitHub repository. This allows you to track changes, observe trends, and analyze the contribution of different authors over time.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract("commit", '$.sha') as commit_sha,
  json_extract("commit", '$.message') as commit_message,
  json_extract("commit", '$.url') as commit_url,
  json_extract(json_extract("commit", '$.author'), '$.user.login') as author,
  json_extract("commit", '$.authored_date') as authored_date,
  json_extract(json_extract("commit", '$.committer'), '$.user.login') as committer,
  json_extract("commit", '$.committed_date') as committed_date,
  json_extract("commit", '$.additions') as additions,
  json_extract("commit", '$.deletions') as deletions,
  json_extract("commit", '$.changed_files') as changed_files
from
  github_branch
where
  repository_full_name = 'turbot/steampipe';
```

### List branch protection information for each protected branch
Gain insights into the protection rules applied to each safeguarded branch within a specific GitHub repository. This can help ensure the repository's integrity by understanding the restrictions and permissions set for each branch.

```sql+postgres
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

```sql+sqlite
select
  name,
  protected,
  json_extract(branch_protection_rule, '$.id') as rule_id,
  json_extract(branch_protection_rule, '$.node_id') as rule_node_id,
  json_extract(branch_protection_rule, '$.allows_deletions') as allows_deletions,
  json_extract(branch_protection_rule, '$.allows_force_pushes') as allows_force_pushes,
  json_extract(branch_protection_rule, '$.creator_login') as rule_creator,
  json_extract(branch_protection_rule, '$.requires_commit_signatures') as requires_signatures,
  json_extract(branch_protection_rule, '$.restricts_pushes') as restricts_pushes
from
  github_branch
where
  repository_full_name = 'turbot/steampipe'
and
  protected = 1;
```