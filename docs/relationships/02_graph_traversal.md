# Graph Traversal

SurrealDB's graph capabilities allow you to traverse relationships between records. SOM provides type-safe access to these features.

## Querying Edges

Find all relationships of a specific type:

```go
// Find all "follows" relationships
follows, err := client.FollowsRepo().Query().All(ctx)
```

## Filtering by Source/Target

### By Source (In)

Find relationships from a specific node:

```go
// Find everyone a user follows
following, err := client.FollowsRepo().Query().
    Where(filter.Follows.In.Equal(user.ID())).
    All(ctx)

// Get the followed users
for _, f := range following {
    fmt.Println("Follows:", f.Out.Name)
}
```

### By Target (Out)

Find relationships to a specific node:

```go
// Find everyone who follows a user
followers, err := client.FollowsRepo().Query().
    Where(filter.Follows.Out.Equal(user.ID())).
    All(ctx)

// Get the follower users
for _, f := range followers {
    fmt.Println("Follower:", f.In.Name)
}
```

## Edge Properties

Query edges by their properties:

```go
// Find recent follows (last 30 days)
recentFollows, err := client.FollowsRepo().Query().
    Where(
        filter.Follows.Since.After(time.Now().AddDate(0, 0, -30)),
    ).
    All(ctx)

// Find mutual follows
mutualFollows, err := client.FollowsRepo().Query().
    Where(filter.Follows.IsMutual.IsTrue()).
    All(ctx)
```

## Combining Filters

```go
// Find mutual follows from a specific user created this year
follows, err := client.FollowsRepo().Query().
    Where(
        filter.Follows.In.Equal(user.ID()),
        filter.Follows.IsMutual.IsTrue(),
        filter.Follows.Since.After(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
    ).
    All(ctx)
```

## Edge CRUD Operations

Edges support the same CRUD operations as nodes:

### Create

```go
follows := &model.Follows{
    Since:    time.Now(),
    IsMutual: false,
}

err := client.FollowsRepo().Relate().
    From(alice).
    To(bob).
    Create(ctx, follows)
```

### Read

```go
follows, exists, err := client.FollowsRepo().Read(ctx, edgeID)
if exists {
    fmt.Printf("%s follows %s since %v\n",
        follows.In.Name, follows.Out.Name, follows.Since)
}
```

### Update

```go
follows.IsMutual = true
err := client.FollowsRepo().Update(ctx, follows)
```

### Delete

```go
err := client.FollowsRepo().Delete(ctx, follows)
```

## Counting Relationships

```go
// Count followers
followerCount, err := client.FollowsRepo().Query().
    Where(filter.Follows.Out.Equal(user.ID())).
    Count(ctx)

// Count following
followingCount, err := client.FollowsRepo().Query().
    Where(filter.Follows.In.Equal(user.ID())).
    Count(ctx)

fmt.Printf("%d followers, following %d\n", followerCount, followingCount)
```

## Checking Relationship Existence

```go
// Check if Alice follows Bob
isFollowing, err := client.FollowsRepo().Query().
    Where(
        filter.Follows.In.Equal(alice.ID()),
        filter.Follows.Out.Equal(bob.ID()),
    ).
    Exists(ctx)

if isFollowing {
    fmt.Println("Alice follows Bob")
}
```

## Live Edge Subscriptions

Subscribe to relationship changes in real-time:

```go
updates, err := client.FollowsRepo().Query().
    Where(filter.Follows.Out.Equal(user.ID())).
    Live(ctx)

for update := range updates {
    if update.Error != nil {
        continue
    }

    switch update.Action {
    case "CREATE":
        fmt.Printf("New follower: %s\n", update.Data.In.Name)
    case "DELETE":
        fmt.Printf("Lost follower: %s\n", update.Data.In.Name)
    }
}
```

## Use Cases

### Social Networks

- Find friends of friends
- Discover mutual connections
- Build follower/following lists
- Recommend connections

```go
// Find who Alice follows that also follows Bob
aliceFollows, _ := client.FollowsRepo().Query().
    Where(filter.Follows.In.Equal(alice.ID())).
    All(ctx)

for _, f := range aliceFollows {
    // Check if this person follows Bob
    follows, _ := client.FollowsRepo().Query().
        Where(
            filter.Follows.In.Equal(f.Out.ID()),
            filter.Follows.Out.Equal(bob.ID()),
        ).
        Exists(ctx)

    if follows {
        fmt.Printf("%s follows both Alice and Bob\n", f.Out.Name)
    }
}
```

### Access Control

- Check group membership
- Verify permissions via relationships
- Build role hierarchies

```go
// Check if user is member of group
isMember, err := client.MemberOfRepo().Query().
    Where(
        filter.MemberOf.In.Equal(user.ID()),
        filter.MemberOf.Out.Equal(group.ID()),
    ).
    Exists(ctx)

// Check if user is admin of group
isAdmin, err := client.MemberOfRepo().Query().
    Where(
        filter.MemberOf.In.Equal(user.ID()),
        filter.MemberOf.Out.Equal(group.ID()),
        filter.MemberOf.Role.Equal("admin"),
    ).
    Exists(ctx)
```

### Content Systems

- Find related articles
- Track content authorship
- Build recommendation graphs

```go
// Find all posts liked by a user
likes, err := client.LikesRepo().Query().
    Where(filter.Likes.In.Equal(user.ID())).
    All(ctx)

for _, like := range likes {
    fmt.Printf("Liked: %s at %v\n", like.Out.Title, like.LikedAt)
}
```

## Ordering Edges

```go
// Most recent followers first
followers, err := client.FollowsRepo().Query().
    Where(filter.Follows.Out.Equal(user.ID())).
    Order(by.Follows.Since.Desc()).
    Limit(10).
    All(ctx)
```

## Performance Tips

- Index frequently queried edge properties
- Use specific filters to narrow traversals
- Consider query complexity for deep traversals
- Use `Limit` for large result sets
- Use `Count` and `Exists` instead of fetching full records when possible
