# Edges

Edges represent relationships between nodes in SurrealDB's graph model. They connect two records
with a typed, directional relationship that can carry its own properties.

## Defining an Edge

An edge is a struct embedding `som.Edge` with exactly one field tagged `som:"in"` and one tagged
`som:"out"`:

```go
package model

import (
    "time"
    "yourproject/gen/som"
)

type MemberOf struct {
    som.Edge

    User  User  `som:"in"`   // source node
    Group Group `som:"out"`  // target node

    Role     string
    JoinedAt time.Time
}
```

Rules enforced by the generator:

- The struct must embed `som.Edge` (provides `ID() string`).
- Exactly one field must be tagged `som:"in"`, one `som:"out"`.
- Both must be node types, used as values (`User`, not `*User`).
- A field named `ID` is not allowed — it is provided by `som.Edge`.
- `som.Timestamps`, `som.OptimisticLock` and `som.SoftDelete` may be embedded.
  `som.Expiry` is **not** supported on edges.

The names of the `in`/`out` fields are up to you — only the tags matter. Those names become the
accessors used when filtering across the edge.

## Registering the Edge on a Node

To create or traverse an edge, the node it starts from must declare it as a field:

```go
type User struct {
    som.Node[som.ULID]

    Name        string
    Memberships []MemberOf  // edge field: generates relate/filter accessors
}
```

The field name (`Memberships`) becomes the accessor name in the generated `relate` and `filter`
packages.

## Edge Properties

Edges can have their own fields, just like nodes, including nested structs:

```go
type MemberOf struct {
    som.Edge
    som.Timestamps

    User  User  `som:"in"`
    Group Group `som:"out"`

    Role string
    Meta EdgeMetadata
}

type EdgeMetadata struct {
    InvitedBy string
    Notes     string
}
```

## Creating Edges

Edges have no repository of their own. They are created through the `Relate()` builder of the
node the edge starts from:

```go
// Both nodes must exist already
membership := &model.MemberOf{
    User:     *user,
    Group:    *group,
    Role:     "admin",
    JoinedAt: time.Now(),
}

err := client.UserRepo().Relate().Memberships().Create(ctx, membership)
```

This executes a `RELATE user:...->member_of->group:...` statement. After a successful call, the
edge carries its generated ID and any database defaults (for example timestamps).

Errors are returned when:

- the edge is `nil`,
- the edge already has an ID set,
- the `in` or `out` node has an empty ID.

> `Update` and `Delete` on edges are not implemented yet and return an error.

## Querying Across Edges

Edges are traversed from a node query using the generated edge accessor on the node's filter.
The accessor takes filters on the edge itself and then exposes the connected node for further
filtering:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.User.
            Memberships(
                filter.MemberOf.Role.Equal("admin"),
            ).
            Group(
                filter.Group.Name.Equal("developers"),
            ),
    ).
    All(ctx)
```

`Memberships(...)` is the node's edge field, `Group(...)` the edge's `out` field. Selecting edge
records directly (as you would with a repository) is not supported — model the traversal from one
of the connected nodes instead.

## Edge Direction

Edges are **directional**. The `in` → `out` direction matters:

```
user:alice ──[member_of]──> group:developers
    (in)                          (out)
```

For bidirectional relationships, create two edges — one per direction.

## Table Naming

The edge table name is the snake_case form of the struct name:

| Struct | Table |
|--------|-------|
| `Follows` | `follows` |
| `MemberOf` | `member_of` |
| `GroupMember` | `group_member` |

## Changefeed on Edges

Like nodes, edges support a changefeed via a tag on the embedded `som.Edge`:

```go
type MemberOf struct {
    som.Edge `som:"changefeed=1d"`

    User  User  `som:"in"`
    Group Group `som:"out"`
}
```

See [Changefeed](../querying/08_changefeed.md).

See [Relationships](../relationships/README.md) for record links, traversal patterns and eager
loading.
