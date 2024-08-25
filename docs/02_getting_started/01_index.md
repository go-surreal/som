# Getting started


*Please note: This package is currently tested against version
[1.0.0](https://surrealdb.com/releases#v1-0-0)
of SurrealDB.*

### Disclaimer

This library is currently considered **HIGHLY EXPERIMENTAL** and under heavy development.

Som is in the stage of (very) early development. While the basic functionality should be working as expected,
there could be unknown and critical bugs. This could theoretically lead to your database and especially
the data within it to get broken. It is not yet recommended for production use, so please be careful!

If you still want to use it for production applications, we do not accept any liability for data loss
or other kind of issues. You have been warned!

Furthermore, it may lack specific features and edge cases you might expect from a fully-featured ORM and
query builder library. Feel free to submit feature requests or pull requests to add additional functionality
to Som. We do ask you to please read our [Contributor Guide](CONTRIBUTING.md).

As long as this library is not considered stable, there will likely be significant and breaking changes to its API.
So please keep in mind, that any new updates might lead to required (manual) migrations or to take similar actions.

But still, please try it out and give us some feedback. We would highly appreciate it. Thank you! 🙏

### Basic usage

Generate the client code:

```
go run github.com/go-surreal/som/cmd/somgen@latest <input_dir> <output_dir>
```

Currently, the generated code does not make use of the official SurrealDB go client.
Instead, it is using a custom implementation called [sdbc](https://github.com/go-surreal/sdbc).
Until the official client is considered stable, this will likely not change.
Final goal would be to make it possible to use both the official client and the custom implementation.
As of now, this might change at any time.

#### Example

Let's say we have the following model at `<root>/model/user.go`:

```go
package model

type User struct {
    ID       string `som:"id"`
    Username string `som:"username"`
    Password string `som:"password"`
    Email    string `som:"email"`
}
```

In order for it to be considered by the generator, it must embed `som.Node`:

```go
package model

import "github.com/go-surreal/som"

type User struct {
    som.Node
    
    // ID string `som:"id"` --> provided by som!
    
    Username string `som:"username"`
    Password string `som:"password"`
    Email    string `som:"email"`
}
```

Now, we can generate the client code:

```
go run github.com/go-surreal/som/cmd/somgen@latest <root>/model <root>/gen/som
```

With the generated client, we can now perform operations on the database:

```go
package main

import (
    "context"
    "log"
    
    "<root>/gen/som"
    "<root>/gen/som/where"
    "<root>/model"
)

func main() {
    ctx := context.Background()

    // create a new client
    client, err := som.NewClient(ctx, som.Config{
        Address:   "ws://localhost:8000",
        Username:  "root",
        Password:  "root",
        Namespace: "test",
        Database:  "test",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    // initialize the model
    user := &model.User{
        Username: "test",
        Password: "test",
        Email:    "test@example.com",
    }
    
    // insert the user into the database
    err = client.UserRepo().Create(ctx, user)
    if err != nil {
        log.Fatal(err)
    }
		
    // query the user by email
    read, err := client.UserRepo().Query().
        Filter(
            where.User.Email.Equal("test@example.com"),
        ).
        First(ctx)

    if err != nil {
        log.Fatal(err)
    }
		
    fmt.Println(read)
}
```
