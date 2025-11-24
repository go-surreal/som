# Graph Traversal

SurrealDB's graph capabilities allow you to traverse relationships between records.

## Querying Edges

Find all relationships of a specific type:

```go
// Find all "follows" relationships
follows, err := client.FollowsRepo().Query().All(ctx)
```

## Filtering by Source/Target

Find relationships from a specific node:

```go
// Find everyone a user follows
follows, err := client.FollowsRepo().Query().
    Filter(where.Follows.In.Equal(userID)).
    All(ctx)
```

Find relationships to a specific node:

```go
// Find everyone who follows a user
followers, err := client.FollowsRepo().Query().
    Filter(where.Follows.Out.Equal(userID)).
    All(ctx)
```

## Edge Properties

Query edges by their properties:

```go
// Find recent follows (last 30 days)
recentFollows, err := client.FollowsRepo().Query().
    Filter(
        where.Follows.Since.GreaterThan(time.Now().AddDate(0, 0, -30)),
    ).
    All(ctx)
```

## Combining Filters

```go
// Find mutual follows created this year
mutualFollows, err := client.FollowsRepo().Query().
    Filter(
        where.Follows.IsMutual.Equal(true),
        where.Follows.Since.GreaterThan(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
    ).
    All(ctx)
```

## Use Cases

### Social Networks
- Find friends of friends
- Discover mutual connections
- Build follower/following lists

### Access Control
- Check group membership
- Verify permissions via relationships

### Content Systems
- Find related articles
- Track content authorship
- Build recommendation graphs

## Performance Tips

- Index frequently queried edge properties
- Use specific filters to narrow traversals
- Consider query complexity for deep traversals
