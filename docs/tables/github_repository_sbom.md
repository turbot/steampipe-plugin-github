---
title: "Steampipe Table: github_repository_sbom - Query GitHub Repositories using SQL"
description: "Allows users to query GitHub Repositories, specifically the Software Bill of Materials (SBOM), providing insights into the components, their versions, and dependencies in a software application."
folder: "Repository"
---

# Table: github_repository_sbom - Query GitHub Repositories using SQL

GitHub Repositories is a feature offered by GitHub that allows developers to store and share their code, manage projects, and collaborate with other developers. It provides a platform for version control and source code management (SCM) functionalities of Git, along with its own features. GitHub Repositories enable developers to maintain a history of project files, track changes, revert to previous versions, and collaborate seamlessly with their team.

## Table Usage Guide

The `github_repository_sbom` table provides insights into the Software Bill of Materials (SBOM) of GitHub Repositories. As a software developer or security analyst, explore the components, their versions, and dependencies in a software application through this table. Utilize it to uncover information about the SBOM, such as the components used in a software application, their versions, and dependencies, which can be crucial for vulnerability management and software maintenance.

**Important Notes**
- You must specify the `repository_full_name` (repository including org/user prefix) column in the `where` or `join` clause to query the table.

## Examples

### List SBOM packages with a specific package version
The query helps to identify software bill of materials (SBOM) packages within a specific GitHub repository that are using a specified version. This can be useful for maintaining version consistency or addressing potential vulnerabilities associated with certain versions.

```sql+postgres
select
  spdx_id,
  spdx_version,
  p ->> 'name' as package_name,
  p ->> 'versionInfo' as package_version,
  p ->> 'licenseConcluded' as package_license
from
  github_repository_sbom,
  jsonb_array_elements(packages) p
where
  p ->> 'versionInfo' = '2.6.0'
  and repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  spdx_id,
  spdx_version,
  json_extract(p.value, '$.name') as package_name,
  json_extract(p.value, '$.versionInfo') as package_version,
  json_extract(p.value, '$.licenseConcluded') as package_license
from
  github_repository_sbom,
  json_each(packages) as p
where
  json_extract(p.value, '$.versionInfo') = '2.6.0'
  and repository_full_name = 'turbot/steampipe';
```

### Find SBOMs conforming to a specific SPDX version
Analyze the settings to understand which Software Bill of Materials (SBOMs) align with a certain SPDX version in a given repository. This can help maintain compliance and compatibility with specific standards.

```sql+postgres
select
  name,
  spdx_version
from
  github_repository_sbom
where
  spdx_version = '2.2'
  and repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  spdx_version
from
  github_repository_sbom
where
  spdx_version = '2.2'
  and repository_full_name = 'turbot/steampipe';
```

### Retrieve SBOMs under a specific data license
Explore which Software Bill of Materials (SBOMs) are under a specific data license in a particular repository. This can help in assessing compliance with licensing requirements and managing intellectual property rights.

```sql+postgres
select
  name,
  data_license
from
  github_repository_sbom
where
  data_license = 'CC0-1.0'
  and repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  name,
  data_license
from
  github_repository_sbom
where
  data_license = 'CC0-1.0'
  and repository_full_name = 'turbot/steampipe';
```

### Find SBOMs created by a specific user or at a specific time
Determine the software bill of materials (SBOMs) created by a specific individual or at a certain date. This is useful for tracking changes and understanding the history of your software development.

```sql+postgres
select
  repository_full_name,
  creation_info
from
  github_repository_sbom
where
  (creation_info ->> 'created_by' = 'madhushreeray30' or creation_info ->> 'created_at' = '2023-11-16')
  and repository_full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  repository_full_name,
  creation_info
from
  github_repository_sbom
where
  (json_extract(creation_info, '$.created_by') = 'madhushreeray30' or json_extract(creation_info, '$.created_at') = '2023-11-16')
  and repository_full_name = 'turbot/steampipe';
```