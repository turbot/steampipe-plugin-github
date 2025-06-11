---
title: "Steampipe Table: github_organization_dependabot_alert - Query GitHub Dependabot Alerts using SQL"
description: "Allows users to query Dependabot Alerts in GitHub, specifically alert details such as status, severity, and package name, providing insights into security vulnerabilities in your GitHub organization's dependencies."
folder: "Dependabot"
---

# Table: github_organization_dependabot_alert - Query GitHub Dependabot Alerts using SQL

GitHub Dependabot is a feature that helps you keep your dependencies up to date. It monitors your project's dependencies and sends you an alert when updates or security vulnerabilities are detected. Dependabot Alerts provide critical information about security vulnerabilities that can affect your project's dependencies.

## Table Usage Guide

The `github_organization_dependabot_alert` table provides insights into Dependabot Alerts within GitHub. As a security analyst or a developer, explore alert-specific details through this table, including alert status, severity, and package name. Utilize it to uncover information about security vulnerabilities in your GitHub organization's dependencies, helping you to keep your projects safe and up to date.

**Important Notes**
- You must specify the `organization` column in `where` or `join` clause to query the table.
- To query this table using Fine-grained access tokens, the following permissions are required(The Fine-Grained access token should be created in Organization level):
  - **"Dependabot alerts" repository permissions (read)** – Required to access the all columns.

## Examples

### List dependabot alerts
Analyze the status and ecosystem of dependency packages in a specific organization using this query. It is particularly useful for identifying potential security vulnerabilities or outdated dependencies within your organization's codebase.

```sql+postgres
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org';
```

```sql+sqlite
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org';
```

### List open dependabot alerts
Discover the segments that have open alerts related to software dependencies within a specific organization. This can be used to identify areas that may be vulnerable or in need of updates, improving security and efficiency.

```sql+postgres
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open';
```

```sql+sqlite
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open';
```

### List open critical dependabot alerts
Determine the areas in which critical security threats are open in your organization's dependabot alerts. This query is useful for prioritizing security issues that need immediate attention.

```sql+postgres
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open'
  and security_advisory_severity = 'critical';
```

```sql+sqlite
select
  organization,
  state,
  dependency_package_ecosystem,
  dependency_package_name
from
  github_organization_dependabot_alert
where
  organization = 'my_org'
  and state = 'open'
  and security_advisory_severity = 'critical';
```