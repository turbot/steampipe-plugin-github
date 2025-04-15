---
title: "Steampipe Table: github_code_owner - Query GitHub Code Owners using SQL"
description: "Allows users to query GitHub Code Owners, providing insights into the users or teams that own certain code in a repository."
folder: "Repository"
---

# Table: github_code_owner - Query GitHub Code Owners using SQL

GitHub Code Owners is a feature within GitHub that allows repository administrators to define specific maintainers for different parts of a codebase. It works by assigning a 'code owner' to specific files or directories in a repository, and these owners are then automatically requested for review when someone opens a pull request that modifies code they own. This feature helps streamline the code review process and ensures that the right people are involved in reviewing changes.

## Table Usage Guide

The `github_code_owner` table provides insights into the code owners within a GitHub repository. As a repository administrator, you can use this table to understand who owns which parts of your codebase and to ensure that all changes are reviewed by the appropriate individuals. This table can also be helpful for developers looking to understand who to contact for specific parts of a codebase, or for project managers looking to understand the distribution of code ownership within a project.

**Important Notes**
- You must specify the `repository_full_name` column in `where` or `join` clause to query the table.

## Examples

### List code owners from rules
Determine the ownership of different sections of code within the 'github/docs' repository. This can help identify who is responsible for certain code patterns, allowing for more efficient collaboration and issue resolution.
For instance, for a given [CODEOWNERS](https://github.com/github/docs/blob/main/.github/CODEOWNERS) file from the [GitHub Docs](https://github.com/github/docs) repository:

```
# Order is important. The LAST matching pattern has the MOST precedence.
# gitignore style patterns are used, not globs.
# https://docs.github.com/articles/about-codeowners
# https://git-scm.com/docs/gitignore

# Engineering
*.js @github/docs-engineering
*.ts @github/docs-engineering
*.tsx @github/docs-engineering
/.github/ @github/docs-engineering
/script/ @github/docs-engineering
/includes/ @github/docs-engineering
/lib/search/popular-pages.json @github/docs-engineering
Dockerfile @github/docs-engineering
package-lock.json @github/docs-engineering
package.json @github/docs-engineering

# Localization
/.github/actions-scripts/msft-create-translation-batch-pr.js @github/docs-engineering
/.github/workflows/msft-create-translation-batch-pr.yml @github/docs-engineering
/translations/ @Octomerger

# Site Policy
/content/site-policy/ @github/site-policy-admins

# Content strategy
/contributing/content-markup-reference.md @github/docs-content-strategy
/contributing/content-style-guide.md @github/docs-content-strategy
/contributing/content-model.md @github/docs-content-strategy
/contributing/content-style-guide.md @github/docs-content-strategy
/contributing/content-templates.md @github/docs-content-strategy

# Requires review of #actions-oidc-integration, docs-engineering/issues/1506
content/actions/deployment/security-hardening-your-deployments/** @github/oidc
```

You can query all of the rules with the following query:

```sql+postgres
select
  line,
  pattern,
  users,
  teams,
  pre_comments,
  line_comment,
  repository_full_name
from
  github_code_owner
where
  repository_full_name = 'github/docs'
order by
  line asc;
```

```sql+sqlite
select
  line,
  pattern,
  users,
  teams,
  pre_comments,
  line_comment,
  repository_full_name
from
  github_code_owner
where
  repository_full_name = 'github/docs'
order by
  line asc;
```

```
+------+-------------------------------------------------------------------+-----------------+-----------------------------------+----------------------------------------------------------------------------------+--------------+----------------------+
| line | pattern                                                           | users           | teams                             | pre_comments                                                                     | line_comment | repository_full_name |
+------+-------------------------------------------------------------------+-----------------+-----------------------------------+----------------------------------------------------------------------------------+--------------+----------------------+
| 7    | *.js                                                              | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 8    | *.ts                                                              | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 9    | *.tsx                                                             | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 10   | /.github/                                                         | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 11   | /script/                                                          | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 12   | /includes/                                                        | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 13   | /lib/search/popular-pages.json                                    | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 14   | Dockerfile                                                        | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 15   | package-lock.json                                                 | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 16   | package.json                                                      | <null>          | ["@github/docs-engineering"]      | ["# Engineering"]                                                                |              | github/docs          |
| 19   | /.github/actions-scripts/msft-create-translation-batch-pr.js      | <null>          | ["@github/docs-engineering"]      | ["# Localization"]                                                               |              | github/docs          |
| 20   | /.github/workflows/msft-create-translation-batch-pr.yml           | <null>          | ["@github/docs-engineering"]      | ["# Localization"]                                                               |              | github/docs          |
| 21   | /translations/                                                    | ["@Octomerger"] | <null>                            | ["# Localization"]                                                               |              | github/docs          |
| 24   | /content/site-policy/                                             | <null>          | ["@github/site-policy-admins"]    | ["# Site Policy"]                                                                |              | github/docs          |
| 27   | /contributing/content-markup-reference.md                         | <null>          | ["@github/docs-content-strategy"] | ["# Content strategy"]                                                           |              | github/docs          |
| 28   | /contributing/content-style-guide.md                              | <null>          | ["@github/docs-content-strategy"] | ["# Content strategy"]                                                           |              | github/docs          |
| 29   | /contributing/content-model.md                                    | <null>          | ["@github/docs-content-strategy"] | ["# Content strategy"]                                                           |              | github/docs          |
| 30   | /contributing/content-style-guide.md                              | <null>          | ["@github/docs-content-strategy"] | ["# Content strategy"]                                                           |              | github/docs          |
| 31   | /contributing/content-templates.md                                | <null>          | ["@github/docs-content-strategy"] | ["# Content strategy"]                                                           |              | github/docs          |
| 34   | content/actions/deployment/security-hardening-your-deployments/** | <null>          | ["@github/oidc"]                  | ["# Requires review of #actions-oidc-integration, docs-engineering/issues/1506"] |              | github/docs          |
+------+-------------------------------------------------------------------+-----------------+-----------------------------------+----------------------------------------------------------------------------------+--------------+----------------------+
```