# Fragments

A fragment is a projection of a node model: a struct holding a subset of the node's fields.
Querying a fragment narrows the `SELECT` list, so only the fields that are actually needed are
read from the database and decoded.

Unlike a [view](07_views.md), a fragment has no table of its own — it is a different Go type
for the same table.

## Defining a Fragment

Embed `som.Fragment[T]`, where `T` is the node model the fragment projects:

```go
package model

import "yourproject/gen/som"

type Person struct {
    som.Node[som.ULID]

    Name    string
    Email   string
    Age     int
    Bio     string
    Address Address
}

// PersonCard projects the fields needed to render a list entry.
type PersonCard struct {
    som.Fragment[Person]

    Name string
    Age  int
}
```

Every field of a fragment must exist on the parent node with the same name and the same type.
The field's database name, its type and its `som` tag are taken over from the node, so a
fragment field must not carry a `som` tag of its own — that keeps the two definitions from
drifting apart.

Code generation rejects a fragment that

- declares a field the node does not have,
- declares a field whose type differs from the node's,
- declares an `ID` field (the record id is always projected),
- projects no field at all.

Fields provided by an embed on the node (`som.Timestamps`, `som.SoftDelete`, `som.Expiry`) are
not projectable.

## Querying a Fragment

Fragments are fetched through the regular query builder. Filters, ordering, `Limit`/`Start`
and `Fetch` stay typed on the parent model; only the terminal call selects the fragment:

```go
cards, err := client.PersonRepo().Query().
    Where(filter.Person.Age.GreaterThan(18)).
    Order(by.Person.Name.Asc()).
    AllAs[model.PersonCard](ctx)
```

This runs:

```surql
SELECT id, name, age FROM person WHERE age > 18 ORDER BY name ASC
```

The available terminals are:

| Terminal | Description |
|----------|-------------|
| `AllAs[F](ctx)` | All matching records as fragments |
| `AllAsAsync[F](ctx)` | Asynchronous `AllAs` |
| `FirstAs[F](ctx)` | First matching record as a fragment |
| `FirstAsAsync[F](ctx)` | Asynchronous `FirstAs` |
| `IterateAs[F](ctx, batchSize)` | Keyset-paged iteration over fragments |
| `LiveAs[F](ctx)` | Live query delivering fragments |
| `DescribeAs[F]()` | The narrowed statement as a string, for debugging |

A fragment of the wrong model is a compile-time error, since `F` is constrained to
`som.FragmentOf[<model>]`:

```go
// PersonCard does not satisfy som.FragmentOf[model.Product]
client.ProductRepo().Query().AllAs[model.PersonCard](ctx)
```

`DescribeAs` shows what is actually sent to the database:

```go
client.PersonRepo().Query().
    Order(by.Person.Name.Asc()).
    DescribeAs[model.PersonCard]()
// SELECT id, name, age FROM person ORDER BY name ASC
```

### Ordering

`ORDER BY` resolves against the projected record, so a sort field that the fragment does not
hold is added to the projection automatically. The same applies to the keyset cursor used by
`IterateAs`.

## Record ID and Expanding

The record id is always projected. It is available as a string through the promoted
`RecordID` method:

```go
card.RecordID() // "person:01JBZ..."
```

To get the full record a fragment was projected from, use `Expand` on the parent repository:

```go
person, exists, err := client.PersonRepo().Expand(ctx, card)
```

`Expand` is generated on the repository of every node that has at least one fragment. It
accepts any fragment of that node.

## Read-Only

A fragment instance holds only a subset of the record, so it is always marked partial
(see [Repository](../api_reference/02_repository.md)):

```go
card.IsPartial() // true
card.Marker()    // som.MarkerLoaded | som.MarkerPartial
```

There is no `Create`, `Update` or `Delete` for fragments — the write methods accept the node
model only. Expand the fragment first if the record needs to be written.

## Caveats

- **Caching.** Fragment queries bypass the in-process cache ([Caching](../api_reference/04_caching.md)),
  which holds full models. `Expand` uses the regular cached read path.
- **Links are not resolved.** A projected record link decodes to a partial model holding just
  its record id, exactly as in a full-model query. Use `Fetch` to resolve it.
- **Live deletes.** For delete notifications SurrealDB only sends the record id, so the
  received fragment holds no field values.
- **No cursor pagination.** `Paginate` returns full models; there is no fragment variant yet.
  Use `IterateAs` for batched walks, or `Limit`/`Start` with `AllAs`.
