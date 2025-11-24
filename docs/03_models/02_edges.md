# Edges

Edges represent relationships between nodes in SurrealDB's graph model. They allow you to connect records with typed, directional relationships that can carry their own properties.

## Defining an Edge

Embed `som.Edge` to create a relationship type:

```go
package model

import (
    "time"
    "github.com/go-surreal/som"
)

type Follows struct {
    som.Edge

    Since time.Time
}
```

## Edge Structure

Every edge automatically has:
- `In` - The source node (where the relationship starts)
- `Out` - The target node (where the relationship points)
- `ID` - Unique identifier for the edge itself

## Edge Properties

Edges can have their own fields, just like nodes:

```go
type MemberOf struct {
    som.Edge

    Role      string
    JoinedAt  time.Time
    IsAdmin   bool
}
```

## Connecting Nodes

Edges connect two node types. The relationship is directional:

```go
// User -[Follows]-> User
// User -[MemberOf]-> Group
```

## Graph Queries

With edges defined, you can traverse the graph:

```go
// Find all users that a specific user follows
followers, err := client.FollowsRepo().Query().
    Filter(where.Follows.In.Equal(userID)).
    All(ctx)
```

See [Relationships](../07_relationships/README.md) for more details on working with edges.
