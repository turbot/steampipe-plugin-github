![image](https://hub.steampipe.io/images/plugins/turbot/github-social-graphic.png)

# GitHub Plugin for Steampipe

* **[Get started →](https://hub.steampipe.io/plugins/turbot/github)**
* Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/github/tables)
* Community: [Slack Channel](https://steampipe.io/community/join)
* Get involved: [Issues](https://github.com/turbot/steampipe-plugin-github/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install github
```

Set your GitHub API key:

```shell
export GITHUB_TOKEN=YOURTOKENHERE
```

Launch the Steampipe REPL:

```shell
steampipe query
```

Run a query:

```sql
select
  name,
  language,
  forks_count,
  stargazers_count
from
  github_my_repository;
```

## Developing

Prerequisites:

* [Steampipe](https://steampipe.io/downloads)
* [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-github.git
cd steampipe-plugin-github
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```sh
make
```

Configure the plugin:

```sh
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/github.spc
```

Try it!

```shell
steampipe query
> .inspect github
```

Further reading:

* [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
* [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-github/blob/main/LICENSE).

`help wanted` issues:

* [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
* [GitHub Plugin](https://github.com/turbot/steampipe-plugin-github/labels/help%20wanted)
