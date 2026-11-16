---
title: "Steampipe Table: github_repository_code_scanning_alert - Query GitHub Code Scanning Alerts using SQL"
description: "Allows users to query code scanning alerts in GitHub repositories, providing insights into vulnerabilities and errors found by static analysis tools such as CodeQL."
folder: "Repository"
---

# Table: github_repository_code_scanning_alert - Query GitHub Code Scanning Alerts using SQL

GitHub code scanning analyzes the code in a repository with CodeQL or a third-party tool that uploads SARIF results, and raises an alert for each security vulnerability or coding error it finds.

## Table Usage Guide

The `github_repository_code_scanning_alert` table provides insights into code scanning alerts within GitHub repositories. As a security engineer, explore alert-specific details through this table, including the rule that raised the alert, its security severity, the tool that found it and where in the code its most recent instance is.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Code scanning alerts (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

A personal access token (classic) needs the `security_events` scope (`public_repo` for public repositories only) ([List code scanning alerts for a repository](https://docs.github.com/en/rest/code-scanning/code-scanning?apiVersion=2022-11-28#list-code-scanning-alerts-for-a-repository)).

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.
- A repository where code scanning has never run returns no rows: GitHub's 404 (`no analysis found`) and its 403 when GitHub Advanced Security or Code Security is not enabled for the repository are both treated as empty. A 403 for a missing token permission is still returned as an error.
- Because of this, no rows does not mean code scanning is enabled. To tell a repository where code scanning is not set up from one with no alerts, check for rows in the `github_repository_code_scanning_analysis` table.

## Examples

### List code scanning alerts
Explore the code scanning alerts raised for a repository, including the rule and tool that raised them.

```sql+postgres
select
  number,
  state,
  rule_id,
  rule_security_severity_level,
  tool_name,
  most_recent_instance_location_path
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo';
```

```sql+sqlite
select
  number,
  state,
  rule_id,
  rule_security_severity_level,
  tool_name,
  most_recent_instance_location_path
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo';
```

### List open critical and high severity alerts
Focus remediation on the open alerts with the highest security severity.

```sql+postgres
select
  number,
  rule_id,
  rule_description,
  rule_security_severity_level,
  html_url
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'open'
  and rule_security_severity_level in ('critical', 'high');
```

```sql+sqlite
select
  number,
  rule_id,
  rule_description,
  rule_security_severity_level,
  html_url
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'open'
  and rule_security_severity_level in ('critical', 'high');
```

### List dismissed alerts with the reason
Review which alerts were dismissed, by whom and why.

```sql+postgres
select
  number,
  rule_id,
  dismissed_by,
  dismissed_reason,
  dismissed_comment,
  dismissed_at
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'dismissed';
```

```sql+sqlite
select
  number,
  rule_id,
  dismissed_by,
  dismissed_reason,
  dismissed_comment,
  dismissed_at
from
  github_repository_code_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'dismissed';
```
