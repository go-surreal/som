# Distinct Values

`query.Distinct` fetches the distinct values of a single field instead of whole records. It is
useful for filter dropdowns, facets and "which categories exist" style queries, and avoids
loading full records just to deduplicate them in Go.

## Usage

```go
import (
    "yourproject/gen/som/field"
    "yourproject/gen/som/query"
)

categories, err := query.Distinct(ctx,
    client.ArticleRepo().Query(),
    field.Article.Category,
)
// categories is []string
```

`Distinct` is a package-level function (Go methods cannot introduce their own type parameters),
taking:

1. the context,
2. a query builder — any filters, joins and options on it apply,
3. a field reference from the generated `field` package.

The return type follows the field type: `field.Article.ViewCount` yields `[]int`,
`field.Article.PublishedAt` yields `[]time.Time`, an enum field yields `[]model.YourEnum`, and so
on.

## With Filters

The builder is used as-is, so filtering works as usual:

```go
categories, err := query.Distinct(ctx,
    client.ArticleRepo().Query().
        Where(filter.Article.Published.True()),
    field.Article.Category,
)
```

## Nested and Slice Fields

The `field` package mirrors the model structure, including nested structs:

```go
cities, err := query.Distinct(ctx,
    client.UserRepo().Query(),
    field.User.Address.City,
)
```

For slice fields the distinct values of the **elements** are returned, not the distinct slices:

```go
// []string of every tag used by any post
tags, err := query.Distinct(ctx,
    client.PostRepo().Query(),
    field.Post.Tags,
)
```

## Notes

- An empty table yields a `nil` slice and no error.
- Ordering is not guaranteed — sort in Go if you need a stable order.
- Only one field can be requested per call.
