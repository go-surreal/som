# Ideas

This document holds new ideas and upcoming additions to the library.
The main usage of this document is brainstorming on whether the listed ideas are useful and actually possible.
For more in-depth talk about single points, a [GitHub discussion](https://github.com/marcbinz/som/discussions)
can be created at any point and linked back to this document.

## References 

- https://pg.uptrace.dev/sql-null-go-zero-values/
- https://ente.io/blog/tech/go-nulls-and-sql/
- 

## Features

### Soft Delete

- https://www.jmix.io/blog/to-delete-or-to-soft-delete-that-is-the-question/
- https://www.brentozar.com/archive/2020/02/what-are-soft-deletes-and-how-are-they-implemented/
- 

### Support custom ID constructs

```
som.Timeseries
```

### Add `som.Slug` type

- https://github.com/gosimple/slug

- is a string internally
- automatically converts its value to a unique slug representation
- 1: check if value is already in slug format (slug it and compare)
- 2: if it is, do nothing more
- 3: if it is not, slug it and check for uniqueness
- 4: if it is unique, do nothing more
- 5: if it is not unique, append or bump a suffix (e.g. "some-slug-1" to "some-slug-2")

### Add automated data seeding capabilities

- https://medium.easyread.co/how-i-seed-my-database-with-go-27488d2e6a75
- https://github.com/kristijorgji/goseeder
- https://medium.com/@thegalang/automatic-data-seeding-on-go-api-server-4ba7eb8c1881
- https://www.reddit.com/r/golang/comments/8o8s7f/how_do_you_seed_your_database/
- https://ieftimov.com/posts/simple-golang-database-seeding-abstraction-gorm/

### Query math functions

```
math::max() - Returns the maximum number in a set of numbers
math::mean() - Returns the mean of a set of numbers
math::median() - Returns the median of a set of numbers
math::min() - Returns the minimum number in a set of numbers
math::product() - Returns the product of a set of numbers
math::sum() - Returns the total sum of a set of numbers
```

### Support Email type

```
som.Email
```

```
parse::email::domain() - Parses and returns an email domain from an email address
parse::email::user() - Parses and returns an email username from an email address
```

### Support URL type

```
url.Url // go standard type
```

```
parse::url::domain() - Parses and returns the domain from a URL
parse::url::fragment() - Parses and returns the fragment from a URL
parse::url::host() - Parses and returns the hostname from a URL
parse::url::path() - Parses and returns the path from a URL
parse::url::port() - Parses and returns the port number from a URL
parse::url::query() - Parses and returns the query string from a URL
```

### More query functions

- https://surrealdb.com/docs/surrealql/functions/string
- https://surrealdb.com/docs/surrealql/functions/time

### Initial value for UUID

```
DEFINE FIELD uuid ON x VALUE $before OR rand::uuid()
```

### Pagination support

Ready-to-use pagination support with direct use via GraphQL.

Both cursor- and page-based.

```go

client.User.Query().
	Filter(...).
	Paginate(cursor, size)

client.User.Query().
    Filter(...).
    LoadPage(page, size)

```

- https://stackoverflow.com/questions/58787039/cursor-based-pagination-for-search-results-without-sequential-unique-ids-eg-lo

### Aggregations

```
-- Group results with aggregate functions
SELECT count() AS total, math::sum(age), gender, country FROM person GROUP BY gender, country;
```

### Selection distinct sub-values directly

```go

var groupNameStrings []string = client.User().Query().
  Distinct(field.User.Groups().Name)
// SELECT groups.name FROM user GROUP BY groups.name

userNameGroupNameTuples := client.User().Query().
  Distinct2(field.User.Name, field.User.Groups().Name)

allGroupNamesOfUser := som.Values(
  field.User.Groups().Name
)

```

### Choose the fields to update specifically

```go

repo.db.user.Update(&userModel, 
  field.User.FirstName,
  field.User.LastName,
)

```

### Better edges with generics

```
package som

type Edge[I, O, P any] struct {
    In I
    Out O
    Props P
}
```

```
type MemberOf som.Edge[User, Group, MemberOfProps]

type MemberOfProps struct {
    som.Timestamps
    
    Roles []Role
}
```

### (Automatic) Migrations

- https://www.edgedb.com/showcase/migrations

### System-versioned temporal tables

- https://learn.microsoft.com/en-us/sql/relational-databases/tables/temporal-tables?view=sql-server-ver16
- MIGHT COME AS NATIVE FEATURE!

### Calculate the complexity (and cost) of a query and issue warnings if too high

- https://www.edgedb.com/blog/why-orms-are-slow-and-getting-slower

### Partial updates

- https://incident.io/blog/code-generation

### ?

```golang

func (r *user) QueryUsersByPartialName(name string) query.User {
	return r.Query().
		Filter(
			where.User.String.FuzzyMatch(name),
		)
}

r.QueryUsersByPartialName("some").
  Filter(
    where.User.Int.GreaterThan(5),
    where.User.Int.LessThan(10),
  ).
  All(ctx)

```

### Allow for union of connected records

```
TYPE record(table1, table2, etc...)
```

```
type SomeRecordUnion interface {
    someRecordUnion()
}

type X struct {
    som.Node
    SomeRecordUnion
}

type Y struct {
    som.Node
    SomeRecordUnion
}

type Z struct {
    som.Node
    XY SomeRecordUnion
}
```

```
DEFINE FIELD xy ON z TYPE record(x, y)
```

## Define custom (type-safe) views

```go

db.Define().View("some_name").
	With(
	    field.User.Name(),	
    ).
	By(
	    ...	
    )

```

Based on the definition a view will be created in the database and 
a type generated that can be consumed by the usual query builder.

## Optimisations

- replace string concatenation in query builder with buffer (faster execution)

## Other links

- https://www.liquibase.com/resources/guides/database-version-control
- https://www.dolthub.com/blog/2022-08-04-database-versioning/
- https://go-rel.github.io/introduction/
- https://hackernoon.com/introducing-bun-a-golang-orm
- https://www.quora.com/What-are-the-basics-of-building-an-object-relationship-mapper-ORM-for-an-SQL-database
