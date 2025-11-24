# Ideas

This document holds new ideas and upcoming additions to the library.
The main usage of this document is brainstorming on whether the listed ideas are useful and actually possible.
For more in-depth talk about single points, a [GitHub discussion](https://github.com/go-surreal/som/discussions)
can be created at any point and linked back to this document.

## References 

- https://pg.uptrace.dev/sql-null-go-zero-values/
- https://ente.io/blog/tech/go-nulls-and-sql/
- 

## Features

### ?

Can I get the total count when doing a paginate query?
So a count where limit and offset are basically ignored, but the result set still takes them into account.

### Views

```
package model

import "github.com/go-surreal/som"

type TestView struct {
	som.View
}
```

### Migrations

Rename:

```
DEFINE FIELD new ON x VALUE $before OR $after.old; // works with value? otherwise:
UPDATE x SET new = old;
DROP FIELD old ON x;
```

Ignore:

```
change.Field.X.Ignore(reason: "must add comment here")
change.Field.X.TODO()
```

### Custom functions

https://surrealdb.com/docs/surrealdb/surrealql/datamodel/closures

Note: SOM might define its own functions in the future.

### Computed Fields

tbd.

### LiveRead model

```
func (r *allFieldTypes) LiveRead(ctx context.Context, id *som.ID) (*model.AllFieldTypes, bool, error) {
	allFieldTypes, exists, err := r.Read(ctx, id) // TODO: mark model as live (similar to fragment) to prevent it from being updated
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, false, nil
	}

	liveRes, err := r.Query().
		Filter(where.AllFieldTypes.ID.Equal(id)).
		Live(ctx)
	if err != nil {
		return nil, false, err
	}

	go func() {
		for {
			select {

			case live := <-liveRes:
				{
				switch res := live.(type) {

				case query.LiveUpdate[model.AllFieldTypes]:
					updatedModel,err := res.Get()
					if err != nil {
						return
					}
					*allFieldTypes = updatedModel

				case query.LiveCreate[model.AllFieldTypes]:
				case query.LiveDelete[model.AllFieldTypes]:
				}
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return allFieldTypes, true, nil
}
```

### Cache

```
// TTL sets the time-to-live for the result set of the query.
// After the given duration, the result set will be invalidated.
// This means that the next query will re-fetch the data from the database.
func (b builder[M, C]) TTL(dur time.Duration) string {
	panic("not implemented")
}
```

### On Delete Cascade

not yet a native feature but might be at some time

```
DEFINE EVENT event_name ON TABLE table_name WHEN ($after == NONE) THEN {
    delete contact where id == $before.contact;
    delete adress where id inside $before.adresses;
};
```

https://github.com/surrealdb/surrealdb/issues/1782
https://github.com/surrealdb/surrealdb/issues/1783

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
(or use patch instead)

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

### Add dirty flag handling

- cache the data itself within the read model (would double the memory usage)
- cache the data in a separate cache (would require a cache invalidation strategy)

### Read through GORM to find more ideas

- https://gorm.io/docs

### Add ability to deep load models

```go

err := db.User.Fetch(ctx, []model.User{}, 
	field.User.Groups(),
	field.User.Groups().Members(),
)

```

### Indexing for schemaless (needed?)

```
DEFINE EVENT example_uid_setting ON TABLE example WHEN $before = null AND $before!=$after THEN {
    LET $next = (SELECT val FROM counter:example) +1 ;
    UPDATE $after SET uid = $next  
    UPDATE counter:example SET val = $next;
};
```

### Caching

add caching to the generation process

### Watch

add a watch command to the generation process
should watch the input dir as well as all other imported dirs

## Optimisations

tbd.

## Other links

- https://www.liquibase.com/resources/guides/database-version-control
- https://www.dolthub.com/blog/2022-08-04-database-versioning/
- https://go-rel.github.io/introduction/
- https://hackernoon.com/introducing-bun-a-golang-orm
- https://www.quora.com/What-are-the-basics-of-building-an-object-relationship-mapper-ORM-for-an-SQL-database
