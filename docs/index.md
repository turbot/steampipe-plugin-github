---
org: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/github.svg"
brand_color: "#191717"
display_name: "GitHub"
name: "github"
description: "Steampipe plugin for querying GitHub Repositories, Organizations, and other resources."
---

# GitHub

The GitHub plugin is used to interact with the many resources in the GitHub API.

## Installation

To download and install the latest github plugin:

```bash
steampipe plugin install github
```


## Credentials

The GitHub plugin uses a personal access token to authenticate to the GitHub APIs.  You must create a [Personal Access Token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and assign the following scopes:
    - `repo` (all)
    - `read:org`
    - `read:user`
    - `gist`
    - `user:email`


## Connection Configuration
Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from ~/.steampipe/config that have a .spc extension. A config file may contain multiple connections.

Installing the latest github plugin will create a connection file (`~/.steampipe/config/github.spc`) with a single default connection named `github`. You must edit this connection to include your API `token`:

  ```hcl
  connection "github" {
    plugin = "github"
    token  = "111222333444555666777888999aaabbbcccddde"
  }
  ```

### Scope

The GitHub plugin query scope is generally the same as the GitHub API. Usually, this means you can list private resources that you have access to, as well as public resources that you own, or that are owned by organizations to which you belong. The same GitHub APIs are used to get information for public resources, but the public items are returned in list calls (because there would be too many). This has an interesting side effect in Steampipe in that you can sometimes query _a specific item_ by _a specific key column or columns_ that does not show up in a list query.

For example, `select * from github_organization` will list details about all the GitHub Organizations to which you belong. `select * from github_organization where login = 'postgres'` will show you the publicly available details about the `postgres` organization, which didn't show up in your first query! It works this way in Steampipe because [that's how it works in the API](https://docs.github.com/en/rest/reference/orgs#list-organizations-for-a-user). While this may seem counter-intuitive at first, it actually can be quite useful.

