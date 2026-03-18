---
title: "Steampipe Table: github_user_contribution_stats - Query GitHub user contributions summary using SQL"
description: "Query GitHub user contribution summaries and calendar data from the GraphQL ContributionsCollection."
folder: "User"
---

## Table: github_user_contribution_stats - Query GitHub user contributions summary using SQL

The `github_user_contribution_stats` table provides access to GitHub's ContributionsCollection data for a user, including total contribution counts and the contribution calendar (weeks/days). This makes it possible to build dashboards and reports similar to a user's public contribution graph.

## Table Usage Guide

The table is scoped to a single user per query. Optionally specify `from_date` and `to_date` to constrain the contribution window, and `max_repositories` to control how many repositories are returned for commit contributions by repository.

## Important Notes

- You must specify the `login` column in the `where` clause.
- The `commit_contributions_by_repository` field returns at most 100 repositories (default 100).

## Examples

### Get contribution summary for a user

```sql+postgres
select
  total_commit_contributions,
  total_issue_contributions,
  total_pull_request_contributions,
  total_pull_request_review_contributions,
  total_repositories_with_contributed_commits
from
  github_user_contribution_stats
where
  login = 'octocat';
```

```sql+sqlite
select
  total_commit_contributions,
  total_issue_contributions,
  total_pull_request_contributions,
  total_pull_request_review_contributions,
  total_repositories_with_contributed_commits
from
  github_user_contribution_stats
where
  login = 'octocat';
```

### Get contribution calendar for a date range

```sql+postgres
select
  contribution_calendar
from
  github_user_contribution_stats
where
  login = 'octocat'
  and from_date = '2025-01-01'
  and to_date = '2025-12-31';
```

```sql+sqlite
select
  contribution_calendar
from
  github_user_contribution_stats
where
  login = 'octocat'
  and from_date = '2025-01-01'
  and to_date = '2025-12-31';
```

### Limit repositories in commit contributions breakdown

```sql+postgres
select
  commit_contributions_by_repository
from
  github_user_contribution_stats
where
  login = 'octocat'
  and max_repositories = 100;
```

```sql+sqlite
select
  commit_contributions_by_repository
from
  github_user_contribution_stats
where
  login = 'octocat'
  and max_repositories = 100;
```
