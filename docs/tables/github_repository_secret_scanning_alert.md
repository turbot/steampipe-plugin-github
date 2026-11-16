---
title: "Steampipe Table: github_repository_secret_scanning_alert - Query GitHub Secret Scanning Alerts using SQL"
description: "Allows users to query secret scanning alerts in GitHub repositories, providing insights into credentials and tokens committed to repository history."
folder: "Repository"
---

# Table: github_repository_secret_scanning_alert - Query GitHub Secret Scanning Alerts using SQL

GitHub secret scanning scans a repository's Git history for known types of secrets, such as API keys and tokens, and raises an alert for each one it finds. Push protection can additionally block pushes that contain a supported secret.

## Table Usage Guide

The `github_repository_secret_scanning_alert` table provides insights into secret scanning alerts within GitHub repositories. As a security engineer, explore alert-specific details through this table, including the secret type, its validity, whether push protection was bypassed, and how the alert was resolved.

The table never returns the detected secret itself. Requests are sent with `hide_secret=true` and the `secret` field of the API response is not read, so no column contains the leaked value.

To query this table using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), the following permissions are required:
- Repository permissions:
  - Secret scanning alerts (Read-only): Required to access all columns.
  - Metadata (Read-only): Required to access general repository metadata.

A personal access token (classic) needs the `repo` or `security_events` scope (`public_repo` for public repositories only). GitHub requires the authenticated user to be an administrator of the repository or of the organization that owns it ([List secret scanning alerts for a repository](https://docs.github.com/en/rest/secret-scanning/secret-scanning?apiVersion=2022-11-28#list-secret-scanning-alerts-for-a-repository)).

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.
- GitHub returns a 404 when secret scanning is disabled for the repository, which this table returns as no rows. To tell a repository with secret scanning disabled from one with no alerts, check `security_and_analysis -> 'secret_scanning' ->> 'status'` in the `github_repository` table.

## Examples

### List secret scanning alerts
Explore the secret scanning alerts raised for a repository, including the type of secret detected and the alert state.

```sql+postgres
select
  number,
  state,
  secret_type_display_name,
  validity,
  created_at
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo';
```

```sql+sqlite
select
  number,
  state,
  secret_type_display_name,
  validity,
  created_at
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo';
```

### List open alerts for secrets that are still active
Find leaked secrets that GitHub reports as still valid, which should be revoked first.

```sql+postgres
select
  number,
  secret_type_display_name,
  html_url
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'open'
  and validity = 'active';
```

```sql+sqlite
select
  number,
  secret_type_display_name,
  html_url
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and state = 'open'
  and validity = 'active';
```

### List alerts where push protection was bypassed
Identify secrets that were pushed even though push protection blocked them, and who bypassed it.

```sql+postgres
select
  number,
  secret_type_display_name,
  push_protection_bypassed_by,
  push_protection_bypassed_at
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and push_protection_bypassed;
```

```sql+sqlite
select
  number,
  secret_type_display_name,
  push_protection_bypassed_by,
  push_protection_bypassed_at
from
  github_repository_secret_scanning_alert
where
  repository_full_name = 'my-org/my-repo'
  and push_protection_bypassed = 1;
```
