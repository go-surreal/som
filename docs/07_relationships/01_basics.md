# Relationship Basics

SurrealDB is a graph database, and SOM provides first-class support for modeling relationships between records.

## Types of Relationships

### Record Links

Direct references from one node to another:

```go
type Post struct {
    som.Node

    Title  string
    Author *User  // Link to User record
}
```

### Edges

Typed relationships with their own properties:

```go
type Follows struct {
    som.Edge

    Since    time.Time
    IsMutual bool
}
```

## Defining Record Links

Link to another node using a pointer:

```go
type Comment struct {
    som.Node

    Content string
    Author  *User  // Single link
    Post    *Post  // Single link
}
```

## Defining Edge Relationships

Create an edge type by embedding `som.Edge`:

```go
type MemberOf struct {
    som.Edge

    Role     string
    JoinedAt time.Time
}
```

Edges connect nodes directionally:
- `In` - Source node
- `Out` - Target node

## When to Use Each

**Record Links** are best for:
- Simple parent-child relationships
- Required references
- When you don't need relationship metadata

**Edges** are best for:
- Many-to-many relationships
- Relationships with properties (role, timestamp, etc.)
- Graph traversal queries
- Social connections (follows, friends, etc.)

## Example: Social Network

```go
// Nodes
type User struct {
    som.Node
    Name string
}

type Group struct {
    som.Node
    Name string
}

// Edges
type Follows struct {
    som.Edge
    Since time.Time
}

type MemberOf struct {
    som.Edge
    Role string
}
```
