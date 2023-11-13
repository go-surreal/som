<br>

<div align="center">
    <img width="400px" src=".github/branding/logo.png" alt="logo">
    <h3>SOM - A SurrealDB object mapper & query builder for Go</h3>
</div>

<hr />

<p align="center">
  <a href="https://go.dev/doc/devel/release">
    <img src="https://img.shields.io/badge/go-1.21.4-informational" alt="Go 1.21.4">
  </a>
  <a href="https://goreportcard.com/report/github.com/go-surreal/som">
    <img src="https://goreportcard.com/badge/github.com/go-surreal/som" alt="Go Report Card">
  </a>
  <a href="https://github.com/go-surreal/som/actions/workflows/pull_request.yml">
    <img src="https://github.com/go-surreal/som/actions/workflows/pull_request.yml/badge.svg" alt="PR">
  </a>
  <a href="https://discord.gg/surrealdb">
    <img src="https://img.shields.io/discord/902568124350599239?label=discord&color=5a66f6" alt="Discord">
  </a>
  <img src="https://img.shields.io/github/contributors/marcbinz/som" alt="Contributors">
</p>

SOM (SurrealDB object mapper) is an ORM and query builder for [SurrealDB](https://surrealdb.com/) with built-in model
mapping and type-safe query operation generator. It provides an easy and sophisticated database access layer.

## What is SurrealDB?

SurrealDB is a cutting-edge database system that offers a SQL-style query language with real-time queries  
and efficient related data retrieval. It supports both schema-full and schema-less data handling.
With its full graph database functionality, SurrealDB enables advanced querying and analysis by allowing
records (or vertices) to be connected with edges, each with its own properties and metadata.
This facilitates multi-table, multi-depth document retrieval without complex JOINs, all within the database.

*(Information extracted from the [official homepage](https://surrealdb.com))*.

## Table of contents

* [Getting started](#getting-started)
  * [Disclaimer](#disclaimer)
  * [Basic usage](#basic-usage)
    * [Example](#example)
* [Development](#development)
  * [Versioning](#versioning)
  * [Compatibility](#compatibility)
  * [Features](#features)
* [How to contribute](#how-to-contribute)
* [FAQ](#faq)
* [Maintainers & Contributors](#maintainers--contributors)
* [References](#references)

## Getting started

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

## Development

### Versioning

In the future this package will follow the [semantic versioning](https://semver.org) specification.

Up until version 1.0 though, breaking changes might be introduced at any time (minor version bumps).

### Compatibility

This go project makes heavy use of generics. As this feature has been introduced with go 1.18, that version is the 
earliest to be supported by this library.

In general, the two latest (minor) versions of go - and within those, only the latest patch - will be supported 
officially. This means that older versions might still work, but could also break at any time, with any new 
release and without further notice.

Deprecating an "outdated" go version does not yield a new major version of this library. There will be no support for 
older versions whatsoever. This rather hard handling is intended, because it is the official handling for the go 
language itself. For further information, please refer to the
[official documentation](https://go.dev/doc/devel/release#policy) or [endoflife.date](https://endoflife.date/go).

### Features

tbd.

[//]: # (## Roadmap)

[//]: # ()
[//]: # (You can find the official roadmap [here]&#40;ROADMAP.md&#41;. As this might not always be the full)

[//]: # (list of all planned changes, please take a look at the issue section on GitHub as well.)

## How to contribute

### Forking

- to contribute please create a fork within your own namespace
- after making your changes, create a pull request that will merge your state into upstream

### Commits & pull requests

- Commit messages follow the [Conventional Commits](https://www.conventionalcommits.org) specification.
- During development (within a feature branch) commits may have any naming whatsoever.
- After a PR has been approved, it will be squashed and merged into `main`.
- The final commit message must adhere to the specification mentioned above.
- A GitHub Workflow will make sure that the PR title and description matches this specification.

### Labels

- there are some different labels for issues and pull requests (e.g. bug and fix)
- the labels on pull requests resemble the conventional commits specification
- the "highest" label should then be used for the final commit message (e.g. feat above fix, or fix above refactor)
  - TODO: create exhaustive list of label order

## FAQ

You can find a separate document for the FAQs [here](FAQ.md).

## Maintainers & Contributors

Please take a look at the [MAINTAINERS.md](MAINTAINERS.md) file.

## References

- [Official SurrealDB documentation](https://surrealdb.com/docs)
