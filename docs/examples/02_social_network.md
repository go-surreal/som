# Social Network Example

This example models a social network with users, posts and relationships using SOM's graph
capabilities.

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

    // Edges starting at this user
    Follows []Follows
    Likes   []Likes
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

    Follower User `som:"in"`
    Followed User `som:"out"`

    FollowedAt time.Time
}

type Likes struct {
    som.Edge

    User User `som:"in"`
    Post Post `som:"out"`

    LikedAt time.Time
}
```

Both edges are declared as fields on `User`, which is what generates the `Relate()` accessors and
the traversal filters.

## Application Code

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "yourproject/gen/som/by"
    "yourproject/gen/som/filter"
    "yourproject/gen/som/repo"
    "yourproject/gen/som/with"
    "yourproject/model"
)

func main() {
    ctx := context.Background()

    client, err := repo.NewClient(ctx, repo.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "social",
        Database:  "network",
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    if err := client.ApplySchema(ctx); err != nil {
        log.Fatal(err)
    }

    // Create users
    alice := &model.User{Username: "alice", Email: "alice@example.com", IsActive: true}
    bob := &model.User{Username: "bob", Email: "bob@example.com", IsActive: true}
    charlie := &model.User{Username: "charlie", Email: "charlie@example.com", IsActive: true}

    for _, user := range []*model.User{alice, bob, charlie} {
        if err := client.UserRepo().Create(ctx, user); err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Created users")

    // Alice follows Bob and Charlie
    for _, target := range []*model.User{bob, charlie} {
        follows := &model.Follows{
            Follower:   *alice,
            Followed:   *target,
            FollowedAt: time.Now(),
        }
        if err := client.UserRepo().Relate().Follows().Create(ctx, follows); err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Alice follows Bob and Charlie")

    // Bob follows Alice (mutual)
    if err := client.UserRepo().Relate().Follows().Create(ctx, &model.Follows{
        Follower:   *bob,
        Followed:   *alice,
        FollowedAt: time.Now(),
    }); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Bob follows Alice")

    // Alice creates a post
    post := &model.Post{Content: "Hello, world!", Author: alice}
    if err := client.PostRepo().Create(ctx, post); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Alice posted: %s\n", post.Content)

    // Bob and Charlie like the post
    for _, user := range []*model.User{bob, charlie} {
        like := &model.Likes{User: *user, Post: *post, LikedAt: time.Now()}
        if err := client.UserRepo().Relate().Likes().Create(ctx, like); err != nil {
            log.Fatal(err)
        }
    }
    fmt.Println("Bob and Charlie liked Alice's post")

    // Query: who does Alice follow?
    following, err := client.UserRepo().Query().
        Where(
            filter.User.Follows().Followed(
                filter.User.Username.Equal(alice.Username),
            ),
        ).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("\nUsers followed by someone: %d\n", len(following))

    // Query: who follows Alice? (users that have a follows edge to alice)
    followers, err := client.UserRepo().Query().
        Where(
            filter.User.Follows().Followed(
                filter.User.ID.Equal(string(alice.ID())),
            ),
        ).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Alice has %d followers\n", len(followers))

    // Query: how many users liked the post?
    likeCount, err := client.UserRepo().Query().
        Where(
            filter.User.Likes().Post(
                filter.Post.ID.Equal(string(post.ID())),
            ),
        ).
        Count(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Post has %d likes\n", likeCount)

    // Query: recent posts with their authors
    recentPosts, err := client.PostRepo().Query().
        Fetch(with.Post.Author()).
        Order(by.Post.CreatedAt.Desc()).
        Limit(10).
        All(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Recent posts: %d\n", len(recentPosts))
}
```

## Live Notifications

Watch for users gaining followers in real time:

```go
func WatchFollowers(ctx context.Context, client repo.Client, user *model.User) {
    updates, err := client.UserRepo().Query().
        Where(
            filter.User.Follows().Followed(
                filter.User.ID.Equal(string(user.ID())),
            ),
        ).
        Live(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for update := range updates {
        switch res := update.(type) {
        case query.LiveCreate[*model.User]:
            follower, err := res.Get()
            if err != nil {
                log.Printf("error: %v", err)
                continue
            }
            fmt.Printf("New follower: %s\n", follower.Username)

        case query.LiveDelete[*model.User]:
            follower, _ := res.Get()
            fmt.Printf("Lost follower: %s\n", follower.Username)

        case query.LiveKilled[*model.User]:
            return
        }
    }
}
```

## Running the Example

```bash
docker run --rm -p 8000:8000 surrealdb/surrealdb:v3.2.0 \
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

Users followed by someone: 2
Alice has 1 followers
Post has 2 likes
Recent posts: 1
```
