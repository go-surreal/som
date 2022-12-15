
[![go 1.19.4](https://img.shields.io/badge/go-1.19.4-informational)](https://go.dev/doc/devel/release)
[![PR](https://github.com/marcbinz/som/actions/workflows/pull_request.yml/badge.svg)](https://github.com/marcbinz/som/actions/workflows/pull_request.yml)

# som

A type-safe database access layer based on code generation for [SurrealDB](https://surrealdb.com).

SurrealDB is a relatively new database approach.
It provides a SQL-style query language with real-time queries and highly-efficient related data retrieval.
Both schemafull and schemaless handling of the data is possible.

With full graph database functionality, SurrealDB enables more advanced querying and analysis.
Records (or vertices) can be connected to one another with edges, each with its own record properties and metadata.
Simple extensions to traditional SQL queries allow for multi-table, multi-depth document retrieval, efficiently 
in the database, without the use of complicated JOINs and without bringing the data down to the client.

(Information extracted from the [official homepage]((https://surrealdb.com)))

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

### Basic usage

Generate the client code:

```
go run github.com/marcbinz/som/cmd/somgen@latest <input_dir> <output_dir>
```

The package `github.com/marcbinz/som` is a dependency for the project in which the generated client code is used.
So it must be added to the `go.mod` file accordingly.

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

- Fully type-safe surrealdb access via generated code.
- Supports most atomic go types: `string`, `int`, `int32`, `int64`, `float32`, `float64`, `bool`
  - Coming soon: `byte`, `[]byte`, `rune`, `uint` ...
- Supports slice values of all atomic types
- Supports pointer fields
- Supports complex types `time.Time` (standard lib) and `uuid.UUID` (google)
  - Maybe future: support any external type with custom encoders and decoders?

## Roadmap

### Before v0.1.0 (first "somewhat stable" non-pre-release)

- [x] Initial implementation.
- [x] Rename project to "som". (#27)
- [x] Add basic GitHub workflow for PR. (#6)
- [ ] Setup golangci-lint with proper config. (#7)
- [ ] Consider reserved (query) keywords. (#18)
- [ ] Add support for pointer fields. (#19)
- [x] Add support for edge (graph) connections. (#20)
- [ ] Add support for `[]byte` (and `byte`?) type.
- [ ] Fix query variable index `rune` not suited for > 26 due to invalid char as variable key.
- [ ] Think about generation based on schema. (#21)
- [ ] How to handle data migrations? (#22)
- [ ] Mark fetched sub-nodes as "invalid to be saved"? (#25)
- [ ] Choose proper licensing for the project. (#11)

### After v0.1.0

- [ ] Make query builder not use pointers, so partial builds and usages are working?
- [ ] Cleanup naming conventions. (#24)
- [ ] Code comments and documentation. (#9)
- [ ] Write tests. (#8)
- [ ] Add "Describe" as query output to get a full description of a generated query. (#17)
- [ ] Generate `sommock` package for easy mocking of the underlying database client.
- [ ] Make casing of database field names configurable.
- [ ] Switch the source code parser to support generics.
- [ ] Add `som.Edge[I, O any]` for defining edges more clearly and without tags (requires generics parser).
- [ ] Support transactions.

### Nice to have (v0.x.x)

- [ ] Add new data type "password". (#16)
- [ ] Add performance benchmarks (and possible optimizations due to it).
- [ ] Integrate external APIs (GraphQL) into the db access layer?

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

### Why are maps not supported?

- With the schemaless database this would be possible.
- Currently, the focus is on structured and deterministic data.
- Might be added in the future though.

## Knowhow

- https://www.appdynamics.com/blog/product/common-application-problems-and-how-to-fix-them-the-select-n-1-problem/
- https://www.sitepoint.com/silver-bullet-n1-problem/
- https://nhibernate.info/doc/howto/various/lazy-loading-eager-loading.html
- https://github.com/mindstand/gogm


## Maintainers & Contributors

- Marc Binz (Author/Owner)

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
