<br>

<div align="center">
    <img width="400px" src=".github/branding/logo.png" alt="logo">
    <h3>SOM - A SurrealDB object mapper & query builder for Go</h3>
</div>

<hr />

<p align="center">
  <a href="https://go.dev/doc/devel/release">
    <img src="https://img.shields.io/badge/go-1.20.2-informational" alt="Go 1.20.2"> 
  </a>
  <a href="https://github.com/marcbinz/som/actions/workflows/pull_request.yml">
    <img src="https://github.com/marcbinz/som/actions/workflows/pull_request.yml/badge.svg" alt="PR">
  </a>
  <a href="https://discord.gg/surrealdb">
    <img src="https://img.shields.io/discord/902568124350599239?label=discord&color=5a66f6" alt="Discord">
  </a>
  <img src="https://img.shields.io/github/contributors/marcbinz/som" alt="Contributors">
</p>

SOM (SurrealDB object mapper) is an ORM and query builder for [SurrealDB](https://surrealdb.com/) with built-in model
mapping and type-safe query operation generator. It provides an easy and sophisticated database access layer.

## What is SurrealDB?

SurrealDB is a relatively new database approach.
It provides a SQL-style query language with real-time queries and highly-efficient related data retrieval.
Both schemafull and schemaless handling of the data is possible.

With full graph database functionality, SurrealDB enables more advanced querying and analysis.
Records (or vertices) can be connected to one another with edges, each with its own record properties and metadata.
Simple extensions to traditional SQL queries allow for multi-table, multi-depth document retrieval, efficiently 
in the database, without the use of complicated JOINs and without bringing the data down to the client.

*(Information extracted from the [official homepage]((https://surrealdb.com)))*

## Table of contents

* [Getting started](#getting-started)
  * [Basic usage](#basic-usage)
  * [Versioning](#versioning)
  * [Compatibility](#compatibility)
  * [Features](#features)
* [Roadmap](#roadmap)
* [How to contribute](#how-to-contribute)
* [Maintainers & Contributors](#maintainers--contributors)
* [References](#references)

## Getting started

*Please note: This package is currently tested against version 
[1.0.0-beta.8](https://github.com/surrealdb/surrealdb/releases/tag/v1.0.0-beta.8)
of SurrealDB.*

### Disclaimer

This library is currently considered **HIGHLY EXPERIMENTAL**.

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
go run github.com/marcbinz/som/cmd/somgen@latest <input_dir> <output_dir>
```

The package `github.com/marcbinz/som` can be considered an invisible dependency for your project. All it does is to
generate code that lives within your project, but the package itself does not need to be added to the `go.mod` file.

The generated code uses the official SurrealDB go client:

```
go get github.com/surrealdb/surrealdb.go
```

### Versioning

In the future this package will follow the [semantic versioning](https://semver.org) specification.

Up until version 1.0 though, breaking changes might be introduced at any time (minor version bumps).

### Compatibility

This go project makes heavy use of generics. As this feature has been introduced with go 1.18, that version is the 
earliest to be supported by this library.

In general, the two latest (minor) versions of go - and within those, only the latest patch - will be supported 
officially. This means that older versions might still work, but could also break at any time and with any new release.

Deprecating an "outdated" go version does not yield a new major version of this library. There will be no support for 
older versions whatsoever. This rather hard handling is intended, because it is the official handling for the go 
language itself. For further information, please refer to the
[official documentation](https://go.dev/doc/devel/release#policy) or [endoflife.date](https://endoflife.date/go).

### Features

- Fully type-safe SurrealDB access via generated code.
- Supports most atomic go types: `string`, `int`, `int32`, `int64`, `float32`, `float64`, `bool`
  - Coming soon: `byte`, `[]byte`, `rune`, `uint` ...
- Supports slice values of all atomic types.
- Supports pointer fields.
- Supports complex types `time.Time` (standard lib) and `uuid.UUID` (google)
  - Maybe future: support any external type with custom encoders and decoders?
- Supports record links (references to other nodes/models).
- Supports graph connections (edges) between nodes/models.

## Roadmap

You can find the official roadmap [here](ROADMAP.md). As this might not always be the full
list of all planned changes, please take a look at the issue section on GitHub as well.

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

You can find a separate document for the FAQ [here](FAQ.md).

## Maintainers & Contributors

- Marc Binz (Author/Owner)

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
- https://github.com/StarlaneStudios/cirql
- https://github.com/uptrace/bun
- https://atlasgo.io/
