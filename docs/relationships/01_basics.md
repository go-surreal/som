# Relationship Basics

SurrealDB is a graph database, and SOM provides first-class support for modeling relationships between records.

## Types of Relationships

### Record Links

Direct references from one node to another:

```go
type Post struct {
    som.Node[som.ULID]

    Title  string
    Author *User  // Link to User record
}
```

### Edges

Typed relationships with their own properties:

```go
type Follows struct {
    som.Edge

    In  *User `som:"in"`   // Source: who is following
    Out *User `som:"out"`  // Target: who is being followed

    Since    time.Time
    IsMutual bool
}
```

## Defining Record Links

Link to another node using a pointer:

```go
type Comment struct {
    som.Node[som.ULID]

    Content string
    Author  *User  // Single link
    Post    *Post  // Single link
}
```

Links are stored as record IDs in the database and can be fetched using the query builder.

## Defining Edge Relationships

Create an edge type by embedding `som.Edge` and specifying `In`/`Out` fields:

```go
type MemberOf struct {
    som.Edge

    In  *User  `som:"in"`   // The user joining
    Out *Group `som:"out"`  // The group being joined

    Role     string
    JoinedAt time.Time
}
```

### Required Edge Fields

Edges must have:

1. `som.Edge` embedding - Provides ID and base functionality
2. `In` field with `som:"in"` tag - Source node of the relationship
3. `Out` field with `som:"out"` tag - Target node of the relationship

Both `In` and `Out` must be pointers to Node types.

### Edge Direction

Edges are directional:

- **In** - The source/origin node
- **Out** - The target/destination node

Think of it as: `In` --[Edge]--> `Out`

For example, in a "Follows" edge:
- `In` = The follower (who follows)
- `Out` = The followed (who is being followed)

## Creating Relationships

### Using Relate Builder

```go
// Alice follows Bob
follows := &model.Follows{
    Since:    time.Now(),
    IsMutual: false,
}

err := client.FollowsRepo().Relate().
    From(alice).  // Sets the In field
    To(bob).      // Sets the Out field
    Create(ctx, follows)
```

### Linking Nodes Directly

For simple record links:

```go
post := &model.Post{
    Title:  "Hello World",
    Author: user,  // Direct reference
}
err := client.PostRepo().Create(ctx, post)
```

## Querying Relationships

### Query by Source (In)

Find all relationships originating from a node:

```go
// Everyone Alice follows
following, err := client.FollowsRepo().Query().
    Where(filter.Follows.In.Equal(alice.ID())).
    All(ctx)
```

### Query by Target (Out)

Find all relationships pointing to a node:

```go
// Everyone following Bob
followers, err := client.FollowsRepo().Query().
    Where(filter.Follows.Out.Equal(bob.ID())).
    All(ctx)
```

### Query by Edge Properties

```go
// Find relationships created this month
recentFollows, err := client.FollowsRepo().Query().
    Where(
        filter.Follows.Since.After(startOfMonth),
    ).
    All(ctx)
```

## When to Use Each

### Record Links

Best for:

- Simple parent-child relationships
- Required references (author, owner)
- When you don't need relationship metadata
- One-to-many relationships

```go
type Post struct {
    som.Node[som.ULID]
    Author   *User     // Required author
    Category *Category // Optional category
}
```

### Edges

Best for:

- Many-to-many relationships
- Relationships with properties (role, timestamp, weight)
- Graph traversal queries
- Social connections (follows, friends, blocks)
- Bidirectional relationships

```go
type Friendship struct {
    som.Edge
    In       *User `som:"in"`
    Out      *User `som:"out"`
    Since    time.Time
    Strength int  // Relationship weight
}
```

## Example: Social Network

```go
// Nodes
type User struct {
    som.Node[som.ULID]
    som.Timestamps

    Name  string
    Email string
}

type Group struct {
    som.Node[som.ULID]

    Name        string
    Description string
    IsPrivate   bool
}

type Post struct {
    som.Node[som.ULID]
    som.Timestamps

    Content string
    Author  *User
}

// Edges
type Follows struct {
    som.Edge

    In  *User `som:"in"`
    Out *User `som:"out"`

    Since    time.Time
    IsMutual bool
}

type MemberOf struct {
    som.Edge

    In  *User  `som:"in"`
    Out *Group `som:"out"`

    Role     string
    JoinedAt time.Time
}

type Likes struct {
    som.Edge

    In  *User `som:"in"`
    Out *Post `som:"out"`

    LikedAt time.Time
}
```

## Fetching Related Records

Use `Fetch` to eager-load related records:

```go
// Get posts with their authors
posts, err := client.PostRepo().Query().
    Fetch(with.Post.Author...).
    All(ctx)

for _, post := range posts {
    fmt.Println(post.Author.Name)  // Already loaded
}
```
