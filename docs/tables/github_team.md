# Table: github_team

The `github_team` table can be used to fetch team information for a given organization. **You must specify the organization** in the where or join clause (`where organization=`, `join github_repository on organization=`).

## Examples

## Get all teams in an organization

```sql
select
  t.id
  t.name
  t.privacy
  t.description
from
    github.github_team
where organization = 'my_org'
```

## Get all teams in your organizations

```sql
select
  o.login,
  t.name,
  t.slug
from github.github_my_organization as o
inner join github.github_team as t
  on t.organization = o.login
```
