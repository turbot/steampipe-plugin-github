---
organization: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/github.svg"
brand_color: "#191717"
display_name: "GitHub"
short_name: "github"
description: "Steampipe plugin for querying GitHub Repositories, Organizations, and other resources."
og_description: "Use SQL to query repos, orgs and more from GitHub. Open source CLI. No DB required."
og_image: "/images/plugins/turbot/github-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# GitHub + Steampipe

[Steampipe](https://steampipe.io) is an open-source zero-ETL engine to instantly query cloud APIs using SQL.

[GitHub](https://www.github.com/) is a provider of Internet hosting for software development and version control using Git. It offers the distributed version control and source code management (SCM) functionality of Git, plus its own features.

For example:

```sql
select
  name,
  primary_language -> 'name' as language,
  fork_count,
  stargazer_count
from
  github_my_repository;
```

```
+-------------------------------+------------+-------------+------------------+
| name                          | language   | forks_count | stargazer_count |
+-------------------------------+------------+-------------+------------------+
| steampipe                     | Go         | 11          | 254              |
| steampipe-plugin-aws          | Go         | 8           | 18               |
| steampipe-plugin-shodan       | Go         | 0           | 9                |
| steampipe-plugin-gcp          | Go         | 0           | 8                |
| steampipe-postgres-fdw        | C          | 0           | 8                |
| steampipe-plugin-azure        | Go         | 1           | 8                |
| steampipe-plugin-sdk          | Go         | 0           | 6                |
+-------------------------------+------------+-------------+------------------+
```

## Documentation

- **[Table definitions & examples →](https://hub.steampipe.io/plugins/turbot/github/tables)**

## Get started

### Install

Download and install the latest GitHub plugin:

```bash
steampipe plugin install github
```

### Credentials

| Item        | Description|
|-------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials | The GitHub plugin uses a personal access token to authenticate to the GitHub APIs.
| Permissions | You must create a [personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and assign the following scopes:<br />&nbsp;&nbsp;&nbsp;&nbsp;- `repo` (all)<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:org`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `gist`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `read:user`<br />&nbsp;&nbsp;&nbsp;&nbsp;- `user:email`
| Radius      | The GitHub plugin query scope is generally the same as the GitHub API. Usually, this means you can list private resources that you have access to, as well as public resources that you own, or that are owned by organizations to which you belong. The same GitHub APIs are used to get information for public resources, but the public items are returned in list calls (because there would be too many). This has an interesting side effect in Steampipe in that you can sometimes query _a specific item_ by _a specific key column or columns_ that does not show up in a list query.<br /><br />For example, `select * from github_my_organization` will list details about all the GitHub Organizations to which you belong. `select * from github_organization where login = 'postgres'` will show you the publicly available details about the `postgres` organization, which didn't show up in your first query! It works this way in Steampipe because [that's how it works in the API](https://docs.github.com/en/rest/reference/orgs#list-organizations-for-a-user). While this may seem counter-intuitive at first, it actually can be quite useful. |
| Resolution  | 1. Credentials in the Steampipe configuration file (`~/.steampipe/config/github.spc`) <br />2. Credentials specified in environment variables, e.g., `GITHUB_TOKEN`.

### Configuration

Installing the latest github plugin will create a config file (`~/.steampipe/config/github.spc`) with a single connection named `github`:

```hcl
connection "github" {
  plugin = "github"

  # The GitHub personal access token to authenticate to the GitHub APIs, e.g., `ghp_3b99b12218f63bcd702ad90d345975ef6c62f7d8`.
  # Please see https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token for more information.
  # Can also be set with the GITHUB_TOKEN environment variable.
  # token = "ghp_J1jzniKzVbFJNB34cJPwFPCmKeFakeToken"

  # GitHub Enterprise requires a base_url to be configured to your installation location.
  # Can also be set with the GITHUB_BASE_URL environment variable.
  # base_url = "https://github.example.com"
}
```

- `token` - [Personal access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) for your GitHub account. This can also be set via the `GITHUB_TOKEN` environment variable.
- `base_url` - GitHub Enterprise users have a custom URL location (e.g. `https://github.example.com`). Not required for GitHub cloud. This can also be via the `GITHUB_BASE_URL` environment variable.
