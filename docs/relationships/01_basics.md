# Relationship Basics

SurrealDB is a graph database, and SOM provides first-class support for modeling relationships
between records. There are two ways to relate records: **record links** and **edges**.

## Record Links

A record link is a direct reference from one node to another, stored as a record ID:

```go
type Post struct {
    som.Node[som.ULID]

    Title  string
    Author *User   // single link, optional
    Tags   []Tag   // slice of links
}
```

Links can be values (`User`), pointers (`*User`), slices (`[]User`, `[]*User`) or pointers to
slices. Only the record ID is stored; the linked record is loaded on demand via
[`Fetch()`](#eager-loading-with-fetch).

## Edges

An edge is a separate table connecting two nodes, with its own ID and optional properties:

```go
type MemberOf struct {
    som.Edge

    User  User  `som:"in"`   // source
    Group Group `som:"out"`  // target

    Role     string
    JoinedAt time.Time
}
```

For the edge to be usable, the source node declares it as a field:

```go
type User struct {
    som.Node[som.ULID]

    Name        string
    Memberships []MemberOf
}
```

See [Edges](../models/02_edges.md) for the full rules.

## Creating Relationships

### Record Links

Assign the linked record and create as usual. The linked record must already exist:

```go
post := &model.Post{
    Title:  "Hello World",
    Author: user,
}
err := client.PostRepo().Create(ctx, post)
```

### Edges

Use the `Relate()` builder of the node the edge starts from:

```go
membership := &model.MemberOf{
    User:     *user,
    Group:    *group,
    Role:     "admin",
    JoinedAt: time.Now(),
}

err := client.UserRepo().Relate().Memberships().Create(ctx, membership)
```

## Querying Relationships

### Through Record Links

A single link field is an accessor returning the target model's filters:

```go
posts, err := client.PostRepo().Query().
    Where(filter.Post.Author().Name.Equal("Alice")).
    All(ctx)
```

A slice of links takes the sub-filters as arguments and returns a slice filter:

```go
posts, err := client.PostRepo().Query().
    Where(
        filter.Post.Tags(
            filter.Tag.Name.Equal("golang"),
        ).NotEmpty(),
    ).
    All(ctx)
```

### Through Edges

The edge accessor takes filters on the edge, then exposes the connected node:

```go
// Users who are admins of the "developers" group
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

Edges have no repository, so there is no `client.MemberOfRepo()`. Every edge query starts from a
node.

## Eager Loading with Fetch

`Fetch()` resolves record links in the same query. Each fetchable relation is a method on the
generated `with` accessor:

```go
posts, err := client.PostRepo().Query().
    Fetch(with.Post.Author()).
    All(ctx)

for _, post := range posts {
    fmt.Println(post.Author.Name)  // already loaded
}
```

Nested relations chain:

```go
client.PostRepo().Query().
    Fetch(with.Post.Author().Organization()).
    All(ctx)
```

> Soft-delete filtering does **not** apply to fetched relations — deleted records are still
> returned. Filter them in application code if needed.

## When to Use Each

### Record Links

Best for:

- Simple parent-child references (author, owner, category)
- One-to-many relationships
- Cases without relationship metadata

### Edges

Best for:

- Many-to-many relationships
- Relationships carrying properties (role, weight, timestamps)
- Graph traversal across multiple hops
- Social connections (follows, friends, blocks)

## Example: Social Network Model

```go
// Nodes
type User struct {
    som.Node[som.ULID]
    som.Timestamps

    Name        string
    Email       string
    Follows     []Follows
    Memberships []MemberOf
}

type Group struct {
    som.Node[som.ULID]

    Name      string
    IsPrivate bool
}

type Post struct {
    som.Node[som.ULID]
    som.Timestamps

    Content string
    Author  *User  // record link
}

// Edges
type Follows struct {
    som.Edge

    From  User `som:"in"`
    To    User `som:"out"`
    Since time.Time
}

type MemberOf struct {
    som.Edge

    User  User  `som:"in"`
    Group Group `som:"out"`

    Role     string
    JoinedAt time.Time
}
```
