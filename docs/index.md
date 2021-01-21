---
org: Turbot
category: ["software development"]
icon_url: "/images/plugins/turbot/github.svg"
brand_color: "#191717"
display_name: "GitHub"
name: "github"
description: "Steampipe plugin for querying Github Repositories, Organizations, and other resources."
---


# Github

The Github plugin is used to interact with the many resources in the Github API.

### Installation
To download and install the latest aws plugin:
```bash
steampipe plugin install github
```

### Configuration

The Github plugin uses a personal access token to authenticate to the Github API:
1. [Create a Personal Access Token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token).  You will need the following scopes:
    - `repo` (all)
    - `read:org`
    - `read:user`
    - `read:email`


2. Set the `GITHUB_TOKEN` environment variable to your access token:
```bash 
export GITHUB_TOKEN=111222333444555666777888999aaabbbcccddde
```
3. Run a query:

```
$ steampipe query
Welcome to Steampipe v0.0.14
Type ".inspect" for more information.
> select name, owner_login, language from github_repository
+------------------------------------+---------------+--------------+
|                name                |  owner_login  |   language   |
+------------------------------------+---------------+--------------+
| my_repo                            | me            | Go           |
| my_other_repo                      | me            | JavaScript   |
+------------------------------------+---------------+--------------+
```

### Scope
The Github plugin query scope is generally the same as the Github API.  Usually, this means you can list private resources that you have access to, as well as public resources that you own, or that are owned by organizations to which you belong.  The same Github APIs are used to get information for public resources, but the public items are returned in list calls (because there would be too many).  This has an interesting side effect in Steampipe in that you can sometimes query *a specific item* by *a specific key column or columns* that does not show up in a list query.

For example,  `select * from github_organization` will list details about all the Github Organizations to which you belong. `select * from github_organization where login = 'postgres'` will show you the publicly available details about the `postgres` organization, which didn't show up in your first query!  It works this way in Steampipe because [that's how it works in the API](https://docs.github.com/en/rest/reference/orgs#list-organizations-for-a-user).  While this may seem counter-intuitive at first, it actually can be quite useful.  