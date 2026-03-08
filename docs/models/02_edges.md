# Edges

Edges represent relationships between nodes in SurrealDB's graph model. They allow you to connect records with typed, directional relationships that can carry their own properties.

## Defining an Edge

Embed `som.Edge` to create a relationship type:

```go
package model

import (
    "time"
    "yourproject/gen/som"
)

type Follows struct {
    som.Edge

    Since time.Time
}
```

## Edge Structure

The `som.Edge` embedding provides an `ID()` method returning the edge's unique identifier. You define the connected nodes by declaring fields with `som:"in"` and `som:"out"` tags.

## Specifying Connected Nodes

Use the `som:"in"` and `som:"out"` tags to specify the node types:

```go
type GroupMember struct {
    som.Edge
    som.Timestamps

    // Specify the connected node types
    In  User  `som:"in"`   // The user
    Out Group `som:"out"`  // The group they belong to

    // Edge properties
    Role     string
    JoinedAt time.Time
}
```

This creates relationships like:
```surql
RELATE user:alice->group_member->group:developers
```

## Edge Properties

Edges can have their own fields, just like nodes:

```go
type MemberOf struct {
    som.Edge

    Role      string        // "admin", "member", "viewer"
    JoinedAt  time.Time
    IsAdmin   bool
    Metadata  *EdgeMetadata // Optional nested data
}

type EdgeMetadata struct {
    InvitedBy string
    Notes     string
}
```

## Creating Edges

Use the generated `Relate()` method:

```go
// Create an edge connecting two nodes
membership := &model.GroupMember{
    Role:     "admin",
    JoinedAt: time.Now(),
}

// The Relate() method handles the RELATE statement
err := client.GroupMemberRepo().Relate().
    From(user).   // In node
    To(group).    // Out node
    Create(ctx, membership)
```

Or create directly with populated In/Out:

```go
membership := &model.GroupMember{
    In:       *user,    // The user
    Out:      *group,   // The group
    Role:     "member",
    JoinedAt: time.Now(),
}

err := client.GroupMemberRepo().Create(ctx, membership)
```

## Querying Edges

### Find by Source Node

```go
// Find all groups a user belongs to
memberships, err := client.GroupMemberRepo().Query().
    Where(filter.GroupMember.In.Equal(user.ID)).
    All(ctx)
```

### Find by Target Node

```go
// Find all members of a group
memberships, err := client.GroupMemberRepo().Query().
    Where(filter.GroupMember.Out.Equal(group.ID)).
    All(ctx)
```

### Filter by Edge Properties

```go
// Find admin memberships
admins, err := client.GroupMemberRepo().Query().
    Where(
        filter.GroupMember.Out.Equal(group.ID),
        filter.GroupMember.Role.Equal("admin"),
    ).
    All(ctx)
```

## Edge Direction

Edges are **directional**. The `In` → `Out` direction matters:

```
User:alice ──[Follows]──> User:bob
   (In)                    (Out)
```

- `In` is where the relationship **starts**
- `Out` is where the relationship **points to**

For bidirectional relationships, create edges in both directions:

```go
// Alice follows Bob
client.FollowsRepo().Create(ctx, &model.Follows{
    In:    alice,
    Out:   bob,
    Since: time.Now(),
})

// Bob follows Alice (separate edge)
client.FollowsRepo().Create(ctx, &model.Follows{
    In:    bob,
    Out:   alice,
    Since: time.Now(),
})
```

## Edge Timestamps

Like nodes, edges support automatic timestamps:

```go
type Follows struct {
    som.Edge
    som.Timestamps  // CreatedAt, UpdatedAt

    Since time.Time
}
```

## Common Patterns

### Self-Referencing (Same Node Type)

```go
// User follows User
type Follows struct {
    som.Edge

    In  User `som:"in"`
    Out User `som:"out"`

    Since    time.Time
    IsMutual bool
}
```

### Different Node Types

```go
// User owns Document
type Owns struct {
    som.Edge

    In  User     `som:"in"`
    Out Document `som:"out"`

    AcquiredAt time.Time
    Permission string
}
```

### Many-to-Many with Metadata

```go
// Student enrolled in Course
type Enrollment struct {
    som.Edge
    som.Timestamps

    In  Student `som:"in"`
    Out Course  `som:"out"`

    Semester    string
    Grade       *float64
    Status      EnrollmentStatus
    CompletedAt *time.Time
}
```

## Example: Social Network

```go
// Nodes
type User struct {
    som.Node[som.ULID]
    Username string
    Email    string
}

type Post struct {
    som.Node[som.ULID]
    som.Timestamps
    Content string
    Author  *User  // Direct link (not an edge)
}

// Edges
type Follows struct {
    som.Edge
    In    User `som:"in"`
    Out   User `som:"out"`
    Since time.Time
}

type Likes struct {
    som.Edge
    In  User `som:"in"`
    Out Post `som:"out"`
}
```

See [Relationships](../relationships/README.md) for more advanced graph queries and traversal patterns.
