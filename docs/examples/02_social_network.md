# Social Network Example

This example demonstrates modeling a social network with users, posts, and relationships using SOM's graph capabilities.

## Models

Create `model/models.go`:

```go
package model

import (
    "time"
    "yourproject/gen/som"
)

// Nodes

type User struct {
    som.Node[som.ULID]
    som.Timestamps

    Username string
    Email    string
    Bio      string
    IsActive bool
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

    In  *User `som:"in"`   // Follower
    Out *User `som:"out"`  // Followed

    FollowedAt time.Time
}

type Likes struct {
    som.Edge

    In  *User `som:"in"`   // Who liked
    Out *Post `som:"out"`  // What was liked

    LikedAt time.Time
}
```

## Application Code

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "yourproject/gen/som"
    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "social",
        Database:  "network",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create users
    alice := &model.User{
        Username: "alice",
        Email:    "alice@example.com",
        Bio:      "Software developer",
        IsActive: true,
    }
    bob := &model.User{
        Username: "bob",
        Email:    "bob@example.com",
        Bio:      "Designer",
        IsActive: true,
    }
    charlie := &model.User{
        Username: "charlie",
        Email:    "charlie@example.com",
        Bio:      "Product manager",
        IsActive: true,
    }

    for _, user := range []*model.User{alice, bob, charlie} {
        if err := client.UserRepo().Create(ctx, user); err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Created users")

    // Alice follows Bob and Charlie
    for _, target := range []*model.User{bob, charlie} {
        follows := &model.Follows{FollowedAt: time.Now()}
        err := client.FollowsRepo().Relate().
            From(alice).
            To(target).
            Create(ctx, follows)
        if err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Alice follows Bob and Charlie")

    // Bob follows Alice (mutual)
    follows := &model.Follows{FollowedAt: time.Now()}
    err = client.FollowsRepo().Relate().
        From(bob).
        To(alice).
        Create(ctx, follows)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Bob follows Alice")

    // Alice creates a post
    post := &model.Post{
        Content: "Hello, world!",
        Author:  alice,
    }
    if err := client.PostRepo().Create(ctx, post); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Alice posted: %s\n", post.Content)

    // Bob and Charlie like the post
    for _, user := range []*model.User{bob, charlie} {
        like := &model.Likes{LikedAt: time.Now()}
        err := client.LikesRepo().Relate().
            From(user).
            To(post).
            Create(ctx, like)
        if err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Bob and Charlie liked Alice's post")

    // Query: Who does Alice follow?
    following, err := client.FollowsRepo().Query().
        Where(filter.Follows.In.Equal(alice.ID())).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nAlice follows %d people:\n", len(following))
    for _, f := range following {
        fmt.Printf("  - %s\n", f.Out.Username)
    }

    // Query: Who follows Alice?
    followers, err := client.FollowsRepo().Query().
        Where(filter.Follows.Out.Equal(alice.ID())).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nAlice has %d followers:\n", len(followers))
    for _, f := range followers {
        fmt.Printf("  - %s\n", f.In.Username)
    }

    // Query: How many likes does the post have?
    likeCount, err := client.LikesRepo().Query().
        Where(filter.Likes.Out.Equal(post.ID())).
        Count(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nPost has %d likes\n", likeCount)

    // Query: Recent posts with most active authors
    activePosts, err := client.PostRepo().Query().
        Order(by.Post.CreatedAt.Desc()).
        Limit(10).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nRecent posts: %d\n", len(activePosts))
}
```

## Live Notifications

Subscribe to new followers:

```go
func WatchFollowers(ctx context.Context, client *som.Client, user *model.User) {
    updates, err := client.FollowsRepo().Query().
        Where(filter.Follows.Out.Equal(user.ID())).
        Live(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for update := range updates {
        if update.Error != nil {
            log.Printf("Error: %v", update.Error)
            continue
        }

        switch update.Action {
        case "CREATE":
            fmt.Printf("New follower: %s\n", update.Data.In.Username)
        case "DELETE":
            fmt.Printf("Lost follower: %s\n", update.Data.In.Username)
        }
    }
}
```

## Running the Example

```bash
docker run --rm -p 8000:8000 surrealdb/surrealdb:v3.0.0 \
    start --user root --pass root
go run main.go
```

## Expected Output

```
Created users
Alice follows Bob and Charlie
Bob follows Alice
Alice posted: Hello, world!
Bob and Charlie liked Alice's post

Alice follows 2 people:
  - bob
  - charlie

Alice has 1 followers:
  - bob

Post has 2 likes

Recent posts: 1
```
