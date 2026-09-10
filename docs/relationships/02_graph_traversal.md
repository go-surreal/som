# Graph Traversal

SurrealDB's graph capabilities let you query across relationships. SOM exposes them through the
generated node filters — every traversal starts at a node query.

The examples below use this model:

```go
type User struct {
    som.Node[som.ULID]

    Name        string
    Follows     []Follows
    Memberships []MemberOf
}

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

    Role string
}
```

## Traversing an Edge

The edge accessor on the node filter takes filters applied to the edge itself, and returns an
accessor for the connected nodes:

```go
// Users who follow someone named "Bob"
users, err := client.UserRepo().Query().
    Where(
        filter.User.
            Follows().
            To(filter.User.Name.Equal("Bob")),
    ).
    All(ctx)
```

## Filtering on the Edge

Pass filters to the edge accessor to constrain the relationship itself:

```go
// Users who started following someone in the last 30 days
users, err := client.UserRepo().Query().
    Where(
        filter.User.
            Follows(
                filter.Follows.Since.After(time.Now().AddDate(0, 0, -30)),
            ).
            To(filter.User.Name.Equal("Bob")),
    ).
    All(ctx)
```

Edge and target filters combine — both must match.

## Combining With Node Filters

Traversal filters are ordinary filters, so they combine with everything else:

```go
users, err := client.UserRepo().Query().
    Where(
        filter.User.Name.StartsWith("A").True(),
        filter.User.
            Memberships(
                filter.MemberOf.Role.Equal("admin"),
            ).
            Group(filter.Group.Name.Equal("developers")),
    ).
    Order(by.User.Name.Asc()).
    All(ctx)
```

## Counting and Existence

Because traversal happens inside a node query, `Count()` and `Exists()` work as usual:

```go
// How many users are admins of this group?
count, err := client.UserRepo().Query().
    Where(
        filter.User.
            Memberships(filter.MemberOf.Role.Equal("admin")).
            Group(filter.Group.Name.Equal("developers")),
    ).
    Count(ctx)

// Does Alice follow Bob?
isFollowing, err := client.UserRepo().Query().
    Where(
        filter.User.Name.Equal("Alice"),
        filter.User.Follows().To(filter.User.Name.Equal("Bob")),
    ).
    Exists(ctx)
```

## Live Queries Over Traversals

Traversal filters keep the builder live-capable:

```go
liveChan, err := client.UserRepo().Query().
    Where(
        filter.User.Memberships().Group(filter.Group.Name.Equal("developers")),
    ).
    Live(ctx)
```

See [Live Queries](../querying/04_live_queries.md) for the event types.

## Limitations

- There is no repository for edges — edge rows cannot be selected, updated or deleted directly.
  Traverse from a node instead.
- Traversal depth is limited to what you express through chained accessors on the generated
  filters. For arbitrary SurrealQL graph expressions, use
  [Raw Queries](../querying/06_raw_queries.md).

## Performance Tips

- Add indexes to fields you filter on frequently (`som:"index"`).
- Narrow the node query before traversing.
- Prefer `Count()` / `Exists()` over fetching full records when you only need a number or a
  boolean.
- Use `Limit()` on large result sets.
