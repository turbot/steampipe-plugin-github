## v0.7.0 [2021-10-12]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v162--2021-10-08) ([#68](https://github.com/turbot/steampipe-plugin-github/pull/68))
- Added the column `files` to `github_my_gist` and `github_gist` tables ([#66](https://github.com/turbot/steampipe-plugin-github/pull/66))
- Reordered the permission scope in the docs/index.md file to match the list in GitHub personal access token page ([#67](https://github.com/turbot/steampipe-plugin-github/pull/67))

## v0.6.1 [2021-09-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v161--2021-09-21) ([#63](https://github.com/turbot/steampipe-plugin-github/pull/63))
- Updated several tables to use the helper function for parsing the `repository_full_name` column ([#55](https://github.com/turbot/steampipe-plugin-github/pull/55))

_Bug fixes_

- Removed all the incorrect references in the documentation (([#58](https://github.com/turbot/steampipe-plugin-github/pull/58))) ([#59](https://github.com/turbot/steampipe-plugin-github/pull/59))
- Minor cleanup in github_my_team table document

## v0.6.0 [2021-09-09]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v150--2021-08-06)
- `PublishedAt` column of `github_release` table will now return `nil` if it does not contain any value
- `github_my_issue` table now uses `RetryHydrate` function from Steampipe SDK for enhanced retry logic ([#19](https://github.com/turbot/steampipe-plugin-github/pull/19))
- The document is now updated with the latest format of the GitHub personal access token ([#47](https://github.com/turbot/steampipe-plugin-github/pull/47))

## v0.5.1 [2021-06-05]

_Bug fixes_

- Fixed: Incorrect reference in documentation.

## v0.5.0 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [github_branch_protection](https://hub.steampipe.io/plugins/turbot/github/tables/github_branch_protection) ([#37](https://github.com/turbot/steampipe-plugin-github/issues/37))
  - [github_community_profile](https://hub.steampipe.io/plugins/turbot/github/tables/github_community_profile) ([#40](https://github.com/turbot/steampipe-plugin-github/issues/40))

_Enhancements_

- Updated: Add `outside_collaborators` and `outside_collaborator_logins` columns to `github_repository` and `github_my_repository` tables ([#38](https://github.com/turbot/steampipe-plugin-github/issues/38))

_Bug fixes_

- Fixed: Remove duplicate column `members_can_create_repositories` from `github_organization` and `github_my_organization` tables ([#36](https://github.com/turbot/steampipe-plugin-github/pull/36))

## v0.4.0 [2021-05-15]

_What's new?_

- New tables added
  - [github_stargazer](https://hub.steampipe.io/plugins/turbot/github/tables/github_stargazer) ([#29](https://github.com/turbot/steampipe-plugin-github/pull/29))
  - [github_tag](https://hub.steampipe.io/plugins/turbot/github/tables/github_tag) ([#30](https://github.com/turbot/steampipe-plugin-github/pull/30))
  - [github_branch](https://hub.steampipe.io/plugins/turbot/github/tables/github_branch) ([#31](https://github.com/turbot/steampipe-plugin-github/pull/31))
  - [github_traffic_view_daily](https://hub.steampipe.io/plugins/turbot/github/tables/github_traffic_view_daily) ([#32](https://github.com/turbot/steampipe-plugin-github/pull/32))
  - [github_traffic_view_weekly](https://hub.steampipe.io/plugins/turbot/github/tables/github_traffic_view_weekly) ([#33](https://github.com/turbot/steampipe-plugin-github/pull/33))


## v0.3.0 [2021-04-30]

_What's new?_

- New tables added
  - [github_commit](https://hub.steampipe.io/plugins/turbot/github/tables/github_commit) ([#22](https://github.com/turbot/steampipe-plugin-github/pull/22))
  - [github_gitignore](https://hub.steampipe.io/plugins/turbot/github/tables/github_gitignore) ([#23](https://github.com/turbot/steampipe-plugin-github/pull/23))
  - [github_release](https://hub.steampipe.io/plugins/turbot/github/tables/github_release) ([#20](https://github.com/turbot/steampipe-plugin-github/pull/20))
  - [github_workflow](https://hub.steampipe.io/plugins/turbot/github/tables/github_workflow) ([#25](https://github.com/turbot/steampipe-plugin-github/pull/25))

_Enhancements_

- Use go v1.16 ([#27](https://github.com/turbot/steampipe-plugin-github/pull/27))

_Bug fixes_

- Cleanup unnecessary logging in github_license ([#24](https://github.com/turbot/steampipe-plugin-github/pull/24))
- Github (lower h) references should be GitHub (capital H) throughout the docs etc ([#26](https://github.com/turbot/steampipe-plugin-github/pull/26))


## v0.2.0 [2021-03-18]

_What's new?_

- New tables added
  - [github_my_gist](https://hub.steampipe.io/plugins/turbot/github/tables/github_my_gist) ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))
  - [github_my_issue](https://hub.steampipe.io/plugins/turbot/github/tables/github_my_issue) ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))
  - [github_my_organization](https://hub.steampipe.io/plugins/turbot/github/tables/github_my_organization) ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))
  - [github_my_repository](https://hub.steampipe.io/plugins/turbot/github/tables/github_my_repository) ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))
  - [github_pull_request](https://hub.steampipe.io/plugins/turbot/github/tables/github_pull_request) ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.4](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v024-2021-03-16)

_Bug fixes_

- Fixed: Renamed table `github_repository_issue` to `github_issue` ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))
- Fixed: Renamed table `github_team` to `github_my_team` ([#16](https://github.com/turbot/steampipe-plugin-github/pull/16))


## v0.1.1 [2021-02-25]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve SDK issues:
  - Fix error for missing required quals [#40](https://github.com/turbot/steampipe-plugin-sdk/issues/42).
  - Queries fail with error socket: too many open files [#190](https://github.com/turbot/steampipe/issues/190)


## v0.1.0 [2021-02-18]

_What's new?_

- Added support for [connection configuration](https://github.com/turbot/steampipe-plugin-github/blob/main/docs/index.md#connection-configuration). You may specify github `token` for each connection in a configuration file.


## v0.0.5 [2021-01-28]

_What's new?_

- Added: `github_repository_issue` table
