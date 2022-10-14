## v0.22.1 [2022-10-14]

_Bug fixes_

- Fixed formatting in `github_code_owner` table doc to prevent breaking the Hub site builds.

## v0.22.0 [2022-10-11]

_What's new?_

- New tables added
  - [github_code_owner](https://hub.steampipe.io/plugins/turbot/github/tables/github_code_owner) ([#200](https://github.com/turbot/steampipe-plugin-github/pull/200)) (Thanks to [@aminvielledebatAtBedrock](https://github.com/aminvielledebatAtBedrock) for adding the table!)

_Enhancements_

- Added retries for `RateLimitError` errors with a reset time of 60 seconds or less. ([#209](https://github.com/turbot/steampipe-plugin-github/pull/209))

## v0.21.0 [2022-09-30]

_What's new?_

- New tables added
  - [github_tree](https://hub.steampipe.io/plugins/turbot/github/tables/github_tree) ([#198](https://github.com/turbot/steampipe-plugin-github/pull/198)) (Thanks to [@asfaltboy](https://github.com/asfaltboy) for adding the table!)

_Enhancements_

- Updated error retry logic to retry on secondary rate limit errors.
- Improved general error retry logic to wait longer between requests, which should result in larger result sets returning more consistently.
- Queries will no retry on rate limit errors since the rate limit reset period is often more than 30 minutes.
- Improved `Quick start` instructions in README.md. ([#204](https://github.com/turbot/steampipe-plugin-github/pull/204)) (Thanks to [@breck7](https://github.com/breck7) for the contribution!)

## v0.20.0 [2022-09-09]

_Enhancements_

- Added `visibility` column as an optional list key column in `github_my_repository` table. ([#121](https://github.com/turbot/steampipe-plugin-github/pull/121))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#197](https://github.com/turbot/steampipe-plugin-github/pull/197))
- Recompiled plugin with Go version `1.19`. ([#195](https://github.com/turbot/steampipe-plugin-github/pull/195))

## v0.19.0 [2022-08-11]

_What's new?_

- New tables added
  - [github_organization_member](https://hub.steampipe.io/plugins/turbot/github/tables/github_organization_member) ([#177](https://github.com/turbot/steampipe-plugin-github/pull/177))

_Enhancements_

- Added column `security` to `github_community_profile` table. ([#180](https://github.com/turbot/steampipe-plugin-github/pull/180))
- Added column `hooks` to `github_my_organization` and `github_my_repository` tables. ([#185](https://github.com/turbot/steampipe-plugin-github/pull/185))
- Added columns `required_conversation_resolution` and `signatures_protected_branch` to `github_branch_protection` table. ([#178](https://github.com/turbot/steampipe-plugin-github/pull/178))
- Added columns `pipeline`, `workflow_file_content` and `workflow_file_content_json` to `github_workflow` table. ([#189](https://github.com/turbot/steampipe-plugin-github/pull/189))
- Added links to GitHub search docs in `github_search_*` table documents. ([#173](https://github.com/turbot/steampipe-plugin-github/pull/173))

_Bug fixes_

- Fixed the `github_commit` table to return an empty row instead of an error when trying to query for commits in an empty repository. ([#191](https://github.com/turbot/steampipe-plugin-github/pull/191))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#192](https://github.com/turbot/steampipe-plugin-github/pull/192))

## v0.18.0 [2022-07-07]

- New tables added (Thanks to [@japborst](https://github.com/japborst) for adding the table and [@ngowraj](https://github.com/ngowraj) for extensively testing it!)
  - [github_audit_log](https://hub.steampipe.io/plugins/turbot/github/tables/github_audit_log) ([#166](https://github.com/turbot/steampipe-plugin-github/pull/166))

## v0.17.1 [2022-06-28]

_Bug fixes_

- Fixed `github_branch_protection` table returning an error when querying unprotected branches.

## v0.17.0 [2022-06-24]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v330--2022-6-22). ([#168](https://github.com/turbot/steampipe-plugin-github/pull/168))

## v0.16.0 [2022-06-16]

_What's new?_

- New tables added (Thanks to [@japborst](https://github.com/japborst) for adding the tables below!)
  - [github_team](https://hub.steampipe.io/plugins/turbot/github/tables/github_team) ([#165](https://github.com/turbot/steampipe-plugin-github/pull/165))
  - [github_team_member](https://hub.steampipe.io/plugins/turbot/github/tables/github_team_member) ([#165](https://github.com/turbot/steampipe-plugin-github/pull/165))
  - [github_team_repository](https://hub.steampipe.io/plugins/turbot/github/tables/github_team_repository) ([#165](https://github.com/turbot/steampipe-plugin-github/pull/165))

_Enhancements_

- Recompiled plugin with [go-github v45.1.0](https://github.com/google/go-github/releases/tag/v45.1.0). ([#165](https://github.com/turbot/steampipe-plugin-github/pull/165))

## v0.15.0 [2022-05-25]

_What's new?_

- Added experimental GitHub Enterprise connection support. To get started, please set the `base_url` config argument in your GitHub connection. For more information, please see [GitHub plugin configuration](https://hub.steampipe.io/plugins/turbot/github#configuration). ([#160](https://github.com/turbot/steampipe-plugin-github/pull/160))

## v0.14.1 [2022-05-12]

_Enhancements_

- Updated `config/github.spc` and index doc with `token` argument environment variable information. ([#161](https://github.com/turbot/steampipe-plugin-github/pull/161))

## v0.14.0 [2022-04-28]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#157](https://github.com/turbot/steampipe-plugin-github/pull/157))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#156](https://github.com/turbot/steampipe-plugin-github/pull/156))

## v0.13.0 [2022-03-23]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#151](https://github.com/turbot/steampipe-plugin-github/pull/151))

## v0.12.0 [2022-02-09]

_What's new?_

- New tables added
  - [github_actions_artifact](https://hub.steampipe.io/plugins/turbot/github/tables/github_actions_artifact) ([#137](https://github.com/turbot/steampipe-plugin-github/pull/137))
  - [github_actions_repository_runner](https://hub.steampipe.io/plugins/turbot/github/tables/github_actions_repository_runner) ([#141](https://github.com/turbot/steampipe-plugin-github/pull/141))
  - [github_actions_repository_secret](https://hub.steampipe.io/plugins/turbot/github/tables/github_actions_repository_secret) ([#143](https://github.com/turbot/steampipe-plugin-github/pull/143))
  - [github_actions_repository_workflow_run](https://hub.steampipe.io/plugins/turbot/github/tables/github_actions_repository_workflow_run) ([#146](https://github.com/turbot/steampipe-plugin-github/pull/146))

## v0.11.1 [2022-02-02]

_Bug fixes_

- Fixed the `github_release` table to set `published_at` column to `nil` for draft releases ([#135](https://github.com/turbot/steampipe-plugin-github/pull/135))

## v0.11.0 [2022-01-12]

_Enhancements_

- Added `repository_full_name` column to `github_search_issue`, `github_search_code`, `github_search_commit`, `github_search_label` and `github_search_pull_request` tables ([#125](https://github.com/turbot/steampipe-plugin-github/pull/125)) ([#130](https://github.com/turbot/steampipe-plugin-github/pull/130))

_Bug fixes_

- Fixed the document of `github_search_commit` table ([#128](https://github.com/turbot/steampipe-plugin-github/pull/128))

## v0.10.1 [2022-01-05]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#120](https://github.com/turbot/steampipe-plugin-github/pull/120))

## v0.10.0 [2021-12-08]

_What's new?_

- New tables added
  - [github_search_code](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_code) ([#107](https://github.com/turbot/steampipe-plugin-github/pull/107))
  - [github_search_commit](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_commit) ([#108](https://github.com/turbot/steampipe-plugin-github/pull/108))
  - [github_search_issue](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_issue) ([#109](https://github.com/turbot/steampipe-plugin-github/pull/109))
  - [github_search_label](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_label) ([#115](https://github.com/turbot/steampipe-plugin-github/pull/115))
  - [github_search_pull_request](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_pull_request) ([#111](https://github.com/turbot/steampipe-plugin-github/pull/111))
  - [github_search_repository](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_repository) ([#114](https://github.com/turbot/steampipe-plugin-github/pull/114))
  - [github_search_topic](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_topic) ([#112](https://github.com/turbot/steampipe-plugin-github/pull/112))
  - [github_search_user](https://hub.steampipe.io/plugins/turbot/github/tables/github_search_user) ([#113](https://github.com/turbot/steampipe-plugin-github/pull/113))

## v0.9.0 [2021-11-23]

_What's new?_

- New tables added
  - [github_my_star](https://hub.steampipe.io/plugins/turbot/github/tables/github_my_star) ([#90](https://github.com/turbot/steampipe-plugin-github/pull/90))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#95](https://github.com/turbot/steampipe-plugin-github/pull/95))
- Updated all the tables which use custom retry logic to use RetryHydrate function from Steampipe plugin SDK ([#86](https://github.com/turbot/steampipe-plugin-github/pull/86))
- Added additional optional key quals, filter support and context cancellation handling across the following tables ([#61](https://github.com/turbot/steampipe-plugin-github/pull/61))
  - `github_branch`
  - `github_commit`
  - `github_issue`
  - `github_my_gist`
  - `github_my_issue`
  - `github_my_organization`
  - `github_my_repository`
  - `github_my_team`
  - `github_pull_request`
  - `github_release`
  - `github_stargazers`
  - `github_tag`
  - `github_workflow`
  - `github_branch`
  - `github_commit`
  - `github_issue`
  - `github_my_issue`
  - `github_my_repository`
  - `github_pull_request`

_Bug fixes_

- `github_repository` table will now return an empty row instead of `not found` error when the repository collaborator details are not available ([#89](https://github.com/turbot/steampipe-plugin-github/pull/89))

## v0.8.1 [2021-10-26]

_Bug fixes_

- Fixed the `member_logins` column of the `github_organization` table to correctly return the user logins that are members of the organization ([#81](https://github.com/turbot/steampipe-plugin-github/pull/81))

## v0.8.0 [2021-10-20]

_What's new?_

- New tables added
  - [github_rate_limit](https://hub.steampipe.io/plugins/turbot/github/tables/github_rate_limit) ([#74](https://github.com/turbot/steampipe-plugin-github/pull/74))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.7.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v170--2021-10-18) ([#77](https://github.com/turbot/steampipe-plugin-github/pull/77))

_Bug fixes_

- `github_organization` table will now return an empty row instead of error when `login` value for a GitHub Organization doesn't exist ([#75](https://github.com/turbot/steampipe-plugin-github/pull/75))

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
