---
title: "Steampipe Table: github_repository_code_scanning_analysis - Query GitHub Code Scanning Analyses using SQL"
description: "Allows users to query code scanning analyses uploaded to GitHub repositories, providing insights into which tools scan a repository and when they last ran."
folder: "Repository"
---

# Table: github_repository_code_scanning_analysis - Query GitHub Code Scanning Analyses using SQL

Each time CodeQL or a third-party tool scans a repository and uploads its results, GitHub records a code scanning analysis for the commit and Git reference that was scanned.

## Table Usage Guide

The `github_repository_code_scanning_analysis` table provides insights into the code scanning analyses of GitHub repositories. As a security engineer, use it to confirm that code scanning is set up for a repository, which tools scan it, and when each last ran.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Code scanning alerts (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

A personal access token (classic) needs the `security_events` scope (`public_repo` for public repositories only) ([List code scanning analyses for a repository](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#list-code-scanning-analyses-for-a-repository)).

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.
- A repository where code scanning has never run returns no rows: GitHub's 404 (`no analysis found`) and its 403 when GitHub Advanced Security or Code Security is not enabled for the repository are both treated as empty. A 403 for a missing token permission is still returned as an error.

## Examples

### Check whether code scanning is set up for a repository
A repository with at least one analysis has code scanning set up.

```sql+postgres
select
  exists (
    select
      1
    from
      github_repository_code_scanning_analysis
    where
      repository_full_name = 'my-org/my-repo'
  ) as code_scanning_set_up;
```

```sql+sqlite
select
  exists (
    select
      1
    from
      github_repository_code_scanning_analysis
    where
      repository_full_name = 'my-org/my-repo'
  ) as code_scanning_set_up;
```

### List the most recent analysis for each tool
Find which tools scan a repository and when each last uploaded results.

```sql+postgres
select
  tool_name,
  max(created_at) as last_analysis_at
from
  github_repository_code_scanning_analysis
where
  repository_full_name = 'my-org/my-repo'
group by
  tool_name;
```

```sql+sqlite
select
  tool_name,
  max(created_at) as last_analysis_at
from
  github_repository_code_scanning_analysis
where
  repository_full_name = 'my-org/my-repo'
group by
  tool_name;
```

### List analyses for the default branch
Review the analyses of a specific Git reference.

```sql+postgres
select
  id,
  tool_name,
  tool_version,
  results_count,
  created_at
from
  github_repository_code_scanning_analysis
where
  repository_full_name = 'my-org/my-repo'
  and ref = 'refs/heads/main';
```

```sql+sqlite
select
  id,
  tool_name,
  tool_version,
  results_count,
  created_at
from
  github_repository_code_scanning_analysis
where
  repository_full_name = 'my-org/my-repo'
  and ref = 'refs/heads/main';
```
