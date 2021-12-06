# Table: github_search_code

The `github_search_code` table helps to search for the specific item you want to find inside of a file. You can search globally across all of GitHub, or scope your search to a particular repository or organization.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. The `query` contains one or more search keywords and qualifiers. Qualifiers allow you to limit your search to specific areas of GitHub.

## Examples

### List searched codes by file name

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

### List searched codes by file location

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables';
```

### List searched codes by extension

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_rate_limit path:docs/tables extension:md';
```

### List searched codes within organization's repositories

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_stargazer org:turbot extension:go';
```

### List searched codes within a user's repositories

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers user:turbot extension:go';
```

### List searched codes within a user's specific repository

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers repo:turbot/steampipe-plugin-github extension:go';
```

### List searched codes by language

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'github_tag org:turbot language:markdown';
```

### List searched codes while file size is greater than 40 KB

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'org:turbot size:>40000 language:markdown';
```

### List searched codes by the file contents or file path

```sql
select
  name,
  query,
  html_url,
  sha
from
  github_search_code
where
  query = 'Stargazers org:turbot in:file,path extension:go';
```

### List text match details

```sql
select
  name,
  jsonb_pretty(match -> 'matches') as matches,
  match ->> 'fragment' as fragment,
  match ->> 'property' as property,
  match ->> 'object_url' as object_url,
  match ->> 'object_type' as object_type
from
  github_search_code,
  jsonb_array_elements(text_matches) as match
where
  query = 'filename:table_github_my_organization RowsRemaining';
```

### List repository details

```sql
select
  name,
  repository -> 'id' as repo_id,
  repository ->> 'name' as repo_name,
  repository ->> 'url' as repo_url,
  jsonb_pretty(repository -> 'owner') as repo_owner
from
  github_search_code
where
  query = 'filename:table_github_my_organization RowsRemaining';
```
