# som

A type-safe database access layer with code generation for [SurrealDB](https://surrealdb.com).

## Getting started

### Usage

Generate the client code:

```
go run github.com/marcbinz/som/cmd/somgen@latest <input_dir> <output_dir>
```

The package `github.com/marcbinz/sdb` is a dependency for the project in which the generated client code is used.
So it must be added to the `go.mod` file accordingly.

### Versioning

In the future this package will follow the [semantic versioning](https://semver.org) specification.

Up until version 1.0 though, breaking changes might be introduced at any time (minor version bumps).

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

### Nice to have (v0.x.x)

- [ ] Add new data type "password" with automatic handling of encryption with salt. (#16)
- [ ] Add data type "email" as alias for string that adds database assertion.
- [ ] Add performance benchmarks (and possible optimizations due to it).
- [ ] Integrate external APIs (GraphQL) into the db access layer?

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
