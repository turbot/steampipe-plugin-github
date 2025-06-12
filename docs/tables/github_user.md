---
title: "Steampipe Table: github_user - Query GitHub Users using SQL"
description: "Allows users to query GitHub Users, specifically the user profiles, providing insights into user activities and contributions."
folder: "User"
---

# Table: github_user - Query GitHub Users using SQL

GitHub Users is a resource within GitHub that represents an individual user account on GitHub. It provides information about the user's public profile, such as the user's username, bio, location, email, and the date the user joined GitHub. GitHub Users also include statistics about the user's activities and contributions to GitHub repositories.

## Table Usage Guide

The `github_user` table provides insights into individual user accounts within GitHub. As a developer or project manager, explore user-specific details through this table, including user profiles, activities, and contributions. Utilize it to uncover information about users, such as their public profiles, the date they joined GitHub, and their contributions to repositories.

**Important Notes**
- You must specify the `login` column in `where` or `join` clause to query the table.

If using a [fine-grained access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token), no permissions are required.

## Examples

### Get information for a user
Explore the details associated with a specific GitHub user to gain insights into their activity and profile. This can be useful for understanding their contributions or for investigating user behavior.

```sql+postgres
select
  *
from
  github_user
where
  login = 'torvalds';
```

```sql+sqlite
select
  *
from
  github_user
where
  login = 'torvalds';
```

### List users that are members of multiple organizations

```sql+postgres
select
  name,
  email,
  created_at,
  bio,
  twitter_username,
  organizations_total_count
from
  github_user
where
  login = 'madhushreeray30'
  and organizations_total_count > 1;
```

```sql+sqlite
select
  name,
  email,
  created_at,
  bio,
  twitter_username,
  organizations_total_count
from
  github_user
where
  login = 'madhushreeray30'
  and organizations_total_count > 1;
```