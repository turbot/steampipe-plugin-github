---
title: "Steampipe Table: github_search_pull_request - Query GitHub Pull Requests using SQL"
description: "Allows users to query GitHub Pull Requests. This table provides extensive details about pull requests across repositories, including the status, creator, and assignee information."
---

# Table: github_search_pull_request - Query GitHub Pull Requests using SQL

GitHub Pull Requests is a feature within GitHub that allows developers to propose changes to a repository. It provides a platform for code review and discussion about the proposed changes before they are merged into the codebase. GitHub Pull Requests helps in maintaining the integrity and quality of the code in a repository by ensuring that all changes are reviewed and approved before they are incorporated.

## Table Usage Guide

The `github_search_pull_request` table provides insights into pull requests within GitHub. As a developer or project manager, explore pull request-specific details through this table, including status, creator, assignee, and associated metadata. Utilize it to monitor the progress of proposed changes, manage code reviews, and ensure the quality of the code in your repositories.

**Important Notes**
- You must always include at least one search term when searching pull requests in the where or join clause using the `query` column. You can narrow the results using these search qualifiers in any combination. See [Searching issues and pull requests](https://docs.github.com/search-github/searching-on-github/searching-issues-and-pull-requests) for details on the GitHub query syntax.

## Examples

### List pull requests by the title, body, or comments
Explore which GitHub pull requests match certain criteria within the title, body, or comments. This can be useful for identifying relevant discussions or changes related to specific topics or keywords.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  repository_full_name,
  url
from
  github_search_pull_request
where
  query = 'github_search_issue in:title in:body in:comments';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  repository_full_name,
  url
from
  github_search_pull_request
where
  query = 'github_search_issue in:title in:body in:comments';
```

### List pull requests in open state assigned to a specific user
Determine the areas in which a specific user has been assigned open pull requests for a particular repository. This is useful for project managers to track individual contributions and progress in a collaborative environment.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:open assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:open assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List pull requests with public visibility assigned to a specific user
Determine the areas in which a specific user has been assigned public visibility pull requests in a particular GitHub repository. This can be useful to track a user's involvement and contribution to publicly visible projects.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:public assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:public assignee:c0d3r-arnab repo:turbot/steampipe-plugin-github';
```

### List pull requests not linked to an issue
Identify instances where open pull requests are not linked to an issue within the 'turbot/steampipe-plugin-github' repository. This can help to uncover potential oversights in issue tracking and ensure all code changes are properly documented.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:open -linked:issue repo:turbot/steampipe-plugin-github';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'is:open -linked:issue repo:turbot/steampipe-plugin-github';
```

### List pull requests with over 50 comments
Determine the areas in which pull requests have sparked significant discussion, by identifying those with over 50 comments. This can provide insights into contentious or complex issues within your organization's GitHub repositories.

```sql+postgres
select
  title,
  id,
  comments,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot comments:>50';
```

```sql+sqlite
select
  title,
  id,
  comments,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot comments:>50';
```

### List open draft pull requests
Explore the open draft pull requests in your GitHub organization. This could be used to identify unfinished work and help prioritize tasks for your development team.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot draft:true state:open';
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot draft:true state:open';
```

### List pull requests that took more than 30 days to close
Determine the areas in which pull requests have taken more than a month to close. This can help in identifying bottlenecks in the code review process and provide insights for improving efficiency.

```sql+postgres
select
  title,
  id,
  state,
  created_at,
  closed_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot state:closed'
  and closed_at > (created_at + interval '30' day);
```

```sql+sqlite
select
  title,
  id,
  state,
  created_at,
  closed_at,
  url
from
  github_search_pull_request
where
  query = 'org:turbot state:closed'
  and closed_at > datetime(created_at, '+30 days');
```