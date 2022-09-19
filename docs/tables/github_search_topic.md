# Table: github_search_topic

The `github_search_topic` table helps to find topics via various criteria. You can search for topics on GitHub, explore related topics, and see how many repositories are associated with a certain topic.

 **You must always include at least one search term when searching source code** in the where or join clause using the `query` column. See [Searching topics](https://docs.github.com/search-github/searching-on-github/searching-topics) for details on the GitHub query syntax.

## Examples

### List topics that are not curated

```sql
select
  name,
  created_at,
  curated
from
  github_search_topic
where
  query = 'target-searching is:not-curated repo:turbot/steampipe-plugin-github';
```

### List featured topics

```sql
select
  name,
  created_at,
  featured
from
  github_search_topic
where
  query = 'markdown is:featured repo:turbot/steampipe-plugin-github';
```

### List topics that have more than 5000 repositories

```sql
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'repositories:>5000 react-redux';
```

### List topics that were created after a specific timestamp

```sql
select
  name,
  created_at,
  created_by,
  featured,
  curated
from
  github_search_topic
where
  query = 'created:>2021-01-01 react-redux';
```
