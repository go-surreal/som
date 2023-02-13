# Ideas

This document holds new ideas and upcoming additions to the library.
The main usage of this document is brainstorming on whether the listed ideas are useful and actually possible.
For more in-depth talk about single points, a [GitHub discussion](https://github.com/marcbinz/som/discussions)
can be created at any point and linked back to this document.

## Features

### Support custom ID constructs

```
som.Timeseries
```

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

## Optimisations

tbd.
