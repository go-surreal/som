# Union Type

The union type represents a record link that may point to any member of a union
of node models. See [Unions](../models/11_unions.md) for how a union is declared.

## Overview

| Property | Value |
|----------|-------|
| Go Type | The union interface (e.g. `Shape`), sealed or open |
| Database Schema | `option<record<square\|circle>>` |
| CBOR Encoding | Direct (as record ID) |
| Sortable | No |
| Pointer | Not allowed — an unset link is the nil interface |

## Definition

```go
type Shape interface {
    som.Union[Shape]

    // Optional: methods every member must implement.
    Area() float64
}

type Square struct {
    som.Node[som.ULID]
    som.Part[Shape]

    Size float64
}

type Circle struct {
    som.Node[som.UUID]
    som.Part[Shape]

    Radius float64
}

type Drawing struct {
    som.Node[som.ULID]

    Name    string
    Subject Shape
}
```

## Schema

```surql
DEFINE FIELD subject ON TABLE drawing TYPE option<record<square|circle>>;
```

## Filter Operations

| Operation | Description | Returns |
|-----------|-------------|---------|
| `Is<Member>()` | Link points to that member's table | Bool filter |
| `<Member>()` | Narrow the link to that member | Member's filter |
| `Nil(true/false)` | Link is (not) set | Bool filter |
| `Equal(table, id)` | Link points to exactly that record | Bool filter |
| `NotEqual(table, id)` | Link points to any other record | Bool filter |
| `IsTable(table)` | Link points to the given table | Bool filter |
| `IsNotTable(table)` | Link points to any other table | Bool filter |
| `InTables(tables)` | Link points to any of the given tables | Bool filter |
| `Exists(true/false)` | Linked record (does not) still exist | Bool filter |

```go
filter.Drawing.Subject().IsSquare()
filter.Drawing.Subject().Square().Size.GreaterThan(3)
filter.Drawing.Subject().Nil(false)
```

Narrowing to a member does not restrict the rows by itself: a link to another
member holds no value for the field, so it never matches the comparison.

## Fetching

```go
drawings, _ := client.DrawingRepo().Query().
    Fetch(with.Drawing.Subject()).
    All(ctx)

switch subject := drawings[0].Subject.(type) {
case *model.Square:
    fmt.Println(subject.Size)
case *model.Circle:
    fmt.Println(subject.Radius)
}
```

## Limitations

- A slice of a union (`[]Shape`) is not supported yet.
- Union fields are not sortable. Sort on a narrowed member's field instead.
