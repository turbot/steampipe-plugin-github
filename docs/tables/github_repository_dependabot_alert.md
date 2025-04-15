---
title: "Steampipe Table: github_repository_dependabot_alert - Query GitHub Dependabot Alerts using SQL"
description: "Allows users to query Dependabot Alerts in GitHub repositories, providing insights into potential security vulnerabilities within project dependencies."
folder: "Dependabot"
---

# Table: github_repository_dependabot_alert - Query GitHub Dependabot Alerts using SQL

GitHub Dependabot is a feature within the GitHub platform that monitors your project dependencies for known security vulnerabilities and automatically opens pull requests to update them to the minimum required version. It provides an automated way to keep your project dependencies up-to-date and secure. GitHub Dependabot helps you maintain the security and reliability of your projects by identifying and suggesting updates for vulnerable dependencies.

## Table Usage Guide

The `github_repository_dependabot_alert` table provides insights into Dependabot alerts within GitHub repositories. As a project maintainer or security engineer, explore alert-specific details through this table, including the dependency name, version, and associated security vulnerabilities. Utilize it to uncover information about potentially insecure dependencies, helping you to maintain the security and integrity of your projects.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List dependabot alerts
Identify the status and type of dependabot alerts for a specific repository to maintain and upgrade dependencies efficiently.

```sql+postgres
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe';
```

### List open dependabot alerts
Discover the segments that have active dependency alerts within a specific GitHub repository. This query is useful for maintaining security and up-to-date dependencies in your projects.

```sql+postgres
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open';
```

```sql+sqlite
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open';
```

### List open critical dependabot alerts
Explore critical alerts in your repository's dependencies that are currently open. This is useful for quickly identifying potential security risks within your project's ecosystem.

```sql+postgres
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open'
  and security_advisory_severity = 'critical';
```

```sql+sqlite
select
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_repository_dependabot_alert
where
  repository_full_name = 'turbot/steampipe'
  and state = 'open'
  and security_advisory_severity = 'critical';
```