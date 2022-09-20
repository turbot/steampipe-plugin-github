# Table: github_tree

A Git tree object creates the hierarchy between files in a Git repository. You
can use the Git tree object to create the relationship between directories and
the files they contain.

A single tree object contains one or more entries, each of which is the SHA-1
hash of a blob or subtree with its associated mode, type, and filename.

The `github_tree` table can be used to query information about any tree, and
**you must specify which repository and tree SHA** in the where or join clause
using the `repository_full_name` and `tree_sha` columns. By default, recursive
entries are not returned, but can be with the `recursive` column.

## Examples

### List tree entries non-recursively

```sql
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a';
```

### List tree entries for a subtree recursively

```sql
select
  tree_sha,
  truncated,
  path,
  mode,
  type,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '5622172b528cd38438c52ecfa3c20ac3f71dd2df'
  and recursive = true;
```

### List executable files

```sql
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = true
  and mode = '100755';
```

### List JSON files

```sql
select
  tree_sha,
  truncated,
  path,
  mode,
  size,
  sha
from
  github_tree
where
  repository_full_name = 'turbot/steampipe'
  and tree_sha = '0f200416c44b8b85277d973bff933efa8ef7803a'
  and recursive = true
  and path like '%.json';
```
