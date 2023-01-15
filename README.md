
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

### Before v0.1.0 (first "somewhat stable" non-pre-release)

- [ ] Implement sub-queries for node and edge types.
- [ ] Add `som.SoftDelete` type with `DeletedAt` timestamp and automated handling throughout som.
- [ ] Mark fetched sub-nodes as "invalid to be saved"? (#25)
- [ ] Consider reserved (query) keywords. (#18)
- [ ] Check for possible security vulnerabilities.
- [ ] Choose proper licensing for the project. (#11)

### After v0.1.0

- [ ] Provide `WithInfo` method.
- [ ] Add support for `[]byte` (and `byte`?) type.
- [ ] How to handle data migrations? (#22)
- [ ] Setup golangci-lint with proper config. (#7)
- [ ] Support (deeply) nested slices? (needed?)
- [ ] Cleanup naming conventions. (#24)
- [ ] Code comments and documentation. (#9)
- [ ] Write tests. (#8)
- [ ] Generate `sommock` package for easy mocking of the underlying database client.
- [ ] Make casing of database field names configurable.
- [ ] Switch the source code parser to support generics.
- [ ] Add `som.Edge[I, O any]` for defining edges more clearly and without tags (requires generics parser).
- [ ] Support transactions.
- [ ] Distinct results (https://stackoverflow.com/questions/74326176/surrealdb-equivalent-of-select-distinct).
- [ ] Integrate external APIs (GraphQL) into the db access layer?
- [ ] Support (deeply) nested slices? (needed?)
- [ ] Unique relations (`DEFINE INDEX __som_unique_relation ON TABLE member_of COLUMNS in, out UNIQUE;`)

### Nice to have (v0.x.x)?

- [ ] Add new data type "password" with automatic handling of encryption with salt. (#16)
- [ ] Add data type "email" as alias for string that adds database assertion.
  - Or provide an API to add custom assertions for types (especially string).
- [ ] Add performance benchmarks (and possible optimizations due to it).

```sql
DEFINE TABLE user SCHEMAFULL 
        PERMISSIONS NONE;
DEFINE FIELD username ON TABLE user
        TYPE string
        ASSERT string::length($value) >= 4
        ASSERT string::length($value) <= 8;
DEFINE FIELD password ON TABLE user
        PERMISSIONS
                FOR SELECT NONE
        TYPE string;
DEFINE FIELD email ON TABLE user
        TYPE string
        ASSERT is::email($value);
DEFINE FIELD num ON TABLE user
        VALUE 42;
```

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

*Disclaimer: Currently those are just questions I asked myself and wanted an answer before the initial public release.
In the future this section will be advanced by topics raised in issues or discussions.*

### Why are maps not supported?

- With the schemaless database this would be possible.
- Currently, the focus is on structured and deterministic data.
- Might be added in the future though.

### Why does a filter like `where.User.Equal(userModel)` not exist?

- This would be an ambiguous case. Should it compare the whole object with all properties or only the ID?
- For this case it is better and more deterministic to just compare the ID explicitly.
- If - for whatever reason - it is required to check the fields, adding those filters one by one makes the purpose of the query clearer.
- Furthermore, we would need to find a way to circumvent a naming clash when the field of a model is named `Equal` (or other keywords).
- On the other hand, this feature is still open for debate. So if anyone can clarify the need for it, we might as well implement it at some point.

## Maintainers & Contributors

- Marc Binz (Author/Owner)

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
