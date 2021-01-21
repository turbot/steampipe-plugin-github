# Table: github_gist

Github Gist is a simple way to share snippets and pastes with others.  The `github_gist` table will list **public** gists that **you own**.  You can query **ANY** gist that you have access to by specifying its `id` explicitly in the where clause with  `where id =`  .

## Examples

### List your public gists

```sql
select
  *
from
  github_gist;
```


### Get details about ANY public gist (by id)

```sql
select
  *
from
  github_gist
where
  id='633175';
```
