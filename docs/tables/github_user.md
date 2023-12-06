---
title: "Steampipe Table: github_user - Query GitHub Users using SQL"
description: "Allows users to query GitHub Users, specifically the user profiles, providing insights into user activities and contributions."
---

# Table: github_user - Query GitHub Users using SQL

GitHub Users is a resource within GitHub that represents an individual user account on GitHub. It provides information about the user's public profile, such as the user's username, bio, location, email, and the date the user joined GitHub. GitHub Users also include statistics about the user's activities and contributions to GitHub repositories.

## Table Usage Guide

The `github_user` table provides insights into individual user accounts within GitHub. As a developer or project manager, explore user-specific details through this table, including user profiles, activities, and contributions. Utilize it to uncover information about users, such as their public profiles, the date they joined GitHub, and their contributions to repositories.

**Important Notes**
- You must specify the `login` column in `where` or `join` clause to query the table.

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

### List of users in your organizations
This example helps you understand who the members of your organization are on Github. It provides insights into their profile details including their name, company, location, and Twitter handle, which can be useful for networking or organizational mapping purposes.

```sql+postgres
select
  u.login,
  o.login as organization,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_organization as o,
  jsonb_array_elements_text(o.member_logins) as member_login
where
  u.login = member_login;
```

```sql+sqlite
select
  u.login,
  o.login as organization,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_organization as o,
  json_each(o.member_logins) as member_login
where
  u.login = member_login.value;
```

### List of users that collaborate on a repository that you own
This query is utilized to identify the collaborators on a specific repository that you manage. It's particularly useful for understanding the team composition, including their affiliations and locations, providing you with a comprehensive view of who is contributing to your project.

```sql+postgres
select
  r.full_name as repository,
  u.login,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_repository as r,
  jsonb_array_elements_text(r.collaborator_logins) as collaborator_login
where
  u.login = collaborator_login
  and r.full_name = 'turbot/steampipe';
```

```sql+sqlite
select
  r.full_name as repository,
  u.login,
  u.name,
  u.company,
  u.location,
  u.twitter_username,
  u.bio
from
  github_user as u,
  github_my_repository as r,
  json_each(r.collaborator_logins) as collaborator_login
where
  u.login = collaborator_login.value
  and r.full_name = 'turbot/steampipe';
```