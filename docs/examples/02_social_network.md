# Social Network Example

This example demonstrates modeling a social network with users, posts, and relationships using SOM's graph capabilities.

## Models

Create `model/models.go`:

```go
package model

import (
    "time"
    "github.com/go-surreal/som"
)

// Nodes

type User struct {
    som.Node
    som.Timestamps

    Username string
    Email    string
    Bio      string
    IsActive bool
}

type Post struct {
    som.Node
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
    "yourproject/gen/som/where"
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
        Filter(where.Follows.In.Equal(alice.ID())).
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
        Filter(where.Follows.Out.Equal(alice.ID())).
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
        Filter(where.Likes.Out.Equal(post.ID())).
        Count(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nPost has %d likes\n", likeCount)

    // Query: Recent posts with most active authors
    activePosts, err := client.PostRepo().Query().
        Filter(where.Post.Author.IsActive.IsTrue()).
        Order(by.Post.CreatedAt.Desc()).
        Limit(10).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nRecent posts from active users: %d\n", len(activePosts))
}
```

## Feed Algorithm

Build a simple chronological feed:

```go
func GetFeed(ctx context.Context, client *som.Client, user *model.User, limit int) ([]*model.Post, error) {
    // Get who the user follows
    following, err := client.FollowsRepo().Query().
        Filter(where.Follows.In.Equal(user.ID())).
        All(ctx)
    if err != nil {
        return nil, err
    }

    // Collect followed user IDs
    var followedIDs []*som.ID
    for _, f := range following {
        followedIDs = append(followedIDs, f.Out.ID())
    }

    // No one followed = empty feed
    if len(followedIDs) == 0 {
        return []*model.Post{}, nil
    }

    // Get posts from followed users
    posts, err := client.PostRepo().Query().
        Filter(where.Post.Author.ID().In(followedIDs...)).
        Order(by.Post.CreatedAt.Desc()).
        Limit(limit).
        All(ctx)

    return posts, err
}
```

## Mutual Followers

Find mutual connections:

```go
func GetMutualFollowers(ctx context.Context, client *som.Client, userA, userB *model.User) ([]*model.User, error) {
    // Get userA's followers
    followersA, err := client.FollowsRepo().Query().
        Filter(where.Follows.Out.Equal(userA.ID())).
        All(ctx)
    if err != nil {
        return nil, err
    }

    // For each of A's followers, check if they also follow B
    var mutuals []*model.User
    for _, f := range followersA {
        followsB, err := client.FollowsRepo().Query().
            Filter(
                where.Follows.In.Equal(f.In.ID()),
                where.Follows.Out.Equal(userB.ID()),
            ).
            Exists(ctx)
        if err != nil {
            return nil, err
        }
        if followsB {
            mutuals = append(mutuals, f.In)
        }
    }

    return mutuals, nil
}
```

## Live Notifications

Subscribe to new followers:

```go
func WatchFollowers(ctx context.Context, client *som.Client, user *model.User) {
    updates, err := client.FollowsRepo().Query().
        Filter(where.Follows.Out.Equal(user.ID())).
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
surreal start --user root --pass root memory
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

Recent posts from active users: 1
```
