---
title: "Steampipe Table: github_actions_repository_runner - Query GitHub Actions Repository Runners using SQL"
description: "Allows users to query GitHub Actions Repository Runners, providing insights into the details of self-hosted runners for a repository."
---

# Table: github_actions_repository_runner - Query GitHub Actions Repository Runners using SQL

GitHub Actions is a CI/CD service provided by GitHub which allows automating software workflows. It enables developers to build, test, and deploy applications on GitHub, making it easier to integrate changes and deliver better software faster. A key feature of GitHub Actions is the ability to use self-hosted runners, which are servers with GitHub Actions runner application installed.

## Table Usage Guide

The `github_actions_repository_runner` table provides insights into the self-hosted runners of GitHub Actions for a specific repository. As a DevOps engineer, you can explore runner-specific details through this table, including runner IDs, names, operating systems, statuses, and associated metadata. Utilize it to manage and monitor your self-hosted runners, ensuring they are functioning properly and are up-to-date.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List runners
Explore which runners are associated with the 'turbot/steampipe' repository in GitHub Actions. This can be beneficial for maintaining and managing the performance and efficiency of your workflows.

```sql+postgres
select
  *
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  *
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe';
```

### List runners with mac operating system
Discover the segments that utilize a Mac operating system within a specific repository. This is beneficial for understanding the distribution of operating systems within your project, which can aid in troubleshooting and optimization efforts.

```sql+postgres
select
  repository_full_name,
  id,
  name,
  os
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe' and os = 'macos';
```

```sql+sqlite
select
  repository_full_name,
  id,
  name,
  os
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe' and os = 'macos';
```

### List runners which are in use currently
This query allows you to identify which runners are currently in use in the specified repository. This can be beneficial for understanding the utilization of resources and managing workflow runs in real-time.

```sql+postgres
select
  repository_full_name,
  id,
  name,
  os,
  busy
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe' and busy;
```

```sql+sqlite
select
  repository_full_name,
  id,
  name,
  os,
  busy
from
  github_actions_repository_runner
where
  repository_full_name = 'turbot/steampipe' and busy;
```