# SDB

A type-safe database access layer with code generation for [SurrealDB](https://surrealdb.com).

## Getting started

### Usage

Generate the client code:

```
go run github.com/marcbinz/sdb/cmd/sdbgen@latest <input_dir> <output_dir>
```

The package `github.com/marcbinz/sdb` is a dependency for the project in which the generated client code is used.
So it must be added to the `go.mod` file accordingly.

### Versioning

In the future this package will follow the [semantic versioning](https://semver.org) specification.

Up until version 1.0 though, breaking changes might be introduced at any time (minor version bumps).

## Roadmap

### TODOs:

- [ ] generate `sdbmock` package for easy mocking of the underlying database client
- [ ] add performance benchmarks (and possible optimizations due to it)
- [ ] make casing of database fields configurable?
- [ ] add support for maps as node fields?
- [ ] add support for `[]byte` (and `byte`?) type
- [ ] Replace `strcase.ToLowerCamel` with `strcase.ToLowerCamel`
- [ ] Switch the source code parsing to support generics
- [ ] add `som.Edge[I, O any]` for defining edges more clearly and without tags (requires generics parser)
- [ ] Make query builder not use pointers, so partial builds and usages are working
- [ ] Fix query variable index `rune` not suited for > 26 due to invalid char as variable key?

## How to contribute

### Commits & pull requests

- Commit messages follow the [Conventional Commits](https://www.conventionalcommits.org) specification.
- During development (within a feature branch) commits may have any naming whatsoever.
- After a PR has been approved, it will be squashed and merged into `main`.
- The final commit message must adhere to the specification mentioned above.
- A GitHub Workflow will make sure that the PR title and description matches this specification.

## FAQ

### Why are maps not supported?

- With the schemaless database this would be possible.
- Currently, the focus is on structured and deterministic data.
- Might be added in the future though.

## Maintainers & Contributors

- Marc Binz (Author/Owner)

## References

- https://surrealdb.com/docs
- https://entgo.io
- https://github.com/d-tsuji/awesome-go-orms
- https://github.com/doug-martin/goqu
- https://github.com/sharovik/orm
