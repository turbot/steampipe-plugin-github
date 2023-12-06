## v0.38.0 [2023-12-06]

_What's new?_

- New tables added
  - [github_repository_sbom](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_sbom) ([#353](https://github.com/turbot/steampipe-plugin-github/pull/353)) (Thanks [@lwakefield](https://github.com/lwakefield) for the contribution!)

_Enhancements_

- Updated the following tables to include support for dynamic GraphQL queries:
  - `github_my_star` ([#369](https://github.com/turbot/steampipe-plugin-github/pull/369))
  - `github_stargazer` ([#370](https://github.com/turbot/steampipe-plugin-github/pull/370))
  - `github_tag` ([#371](https://github.com/turbot/steampipe-plugin-github/pull/371))
  - `github_rate_limit` ([#368](https://github.com/turbot/steampipe-plugin-github/pull/368))
  - `github_community_profile` ([#367](https://github.com/turbot/steampipe-plugin-github/pull/367))
  - `github_license` ([#366](https://github.com/turbot/steampipe-plugin-github/pull/366))
  - `github_organization_member` ([#364](https://github.com/turbot/steampipe-plugin-github/pull/364))
  - `github_team_member` ([#364](https://github.com/turbot/steampipe-plugin-github/pull/364))
  - `github_user` ([#364](https://github.com/turbot/steampipe-plugin-github/pull/364))
  - `github_my_team` ([#363](https://github.com/turbot/steampipe-plugin-github/pull/363))
  - `github_team` ([#363](https://github.com/turbot/steampipe-plugin-github/pull/363))
  - `github_commit` ([#362](https://github.com/turbot/steampipe-plugin-github/pull/362))
  - `github_my_organization` ([#361](https://github.com/turbot/steampipe-plugin-github/pull/361))
  - `github_organization` ([#361](https://github.com/turbot/steampipe-plugin-github/pull/361))
  - `github_organization_external_identity` ([#361](https://github.com/turbot/steampipe-plugin-github/pull/361))
  - `github_branch` ([#360](https://github.com/turbot/steampipe-plugin-github/pull/360))
  - `github_branch_protection` ([#360](https://github.com/turbot/steampipe-plugin-github/pull/360))
  - `github_repository_collaborator` ([#365](https://github.com/turbot/steampipe-plugin-github/pull/365))
  - `github_repository_deployment` ([#365](https://github.com/turbot/steampipe-plugin-github/pull/365))
  - `github_repository_environment` ([#365](https://github.com/turbot/steampipe-plugin-github/pull/365))
  - `github_repository_vulnerability_alert` ([#365](https://github.com/turbot/steampipe-plugin-github/pull/365))
  - `github_issue` ([#359](https://github.com/turbot/steampipe-plugin-github/pull/359))
  - `github_issue_comment` ([#359](https://github.com/turbot/steampipe-plugin-github/pull/359))
  - `github_pull_request` ([#359](https://github.com/turbot/steampipe-plugin-github/pull/359))
  - `github_pull_request_comment` ([#359](https://github.com/turbot/steampipe-plugin-github/pull/359))
  - `github_pull_request_review` ([#359](https://github.com/turbot/steampipe-plugin-github/pull/359))

## v0.37.1 [2023-11-16]

_Bug fixes_

- Fixed the `GetConfig` of `github_team_repository` table to include support for dynamic GraphQL queries. ([#379](https://github.com/turbot/steampipe-plugin-github/pull/379))
- Fixed the example queries in `github_commit` doc file. ([#377](https://github.com/turbot/steampipe-plugin-github/pull/377))
- Fixed the example queries in `github_search_issue` doc file to filter out results from the API. ([#378](https://github.com/turbot/steampipe-plugin-github/pull/378))

## v0.37.0 [2023-11-10]

_Enhancements_

- Added the `run_started_at` column to `github_actions_repository_workflow_run` table. ([#358](https://github.com/turbot/steampipe-plugin-github/pull/358)) (Thanks [@mridang](https://github.com/mridang) for the contribution!)

## v0.36.1 [2023-10-27]

_Bug fixes_

- Fixed the required quals of `github_issue` and `github_pull_request` tables to correctly return data instead of an error. ([#355](https://github.com/turbot/steampipe-plugin-github/pull/355))

## v0.36.0 [2023-10-24]

_What's new_

- Updated `github_issue`, `github_my_issue`, `github_pull_request`, `github_search_issue`, and `github_search_pull_request` tables to only include nested and user permission columns in GraphQL request when requested. This should result in faster queries and large scale queries completing more consistently. ([#342](https://github.com/turbot/steampipe-plugin-github/pull/342))

## v0.35.1 [2023-10-04]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#346](https://github.com/turbot/steampipe-plugin-github/pull/346))

## v0.35.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#341](https://github.com/turbot/steampipe-plugin-github/pull/341))

## v0.34.1 [2023-09-21]

_Bug fixes_

- Empty values will no longer be cached incorrectly for the `github_my_repository`, `github_repository`, and `github_search_repository` tables. ([#340](https://github.com/turbot/steampipe-plugin-github/pull/340))
- Fixed `github_team_repository table` to include support for dynamic GraphQL queries. ([#339](https://github.com/turbot/steampipe-plugin-github/pull/339))

## v0.34.0 [2023-09-20]

_What's new_

- Updated `github_my_repository`, `github_repository`, and `github_search_repository` tables to only include requested columns in GraphQL request. This should result in faster queries and large scale queries completing more consistently. ([#338](https://github.com/turbot/steampipe-plugin-github/pull/338))

_Dependencies_

- Recompiled plugin with Go 1.21. ([#338](https://github.com/turbot/steampipe-plugin-github/pull/338))

## v0.33.1 [2023-09-19]

_Bug fixes_

- Fixed `github_search_repository` table queries failing when selecting the `has_downloads`, `has_pages`, `hooks`, `network_count`, `subscribers_count`, or `topics` columns. ([#337](https://github.com/turbot/steampipe-plugin-github/pull/337))

## v0.33.0 [2023-09-13]

_Breaking changes_

- Removed the `security_advisory_cwes_cweid` and `security_advisory_cwes_name` columns from `github_organization_dependabot_alert` and `github_repository_dependabot_alert` tables. ([#332](https://github.com/turbot/steampipe-plugin-github/pull/332))

_Enhancements_

- Added the `security_advisory_cwes` column to `github_organization_dependabot_alert` and `github_repository_dependabot_alert` tables. ([#332](https://github.com/turbot/steampipe-plugin-github/pull/332))
- Added the `actor`, `actor_login`, `triggering_actor`, and `triggering_actor_login` columns to `github_actions_repository_workflow_run` table. ([#332](https://github.com/turbot/steampipe-plugin-github/pull/332))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v551-2023-07-26). ([#325](https://github.com/turbot/steampipe-plugin-github/pull/325))
- Recompiled plugin with `golang.org/x/oauth2 v0.12.0`. ([#326](https://github.com/turbot/steampipe-plugin-github/pull/326))
- Recompiled plugin with `Github.com/argonsecurity/pipeline-parser v0.3.3`. ([#330](https://github.com/turbot/steampipe-plugin-github/pull/330))
- Recompiled plugin with `github.com/turbot/go-kit v0.7.0`. ([#328](https://github.com/turbot/steampipe-plugin-github/pull/328))
- Recompiled plugin with `github.com/google/go-github v55.0.0`. ([#332](https://github.com/turbot/steampipe-plugin-github/pull/332))

## v0.32.0 [2023-09-07]

_Breaking changes_

- Removed the `temp_clone_token` column from `github_my_repository` and `github_repository` tables to allow queries with fine-grained access tokens. ([#321](https://github.com/turbot/steampipe-plugin-github/pull/321))

_What's new?_

- New tables added
  - [github_repository_vulnerability_alert](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_vulnerability_alert) ([#318](https://github.com/turbot/steampipe-plugin-github/pull/318))

_Enhancements_

- The plugin has been updated to use `DefaultRetryConfig` rather than `plugin.RetryHydrate` in each table. ([#322](https://github.com/turbot/steampipe-plugin-github/pull/322))

## v0.31.0 [2023-08-17]

_What's new?_

- New tables added
  - [github_pull_request_review](https://hub.steampipe.io/plugins/turbot/github/tables/github_pull_request_review) ([#313](https://github.com/turbot/steampipe-plugin-github/pull/313)) (Thanks [@jramosf](https://github.com/jramosf) for the contribution!)

## v0.30.2 [2023-07-28]

_Bug fixes_

- Fixed the `base_url` config argument to correctly parse the GitHub Enterprise graphql url to avoid queries returning `404` error. ([#307](https://github.com/turbot/steampipe-plugin-github/pull/307))

## v0.30.1 [2023-07-24]

_Bug fixes_

- Fixed the `labels` column of `github_issue` and `github_my_issue` tables to correctly return data instead of an error. ([#303](https://github.com/turbot/steampipe-plugin-github/pull/303))

## v0.30.0 [2023-07-20]

_Enhancements_

- Added the `labels_src` and `labels` columns to `github_issue`, `github_my_issue`, and `github_pull_request` tables. ([#299](https://github.com/turbot/steampipe-plugin-github/pull/299))

## v0.29.0 [2023-07-07]

_Breaking changes_

- Renamed column `user` to `user_login` in `github_audit_log` table to avoid incorrectly returning database username instead of the GitHub user, since it is a reserverd keyword. ([#289](https://github.com/turbot/steampipe-plugin-github/pull/289))
- Renamed column `user` to `user_detail` in `github_stargazer` table to avoid incorrectly returning database username instead of the GitHub user, since it is a reserverd keyword. ([#289](https://github.com/turbot/steampipe-plugin-github/pull/289))
- Removed column `outside_collaborators_total_count` column from `github_repository` table. The data is now available in `github_repository_collaborator` table. ([#292](https://github.com/turbot/steampipe-plugin-github/pull/292))

_What's new?_

- New tables added
  - [github_organization_external_identity](https://hub.steampipe.io/plugins/turbot/github/tables/github_organization_external_identity) ([#290](https://github.com/turbot/steampipe-plugin-github/pull/290))

_Enhancements_

- Added column `name` to `github_user` table. ([#294](https://github.com/turbot/steampipe-plugin-github/pull/294)) (Thanks [@ahlfors](https://github.com/ahlfors) for the contribution!)
- Added column `affiliation` to `github_repository_collaborator` table. ([#292](https://github.com/turbot/steampipe-plugin-github/pull/292))

_Bug fixes_

- Fixed the GraphQL API URLs when using base_url config argument.

## v0.28.1 [2023-06-21]

_Bug fixes_

- Fixed `saml_identity_provider` column errors in `github_my_organization` and `github_organization` tables for organizations with SAML SSO enabled. ([#287](https://github.com/turbot/steampipe-plugin-github/pull/287))

## v0.28.0 [2023-06-21]

_Breaking changes_

This release updates 25 tables to use [GitHub's GraphQL API](https://docs.github.com/en/graphql/overview/about-the-graphql-api) in an effort to optimize the number of outgoing API calls, increase query speed, and make the latest data available.

Due to the significant differences between the GitHub REST and GraphQL APIs, most updated tables have a large number of column breaking changes (removed or renamed columns, column data value changes).

Below is a list of updated tables along with all column changes:
- [github_branch](https://github.com/turbot/steampipe-plugin-github/issues/244)
- [github_branch_protection](https://github.com/turbot/steampipe-plugin-github/issues/251)
- [github_commit](https://github.com/turbot/steampipe-plugin-github/issues/262)
- [github_community_profile](https://github.com/turbot/steampipe-plugin-github/issues/266)
- [github_issue](https://github.com/turbot/steampipe-plugin-github/issues/260)
- [github_license](https://github.com/turbot/steampipe-plugin-github/issues/236)
- [github_my_issue](https://github.com/turbot/steampipe-plugin-github/issues/260)
- [github_my_organization](https://github.com/turbot/steampipe-plugin-github/issues/259)
- [github_my_repository](https://github.com/turbot/steampipe-plugin-github/issues/257)
- [github_my_star](https://github.com/turbot/steampipe-plugin-github/issues/270)
- [github_my_team](https://github.com/turbot/steampipe-plugin-github/issues/241)
- [github_organization](https://github.com/turbot/steampipe-plugin-github/issues/258)
- [github_organization_member](https://github.com/turbot/steampipe-plugin-github/issues/279)
- [github_pull_request](https://github.com/turbot/steampipe-plugin-github/issues/276)
- [github_repository](https://github.com/turbot/steampipe-plugin-github/issues/256)
- [github_search_issue](https://github.com/turbot/steampipe-plugin-github/issues/268)
- [github_search_pull_request](https://github.com/turbot/steampipe-plugin-github/issues/268)
- [github_search_repository](https://github.com/turbot/steampipe-plugin-github/issues/268)
- [github_search_user](https://github.com/turbot/steampipe-plugin-github/issues/268)
- [github_stargazer](https://github.com/turbot/steampipe-plugin-github/issues/270)
- [github_tag](https://github.com/turbot/steampipe-plugin-github/issues/238)
- [github_team](https://github.com/turbot/steampipe-plugin-github/issues/241)
- [github_team_member](https://github.com/turbot/steampipe-plugin-github/issues/243)
- [github_team_repository](https://github.com/turbot/steampipe-plugin-github/issues/254)
- [github_user](https://github.com/turbot/steampipe-plugin-github/issues/248)

_What's new_

- New tables added
  - [github_issue_comment](https://hub.steampipe.io/plugins/turbot/github/tables/github_issue_comment) ([#281](https://github.com/turbot/steampipe-plugin-github/pull/281))
  - [github_pull_request_comment](https://hub.steampipe.io/plugins/turbot/github/tables/github_pull_request_comment) ([#281](https://github.com/turbot/steampipe-plugin-github/pull/281))
  - [github_rate_limit_graphql](https://hub.steampipe.io/plugins/turbot/github/tables/github_rate_limit_graphql) ([#280](https://github.com/turbot/steampipe-plugin-github/pull/280))
  - [github_repository_collaborator](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_collaborator) ([#280](https://github.com/turbot/steampipe-plugin-github/pull/280))
  - [github_repository_deployment](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_deployment) ([#282](https://github.com/turbot/steampipe-plugin-github/pull/282))
  - [github_repository_environment](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_environment) ([#282](https://github.com/turbot/steampipe-plugin-github/pull/282))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v550-2023-06-16) which includes pending cache improvements that improve query speed and decrease the number of unnecessary API calls. ([#280](https://github.com/turbot/steampipe-plugin-github/pull/280))


## v0.27.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#250](https://github.com/turbot/steampipe-plugin-github/pull/250))

## v0.26.0 [2023-03-22]

_Enhancements_

- Added column `parent` to `github_team` and `github_my_team` tables. ([#232](https://github.com/turbot/steampipe-plugin-github/pull/232))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#233](https://github.com/turbot/steampipe-plugin-github/pull/233))

## v0.25.1 [2023-02-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#225](https://github.com/turbot/steampipe-plugin-github/pull/225))

## v0.25.0 [2022-12-27]

_What's new?_

- New tables added
  - [github_organization_dependabot_alert](https://hub.steampipe.io/plugins/turbot/github/tables/github_organization_dependabot_alert) ([#215](https://github.com/turbot/steampipe-plugin-github/pull/215)) (Thanks [@francois2metz](https://github.com/francois2metz) for the contribution!)
  - [github_repository_dependabot_alert](https://hub.steampipe.io/plugins/turbot/github/tables/github_repository_dependabot_alert) ([#215](https://github.com/turbot/steampipe-plugin-github/pull/215)) (Thanks [@francois2metz](https://github.com/francois2metz) for the contribution!)

_Enhancements_

- Added columns `base_ref` and `head_ref` to `github_pull_request` table. ([#223](https://github.com/turbot/steampipe-plugin-github/pull/223))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#219](https://github.com/turbot/steampipe-plugin-github/pull/219))

## v0.24.0 [2022-11-11]

_Breaking changes_

- Updated `github_organization_member` table to use GitHub GraphQL API. Based on the available fields, the table has been updated to include the following columns: ([#220](https://github.com/turbot/steampipe-plugin-github/pull/220))
  - `has_two_factor_enabled`
  - `login`
  - `organization`
  - `role`

## v0.23.0 [2022-11-10]

_Enhancements_

- Added `members_can_fork_private_repos` column to `github_my_organization` and `github_organization` tables. ([#214](https://github.com/turbot/steampipe-plugin-github/pull/214))

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
