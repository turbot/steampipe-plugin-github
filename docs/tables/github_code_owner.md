# Table: github_code_owner

All of your code owners rules defined in your repository, in CODEOWNERS files.

The `github_code_owner` table can be used to query information about **ANY** repository, and **you must specify which repository** in the where or join clause (`where repository_full_name=`, `join github_code_owner on repository_full_name=`).

## Examples

### Get All your CodeOwners rules about a specific repository

```sql
select
  line,
  repository_full_name,
  pattern,
  users,
  teams,
  pre_comments,
  line_comment 
from
  github_code_owner 
where
  repository_full_name = 'github/docs' 
order by
  line asc
```

**CODEOWNERS file content**

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

**ResultSet**

| line | repository\_full\_name | pattern                                                             | users             | teams                               | pre\_comments                                                                      | line\_comment |
| ---- | ---------------------- | ------------------------------------------------------------------- | ----------------- | ----------------------------------- | ---------------------------------------------------------------------------------- | ------------- |
| 7    | github/docs            | \*.js                                                               | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 8    | github/docs            | \*.ts                                                               | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 9    | github/docs            | \*.tsx                                                              | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 10   | github/docs            | /.github/                                                           | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 11   | github/docs            | /script/                                                            | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 12   | github/docs            | /includes/                                                          | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 13   | github/docs            | /lib/search/popular-pages.json                                      | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 14   | github/docs            | Dockerfile                                                          | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 15   | github/docs            | package-lock.json                                                   | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 16   | github/docs            | package.json                                                        | <null>            | \["@github/docs-engineering"\]      | \["# Engineering"\]                                                                |               |
| 19   | github/docs            | /.github/actions-scripts/msft-create-translation-batch-pr.js        | <null>            | \["@github/docs-engineering"\]      | \["# Localization"\]                                                               |               |
| 20   | github/docs            | /.github/workflows/msft-create-translation-batch-pr.yml             | <null>            | \["@github/docs-engineering"\]      | \["# Localization"\]                                                               |               |
| 21   | github/docs            | /translations/                                                      | \["@Octomerger"\] | <null>                              | \["# Localization"\]                                                               |               |
| 24   | github/docs            | /content/site-policy/                                               | <null>            | \["@github/site-policy-admins"\]    | \["# Site Policy"\]                                                                |               |
| 27   | github/docs            | /contributing/content-markup-reference.md                           | <null>            | \["@github/docs-content-strategy"\] | \["# Content strategy"\]                                                           |               |
| 28   | github/docs            | /contributing/content-style-guide.md                                | <null>            | \["@github/docs-content-strategy"\] | \["# Content strategy"\]                                                           |               |
| 29   | github/docs            | /contributing/content-model.md                                      | <null>            | \["@github/docs-content-strategy"\] | \["# Content strategy"\]                                                           |               |
| 30   | github/docs            | /contributing/content-style-guide.md                                | <null>            | \["@github/docs-content-strategy"\] | \["# Content strategy"\]                                                           |               |
| 31   | github/docs            | /contributing/content-templates.md                                  | <null>            | \["@github/docs-content-strategy"\] | \["# Content strategy"\]                                                           |               |
| 34   | github/docs            | content/actions/deployment/security-hardening-your-deployments/\*\* | <null>            | \["@github/oidc"\]                  | \["# Requires review of #actions-oidc-integration, docs-engineering/issues/1506"\] |